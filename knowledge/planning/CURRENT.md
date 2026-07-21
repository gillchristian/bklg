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

### TASK-013 — Dashboard-mode resolution
**Source:** BACKLOG v3 (dashboard adapter; ADR-0004)
**Acceptance criteria:**
- [ ] Manifest lookup tries `README.md` then `index.md`: a KB whose only manifest is `index.md` carrying a `## Locations` `dashboard:` key resolves to dashboard mode. *(decider: unit test — `Resolve` on a `testdata` dashboard KB returns `Areas.DashboardFile` = the resolved file, no error though there is no `planning/`; `go test ./...` green.)*
- [ ] A `dashboard: <path>` Locations key sets `Areas.DashboardFile = <path>/<value>` and lifts the planning-dir requirement. *(decider: same unit test asserts the exact path.)*
- [ ] `--dashboard <file>` flag forces dashboard mode with no Locations block, resolving `<file>` against `[path]`. *(decider: unit test for the resolve helper + smoke — `bklg --dashboard knowledge/work/index.md ~/dev/Pinata-dev/Pinata` prints a `dashboard: …/work/index.md` startup line and `curl /` → 200.)*
- [ ] A missing dashboard file exits non-zero with a clear message. *(decider: smoke — `bklg --dashboard nope.md .`; exit 1, stderr contains `no dashboard file`.)*
- [ ] Framework mode unchanged. *(decider: the existing `resolve_test.go` cases still pass; smoke — `bklg .` on this repo prints the same `knowledge/planning/progress` line as before.)*
**Notes:** Resolution slice only (mirrors TASK-002). The dashboard *parser* is TASK-014 — in this slice the parser dispatches on `Areas.DashboardFile` and returns an empty board, so dashboard mode starts and serves 200 with empty columns (parsing pending). Flips ADR-0004 to accepted (its own condition: "when TASK-013 is promoted"). Linear base URL default `https://linear.app/gopinata/issue/` confirmed by user 2026-07-21.
