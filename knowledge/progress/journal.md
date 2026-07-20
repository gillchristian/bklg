# Journal

Append-only. Newest at the bottom. Each entry is a snapshot for future-me
with no memory of this session.

## Entry format

---
## YYYY-MM-DD HH:MM — <short heading>
   (HH:MM may be dropped only if the heading alone makes the entry findable
   from a DONE.md pointer)

**Task:** TASK-NNN (or "scaffolding" / "exploration" / "blocker triage")
**What I did:** 1–3 sentences.
**What I verified:** which gates ran, with literal command output worth keeping.
**What changed:** files touched + the delivery record per my profile.
**What I learned:** anything that would surprise future-me. Non-obvious only.
**Next:** the very next thing to do on resume.

## Entries

---
## 2026-07-20 13:43 — Framework init + spec → backlog + bootstrap

**Task:** scaffolding
**What I did:** Ran SETUP.md adoption against this project. Created the five
instance areas and their files from the skeletons; wrote the manifest
(`knowledge/README.md`) declaring **delivery: pr**, the Locations block (default
`knowledge/…` paths), attribution = user identity only (no agent attribution —
user decision), and the session envelope = ship the MVP (TASK-001…007). Wrote
the project `CLAUDE.md`, `.gitignore`, root `README.md`, MIT `LICENSE`. Filled
the brief and glossary from the spec; recorded the real Go local-CI commands in
`reference/local-ci.md`. Relocated the spec to
`reference/specs/bklg-spec.md` (its own closing note says it belongs there).
Processed the spec's §15 into a 7-task `BACKLOG.md`; pulled **TASK-001** into
`CURRENT.md` with acceptance criteria (each naming its decider). Initialized git
and made the one allowed bootstrap commit on `main`.

**What I verified:** No code yet — nothing to run. Framework guard checked
before starting: no pre-existing `planning/CURRENT.md`, so this was a genuine
first-time adoption, not a re-run against live state. Bootstrap commit + push
output quoted in the git-bootstrap step below.

