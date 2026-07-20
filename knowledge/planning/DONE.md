# Done

Tasks that passed all verification gates. Newest at top.

Each entry: id, title, completion date, a summary (a sentence for trivial
tasks, a paragraph when the detail earns its keep — the journal holds the
full story), the delivery record (``PR #N, merged `sha` `` under the `pr`
profile), and a journal pointer.

- TASK-NNN — title — YYYY-MM-DD — summary — PR #N, merged `sha`. See journal YYYY-MM-DD [HH:MM]

## Completed

- TASK-002 — Area resolution — 2026-07-20 — `internal/backlog.Resolve`: dereferences the manifest's `## Locations` block (repo-root-relative), falls back to `base/planning`+`base/progress`, and detects a multi-system root manifest (`systems/<name>` index) to list per-system invocations rather than erroring blankly; `main` echoes resolved paths and exits non-zero with a clear message (no server) on any failure. 8 unit tests + 5 `testdata/resolve` fixtures. Fresh-context review passed all 5 AC, no findings; two notes addressed (documented key matching, added a partial-block test). Dogfood: `bklg .` resolves this repo via its real Locations block. — PR #3, merged `4cf04c1`. See journal 2026-07-20 14:19.
- TASK-001 — CLI + server skeleton — 2026-07-20 — `splitArgs`-based CLI (positional in any position; `--flag v` and `--flag=v`), flags `--port`/`--dir`, loopback-only bind with clean port-in-use failure, and a `GET /{$}` handler returning a 200 HTML placeholder; startup block matches spec §9 verbatim. Fresh-context review passed all 5 AC; its two non-blocking notes (silent extra-positional drop, `-h` exit code) were fixed and covered by table-driven `splitArgs`/`joinDisplay` tests. — PR #1, merged `288a814`. See journal 2026-07-20 14:00.
