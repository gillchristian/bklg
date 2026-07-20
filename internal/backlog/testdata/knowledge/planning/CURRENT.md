# Current task

> One task at a time. When this file is empty, pull the next item from `BACKLOG.md`.

## Entry template

### TASK-NNN — <title>
**Source:** BACKLOG / parking lot / user request
**Acceptance criteria:**
- [ ] criterion (this template section must be ignored by the parser)
**Notes:** ignored.

## Active

### DEMO-1 — Wire up the widget pipeline
**Source:** BACKLOG (spec §4)
**Acceptance criteria:**
- [x] AC1 — the pipeline builds (go build exits 0)
- [ ] AC2 — the widget renders on the demo page
**Notes:** Depends on the exporter landing first. Kept small on purpose.
**Delivery override:** may commit directly to the branch — user, 2026-07-19

### DEMO-2 — Investigate the flaky importer
**Source:** user request
**Notes:** No acceptance criteria yet — should raise the no-ac hygiene badge, and
this task is the target of an open blocker.
