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

### TASK-005 — Board template + badges + Tailwind
**Source:** BACKLOG (spec §15.5; rendering §7, styling §8, badges §6)
**Acceptance criteria:**
- [ ] AC1 — Three columns: `GET /` renders the three column headings **Backlog**, **In Progress**, **Done**. (Decider: `curl -s /` contains all three heading strings.)
- [ ] AC2 — A card per task: one card block per parsed card, each showing id + title + its badges; a card with an id links to `/<id>`. (Decider: `curl -s /` on the fixture contains each card's id and title; the number of card blocks equals the parsed card count (7).)
- [ ] AC3 — Badges render (§6): the blocked card shows a `blocked` badge; `parking`/`override`/`no-ac`/`namespace` chips render on the right cards. (Decider: `curl -s /` shows `blocked` on DEMO-2's card and the namespace chip `DEMO`.)
- [ ] AC4 — Diagnostics banner + header (§7): a header line shows the resolved planning path; when `len(Warnings)>0` a subtle banner links to `/_diag`. (Decider: `curl -s /` on the fixture contains the planning path and a link to `/_diag`; a clean instance shows no banner.)
- [ ] AC5 — `html/template` + `go:embed` + Tailwind Play CDN, no injection (§2/§7/§8): templates are embedded; captured repo text is auto-escaped (a `<` in a title renders escaped, not as a tag); empty columns show a muted placeholder; the Tailwind Play CDN `<script>` is present. (Decider: `curl -s /` contains the `cdn.tailwindcss.com` script tag; feed a fixture card a title with `<b>` and confirm it appears escaped in the HTML.)

**Notes:**
- Templates live at `internal/backlog/templates/*.html` (embedded via `go:embed` from the `backlog` package — `go:embed` can't traverse `..`, so co-locate rather than a repo-root `templates/`; minor, justified deviation from spec §10). A `layout` partial + `board` template (a `task` template comes in TASK-006).
- Build a small view model (columns → cards) from the `Board`; render with `html/template` (auto-escaping is the injection defense — spec §2, "do not render text as HTML"). Badge → CSS class map for colour (blocked red, parking slate, override amber, no-ac yellow, namespace gray — spec §6 suggested styles).
- Replace `server.go`'s `handleBoard` placeholder with real rendering; keep per-request parse. Header shows `DisplayPath(planning)`; banner only when warnings exist, linking `/_diag`.
- Cards with an id are `<a href="/{id}">`; parking/no-id cards are muted and not linked (no detail page — TASK-006). Empty column → muted placeholder.
- Add a smoke asserting the escaping (a `<` in repo text must not become a live tag).