**What changed:** New: `CLAUDE.md`, `.gitignore`, `README.md`, `LICENSE`,
`knowledge/README.md`, `knowledge/planning/{CURRENT,BACKLOG,DONE}.md`,
`knowledge/progress/{journal,blockers}.md`, `knowledge/decisions/INDEX.md`,
`knowledge/reference/{project-brief,glossary,local-ci}.md`,
`knowledge/reference/specs/bklg-spec.md` (moved), `knowledge/whiteboard/README.md`.
Delivery record: **bootstrap commit on `main`** (pre-framework, covered by the
manifest's Rule-1 exception — the only direct-to-main commit this project gets).

**What I learned:** The GitHub repo `gillchristian/bklg` was pre-created by the
user — **empty, public**, default branch not yet set — so the bootstrap push
establishes `main` as the default branch, after which every change goes via PR.
`gh` is authed as `gillchristian` over SSH. The framework's upstream repo/path
is not recorded (it was pre-placed at `knowledge/framework/`); logged as a
non-blocking open question in the brief. No remote CI exists, so delivery gate
D3 is vacuous until GitHub Actions is added (parked).

**Next:** Start TASK-001 (CLI + server skeleton): branch `feat/task-001-…`,
`go mod init github.com/gillchristian/bklg`, implement `cmd/bklg/main.go`
(`splitArgs`, flags, 127.0.0.1 bind, `GET /{$}`→200, startup block), local CI,
PR, fresh-context review, squash-merge, close PR.

---
## 2026-07-20 14:00 — TASK-001 CLI + server skeleton

**Task:** TASK-001
**What I did:** Built `cmd/bklg/main.go` (module `github.com/gillchristian/bklg`,
`go 1.22`): `splitArgs` pre-splits argv (optional `[path]` in any position;
`--flag v` and `--flag=v`), flags `--port` (1235) / `--dir` (knowledge),
loopback-only listener created before the readiness echo, and a `GET /{$}`
handler returning a 200 HTML placeholder. Startup block matches spec §9. After
the fresh-context review I hardened arg handling (reject extra positionals at
the boundary; `-h` exits 0) and added table-driven tests.

**What I verified:** Local CI green —
`go build ./... && go vet ./... && gofmt -l (empty) && go test ./...` → exit 0;
`ok github.com/gillchristian/bklg/cmd/bklg`. All 5 AC by running the binary:
startup block verbatim `Running Backlog on port 9041` + resolution echo +
`http://localhost:9041`; all three arg forms → `curl / = 200`; `lsof` showed the
listener bound to `127.0.0.1:*` (reviewer's check); port-in-use → exit 1 with
`bklg: cannot listen on 127.0.0.1:9041: … address already in use` and no
"Running" line; extra positional → `bklg: unexpected argument(s): b …` exit 2.
Fresh-context review (subagent, diff+AC only): **PASS**, no blocking findings.

**What changed:** `cmd/bklg/main.go`, `cmd/bklg/main_test.go`, `go.mod` (new);
`knowledge/planning/CURRENT.md` (scope-refinement note re: deferring the embed
seam to TASK-005). Delivery record: **PR #1, merged `288a814`** (squash).

**What I learned:** (1) The Bash tool here runs **`/bin/zsh`**, which does *not*
word-split unquoted `$var` — a loop of `bklg $form` passed each whole string as
one arg and misfired. Pass separate args / use `"$@"`. (2) macOS first-exec of a
freshly built binary can outrun a 0.4s startup sleep (Gatekeeper-ish check);
give warm-up or ~0.6s. (3) stdlib `flag` *ignores* (doesn't reject) extra
positionals, and since it stops at the first non-flag token, a stray positional
before a flag silently drops that flag — hence the boundary rejection.

**Next:** TASK-002 — area resolution (`internal/backlog/resolve.go`): Locations
dereference, default fallback, root-manifest system list; wire into `main`, echo
resolved paths; add `testdata/` resolution fixtures. AC in `CURRENT.md`.

---
## 2026-07-20 14:19 — TASK-002 area resolution

**Task:** TASK-002
**What I did:** Added `internal/backlog/resolve.go` (`Resolve(path,dir)` →
`Areas{KnowledgeDir,PlanningDir,ProgressDir}`): parse `base/README.md`'s
`## Locations` block (keys `planning`/`progress`, values joined **against
`path`** = repo-root-relative), fall back to `base/planning`+`base/progress`,
and on a missing planning dir detect a `systems/<name>` root manifest
(`RootManifestError`) vs a plain "no planning area". `main` now resolves before
binding, exits non-zero with a clear message (no server) on failure, prints the
per-system invocation list for a root manifest, and echoes resolved paths via
`DisplayPath` (restores the `./`). Removed the skeleton's `joinDisplay`.

**What I verified:** Local CI green — `go build/vet ./...`, `gofmt -l` empty,
`go test ./...` → exit 0; `ok …/internal/backlog`, `ok …/cmd/bklg`. 9 backlog
unit tests PASS (Locations dereference w/ **non-default** `alt/planning` so it
can't pass coincidentally; default fallback; root-manifest dedup+order; no
planning area ≠ root manifest; path-not-dir; parseLocations incl. partial block;
parseSystems; DisplayPath). Binary smokes: root manifest → exit 1 listing
`bklg … --dir systems/alpha/knowledge` + `…beta…`; `empty` → exit 1
`no planning area at …/empty/knowledge/planning`; file & missing path → exit 1
`path is not a directory`. **Dogfood** `bklg .` → `planning: ./knowledge/planning`
via the real manifest, HTTP 200. Fresh-context review: **PASS**, no findings.

**What changed:** New `internal/backlog/{resolve.go,resolve_test.go}` +
`testdata/resolve/*` fixtures; `cmd/bklg/main.go` (+`errors`, backlog import,
resolution wiring, −`joinDisplay`); `cmd/bklg/main_test.go` (−`TestJoinDisplay`).
Delivery record: **PR #3, merged `4cf04c1`** (squash).

**What I learned:** Go tooling ignores `testdata/` dirs, so package-local
`internal/backlog/testdata/` is the clean home for fixtures (tests run with CWD =
package dir → simple relative paths); minor, justified deviation from spec §10's
repo-root `testdata/`. Review flagged that a present-but-**unreadable** manifest
silently falls back to default (only `err==nil` checked) — accepted for now; the
right fix is a `/_diag` warning once the Warnings system exists (**TASK-004**).

**Next:** TASK-003 — planning parser → `[]Card` + dedup + the three
reconciliation warnings; build the rich `testdata/knowledge` §11 fixture. AC in
`CURRENT.md`.
