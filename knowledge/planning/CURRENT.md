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

### TASK-006 — Task detail + 404
**Source:** BACKLOG (spec §15.6; detail rendering §7, route key D2)
**Acceptance criteria:**
- [ ] AC1 — Detail for a known id: `GET /<id>` (id matched **case-insensitively**, D2) renders id, title, namespace, current column, and badges. (Decider: `curl -s /DEMO-1` shows `DEMO-1`, its title, `In Progress`, and its badges.)
- [ ] AC2 — State-appropriate fields (§7): an In-Progress card shows Source, Acceptance criteria as a checklist (checked/unchecked), Notes, Delivery override; a Done card shows date, summary, delivery record, journal pointer (plain text — the journal isn't served in v1). (Decider: `curl -s /DEMO-1` shows its 2 criteria + the override; `curl -s /DEMO-4` shows its date/summary/delivery/journal-pointer.)
- [ ] AC3 — Referencing blockers (§7): the page lists blockers whose `Task affected` is this id, **open first then resolved**. (Decider: `curl -s /DEMO-2` shows the open BLOCKER-001; `curl -s /DEMO-1` shows the resolved BLOCKER-002.)
- [ ] AC4 — Collapsed source block (§7): the page always shows a collapsed block containing the card's `Raw` source (escaped). (Decider: `curl -s /DEMO-1` contains a `<details>` with the raw `### DEMO-1 …` block.)
- [ ] AC5 — Unknown id → 404; parking/id-less cards have no detail page (§5); case-insensitive lookup works. (Decider: `curl -o /dev/null -w '%{http_code}' /NOPE-999` → `404`; `/demo-1` → `200`.)

**Notes:**
- Add `task.html` (defines its own `content`) + a `taskTmpl` = layout + task (separate set from `boardTmpl` so both can define `content`). Reinstate `Board.CardByRawID(raw)` (case-insensitive; skips id-less cards) — removed in TASK-003 as premature, now needed.
- `handleTask`: `id := r.PathValue("id")`; look up via `CardByRawID`; 404 (`http.NotFound`) if missing. Collect referencing blockers by matching `parseID(ToUpper(TaskRaw))` to the card id, order open-first. Buffered render like `handleBoard`. Route `GET /{id}` (literal `/_v`,`/_diag`,`/{$}` beat it — spec §7).
- Detail view model carries `PlanningDir`+`Warnings` (so the shared layout header/banner render) plus the `Card` and its referencing `Blockers`.
- Everything escaped via `html/template` (Raw block included). Acceptance rendered as a checklist (☑/☐ or checked styling).
