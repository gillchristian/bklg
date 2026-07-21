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

### TASK-016 — Dashboard detail + links out
**Source:** BACKLOG v3 (dashboard adapter; ADR-0004) — final slice
**Acceptance criteria:**
- [ ] Each dashboard card gets a stable url-safe slug from its title; duplicate titles get distinct slugs; an empty slug falls back. *(decider: unit test — `slugify` output + `assignSlugs` uniqueness on a constructed duplicate-title slice.)*
- [ ] `GET /<slug>` renders a dashboard card's detail: title, column, linked ticket chips (to the Linear base), Material, Status, and the collapsed raw block. *(decider: server test — GET the Alpha card's slug → 200 containing the `PINATA-100` Linear link + Material + Status; smoke on the real Pinata KB.)*
- [ ] The board links each dashboard card's title to its `/<slug>` detail page. *(decider: board render test — a dashboard card's HTML contains `href="/<slug>"`.)*
- [ ] Unknown slug → 404. *(decider: server test — `GET /no-such-slug` → 404.)*
- [ ] Framework-mode detail (`/<ID>`) still works unchanged. *(decider: existing `detail_test.go` green; smoke — `bklg .` `/<id>` → 200.)*
**Notes:** Add `Card.Slug` (computed in `parseDashboard` via `assignSlugs`, unique across the board) + `Card.RouteKey()` (value receiver: `ID.Raw` for framework, `Slug` for dashboard) + `Board.CardByRoute` (case-insensitive); `handleTask` looks up by route key. `task.html` gains a `{{if .Card.Dashboard}}` block (subtitle, ticket chips via the now-used `taskVM.LinearBase`, Material, Status); framework blocks stay guarded by their own fields. `board.html` links the dashboard card title to `/<slug>`. This completes the adapter.
