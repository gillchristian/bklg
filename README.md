# bklg

A read-only, localhost **backlog-board viewer** for a
[knowledge-framework](knowledge/framework/README.md) instance. Point it at a
repo and get a live 3-column kanban — **CURRENT / BACKLOG / DONE** — with badges
for cross-cutting state (blocked, parking-lot, override, namespace) and a detail
page per task.

It parses the semi-structured markdown a framework instance already maintains
(`CURRENT.md`, `BACKLOG.md`, `DONE.md`, `blockers.md`), writes nothing, mutates
no VCS state, and re-reads on every request — so the board tracks a live agent
session.

## Install

```sh
go install github.com/gillchristian/bklg/cmd/bklg@latest
```

Requires **Go 1.22+**. Zero Go module dependencies (stdlib only). Tailwind is
loaded client-side via the Play CDN (the only external artifact); the board
needs network at view time for styling.

## Usage

```
bklg [path] [--port N] [--dir D]

  [path]    path to the repo root                 (default ".")
  --port N  port to listen on                      (default 1235)
  --dir  D  knowledge dir, relative to [path]      (default "knowledge")
```

```sh
$ bklg
Running Backlog on port 1235
  knowledge: ./knowledge   planning: ./knowledge/planning   progress: ./knowledge/progress
  http://localhost:1235
```

The positional path works in any position, and `--port 9000` / `--port=9000`
both parse — so `bklg systems/trail --port 9000` behaves as expected. The server
binds `127.0.0.1` only.

For a single system inside a monorepo, point `--dir` at that system's knowledge
tree: `bklg . --dir systems/trail/knowledge`. Pointed at a monorepo **root**
manifest (many systems, no root planning area), bklg lists the per-system
invocations instead of erroring blankly.

## Routes

| Route | What |
|---|---|
| `/` | the board — three columns, a card per task |
| `/<id>` | task detail (e.g. `/MONO-006`); unknown id → 404 |
| `/_v` | max mtime across parsed files (drives ~3s live-reload poll) |
| `/_diag` | parse warnings, verbatim |

## Status

Built with the knowledge framework it views (see [`knowledge/`](knowledge/)) —
the in-repo instance doubles as the dogfood fixture. Build plan: the seven
executable slices in
[`knowledge/reference/specs/bklg-spec.md`](knowledge/reference/specs/bklg-spec.md)
§15.

## License

MIT (see `LICENSE`).
