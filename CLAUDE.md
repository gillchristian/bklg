# CLAUDE.md — bklg

**Read these, in order, before doing anything else:** `knowledge/README.md`
(the manifest: delivery mode + Locations + project rules),
`knowledge/framework/README.md` (the working system), then the enabled `pr`
profile in `knowledge/framework/delivery.md`. The summary below is just the
headline rules so you can't accidentally violate them while still loading the
rest.

## Non-negotiables

1. **One task at a time.** Pull from `knowledge/planning/CURRENT.md`; if empty,
   promote the top unchecked item of `knowledge/planning/BACKLOG.md`.
   Acceptance criteria (each naming its decider) before code.
2. **Delivery: pr.** `main` is sacrosanct — NEVER commit on or push to it
   directly (the one bootstrap commit is the documented exception). Every change
   lands via a PR I open, get a fresh-context review on, and merge myself with
   `--squash --delete-branch`. Read-only git is fine.
3. **Attribution:** commits/PRs carry `gillchristian <gillchristiang@gmail.com>`
   only — no `Co-Authored-By`, no "generated with" footer.
4. **Verify before declaring done.** Gates in
   `knowledge/framework/verification.md`. Run the program, quote actual output;
   don't confuse "compiles" with "works." Local CI: `knowledge/reference/local-ci.md`.
5. **Journal everything.** Append to `knowledge/progress/journal.md` after every
   task.
6. **When stuck, follow `knowledge/framework/when-stuck.md`.** Don't ask the
   user; log real blockers to `knowledge/progress/blockers.md`, then pivot.

## This project

`bklg` is a zero-dependency (stdlib-only) Go 1.22+ tool: a read-only localhost
kanban viewer for a framework knowledge instance. Spec:
`knowledge/reference/specs/bklg-spec.md`. Build order = the seven slices in the
spec's §15 (TASK-001…007). Tailwind via Play CDN; templates via `go:embed`.
