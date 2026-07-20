# bklg — backlog board viewer — spec

**Status:** draft, for review
**Purpose:** serve a live 3-column board of a framework knowledge instance's backlog
(CURRENT / BACKLOG / DONE) with badges for cross-cutting state (blocked, parking-lot,
override), plus a detail page per task.
**Delivery target:** standalone Go module, `go install`-able as `bklg`; **Go 1.22+**;
**zero Go module dependencies** (stdlib only); Tailwind via Play CDN (the only external
artifact, and it is client-side).
**Reads:** one framework instance's *planning* and *progress* areas, resolved through the
manifest's Locations block; falls back to `<dir>/planning` + `<dir>/progress`.
**One line:** point it at a repo, get a localhost kanban with blocked badges and a page per task.

---

## 1. What it is

A read-only viewer. It parses the semi-structured markdown a framework instance already
maintains (`CURRENT.md`, `BACKLOG.md`, `DONE.md`, `blockers.md`) and renders it as a kanban
board plus per-task pages. It writes nothing, mutates no VCS state, and re-reads on every
request so the board tracks a live agent session. It is generic — it works against any repo
that follows the framework file conventions, so it belongs to no single product and ships as
its own tool (`github.com/<you>/bklg`).

Scope is deliberately narrow: **one planning area at a time.** The monorepo-root case (many
systems, no root planning) is handled with a helpful error + system list in v1, and a real
multi-system picker in v2 (§13).

---

## 2. Input contract (the crux)

Correctness is bounded by how faithfully the target repo follows the skeletons in
`framework/SETUP.md`. Those files are **prose-with-conventions, not a schema.** Two
consequences drive the whole design:

1. **Parse defensively.** A malformed entry must never crash the server or blank the board.
   The parser captures what it can, skips what it can't, and records a warning. Warnings
   surface in a subtle banner linking to `/_diag` — honest about what it couldn't read.
2. **The patterns below *are* the contract.** A repo that diverges (`* [ ]` instead of
   `- [ ]`, a `DONE` line missing its ` — ` separators) degrades to best-effort and shows up
   in `/_diag`, rather than silently corrupting the board.

The parser is line-oriented and keyed to these exact conventions (all em-dashes are U+2014;
IDs match `[A-Z]+-\d+`; dates match `\d{4}-\d{2}-\d{2}`):

| File (in planning area) | Section parsed | Ignored | Item shape | Fields extracted |
|---|---|---|---|---|
| `CURRENT.md` | everything under `## Active` | `## Entry template` and prose above it | `### <ID> — <title>` | `**Source:**`, `**Acceptance criteria:**` → following `- [ ]` / `- [x]` list, `**Notes:**`, `**Delivery override:**` |
| `BACKLOG.md` | `## Active`, `## Parking lot` | header/conventions blockquote | `- [ ] <ID> — <title>` / `- [x] …` (parking-lot items may be freeform, id optional) | id (optional), title, checked-state, which section |
| `DONE.md` | `## Completed` | template prose | `- <ID> — <title> — <date> — <summary> — <delivery record>. See journal …` | id, title, date, summary, delivery record, journal pointer |

Plus the **progress area**:

| File (in progress area) | Section parsed | Ignored | Item shape | Fields extracted |
|---|---|---|---|---|
| `blockers.md` | `## Open`, `## Resolved` | `## Format` example | `## <BLOCKER-ID> — <title> — opened <ts>` | id, title, opened, `**Task affected:**` (the join key), body, status = open/resolved |

**Ambiguity the parser must resolve:** in `blockers.md`, blocker headings and section headings
are *both* `##`. Distinguish by content — a blocker heading matches `^BLOCKER-\d+`; a section
heading is one of `Format` / `Open` / `Resolved`. Track the current section as the last section
heading seen; assign each blocker to it. Skip everything under `## Format`.

