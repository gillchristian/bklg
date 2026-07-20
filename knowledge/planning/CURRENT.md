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

### TASK-011 — Make /_diag actionable
**Source:** BACKLOG (v2 feedback; `../whiteboard/trail-instance-findings.md` #5)
**Acceptance criteria:**
- [ ] AC1 — Structured warnings: `Board.Warnings` is `[]Warning{Kind, Message, TaskRaw}`; `readArea`/`parseDone`/`reconcile` populate `Kind` (e.g. `shipped-missing-done`, `done-not-ticked`, `current-multiple`, `read-error`, `malformed-done`) and `TaskRaw` (the id, when one is involved). (Decider: unit test asserts the fixture's warnings carry the right Kinds + TaskRaw.)
- [ ] AC2 — `/_diag` renders an HTML page (via the layout) that **groups** warnings by Kind, with a **count** per group and a one-line **explanation** of what it means and what to do. (Decider: `curl -s /_diag` on the fixture shows the three groups with counts + explanations.)
- [ ] AC3 — Actionable links: every warning referencing a task id links to that task's detail page. (Decider: `curl -s /_diag` contains `href="/DEMO-5"`, `href="/DEMO-6"`, etc.)
- [ ] AC4 — "missing from DONE.md" reframed as informational: its explanation notes it's expected when a repo keeps the full record inline on the `[x]` line rather than duplicating into DONE.md (not a hard error). (Decider: the explanation text is present in `/_diag`.)
- [ ] AC5 — Real-repo effect: on `trail --dir systems/track/knowledge`, `/_diag` shows the 14 shipped-missing warnings as **one grouped section (14) with linked ids + the explanation**, not 14 alarming raw lines; the header banner still links `/_diag`. (Decider: trail smoke — one group, 14 links, explanation present.)

**Notes:**
- Model: new `Warning{Kind, Message, TaskRaw string}`; `Board.Warnings []Warning`. Update the few producers (`readArea` → read-error; `parseDone` → malformed-done; `reconcile` → shipped-missing-done / done-not-ticked / current-multiple with `TaskRaw`). The layout banner just needs `len(.Warnings)` — unchanged.
- `/_diag`: build a `diagVM` grouping `[]Warning` by Kind → `{Title, Explanation, []Warning}`; render an HTML page (reuse `layout`; a new `diag.html` defining `content`). Each `TaskRaw` → `<a href="/{id}">`. A small `explain(kind)` map supplies titles + how-to-fix text. Keep it escaped/safe (ids via the same route key; text via auto-escape).
- Don't add a fragile suppression heuristic for shipped-missing — grouping + explanation is the fix (honest + actionable). If a repo wants the check off entirely, that's a future config flag (parking), not this task.
- Update `TestParseWarnings` (assert on `w.Kind`/`w.Message`) and `TestDiagRoute` (assert the grouped HTML + links) for the new shape.
