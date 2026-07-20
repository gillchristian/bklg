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

### TASK-004 — Blocker parse + blocked join
**Source:** BACKLOG (spec §15.4; blocker contract §2, badges §6, `/_diag` route §7)
**Acceptance criteria:**
- [ ] AC1 — Blocker parse: `blockers.md` → `[]Blocker`. Disambiguate blocker headings (`^BLOCKER-\d+`) from section headings (`Format`/`Open`/`Resolved`); track the current section; **skip everything under `## Format`**; each blocker captures id, title, `opened` ts, `**Task affected:**` (the join key), body, and `Open` = (section is `Open`). (Decider: unit test on the fixture asserts BLOCKER-001 is Open affecting DEMO-2, BLOCKER-002 is Resolved affecting DEMO-1, and the `## Format` example is ignored.)
- [ ] AC2 — `blocked` badge join (the load-bearing badge, §6): a card whose id is the `Task affected` of an **open** blocker carries a `blocked` badge; a card affected only by a **resolved** blocker does not. (Decider: test asserts DEMO-2 has `blocked`, DEMO-1 does not.)
- [ ] AC3 — Other badges (§6): `parking` (parking-lot card), `override` (card with a `DeliveryOverride`), `no-ac` (In-Progress card with zero acceptance criteria), and the `namespace` chip (the id's namespace). (Decider: test asserts DEMO-2 → {blocked, no-ac, namespace:DEMO}; DEMO-1 → {override, namespace:DEMO}; the parking card → {parking}.)
- [ ] AC4 — `/_diag` route (§7): `GET /_diag` returns the board's warnings verbatim, one per line, wired into `main` with the board built **per request**; the board build now also parses blockers from the progress area. (Decider: `curl /_diag` against the fixture instance lists the warnings.)
- [ ] AC5 — Zero **unexpected** warnings (spec §15.4): against the fixture, `/_diag` shows only the three seeded reconciliation warnings — no blocker-parse or read warnings. (Decider: `curl -s /_diag` equals exactly the three known lines.)

**Notes:**
- Extend `lineParser.Parse` to also read `blockers.md` from `Areas.ProgressDir` (making `Parse` produce the **complete** board) and to compute `Card.Badges` in one place — `blocked` needs the join, the rest are planning-only. New `parseBlockers` in `parse.go` (or `blockers.go`).
- Section/heading disambiguation per §2: a `##` line matching `^BLOCKER-\d+` is a blocker; `Format`/`Open`/`Resolved` are sections; assign each blocker to the last section seen; skip `Format`.
- Wire the board into `server` (build per request — spec §7 freshness). Add `GET /_diag`. `/` stays the placeholder until TASK-005; `Board.Meta` dirs feed the `/_diag` header line if useful.
- Escaping: `/_diag` is plain text (warnings verbatim); no HTML rendering of repo content (spec §2).
- Consider the review's earlier note: a present-but-unreadable manifest/file could warrant a `/_diag` warning — optional here, don't over-build.