**Do not render text as HTML.** Captured field text is displayed escaped (via `html/template`
auto-escaping), with line breaks preserved. This keeps the tool stdlib-only and immune to
injection from repo content. Rendering a safe markdown subset is a v2 opt-in (needs a
sanitizer; §13).

---

## 3. Resolution & discovery

```
path := <positional arg, default ".">
dir  := <--dir, default "knowledge">
base := path/dir                         # the knowledge dir
```

Area resolution, in order:

1. If `base/README.md` exists and contains a `## Locations` block with `planning:` /
   `progress:` keys, dereference them. Locations paths are **repo-root-relative**, so resolve
   them against `path`: `planningDir = path/<locations.planning>`. (This is the framework-native
   path, and it composes for a single system: `bklg . --dir systems/trail/knowledge` →
   Locations `planning: systems/trail/knowledge/planning` → resolved against repo root.)
2. Else default: `planningDir = base/planning`, `progressDir = base/progress`.
3. If the resolved `planningDir` does not exist, this is probably a **multi-system root
   manifest**. Parse the system-index table (rows containing `systems/<name>`); if systems are
   found, exit with a message listing them and the exact `bklg` invocation for each. If none are
   found, exit: `no planning area at <planningDir>`.

Locations-block parse: on the line `## Locations`, enter the block; until the next line
beginning `## `, split each non-empty line on the first `:` into key/value and trim. Keep
`planning` and `progress`; ignore the rest.

System-index parse: scan `|`-delimited table rows for a cell matching `systems/\S+`; collect the
distinct directory names.

---

## 4. Data model

```go
type ID struct { Namespace string; Number int; Raw string }   // {"MONO", 6, "MONO-006"}

type Column int
const ( ColBacklog Column = iota; ColInProgress; ColDone )

type Criterion struct { Text string; Checked bool }
type DoneRecord struct { Date, Summary, DeliveryRecord, JournalPointer string }

type Card struct {
    ID       *ID          // nil for parking-lot items without an id
    Title    string
    Column   Column
    Badges   []Badge

    Source           string       // from CURRENT
    Notes            string
    DeliveryOverride string
    Acceptance       []Criterion  // from CURRENT
    Done             *DoneRecord  // from DONE
    ParkingLot       bool
    Raw              string       // the source block, shown on the detail page
}

type Blocker struct {
    ID      ID           // BLOCKER-003
    Title, Opened, Body string
    TaskRaw string       // affected task id, raw ("TRAIL-007")
    Open    bool
}

type Board struct {
    Cards    []Card
    Blockers []Blocker
    Warnings []string    // parse diagnostics -> /_diag
    Meta     struct{ KnowledgeDir, PlanningDir, ProgressDir string; LatestMTime time.Time }
}
```

Everything on a card's detail page derives from joining that card's id across the parsed files;
`Card.Raw` keeps the original block so the detail page can show source verbatim.

---

## 5. Column mapping & reconciliation

| Column | Source | Order |
|---|---|---|
| **Backlog** | `BACKLOG.md ## Active` unchecked `[ ]` + all `## Parking lot` items | file order (top = next); parking-lot grouped below, muted |
| **In Progress** | `CURRENT.md ## Active` tasks (framework: exactly one; tool handles 0..N) | as written |
| **Done** | `DONE.md ## Completed` entries | file order (newest at top) |

A task id can legitimately appear in more than one file (an in-progress task is usually still
listed unchecked in `BACKLOG`; a shipped task is checked in `BACKLOG` *and* present in `DONE`).
**Dedup rule — most-advanced state wins:** `Done > InProgress > Backlog`. A given id renders as
exactly one card, in its furthest column. Corollary hygiene checks, each emitted as a warning
(not a hard error):

- `[x]` in `BACKLOG` with **no** matching `DONE` entry → still shown in Done, warning "shipped
  item missing from DONE.md".
- `DONE` entry whose id is **not** checked in `BACKLOG` → warning "DONE item not ticked in
  BACKLOG".
