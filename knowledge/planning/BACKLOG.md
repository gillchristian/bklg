# Backlog

Ordered. The top *unchecked* item is next. Promote into `CURRENT.md` when started.

Conventions:
- Completed tasks stay here, checked `[x]` with their delivery record — the
  list doubles as a one-line shipping index. The full record lives in `DONE.md`.
- TASK ids are global and monotonically increasing across all of `planning/`:
  next id = highest id appearing anywhere + 1.

Derived from the spec's §15 "Task breakdown (executable slices)"
(`../reference/specs/bklg-spec.md`). The seven slices below **are** the v1 MVP
(session envelope). Each is one verifiable vertical slice. Test fixtures (§11)
are built inside the tasks that first need them — resolution fixtures in
TASK-002, the rich parser/demo instance in TASK-003 — not as a separate task.

## Active

- [x] TASK-001 — CLI + server skeleton — `splitArgs` (positional-in-any-position, `--flag v`/`--flag=v`), `flag` parsing, `127.0.0.1` bind, `GET /{$}` → 200, startup block. Creates `go.mod` (module `github.com/gillchristian/bklg`), `cmd/bklg/main.go`. — **PR #1, merged `288a814`**
- [x] TASK-002 — Area resolution — Locations-block dereference (repo-root-relative), default `base/planning`+`base/progress` fallback, root-manifest system-index list + helpful exit. Adds resolution fixtures under `testdata/`. — **PR #3, merged `4cf04c1`**
- [x] TASK-003 — Planning parser → model — line-oriented parser for `CURRENT.md`/`BACKLOG.md`/`DONE.md` → `[]Card` with the most-advanced-state dedup and the three reconciliation warnings. Ships the rich `testdata/knowledge/` demo/test instance (§11). — **PR #5, merged `42ff5c9`**
- [x] TASK-004 — Blocker parse + blocked join — parse `blockers.md` (§/blocker heading disambiguation), attach the `blocked` badge to open-blocker tasks; `/_diag` shows warnings. — **PR #7, merged `835363f`**
- [x] TASK-005 — Board template + badges + Tailwind — render `/` (three columns, card per task, badge markup, namespace chip, diagnostics banner) via `html/template` + Tailwind Play CDN. — **PR #9, merged `a9a1fc3`**
- [x] TASK-006 — Task detail + 404 — render `/{id}` (state-appropriate fields, referencing blockers, collapsed Raw block); unknown id → 404. — **PR #11, merged `cd00f36`**
- [x] TASK-007 — Live reload — `/_v` returns max mtime across parsed files; page polls ~3s and `location.reload()`s on change. — **PR #13, merged `14b1df6`**

## Parking lot

_(deferred — the spec's §13 non-goals, v2+, out of the MVP envelope)_

- Multi-system board — read the root manifest's system index, render a switcher / aggregate board with per-card system chip. (Spec §13.1 — the natural v2.)
- Self-contained release — `tailwindcss` build + `go:embed` the CSS; offline single binary. (Spec §13.2 / decision D1's other branch.)
- Markdown rendering of captured field text (safe subset + sanitizer) instead of escaped plain text. (Spec §13.3.)
- Journal deep-links — parse `journal.md` so DONE journal pointers / detail pages link to the entry. (Spec §13.4.)
- JSON API (`/api/board.json`) for external consumers. (Spec §13.5.)
- Live push — swap the mtime poll for SSE/`fsnotify` if the poll feels laggy. (Spec §13.6.)
- GitHub Actions CI — build/vet/test on push; would make delivery gate D3 (remote check) meaningful. (Not in the spec; would let the framework's remote-check gate bite.)
