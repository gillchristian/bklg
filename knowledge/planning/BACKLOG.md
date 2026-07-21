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
- [x] TASK-009 — Parser robustness for real instances — parse `### <ID> — <title>` heading-style `DONE.md` entries (a `**Completed:** … **PR:** …` line + prose body), *and* recognize task ids wrapped in `**…**`/decorated. Widens the §2 input contract (ADR-0001). NB: on trail this fixed the DONE **format** parsing (TRACK-000) but only 1/15 warnings — the other 14 are the check correctly flagging that trail keeps records in the `[x]` line, not DONE.md → TASK-011. (Feedback: /_diag noise + id-less cards.) — **PR #18, merged `d0e33a5`**
- [x] TASK-010 — Render a safe markdown subset — bold / italic / inline-code / links in card + detail text instead of literal `**`/`` ` ``. Hand-rolled escape-first inline renderer, stdlib-only (ADR-0002); asterisk-only emphasis (no identifier mangling). Passed an adversarial security review (no injection). — **PR #20, merged `ff2ea83`**
- [x] TASK-011 — Make `/_diag` actionable — structured `[]Warning{Kind,Message,TaskRaw}`; `/_diag` is an HTML page grouping warnings by kind with counts, explanations, and links to task detail pages; "missing from DONE.md" reframed as informational. (Feedback: "what's the purpose of /_diag?") — **PR #22, merged `ce95127`**
- [x] TASK-012 — Multi-system board (aggregate + filter) — at a monorepo **root** manifest, aggregate every `systems/<name>` instance into one board with a per-card **system chip** and a server-side `?system=` filter bar (lists every system); global detail lookup; `/_v` across systems; unresolvable systems skipped with a warning. Replaces v1's error. ADR-0003. On real trail: 5 systems, ~105 cards. (Feedback #1 + spec §13.1.) — **PR #24, merged `95acdf7`**

### v3 — dashboard adapter (Pinata-shape KBs) (2026-07-21)

Grounded in this session's investigation of the real Pinata KB
(`~/dev/Pinata-dev/Pinata/knowledge`): it runs a full Active/Backlog/Done
lifecycle but in a shape the framework parser can't read — one file
(`work/index.md`), Active/Done as pipe **tables**, Backlog as bullet groups
under bold subheads, identity via inline Linear ids (`PINATA-\d+`, 0..N per row)
not one structured id per card, and "blocked" as prose (`⛔`). This batch
teaches bklg a second input convention (a "dashboard adapter") reading that
shape; it pairs with a light dashboard-format contract the target KB follows
([`../reference/specs/dashboard-format.md`](../reference/specs/dashboard-format.md)).
Two-sided by design: the parser is defensive, but a useful board needs the input
to stay regular. ADR-0004 records the second-convention decision (status
**proposed** — flip to accepted when TASK-013 is promoted). AC written when each
is promoted to `CURRENT.md`.

- [x] TASK-013 — Dashboard-mode resolution — `Resolve` tries `README.md` then `index.md` as the manifest; a new `dashboard:` Locations key (single file, repo-root-relative) plus a `--dashboard <file>` escape-hatch flag select dashboard mode; in that mode the `planning/` dir requirement is lifted and the target is the one dashboard file. Framework mode unchanged. (ADR-0004.) — **PR #27, merged `6cea51f`**
- [x] TASK-014 — Dashboard parser → model — parse `## Active`/`## Done` pipe tables (split on unescaped `|`; `\|` literal) + `## Backlog` bullet groups (a `**Group:**` subhead sets the group label) into `[]Card`: title = leading `**bold**` (short/subtitle split on ` — `, U+2014), column from section, every inline `[A-Z]+-\d+` → `Card.Tickets` (0..N; `#\d+` PR refs ignored), Material link captured; defensive → warnings, never panics. Ships `testdata/dashboard/` mirroring the Pinata shapes. `Card` grows an optional multi-ticket slice + a dashboard flag. — **PR #29, merged `2697702`**
- [ ] TASK-015 — Dashboard badges + board render — `blocked` from a **leading** `⛔`/`**Blocked**` status marker (there is no `blockers.md` in this mode); Linear ticket chips with a configurable base (`linear:` Locations key / `--linear-base`, default `https://linear.app/gopinata/issue/`); Backlog group chip; board template tolerates AC-less, multi-ticket cards.
- [ ] TASK-016 — Dashboard detail + links out — `/{slug}` (title-slugged, collisions disambiguated) → detail page: title, column, linked tickets, Material link, and the raw row block; unknown slug → 404.

## Parking lot

_(deferred — the spec's §13 non-goals, v2+, out of the MVP envelope. Multi-system
board and markdown rendering were promoted to Active as TASK-012 and TASK-010.)_

- Self-contained release — `tailwindcss` build + `go:embed` the CSS; offline single binary. (Spec §13.2 / decision D1's other branch.)
- Journal deep-links — parse `journal.md` so DONE journal pointers / detail pages link to the entry. (Spec §13.4.)
- JSON API (`/api/board.json`) for external consumers. (Spec §13.5.)
- Live push — swap the mtime poll for SSE/`fsnotify` if the poll feels laggy. (Spec §13.6.)
- GitHub Actions CI — build/vet/test on push; would make delivery gate D3 (remote check) meaningful. (Not in the spec; would let the framework's remote-check gate bite.)
- Dashboard auto-detect — infer dashboard mode when a KB has no `planning/` but a `work/index.md`, dropping the explicit `dashboard:` Locations key. (Deferred from ADR-0004: explicit config is safer for v1.)
- Linear status sync — cross-check dashboard Active rows against live Linear status on read (mirrors the Pinata KB's own sync-on-read habit). Needs network + auth; likely stays a non-goal for a zero-dep localhost tool.
- Tighten dashboard ticket detection — the `[A-Z]+-\d+` scan (TASK-014) also matches `UTF-8`/`SHA-256`-style tokens in freeform prose, which would render as (dead) Linear chips. Per the `dashboard-format.md` contract for now; revisit with a project-prefix allowlist or a word-boundary guard if it proves noisy once chips render (TASK-015). (Fresh-context review, TASK-014.)
