# Backlog

Ordered. The top *unchecked* item is next. Promote into `CURRENT.md` when started.

Conventions:
- Completed tasks stay here, checked `[x]` with their delivery record — the
  list doubles as a one-line shipping index. The full record lives in `DONE.md`.
- TASK ids are global and monotonically increasing across all of `planning/`:
  next id = highest id appearing anywhere + 1.

Derived from the spec's §15 "Task breakdown (executable slices)"
(`../reference/specs/bklg-spec.md`). The **first seven** slices below **are** the
v1 MVP (shipped). Below them is the **v2 batch** from user feedback
(2026-07-20). Each is one verifiable vertical slice. Test fixtures (§11) are
built inside the tasks that first need them.

## Active

- [x] TASK-001 — CLI + server skeleton — `splitArgs` (positional-in-any-position, `--flag v`/`--flag=v`), `flag` parsing, `127.0.0.1` bind, `GET /{$}` → 200, startup block. Creates `go.mod` (module `github.com/gillchristian/bklg`), `cmd/bklg/main.go`. — **PR #1, merged `288a814`**
- [x] TASK-002 — Area resolution — Locations-block dereference (repo-root-relative), default `base/planning`+`base/progress` fallback, root-manifest system-index list + helpful exit. Adds resolution fixtures under `testdata/`. — **PR #3, merged `4cf04c1`**
- [x] TASK-003 — Planning parser → model — line-oriented parser for `CURRENT.md`/`BACKLOG.md`/`DONE.md` → `[]Card` with the most-advanced-state dedup and the three reconciliation warnings. Ships the rich `testdata/knowledge/` demo/test instance (§11). — **PR #5, merged `42ff5c9`**
- [x] TASK-004 — Blocker parse + blocked join — parse `blockers.md` (§/blocker heading disambiguation), attach the `blocked` badge to open-blocker tasks; `/_diag` shows warnings. — **PR #7, merged `835363f`**
- [x] TASK-005 — Board template + badges + Tailwind — render `/` (three columns, card per task, badge markup, namespace chip, diagnostics banner) via `html/template` + Tailwind Play CDN. — **PR #9, merged `a9a1fc3`**
- [x] TASK-006 — Task detail + 404 — render `/{id}` (state-appropriate fields, referencing blockers, collapsed Raw block); unknown id → 404. — **PR #11, merged `cd00f36`**
- [x] TASK-007 — Live reload — `/_v` returns max mtime across parsed files; page polls ~3s and `location.reload()`s on change. — **PR #13, merged `14b1df6`**

### v2 — from user feedback (2026-07-20)

Grounded in [`../whiteboard/trail-instance-findings.md`](../whiteboard/trail-instance-findings.md)
(running bklg against the real `trail` monorepo). Ordered: quick fixes that
repair the visibly-broken real-instance experience first, the big multi-system
feature last. AC written when each is promoted to `CURRENT.md`.

- [x] TASK-008 — Trim card text — the board card shows a concise, one-line title (trimmed/ellipsized); the **full** text moves to the detail page. Real instances put whole paragraphs on `[x]`/parking lines, so cards are currently walls of text. (Feedback: card screenshots.) — **PR #16, merged `08c4053`**
- [ ] TASK-009 — Parser robustness for real instances — parse `### <ID> — <title>` heading-style `DONE.md` entries (a `**Completed:** … **PR:** …` line + prose body), *and* recognize task ids wrapped in `**…**`/decorated. Kills trail's 15 false "shipped item missing from DONE.md" warnings, gives Done cards real titles and id-less cards their detail links. Widens the §2 input contract → record an ADR; revisit the D4 goldmark seam. (Feedback: /_diag noise + id-less cards.)
- [ ] TASK-010 — Render a safe markdown subset — bold / italic / inline-code / links / lists in card + detail text instead of literal `**`/`` ` ``. Needs a sanitizer or a strict hand-rolled inline subset → ADR on relaxing zero-dep vs stdlib-only. (Feedback: "Render markdown.")
- [ ] TASK-011 — Make `/_diag` actionable — each warning links to the offending file/card and says what to do; group/dedupe by type; show counts. (TASK-009 removes most of today's noise at the source.) (Feedback: "what's the purpose of /_diag?")
- [ ] TASK-012 — Multi-system board (aggregate + filter) — at a monorepo **root** manifest, instead of the per-system error, read every `systems/<name>` instance, aggregate cards into one board with a per-card **system chip**, and let the user **filter** to any project (e.g. `?system=track` + UI toggles). Handle cross-system id collisions in the detail route. Replaces v1's helpful-error behaviour. Has real UX/CLI design choices — worth a design pass before building. (Feedback #1 + spec §13.1.)

## Parking lot

_(deferred — the spec's §13 non-goals, v2+, out of the MVP envelope. Multi-system
board and markdown rendering were promoted to Active as TASK-012 and TASK-010.)_

- Self-contained release — `tailwindcss` build + `go:embed` the CSS; offline single binary. (Spec §13.2 / decision D1's other branch.)
- Journal deep-links — parse `journal.md` so DONE journal pointers / detail pages link to the entry. (Spec §13.4.)
- JSON API (`/api/board.json`) for external consumers. (Spec §13.5.)
- Live push — swap the mtime poll for SSE/`fsnotify` if the poll feels laggy. (Spec §13.6.)
- GitHub Actions CI — build/vet/test on push; would make delivery gate D3 (remote check) meaningful. (Not in the spec; would let the framework's remote-check gate bite.)
