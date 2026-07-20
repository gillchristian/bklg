# Done

Tasks that passed all verification gates. Newest at top.

Each entry: id, title, completion date, a summary (a sentence for trivial
tasks, a paragraph when the detail earns its keep — the journal holds the
full story), the delivery record (``PR #N, merged `sha` `` under the `pr`
profile), and a journal pointer.

- TASK-NNN — title — YYYY-MM-DD — summary — PR #N, merged `sha`. See journal YYYY-MM-DD [HH:MM]

## Completed

- TASK-001 — CLI + server skeleton — 2026-07-20 — `splitArgs`-based CLI (positional in any position; `--flag v` and `--flag=v`), flags `--port`/`--dir`, loopback-only bind with clean port-in-use failure, and a `GET /{$}` handler returning a 200 HTML placeholder; startup block matches spec §9 verbatim. Fresh-context review passed all 5 AC; its two non-blocking notes (silent extra-positional drop, `-h` exit code) were fixed and covered by table-driven `splitArgs`/`joinDisplay` tests. — PR #1, merged `288a814`. See journal 2026-07-20 14:00.