- More than one task in `CURRENT ## Active` → warning "CURRENT holds >1 active task" (the
  framework's one-task invariant).

Parking-lot items without ids never dedup (they live only in `BACKLOG`) and get no detail page.

---

## 6. Badges

Derived by joining planning + progress. Badges are the whole reason the board beats reading the
files raw.

| Badge | Condition | Suggested style |
|---|---|---|
| `blocked` | id is the `Task affected` of an **open** blocker | red |
| `parking` | item under `## Parking lot` | slate / muted |
| `override` | card has a `**Delivery override:**` field | amber |
| `no-ac` *(opt)* | In-Progress card with zero acceptance criteria (framework hygiene miss) | yellow |
| namespace tag | the id's namespace (`MONO`, `TRAIL`, …), rendered as a neutral chip | gray |

The namespace chip is cheap and useful in mixed-namespace planning areas (trail's holds both
`TRAIL-` and `MONO-`). `blocked` is the load-bearing one the user called out.

---

## 7. Routes & rendering

Go 1.22 `ServeMux` pattern routing; literal segments beat the `{id}` wildcard, so `/_v` and
`/_diag` win over `/{id}` without ordering tricks:

```go
mux.HandleFunc("GET /{$}",    boardHandler)   // exact "/"  -> the board
mux.HandleFunc("GET /_v",     versionHandler) // latest mtime, for live reload
mux.HandleFunc("GET /_diag",  diagHandler)    // parse warnings
mux.HandleFunc("GET /{id}",   taskHandler)    // /MONO-006  -> detail; unknown id -> 404
```

- **Board (`/`)** — three columns, a card per task. Card shows id + title + badges. Cards with
  an id link to their detail page. Empty columns render a muted placeholder. Above the board: a
  header line with the resolved planning path and, when `len(Warnings)>0`, the diagnostics banner.
- **Detail (`/<id>`)** — id, title, namespace, current column, badges; then, by state: In-Progress
  → source, acceptance criteria (as a checklist), notes, override; Done → date, summary, delivery
  record, journal pointer (plain text — the journal isn't served in v1). Always: any blockers
  referencing the id (open first, then resolved), and a collapsed "source" block showing `Raw`.
  Unknown id → `404`.
- **`/_v`** — returns the max mtime across the parsed files as a bare integer. The page polls it
  every ~3s and calls `location.reload()` on change (≈8 lines of vanilla JS), so the board stays
  live during a session without full-page meta-refresh churn.
- **`/_diag`** — the parse warnings, verbatim, one per line.

**Freshness:** parse per request. The files are tiny; correctness-while-editing beats caching.
(An mtime-keyed cache is a trivial later optimization; not needed for v1.)

**Templates:** `html/template`, embedded with `go:embed`. A `layout` partial + `board` + `task`.
No frontend framework. No client state beyond the reload poll.

---

## 8. Styling — the one external-dependency fork

The user wants Tailwind, no frontend framework, "all served by Go," and *simple*. Those pull in
slightly different directions:

- **Play CDN** (`<script src="https://cdn.tailwindcss.com">`): zero build, **no Node toolchain
  ever**, `go install` just works. Cost: the CSS/JS comes from the CDN (so not *strictly* "all
  served by Go") and needs network at view time. For a localhost dev tool this is a non-issue.
- **Prebuilt + embedded**: run `tailwindcss` CLI once, `go:embed` the generated `styles.css`,
  serve it from Go. Fully self-contained binary, offline, everything Go-served. Cost: adds a Node
  build step to the tool's own release flow.

**Recommendation for v1: Play CDN.** It is the only option that keeps the tool a pure-stdlib,
Node-free `go install`. The embedded-CSS path is documented as the "self-contained release"
upgrade (§13) for whoever wants an offline single binary. *(This is decision D1 in §12 — flip it
if "all served by Go" is a hard line for you.)*

---

## 9. CLI

```
bklg [path] [--port N] [--dir D]

  [path]    optional path to the repo root         (default ".")
  --port N  port to listen on                        (default 1235)
  --dir  D  knowledge dir, relative to [path]        (default "knowledge")
```

Startup output — first line matches the requested format, then the resolution is echoed so the
user can see what it's reading (honest + debuggable):

```
$ bklg
Running Backlog on port 1235
  knowledge: ./knowledge   planning: ./knowledge/planning   progress: ./knowledge/progress
  http://localhost:1235
```

**Arg ordering:** stdlib `flag` stops at the first non-flag token, which would break
`bklg systems/trail --port 9000`. Rather than take a dependency (pflag) for two flags, pre-split
argv so the positional works in any position and both `--port 9000` and `--port=9000` parse.
~15 lines, zero deps:

```go
func splitArgs(argv []string) (path string, flagArgs []string) {
    path = "."
    takesValue := map[string]bool{"--port": true, "-port": true, "--dir": true, "-dir": true}
    seenPath := false
    for i := 0; i < len(argv); i++ {
        a := argv[i]
        if strings.HasPrefix(a, "-") {
            flagArgs = append(flagArgs, a)
            if takesValue[a] && !strings.Contains(a, "=") && i+1 < len(argv) {
                i++; flagArgs = append(flagArgs, argv[i]) // consume its value
            }
            continue
        }
        if !seenPath { path, seenPath = a, true } else { flagArgs = append(flagArgs, a) }
    }
    return
}
// main: p, fa := splitArgs(os.Args[1:]); fs.Parse(fa)
```

**Bind to `127.0.0.1`** by default, not `0.0.0.0` — it's a personal dev tool, keep it off the
network. The `/<id>` route is looked up in the parsed model, never used as a filesystem path, so
there's no traversal surface.

**Failure modes** (exit non-zero, clear message, no server): planning area not found (→ list
systems if it's a root manifest, §3); port in use; `path` not a directory.

---

## 10. Suggested layout (the tool's own code)

```
cmd/bklg/main.go        CLI: splitArgs, flags, resolution, server bootstrap
internal/backlog/
  model.go              types from §4
  parse.go              the line-oriented parser (§2) — behind a small interface
  resolve.go            Locations + system-index resolution (§3)
  server.go             handlers + routes (§7)
templates/*.html        layout, board, task  (go:embed)
testdata/knowledge/…    fixtures exercising every case (see §11)
```

Flat `package main` with a few files is also fine — the tool is small. The one seam worth keeping
is `parse.go` behind an interface, so a goldmark-backed parser can replace the line-scanner later
if real-world files prove messier than the skeletons (decision D4, §12).

---

## 11. Test fixtures

Ship a `testdata/knowledge/` instance whose planning + progress files exercise, at minimum: an
unchecked backlog item with an id; one without an id (parking lot); a checked backlog item that
*is* in DONE and one that *isn't*; a CURRENT task with acceptance criteria + a delivery override;
a CURRENT task with **no** criteria (`no-ac`); an open blocker naming a CURRENT task; a resolved
blocker. This doubles as parser test input (assert on counts, ids, states, badge sets — quoted in
the journal) and as the live demo the board renders.

---

## 12. Open decisions (confirm or flip before build)

- **D1 — Tailwind delivery.** v1 = Play CDN (Node-free `go install`); embedded prebuilt CSS is the
  self-contained-release upgrade. §8.
- **D2 — Route key = full task id** (`/MONO-006`, case-insensitive), not a bare number — bare
  numbers collide across namespaces (`TRAIL-007` vs `GW-007`) even within one planning area. Your
  `/<task-number>` is read as `/<task-id>`.
- **D3 — Scope = one planning area in v1.** Root-of-monorepo → helpful error + `bklg` invocation
  per system. Real multi-system board is v2 (§13). Everyday monorepo use is `bklg . --dir
  systems/trail/knowledge` (or run one instance per system).
- **D4 — Zero-dep line-scanner over goldmark** for v1, behind a swappable interface. The
  conventions are fixed and heading/checkbox-based, and the domain interpretation is needed either
  way, so goldmark buys little while costing the pure-stdlib story. §10.

---

## 13. Non-goals (v1) / follow-ups

Out of scope for v1, in rough priority for later:

1. **Multi-system board** — read the root manifest's system index, render a system switcher (or an
   aggregate board with a per-card system chip). This is the natural v2 and what makes it sing on
   the actual monorepo.
2. **Self-contained release** — `tailwindcss` build + `go:embed` the CSS, offline single binary
   (D1's other branch).
3. **Markdown rendering** of field text (safe subset + sanitizer) instead of escaped plain text.
4. **Journal deep-links** — parse the progress area's `journal.md` so DONE journal pointers and
   detail pages link to the actual entry.
5. **JSON API** (`/api/board.json`) for external consumers; today only `/_v` (mtime) exists.
6. **Live push** — swap the mtime poll for SSE/`fsnotify` if the poll ever feels laggy.

Explicitly *not* planned: editing/mutating any file, any auth, any DB, any bundler.

---

## 14. Build & delivery

- Module `github.com/<you>/bklg`; command `bklg`. `go install github.com/<you>/bklg/cmd/bklg@latest`.
- **Go 1.22+** (required for method + `{id}` / `{$}` mux patterns).
- **Zero Go module dependencies.** stdlib only: `net/http`, `html/template`, `embed`, `flag`
  (or the hand-split argv), `os`, `time`, `strings`, `regexp`, `bufio`.
- The tool follows the same framework it views: if you build it under an agent session, the tasks
  in §15 are the plan; each acceptance criterion below names its decider.

## 15. Task breakdown (executable slices)

Ordered; each is one verifiable slice with a fresh-checkable criterion.

1. **CLI + server skeleton.** `splitArgs`, flags, `127.0.0.1` bind, `GET /{$}` returns 200.
   *AC:* `bklg --port 9001 .` prints the startup block with port 9001; `curl -s localhost:9001/`
   → HTTP 200. `bklg . --port 9001` and `bklg --port=9001 .` behave identically.
2. **Area resolution.** Locations-block dereference, default fallback, root-manifest system list.
   *AC:* against fixtures, resolves the correct `planning`/`progress` dirs (assert the paths);
   pointed at a root manifest with no planning, exits listing the discovered `systems/*`.
3. **Planning parser → model.** Parse CURRENT/BACKLOG/DONE into `[]Card` with dedup (§5).
   *AC:* on `testdata`, card count + per-card `{id, column, len(Acceptance), ParkingLot}` match an
   expected table quoted in the journal; the three reconciliation warnings fire on the seeded
   inconsistencies.
4. **Blocker parse + blocked join.** Parse `blockers.md`; attach the `blocked` badge.
   *AC:* the fixture's open-blocker task carries `blocked`; the resolved-blocker task does not;
   `/_diag` lists zero *unexpected* warnings.
5. **Board template + badges + Tailwind.** Render `/` with three columns, a card per task, badge
   markup, the diagnostics banner.
   *AC:* `curl -s localhost:PORT/` contains the three column headings, one card block per parsed
   task, and a `blocked` badge on the blocked card.
6. **Task detail + 404.** Render `/<id>` (state-appropriate fields + referencing blockers + raw
   block); unknown id → 404.
   *AC:* `/<known-id>` shows its title and acceptance/summary and any blocker; `/NOPE-999` → HTTP
   404.
7. **Live reload.** `/_v` returns max mtime; page polls and reloads on change.
   *AC:* `curl -s localhost:PORT/_v` changes value after `touch`-ing a parsed file; loading `/` and
   editing a file reloads the board within the poll interval.

---

*This spec is written to slot into `reference/specs/` of a framework instance. If the tool is
built inside this monorepo's world it is shared-tier tooling, not a product system — but it more
naturally lives in its own repo, since its input contract is the framework itself, not any one
project.*
