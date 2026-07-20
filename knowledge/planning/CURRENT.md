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

### TASK-007 — Live reload
**Source:** BACKLOG (spec §15.7; `/_v` + poll §7)
**Acceptance criteria:**
- [ ] AC1 — `/_v` returns the max mtime as a bare integer: `GET /_v` returns the maximum modification time across the parsed files (CURRENT/BACKLOG/DONE/blockers.md) as a bare integer, `text/plain`; after `touch`-ing a parsed file the value strictly increases. (Decider: `curl -s /_v` is an integer; `touch` a parsed file, `curl -s /_v` again → larger value.)
- [ ] AC2 — `Board.Meta.LatestMTime` computed: the parser sets it to the max mtime across the parsed files (was stubbed zero since TASK-003). (Decider: unit test — parse the fixture, assert `LatestMTime` is non-zero and equals the newest of the four files' mtimes.)
- [ ] AC3 — Page polls + reloads: `/` (and `/<id>`) embed ~8 lines of vanilla JS that fetch `/_v` every ~3s and call `location.reload()` when the value differs from the value baked into the page. (Decider: `curl -s /` contains a script referencing `/_v`, `location.reload`, and the current version value.)
- [ ] AC4 — End-to-end freshness: with the board loaded, editing a parsed file changes `/_v` so the next poll reloads. (Decider: scripted — record `/_v`, touch a parsed file, confirm `/_v` changed; note the manual reload observation.)
- [ ] AC5 — `/_v` precedence + shape: `/_v` is a literal route (beats `/{id}`), returns 200 with just the integer (no HTML). (Decider: `curl -s -o /dev/null -w '%{http_code}' /_v` → 200; body matches `^-?\d+$`.)

**Notes:**
- `parse.go`: add `latestMTime(paths...) time.Time` (stat each, ignore missing, return max); set `b.Meta.LatestMTime` in `Parse` over CURRENT/BACKLOG/DONE + blockers.md. Add `time` import.
- `server.go`: `handleVersion` → `versionString(board)` = `strconv.FormatInt(mtime.UnixNano(),10)` (or `"0"` when zero); route `GET /_v`. Same `versionString` feeds the page so the baked-in value and `/_v` match exactly.
- `render.go`: add `Version string` to `boardVM` **and** `taskVM`; `viewModel` sets it; `handleTask` sets it. `layout.html` gets the poll `<script>` using `{{.Version}}` (html/template JS-context escaping; integer value is safe).
- Poll interval ~3s; `location.reload()` on change; wrap in try/catch-ish `.catch` so a transient fetch error doesn't spam. This is the only client JS in the tool (spec §7).
- Verify the increase with nanosecond mtime (macOS APFS); if a coarse FS makes it flaky, write changed content rather than bare `touch`.
