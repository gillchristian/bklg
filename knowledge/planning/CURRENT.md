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

### TASK-001 — CLI + server skeleton
**Source:** BACKLOG (spec §15.1)
**Acceptance criteria:**
- [ ] AC1 — Builds clean: `go build ./...` exits 0 and `go build -o bklg ./cmd/bklg` produces the binary. (Decider: exit code 0.)
- [ ] AC2 — Startup block: `bklg --port 9001 .` prints a block whose first line is exactly `Running Backlog on port 9001`, then the `knowledge:` / `planning:` / `progress:` resolution echo, then `http://localhost:9001`. (Decider: capture stdout; assert the first line string and the `http://localhost:9001` line are present.)
- [ ] AC3 — Serves `/`: with the server up on 9001, `curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9001/` prints `200`. (Decider: HTTP status == 200.)
- [ ] AC4 — Arg-order robustness: `bklg . --port 9001`, `bklg --port 9001 .`, and `bklg --port=9001 .` all bind 9001 and serve 200 — positional path parses in any position, and both space- and `=`-form flags parse. (Decider: run all three; each startup line shows `port 9001` and a follow-up curl returns 200.)
- [ ] AC5 — Loopback only: the listener binds `127.0.0.1`, not `0.0.0.0` (spec §9). (Decider: code passes `127.0.0.1:<port>` to the server; `curl http://127.0.0.1:9001/` → 200 confirms loopback reachability.)

**Notes:**
- Creates `go.mod` (`module github.com/gillchristian/bklg`, `go 1.22`) and `cmd/bklg/main.go`. **Scope refinement (2026-07-20):** the `go:embed` templates seam is deferred to TASK-005, where real templates arrive — planting a placeholder now would only be churned then, and TASK-001's AC don't need it. `/` returns a minimal static HTML placeholder for now.
- Flag defaults per spec §9: `--port` 1235, `--dir` knowledge. First startup line is literally `Running Backlog on port N` (§9 wording — "Backlog", not "bklg").
- Resolution here is the **skeleton** only: echo `path/dir`, `path/dir/planning`, `path/dir/progress` as plain defaults, with **no** Locations dereference and **no** existence checks — real resolution (Locations, fallback, root-manifest system list) is TASK-002.
- `splitArgs` per spec §9 (pre-split argv, ~15 lines, zero deps) so `flag` doesn't choke on a positional after flags.
