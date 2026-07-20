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

### TASK-002 — Area resolution
**Source:** BACKLOG (spec §15.2, detail in §3)
**Acceptance criteria:**
- [ ] AC1 — Locations dereference: for a fixture repo whose `knowledge/README.md` has a `## Locations` block, resolution returns the `planning:`/`progress:` paths from that block, resolved **repo-root-relative** against `path` (spec §3 step 1). (Decider: a `resolve` unit test asserts `PlanningDir`/`ProgressDir` equal the expected fixture paths; the startup block echoes the same resolved paths.)
- [ ] AC2 — Default fallback: for a fixture with a `knowledge/` dir but **no** `## Locations` block (or no manifest), resolution returns `base/planning` and `base/progress` where `base = path/dir` (spec §3 step 2). (Decider: unit test asserts the two resolved paths.)
- [ ] AC3 — Root-manifest system list: pointed at a fixture **root manifest** (its resolved `planningDir` does not exist, but the manifest has a table with `systems/<name>` rows), `bklg` exits **non-zero** and prints each discovered system with the exact `bklg` invocation to board it (spec §3 step 3 / §9). (Decider: run it; assert exit≠0 and stdout lists the fixture's `systems/*` names + a `bklg … --dir systems/<name>/knowledge` line each.)
- [ ] AC4 — No planning area and not a root manifest: exits non-zero with `no planning area at <planningDir>` (spec §3 step 3 tail). (Decider: run against a dir with neither a planning area nor a systems table; assert the message + non-zero exit.)
- [ ] AC5 — `path` not a directory: `bklg <missing-or-file>` exits non-zero with a clear message before starting a server (spec §9 failure mode). (Decider: run against a nonexistent path and against a regular file; assert non-zero exit + message.)

**Notes:**
- New file `internal/backlog/resolve.go`: `Resolve(path, dir string) (Areas, error)` returning `{KnowledgeDir, PlanningDir, ProgressDir}` (+ a typed "root manifest → systems" signal), unit-tested. `main.go` wires it in and the startup block echoes the **resolved** paths (replacing TASK-001's naive echo). On a resolution error, print the message to stderr and exit non-zero — no server.
- Locations-block parse (spec §3): on `## Locations`, enter block; until the next line starting `## `, split each non-empty line on the **first** `:` into key/value, trim; keep `planning`/`progress`, ignore the rest.
- System-index parse (spec §3): scan `|`-delimited rows for a cell matching `systems/\S+`; collect **distinct** dir names.
- Resolution uses real filesystem paths (`filepath`); keep a display form for the echo. Fixtures needed: (a) a repo with a Locations manifest; (b) a repo with a default (no-Locations) `knowledge/`; (c) a root manifest with a `systems/*` table and no planning area; (d) a bare dir with neither. Put them under `testdata/` (the rich §11 parser instance is TASK-003).
