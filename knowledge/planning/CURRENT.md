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

### TASK-015 — Dashboard badges + board render
**Source:** BACKLOG v3 (dashboard adapter; ADR-0004)
**Acceptance criteria:**
- [ ] Dashboard `Blocked` cards show a red `blocked` badge on the board and detail page; non-blocked ones don't. *(decider: render/unit test + smoke — the leading-`⛔` Active row on the real Pinata KB renders a `blocked` badge; grep the HTML.)*
- [ ] Each card's `Tickets` render as chips linking to `<linear-base><ID>` (default `https://linear.app/gopinata/issue/PINATA-602`). *(decider: smoke — board HTML contains `<a href="https://linear.app/gopinata/issue/PINATA-...">`; unit test on the href builder.)*
- [ ] Linear base configurable via a `--linear-base` flag and a `linear:` Locations key (flag wins), trailing slash tolerated. *(decider: smoke — `--linear-base https://example.com/i/` makes ticket hrefs use that base; unit test.)*
- [ ] Backlog cards show their `Group` as a chip. *(decider: smoke/render — a `Product / code` chip appears on the relevant card.)*
- [ ] Board tolerates AC-less, id-less, multi-ticket dashboard cards: no spurious `no-ac` badge, no broken `/{id}` link (dashboard cards have no single ID). *(decider: render test / smoke — dashboard In-Progress cards carry no yellow `no-ac` badge and no `/<id>` anchor.)*
**Notes:** Rendering slice — the model (Blocked/Tickets/Group) already lands from TASK-014. Add a `group` badge kind + a ticket-chip block (anchors) to `board.html`/`task.html`; thread the Linear base through `Meta`/view models (resolved in `Resolve`/`ResolveDashboard` + `main`). Dashboard cards compute their own badges (blocked + group) — do NOT run the framework `computeBadges` (its `no-ac`/blocker join don't apply). Detail routing for dashboard cards (title slug) is TASK-016.
