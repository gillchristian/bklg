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

### TASK-012 — Multi-system board (aggregate + filter)
**Source:** BACKLOG (v2 feedback #1; spec §13.1; `../whiteboard/trail-instance-findings.md` #6)
**Acceptance criteria:**
- [ ] AC1 — Root manifest aggregates instead of erroring: `bklg <repo>` at a monorepo root (planning area absent, a `systems/<name>` index present) resolves **each** system's instance and serves **one** board combining all their cards, rather than the v1 exit-with-list. (Decider: `bklg /Users/bb8/dev/trail` → HTTP 200 board; cards from ≥2 systems present.)
- [ ] AC2 — Per-card system chip: in aggregate mode each card shows which system it came from (e.g. `track`). Single-system mode is unchanged (no chip). (Decider: board HTML shows a `track` chip on a track card; the single-system fixture board has no system chip.)
- [ ] AC3 — Filter to one project: `GET /?system=<name>` shows only that system's cards; a filter bar links each system (+ "All"). (Decider: `/?system=track` shows only track cards; `/` shows all; the bar lists the systems.)
- [ ] AC4 — Detail + `/_v` across systems: `/<id>` resolves against the aggregated cards (renders, with the system chip); `/_v` is the max mtime across **all** systems' files. (Decider: a known aggregated id → 200 detail; `/_v` changes when any system's file changes.)
- [ ] AC5 — ADR-0003 records the design: root-manifest auto-aggregates; server-side `?system=` filter (no client JS); per-card `System` tag/chip; detail global lookup (cross-system id-collision caveat noted). A system in the index that fails to resolve is skipped with a `/_diag` warning, not a crash. (Decider: `decisions/0003-*.md` exists + indexed; a bogus system entry → warning, board still serves.)

**Notes:**
- Reuse `Resolve(path, sys+"/knowledge")` per system from `RootManifestError.Systems`. Add `Card.System string` (empty single-system → no chip/filter). New server mode: `NewMultiServer(path, systems)` resolving each; `board()` parses each, tags `System`, concatenates (no cross-system dedup — namespaces differ per system; each system's board is already deduped). `main`: on `RootManifestError`, build a multi-server instead of exiting; startup block notes "aggregate: N systems".
- Filter server-side: `?system=` filters the aggregated cards at render; a `systemFilter` bar in `boardVM` (name, count, active, href). `/_diag` aggregates warnings (tag each with its system? keep simple: prefix Message or add a system field — minimal). `/_v` = max over all systems.
- Detail: search aggregated cards for the id (case-insensitive). Cross-system collision (same namespace+number in two systems) is unlikely (namespaces differ); first match wins — note in the ADR.
- A system whose instance won't resolve (missing planning) → skip + warning, don't fail the whole board.
- Biggest task; if it grows past one clean slice, split (e.g. aggregate-core vs filter-UI) — but aim to land it whole. Record ADR-0003 with the design forks (this is the "sane defaults + ADR" the user opted into).
