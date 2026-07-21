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

---
## 2026-07-20 14:42 — TASK-003 planning parser → model

**Task:** TASK-003
**What I did:** Added `internal/backlog/model.go` (§4 types) and `parse.go`: a
line-oriented `Parser` (behind an interface per D4). `splitSections` splits on
`## ` (so `### ` task headers survive; `## Entry template` + preamble ignored);
`parseCurrent`/`parseBacklog`/`parseDone` extract the §2 fields; `reconcile`
places one `Card` per id in its furthest column (Done>InProgress>Backlog) and
emits the three §5 warnings. DONE lines split on `" — "` with date=field 3 and
delivery=last field, so a summary with an embedded em-dash round-trips. Built
the rich `testdata/knowledge` demo instance seeding every case.

**What I verified:** Local CI green (`go build/vet ./...`, `gofmt -l` empty,
`go test ./...` exit 0). Dumped the actual parsed board (throwaway test, since
no HTTP surface yet) and quoted it:
```
CARDS (7): DEMO-4 Done done=true | DEMO-6 Done done=true |
  DEMO-5 Done done=false (backlog [x], no DONE) | DEMO-1 InProgress ac=2 |
  DEMO-2 InProgress ac=0 | DEMO-3 Backlog | (no id) Backlog parking=true
WARNINGS (3): CURRENT holds >1 active task | DONE item not ticked in BACKLOG:
  DEMO-6 | shipped item missing from DONE.md: DEMO-5
```
Tests: card table, dedup, exactly-3-warnings, field extraction (incl. em-dash
summary), defensive (malformed line + missing area). Fresh-context review:
**PASS** all 5 AC (it probed 2-em-dash summaries, NUL bytes, empty input — no
panics). Two [low] findings fixed (parking-with-id dedup; shipped-in-CURRENT
keeps detail) + dup-warning/symmetry notes; fixture board unchanged.

**What changed:** New `internal/backlog/{model.go,parse.go,parse_test.go}` +
`testdata/knowledge/*`. Delivery record: **PR #5, merged `42ff5c9`** (squash).

**What I learned:** The DONE " — "-split is robust *because* date is pinned to
field 3 and delivery to the last field — never trust field count alone when a
free-text field sits in the middle. Also: `Resolve(path,dir)` joins `path/dir`,
so a fixture instance at `testdata/knowledge` is reached with
`Resolve("testdata","knowledge")`, not `Resolve("testdata/knowledge",…)` (cost
me one red test). This is the 3rd close PR — the close mechanics (branch/commit/
push/create/merge/sync) are now scriptable; will add `scripts/close-task.sh` for
TASK-004+.

**Next:** TASK-004 — parse `blockers.md` (§/blocker heading disambiguation),
compute the `blocked` join + all badges, wire the board into the server, add
`GET /_diag`. AC in `CURRENT.md`.

---
## 2026-07-20 14:57 — TASK-004 blockers + badges + /_diag

**Task:** TASK-004
**What I did:** `parseBlockers` reads `blockers.md`; a `## ` line matching
`^BLOCKER-\d+` is a blocker (assigned to the last section seen), else it's a
section — everything under `## Format` is skipped; `Open` = under `## Open`.
`computeBadges` joins planning+progress in one place: `blocked` (open blocker's
affected id), `parking`, `override`, `no-ac`, `namespace`. New `server.go`:
`Server` re-parses per request, serves a `/` placeholder + `GET /_diag`
(warnings verbatim, text/plain). `main` now serves `srv.Routes()`.

