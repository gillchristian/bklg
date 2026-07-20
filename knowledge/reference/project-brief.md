# Project brief — bklg

*The brief wins conflicts with planning. Full detail: `specs/bklg-spec.md`.*

## What it is

`bklg` — a read-only, localhost **backlog-board viewer** for a knowledge-framework
instance. Point it at a repo; it parses the semi-structured markdown a framework
instance already maintains (`CURRENT.md`, `BACKLOG.md`, `DONE.md`, `blockers.md`)
and serves a live 3-column kanban (CURRENT / BACKLOG / DONE) with cross-cutting
badges (blocked, parking, override, namespace) plus a detail page per task.
One line: *point it at a repo, get a localhost kanban with blocked badges and a
page per task.*

## Why it exists

Reading the raw planning/progress files doesn't show cross-cutting state at a
glance — especially which tasks are **blocked** (the load-bearing badge). The
board joins planning + progress so that state is visible, and re-reads on every
request so it tracks a live agent session.

## Hard constraints

- **Go 1.22+** (needs method + `{id}`/`{$}` `ServeMux` patterns).
- **Zero Go module dependencies** — stdlib only (`net/http`, `html/template`,
  `embed`, `flag`, `os`, `time`, `strings`, `regexp`, `bufio`).
- `go install github.com/gillchristian/bklg/cmd/bklg@latest`-able; command `bklg`.
- **Writes nothing, mutates no VCS state.** Read-only viewer.
- **Never render repo text as HTML** — captured field text is escaped via
  `html/template` auto-escaping (immune to injection from repo content).
- **Parse defensively** — a malformed entry must never crash the server or blank
  the board; capture what's parseable, skip the rest, record a warning → `/_diag`.
- Bind **127.0.0.1** by default (personal dev tool, off the network).
- Tailwind via **Play CDN** (the only external artifact, client-side) — keeps the
  tool a Node-free, pure-stdlib `go install`. (Decision D1; embedded-CSS is a v2
  self-contained-release upgrade.)

## Out of scope (v1)

Editing/mutating any file; auth; DB; bundler. And the §13 non-goals (v2+):
multi-system board, self-contained embedded-CSS release, markdown rendering of
field text, journal deep-links, JSON API, live push (SSE/fsnotify).

## Stack

Go stdlib. `html/template` templates embedded with `go:embed`. Tailwind Play CDN
client-side. Go 1.22 `ServeMux` pattern routing (`GET /{$}`, `/_v`, `/_diag`,
`/{id}`).

## Success criteria

The seven executable slices of spec §15 (TASK-001…007), each with its
fresh-checkable AC, all shipped: CLI+server skeleton; area resolution; planning
parser+dedup; blocker parse+blocked join; board+badges+Tailwind; task detail+404;
live reload. That set = the v1 MVP.

## Open questions

- Framework upstream repo/path for `framework/` is unrecorded (pre-placed by the
  user). Non-blocking; note in the manifest, fill in when known.
- Spec decisions D1–D4 are all confirmed as the spec recommends (Play CDN;
  route key = full task id; one planning area in v1; zero-dep line-scanner behind
  a swappable interface). No flips requested.

## Raw notes (user words)

- "See the SETUP.md file. Run the initialization of the framework. The spec for
  the project is bklg-spec.md. Once done, process the spec into tasks. Once
  that's done, work on the tasks until the project is done. Go with the PR way of
  working (ie. autonomous working)."
- "Goal set: Init the framework, split spec into tasks, ship all the tasks which
  make up MVP"
- Attribution (2026-07-20): "Your identity only" — no agent attribution on
  commits/PRs.
