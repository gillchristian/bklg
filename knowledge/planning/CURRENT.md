# Current task

> One task at a time. When this file is empty, pull the next item from `BACKLOG.md`.

## Entry template

### TASK-NNN — <title>
**Source:** BACKLOG / parking lot / user request
**Acceptance criteria:**
- [ ] criterion (how it will be verified)
**Notes:** scope cuts, links, anything decided while planning.
(Add `**Delivery override:** … — user, YYYY-MM-DD` only when the user grants
one; see framework/delivery.md.)

## Active

### TASK-014 — Dashboard parser → model
**Source:** BACKLOG v3 (dashboard adapter; ADR-0004)
**Acceptance criteria:**
- [ ] `## Active` / `## Done` pipe tables parse into `[]Card`: on the `testdata/dashboard` fixture, per-column card counts and each card's `{title, column, tickets, group, blocked, material}` match an expected table quoted in the journal. *(decider: unit test `TestParseDashboard`; `go test ./...` green.)*
- [ ] A literal `\|` inside a cell is not a column split: a row whose status contains `a \|\| b` keeps that text intact in the card. *(decider: unit test asserts the card's text contains `||`.)*
- [ ] `## Backlog` bullet groups: a `**Group:**` subhead sets `Card.Group` on the bullets under it; ungrouped bullets have empty Group. *(decider: unit test on the fixture.)*
- [ ] Title = leading `**bold**` (subtitle after ` — `, U+2014); every inline `[A-Z]+-\d+` → `Card.Tickets` (0..N), `#\d+` PR refs ignored; the Material/Record link is captured. *(decider: unit test asserting tickets slice + material for a multi-ticket row and a zero-ticket row.)*
- [ ] Defensive: a malformed row (no bold title / stray unescaped pipe) becomes a `/_diag` warning and is skipped; the rest of the board still renders; never panics. *(decider: unit test with a seeded bad row asserts a warning + a non-empty board; `go test -race`.)*
- [ ] End-to-end: `bklg --dashboard knowledge/work/index.md ~/dev/Pinata-dev/Pinata` now shows populated columns (card counts > 0 in Active/Backlog/Done). *(decider: smoke — `curl /` card-count grep quoted in the journal.)*
**Notes:** `Card` grows `Tickets []ID`, `Group string`, `Material string`, and a `Dashboard bool` (or reuse an existing flag) — kept optional so framework-mode cards are unaffected. Ticket chips + blocked badge + group chip + Linear links are TASK-015 (rendering); this slice is parse→model only, asserted by unit tests. Ship the rich `testdata/dashboard/` fixture (table Active/Done, bullet-group Backlog, an escaped-pipe row, a leading-`⛔` row, a multi-ticket row, a zero-ticket row, a malformed row).