**What I verified:** Local CI green. Tests: `TestParseBlockers` (BLOCKER-001
Open/DEMO-2, BLOCKER-002 Resolved/DEMO-1, Format example skipped), `TestBadges`
(DEMO-2 blocked+no-ac+namespace; DEMO-1 override, not blocked; parking chip),
`TestDiagRoute` (exactly 3 warning lines, text/plain), `TestBoardRoute` (200).
Binary: `/_diag` on the fixture →
```
CURRENT holds >1 active task (framework one-task invariant)
DONE item not ticked in BACKLOG: DEMO-6
shipped item missing from DONE.md: DEMO-5
```
**Dogfood** `bklg . /_diag` → `no warnings` (this repo's real instance is clean).
Fresh-context review: **PASS** all 5 AC (probed case-insensitive join, live
reparse, nil-ID safety). Its [low] finding — em-dash in a blocker title — and a
trailing-text join note were fixed (anchor on `opened <ts>`; normalize the join
id) and pinned with tests; fixture unchanged.

**What changed:** `internal/backlog/parse.go` (+parseBlockers/parseBlockerHead/
computeBadges, extended Parse), new `internal/backlog/server.go`,
`internal/backlog/{blockers_test,server_test}.go`, `cmd/bklg/main.go` (serves
via Server; dropped the inline mux/skeleton). Delivery: **PR #7, merged
`835363f`**.

**What I learned:** Two functions now anchor a middle free-text field by pinning
the fixed fields at the ends (DONE: date=3rd, delivery=last; blocker: title
between id and `opened`) — the robust pattern for em-dash-delimited lines.
Also: first use of `scratchpad/close-task.sh` for this close PR's mechanics.

**Next:** TASK-005 — board template (`/`): three columns, card per task, badge
markup, diagnostics banner; `html/template` + `go:embed` + Tailwind Play CDN.
AC in `CURRENT.md`.

---
## 2026-07-20 15:11 — TASK-005 board template + badges + Tailwind

**Task:** TASK-005
**What I did:** Embedded `html/template` templates (`templates/layout.html` +
`board.html`, `go:embed` from the `backlog` package — embed can't traverse `..`)
and `render.go`: `viewModel` splits the `Board` into Backlog / In Progress /
Done columns; `badgeClass`/`badgeText` map the five badge kinds to Tailwind
colours (blocked red, parking slate, override amber, no-ac yellow, namespace
gray). `handleBoard` renders into a buffer (clean 500 on template error) and the
layout shows the planning-path header + a `/_diag` banner when warnings exist.

**What I verified:** Local CI green. `TestBoardRender` (3 columns, 6 id-cards +
links, blocked badge, namespace chip, header path, `/_diag`, Tailwind CDN) and
`TestBoardEscaping` (a `<script>` in card text → `&lt;script&gt;`, never a live
tag) pass. Binary smoke on the fixture: 200; DEMO-2 card shows red `blocked` +
yellow `no-ac` + gray `DEMO` chips; parking card muted + unlinked; empty column
→ "nothing here". **Dogfood** `bklg .` → this repo's board, 200, 3 columns, all
7 `TASK-00N` cards linked. Fresh-context review: **PASS**, **zero findings** —
included a hostile-input injection probe (titles/path/labels/pathological id in
`href`) all escaped; `/_diag` stays text/plain; empty board safe.

**What changed:** New `internal/backlog/render.go`,
`templates/{layout,board}.html`; `server.go` (buffered board render, −placeholder);
`server_test.go` (+render/escape tests). Delivery: **PR #9, merged `a9a1fc3`**.

**What I learned:** Separate `template.Template` sets per page (boardTmpl now,
taskTmpl in TASK-006) let each define `content` without collision — the clean
way to share a `layout` across pages in `html/template`. Buffered render is the
right default: a template error becomes a 500 instead of a half-written page.

**Next:** TASK-006 — task detail (`/{id}`, case-insensitive) + 404;
state-appropriate fields, referencing blockers (open first), collapsed Raw
block. Reinstate `Board.CardByRawID`. AC in `CURRENT.md`.

---
## 2026-07-20 15:24 — TASK-006 task detail + 404

**Task:** TASK-006
**What I did:** Added `GET /{id}` → `handleTask`: `CardByRawID` (case-insensitive
via `EqualFold`), `http.NotFound` for unknown/id-less cards. Renders
id/title/namespace/column/badges, state fields (In-Progress vs Done),
referencing blockers (open-first via a `[]bool{true,false}` pass), and a
collapsed `<details>` Raw block. New `task.html` + separate `taskTmpl`;
reinstated `Board.CardByRawID`.

**What I verified:** Local CI green. `detail_test.go` (In-Progress fields +
checklist ☑/☐ + override + notes + raw; Done date/summary/delivery/journal;
open vs resolved blockers on DEMO-2/DEMO-1; `/NOPE-999`→404, `/demo-1`→200).
Binary: `/DEMO-1`→200, `/demo-1`→200, `/NOPE-999`→404, `/_diag`→200 (literal
still wins). Fresh-context review: **PASS**, **zero findings** — injected
`<script>`+quotes into all 16 rendered fields (all escaped to `&lt;`), path
traversal (`/..`, `%2e%2e`, `/DEMO-1/extra`) → 301/404 no file-read/panic,
blocker open-first ordering proven with a seeded resolved-before-open case. Its
one note (stale `Routes()` comment mentioning unregistered `/_v`) was fixed.

**What changed:** `internal/backlog/{server,model,render}.go`, new
`templates/task.html`, new `detail_test.go`. Delivery: **PR #11, merged
`cd00f36`** (+ comment-fix commit).

**What I learned:** Go 1.22 mux precedence is genuinely "no ordering tricks" —
`/_diag` (literal) beats `/{id}` (wildcard) regardless of registration order,
and `/DEMO-1/extra` doesn't match `/{id}` (single segment) so it 404s cleanly.
Escaping is entirely `html/template`'s doing — every field uses plain `{{}}`,
zero `template.HTML`.

**Next:** TASK-007 (final MVP task) — live reload: `/_v` max-mtime endpoint +
~3s poll → `location.reload()`; compute `Meta.LatestMTime`. AC in `CURRENT.md`.

---
## 2026-07-20 15:35 — TASK-007 live reload

**Task:** TASK-007
**What I did:** `areaMTime(Areas)` = max mtime across CURRENT/BACKLOG/DONE/
blockers (4 stats); `Parse` sets `Board.Meta.LatestMTime`. `GET /_v` →
`versionString(areaMTime)` (bare UnixNano int, text/plain) — stats, doesn't
reparse, so polling is cheap. Both view models carry `Version`; `layout.html`
embeds ~8 lines of vanilla JS: `fetch("/_v",{cache:"no-store"})` every 3s →
`location.reload()` when the value differs from the baked-in one.

**What I verified:** Local CI green. `version_test.go`: `TestMetaLatestMTime`
(== max of 4 files), `TestVersionRoute` (200/text-plain/`^-?\d+$`),
`TestPollScript` (script present, baked `var v` == `/_v`),
`TestVersionChangesOnModify` (deterministic via `Chtimes` to a future time).
Binary smoke on a temp copy: `/_v` `…150491…` → after editing BACKLOG →
`…151152…` (changed: YES); poll script present. Fresh-context review: **PASS**,
**zero findings** — `-race` clean, 40 concurrent `/_v` consistent, XSS probe on
the baked version (JS-context escaping holds), non-parsed-file edits ignored,
empty instance → `0`.

**What changed:** `internal/backlog/{parse,render,server}.go`,
`templates/layout.html`, new `version_test.go`. Delivery: **PR #13, merged
`14b1df6`**.

**What I learned:** `/_v` deliberately stats rather than reparsing — the poll
path must stay cheap. String **inequality** (not strict-greater) as the reload
trigger also catches an mtime rewind (git checkout). Sharing one `versionString`
between the baked value and `/_v` guarantees they're the same computation, so
equality reliably means "unchanged."

**Next:** MVP complete — run the end-of-session sweep and stop.

---
## 2026-07-20 15:35 — v1 MVP complete (session wrap)

**Task:** session summary
**What I did:** Shipped all seven §15 slices, each as a feature PR + a fresh-
context review + a close PR: TASK-001 CLI/server skeleton (#1), TASK-002 area
resolution (#3), TASK-003 planning parser + dedup (#5), TASK-004 blockers +
badges + `/_diag` (#7), TASK-005 board template + Tailwind (#9), TASK-006 task
detail + 404 (#11), TASK-007 live reload (#13). 13 feature/close PR pairs total,
all squash-merged; `main` is the only branch.

**What I verified:** Every task passed all its acceptance criteria and an
independent fresh-context review (all PASS; findings fixed or rebutted in-PR).
End-to-end **dogfood**: `bklg .` serves this repo's own board (7 cards across
Backlog/In-Progress/Done, `/_diag` = "no warnings" — the real instance is
clean), and every route works (`/`, `/{id}`, `/_v`, `/_diag`, 404s). Zero Go
module deps; stdlib only. Full end-of-session sweep quoted below the DONE entries.

**What I learned:** The close-PR-per-task rhythm is heavy (14 PRs for 7 tasks)
but keeps `main` always-green and every change reviewed; the close mechanics
became a script (`scratchpad/close-task.sh`) after the 3rd repetition. The tool
viewing its own build instance (dogfooding) caught nothing broken but is a
genuinely useful confidence check.

**Next:** Session envelope spent (MVP shipped). Stop. Future work: promote a
parking-lot item (spec §13) in a new session — GitHub Actions CI would be a good
first one (it makes delivery gate D3 real).

---
## 2026-07-20 16:02 — v2 feedback triaged into TASK-008…012

**Task:** blocker triage / planning
**What I did:** User reviewed the shipped MVP (happy) and gave feedback after
running bklg against the real `trail` monorepo. Investigated
`trail/systems/track/knowledge/planning`, recorded the concrete divergences in
`whiteboard/trail-instance-findings.md`, and turned the feedback into five v2
tasks: TASK-008 trim card text, TASK-009 parser robustness (real-instance
formats), TASK-010 render markdown, TASK-011 make /_diag actionable, TASK-012
multi-system board (aggregate + filter). Promoted the multi-system + markdown
parking items to Active.

**What I verified:** Read the real trail files (quoted in the whiteboard entry):
`DONE.md` has **0** `- ` bullets — it uses `### TRACK-NNN — …` headings — so our
parser finds no DONE entries → all 15 `[x]` backlog items fire the false
"shipped item missing from DONE.md" warnings the user saw; parking ids look like
`**WI-8 — …`; `[x]`/parking lines carry full paragraphs (the wall-of-text cards).

**What changed:** `BACKLOG.md` (+5 v2 tasks, parking pruned),
`whiteboard/{trail-instance-findings.md, README.md}`. Delivery: this docs PR.

**What I learned:** v1's spec was keyed to the framework **skeleton**; real
instances diverge, which is exactly what bklg must handle. The v2 parser work
widens the §2 input contract — ADRs due (format-widening; markdown vs the
zero-dep constraint). This is the trigger the D4 goldmark seam anticipated.

**Next:** Asked the user whether to ship the v2 batch autonomously now (same
PR-per-task flow), do only the quick trail-fixes and pause before the
design-heavy multi-system feature, or leave it queued. Await steer.

---
## 2026-07-20 16:30 — TASK-008 trim card text

**Task:** TASK-008
**What I did:** Added `truncate(s, max)` template helper (rune-based, word-
boundary backoff, `…`) in `render.go`; `board.html` uses `{{truncate .Title 140}}`
for the card title; `task.html` keeps the full title. Display-only.

**What I verified:** Local CI green. `TestTruncate` (short unchanged, long→`…`
≤141 runes, multibyte valid), `TestCardTitleBoardVsDetail` (board truncates,
detail full). **Real-repo smoke** `bklg /Users/bb8/dev/trail --dir
systems/track/knowledge` → 18 cards now end in `…` (~140-char titles, were full
paragraphs). Fresh-context review: **PASS** all 4 AC, no findings (ruled out
index-out-of-range even on invalid UTF-8; escaping intact — `+`→`&#43;`).

**What changed:** `render.go` (+truncate), `board.html` (card title),
`trim_test.go` (new). Delivery: **PR #16, merged `08c4053`**.

**What I learned:** User chose "ship all five" v2 tasks autonomously. Trimming
is the raw title now; TASK-010 markdown will make it emphasis-aware. The huge
Done-card titles will also shrink once TASK-009 parses the `### `-heading DONE
entries (dedup picks the short heading title over the fat `[x]` backlog stub).

**Next:** TASK-009 — parse `### <ID>`-heading DONE entries + emphasized ids;
record ADR-0001 (contract widening). AC in `CURRENT.md`.

---
## 2026-07-20 16:47 — TASK-009 parser robustness (real instances)

**Task:** TASK-009
**What I did:** Extended `parseDone` to also parse `### <ID> — <title>`
heading-style entries (`parseDoneHeading` + `fieldAfter`/`doneSummary`/
`prPrefixed`), kept the bullet format; `parseID` now strips leading `* ` `` ` ``
`_`. ADR-0001 records the contract widening (keep line-scanner, defer goldmark/D4).

**What I verified:** Local CI green. Tests: `TestParseDoneHeadingFormat`,
`TestParseIDEmphasis`, `TestParseDoneFormatExclusivity`, and (after review)
`TestParseIDTitleDecorated`. Real trail: `/TRACK-000` now shows Date 2026-06-25 /
PR #161 / Journal / Summary; `/_diag` no longer lists TRACK-000. Fresh-context
review: **PASS** all 5 AC; caught a **latent [low] bug** — `parseIDTitle`'s
no-separator branch sliced `s[len(id.Raw):]` from index 0, wrong once `parseID`
strips decoration (`**WI-8**` → `"8**"`). Fixed (locate id in stripped string +
trim) + added the missing test.

**What I learned (correction):** My AC4 hypothesis "15 → 0 warnings" was WRONG.
trail's DONE.md holds **only** TRACK-000; TRACK-001…014 live solely in the rich
`[x]` BACKLOG lines. So 14/15 warnings are the check *correctly* flagging a
convention divergence — a **policy** question for TASK-011, not a parse bug. I
revised AC4 honestly (didn't fake it) and corrected the whiteboard. Lesson:
verify the *data* before assuming a warning is a parse failure.

**What changed:** `parse.go` (+heading DONE, +decoration-strip, parseIDTitle
fix), `donefmt_test.go`, `decisions/0001-*.md` + INDEX, whiteboard correction.
Delivery: **PR #18, merged `d0e33a5`**.

**Next:** TASK-010 — safe inline markdown (escape-first, stdlib-only, ADR-0002);
the one place `template.HTML` is used → paranoid security review. AC in `CURRENT.md`.

---
## 2026-07-20 17:11 — TASK-010 safe inline markdown

**Task:** TASK-010
**What I did:** `markdown.go` `renderMarkdown` — escape-first, then whitelist
inline patterns (`**`/`*`/`` ` ``/links) → `template.HTML`; `safeURL` scheme-check
(rejects javascript/data/protocol-relative/backslash); code spans placeheld;
NUL stripped; asterisk-only emphasis. Registered `md`; templates render free-text
fields via `{{md}}`, Raw block stays auto-escaped. ADR-0002 (keep zero-dep vs
goldmark+sanitizer).

**What I verified:** Local CI green + `-race`. 7 markdown tests (render, escape-
first, unsafe-link, code-protected, unclosed, no-underscore, safeURL, NUL). Real
trail: `/TRACK-000` renders `<strong>`/6×`<code>`; board 0 live `<script>`.
**Adversarial security review: PASS, no injection** (~90 hostile inputs, `"`-break
attempt, placeholder spoof, structural fuzzing). Its two non-XSS notes hardened
(reject `//`/backslash links; strip NUL).

**What I learned:** Escape-first + tag-whitelist + href-scheme-check is safe
*without* a sanitizer precisely because repo HTML never reaches output — the
sanitizer's job never arises. That's the whole justification for staying
zero-dep (ADR-0002). Underscore emphasis is a footgun on technical text
(snake_case, dunders) → asterisk-only.

**What changed:** `markdown.go`+`markdown_test.go` (new), `render.go` (funcs),
`templates/{board,task}.html`, `decisions/0002-*.md`+INDEX. Delivery: **PR #20,
merged `ff2ea83`**.

**Next:** TASK-011 — make `/_diag` actionable: structured `[]Warning` grouped by
kind, counts, explanations, id links; reframe "missing from DONE.md". AC in
`CURRENT.md`.

---
## 2026-07-20 17:30 — TASK-011 actionable /_diag

**Task:** TASK-011
**What I did:** Structured `Board.Warnings` into `[]Warning{Kind,Message,TaskRaw}`
(readArea→read-error, parseDone→malformed-done, reconcile→shipped-missing-done/
done-not-ticked/current-multiple). `/_diag` is now an HTML page (`diag.html` via
the layout) grouping by kind with counts + explanations + id links;
`buildDiagVM` in render.go. shipped-missing reframed as informational.

**What I verified:** Local CI green (48 tests). `TestWarningKinds`, rewritten
`TestDiagRoute` (HTML, groups, links, escaped message). Real trail: 14
shipped-missing → one "Shipped, but not in DONE.md (14)" group with 14 TRACK-
links + explanation. Fresh-context review PASS (Message auto-escaped, href
URL-normalized); fixed a latent map-order nondeterminism (sort unknown kinds).

**What changed:** `model.go` (Warning type), `parse.go` (structured producers),
`render.go` (diagVM/buildDiagVM/diagTmpl + VM Warnings type + sort), `server.go`
(handleDiag HTML), `templates/diag.html` (new), tests. Delivery: **PR #22,
merged `ce95127`**.

**What I learned:** Structured warnings pay off immediately — grouping +
explanation turned 15 scary lines into one informational section. The `>`-in-
message escaping tripped a test assertion (literal `>` vs rendered `&gt;`) —
good reminder that auto-escaping applies to my own diagnostic strings too.

**Next:** TASK-012 (final) — multi-system board: aggregate all `systems/*` at a
root manifest + `?system=` filter + per-card system chip; ADR-0003. AC in
`CURRENT.md`.

---
## 2026-07-20 17:52 — TASK-012 multi-system board

**Task:** TASK-012
**What I did:** Aggregate mode. `Card.System`; `NewMultiServer(root, systems)`;
`Server.board()` resolves+parses each `systems/<name>`, tags System,
concatenates; `main` builds it on `RootManifestError` (startup `aggregate: N
systems`). Server-side `?system=` filter; filter bar lists **all** discovered
systems (incl. empty, so any project is filterable); per-card system chip on
board + detail; global `/{id}` lookup; `/_v` = max mtime across systems;
unresolvable system → skip + warning. ADR-0003.

**What I verified:** Local CI green (53→ now more tests). Controlled `testdata/
mono` (alpha+beta) tests + `TestSingleModeIgnoresSystemFilter`. **Real trail**:
`bklg /Users/bb8/dev/trail` → `aggregate: 5 systems — trail, cadence, gateway,
track, reflect`, `/`=200 ~105 cards, chips, `?system=track`→15 no-leak, all 5 in
the bar (cadence/reflect empty), `/TRACK-000`→200, `/_v` live. Fresh-context
review PASS all 5 AC; fixed a [low] contract bug (single mode didn't ignore
`?system=` → blanked the board; now zeroed when no systems).

**What I learned:** cadence/reflect showed 0 cards not from a parse bug but
because their only `### ` heading sits under `## Entry template` (correctly
ignored) — verified the raw files before assuming a bug (the TASK-009 lesson,
reapplied). Listing empty systems in the bar makes "filter to any project"
literally true and matches the startup count.

**What changed:** `model.go` (Card.System), `server.go` (multi-server +
aggregate + Systems + latestMTime), `render.go` (viewModel filter + systemChip),
`main.go` (root→aggregate), `templates/{board,task}.html`, `testdata/mono/*`,
`multi_test.go`, `decisions/0003-*.md`+INDEX. Delivery: **PR #24, merged `95acdf7`**.

**Next:** v2 batch (TASK-008…012) complete — CURRENT empty. End-of-session sweep,
then stop. Future: parking-lot items + the v2 ADR follow-ups (block markdown,
per-system detail routes, CI).

---
## 2026-07-21 12:31 — Plan the dashboard adapter (Path B) for Pinata-shape KBs

**Task:** none — planning/grooming (no code). Prompted by a user report that
`bklg .` in `~/dev/Pinata-dev/Pinata` exits `no planning area at
knowledge/planning`.

**What I did:** Investigated the real Pinata KB. Diagnosis: it's neither
"shifted slightly" nor "no planning" — it runs a full Active/Backlog/Done
lifecycle in a *different, incompatible shape*: manifest is `index.md` (not
`README.md`) with no `## Locations` block; no `planning/`/`progress/` dirs; the
whole board is one file, `work/index.md`, with Active/Done as pipe **tables**,
Backlog as bullet groups under bold subheads, identity via inline Linear ids
(`PINATA-\d+`, 0..N per row) not one structured id per card, and "blocked" as
prose (`⛔`); Linear is the source of truth for status. So the mismatch is at
the **data-model** level, not just paths. Laid out two paths — A: conform the KB
to the skeleton (rejected: hand-maintained duplicate of `work/index.md` +
Linear); B: teach bklg a second "dashboard adapter" convention. Drafted Path B
as backlog **v3 (TASK-013…016)**, wrote **ADR-0004** (status *proposed*), a
canonical contract at `reference/specs/dashboard-format.md`, and a standalone
handoff copy at repo-root `pinata.md` for the Pinata coding agent.

**What I verified:** Local CI green (`go build`/`vet`/`gofmt -l`/`go test`) — no
Go changed, but confirmed the docs branch doesn't break the build. This is
docs-only; no runtime behaviour to smoke-test yet (the adapter is unbuilt).

**What I learned:** bklg's "the patterns are the contract" (§2) makes Path B
inherently two-sided — a defensive parser over irregular input yields a `/_diag`
dump, so the adapter must pair with a format contract the target KB follows.
Two deliberate calls (in ADR-0004): mode is an explicit `dashboard:` config key
(no auto-detect in v1), and blocked is a **leading** `⛔` only (Pinata already
uses `⛔` mid-prose as a decision marker, so "contains" would false-positive).

**What changed:** `planning/BACKLOG.md` (v3 batch + 2 parking-lot follow-ups),
`decisions/0004-dashboard-adapter.md` + INDEX, `reference/specs/dashboard-format.md`,
this journal. Repo-root `pinata.md` written as an untracked handoff artifact
(not part of the PR). Delivery: docs PR (opened; merge left to the user — it
proposes un-greenlit work).

**Next:** user greenlights (or edits) Path B; on go-ahead, promote TASK-013 to
`CURRENT.md` with AC and flip ADR-0004 to accepted. The Pinata KB update is
handed off separately via `pinata.md`.

---
## 2026-07-21 15:20 — TASK-013 dashboard-mode resolution

**Task:** TASK-013 (BACKLOG v3, ADR-0004). User greenlit Path B.

**What I did:** Added dashboard-mode resolution — the first slice of the
dashboard adapter. `Resolve` manifest lookup now tries `README.md` then
`index.md` (first hit wins); a `dashboard:` key in the `## Locations` block, or a
new `--dashboard FILE` flag, selects the single-file adapter and short-circuits
the planning/progress + systems resolution. Dashboard paths resolve against the
repo root (`path`), like every other Locations value. `Areas` gained
`DashboardFile`; `NewParser()` is now a `defaultParser` that dispatches on it —
`parseDashboard` returns an empty board for this slice (real parsing is
TASK-014). `areaMTime` watches the single file in dashboard mode so `/_v` stays
live. `main` grew the `dashboardServer` helper + the flag in `splitArgs`
`takesValue`. Flipped ADR-0004 proposed→accepted.

**What I verified:** Local CI green (build/vet/gofmt/test). Smoke (all quoted in
the PR): missing file → `bklg: no dashboard file at nope.md` exit 1; `bklg .`
framework startup line unchanged, `/`→200; `--dashboard knowledge/work/index.md
~/dev/Pinata-dev/Pinata` → startup `dashboard: …/work/index.md`, `/`→200; the
Locations-key fixture via `main` → `/`→200 and `/_v` a real (non-zero) mtime.
Unit tests: both resolution paths + missing-file (flag AND Locations key) +
non-dir + framework-unaffected. Fresh-context review (diff+AC only): PASS all 5
AC, zero correctness findings; its coverage nit (missing-file via Locations)
fixed with a test + `dashboard-missing` fixture, its other two nits acknowledged
(`Meta.PlanningDir` reuse → TASK-014; manifest widening → intended per ADR-0004).

**What I learned:** The fixture's `dashboard:` value deliberately includes the
`knowledge/` prefix (`knowledge/work/index.md`), so it only resolves correctly
under repo-root joining — the fixture genuinely proves the repo-root-relative
rule rather than accidentally passing under a base-relative join.

**What changed:** `resolve.go` (manifest loop, `dashboard:` branch,
`ResolveDashboard`/`dashboardAreas`, `parseLocations` key), `parse.go`
(`defaultParser` dispatch, `parseDashboard` stub, dashboard-aware `areaMTime`),
`main.go` (`--dashboard` flag, mode switch, `dashboardServer`),
`resolve_test.go` + `testdata/resolve/dashboard{,-missing}`, planning/ADR docs.
Delivery: **PR #27, merged `6cea51f`**. D3 remote check vacuous (no CI).

**Next:** TASK-014 — dashboard parser → model (Active/Done tables + Backlog
bullet groups → `[]Card` with Tickets/Group/Material; escaped-pipe handling;
defensive warnings). AC in `CURRENT.md`.
