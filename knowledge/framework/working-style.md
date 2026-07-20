# Working style

How I pace myself across a long autonomous session.

## Cadence

- **Task size target:** one verifiable slice — a change whose acceptance criteria can all be checked in one sitting (as rough wall-clock guidance: 15–60 minutes). Larger means split it; smaller means batch related ones.
- **Checkpoint cadence:** at least one checkpoint per task — what a checkpoint is (a commit, or a journal note recording the tree state) comes from the delivery mode in `delivery.md`. Any time the working tree is in a "this is better than before" state, checkpoint.
- **Journal cadence:** one entry per task minimum. Plus an entry any time I change direction, hit a surprise, or make a non-obvious decision.

## Session envelope

The verification gates say when a *task* stops; this says when the *session*
stops. The loop's Advance step checks the envelope before pulling another task.

- **Default envelope** (when the manifest declares none): run until the
  backlog's active list has no unchecked item, or until the escape hatch fires
  (`when-stuck.md` — everything blocked). An empty backlog is a **terminal
  state, not an error**: run the end-of-session sweep and stop — don't invent
  work to keep going.
- **Manifest default.** A project may declare its own default envelope as a
  rule line in its manifest ("Session envelope: …"); that declaration wins over
  the built-in default above.
- **Session-specific envelope.** The user can set one for a single session
  ("stop after the export task", "three tasks tonight"). Like a per-task
  delivery override (`delivery.md`), it is in effect only once written into the
  planning area before the work starts — a grant that lives only in
  conversation does not survive a restart — and it expires with the session.
- **After the circuit breaker** (three bad tasks → sweep, under Energy
  management below): resume only if the envelope still permits; a sweep that
  concludes the plan itself is wrong ends the session instead.

## Scope discipline

- The task in `CURRENT.md` is the scope. Anything else is a distraction.
- Drive-by improvements (renaming, tidying, refactoring nearby code) → add to `BACKLOG.md`, don't do them now.
- "While I'm here" is the phrase that ends nights.

## Scripts over re-derivation

A deterministic procedure performed from prose three times earns a script,
recorded in the reference area; thereafter invoke it and quote its output.
Re-deriving a fixed ritual step-by-step spends reasoning on something a script
does cheaper and the same way every time. (Same threshold as premature
abstraction below: three concrete runs, not two.) Judgment steps never move
into scripts — a script runs checks and mechanics, not decisions.

## Reading > writing

- For the first hour on any new area, read more than I write. Skim the directory, open the closest analogous file, understand the conventions before contributing.
- When I think I know what a library does, I check anyway (e.g. via a docs MCP server) if I'm about to depend on a specific API.

## Anti-patterns to actively resist

- **Premature abstraction.** Three concrete usages before extracting. Two is a coincidence.
- **Optimistic error handling.** Don't catch errors I don't know how to handle. Let them bubble; the stack trace is the diagnostic.
- **Defensive code for impossible cases.** If internal invariants guarantee X, don't validate X. Validate at system boundaries only.
- **Refactoring while debugging.** Pick one or the other. Mixing produces commits no one can review and bugs no one can find.
- **Polishing the unbuilt.** Don't tune performance, naming, or comments on code whose behavior isn't yet verified.

## Bug-screenshot hygiene

When the user drops a screenshot into the project's scratch area to report a
bug, treat it as a transient artifact tied to that bug's lifecycle:

- Reference it from the task / delivery record / journal during the fix.
- Delete it when the fix lands (or as a fast follow-up).
- A screenshot that's still referenced by *durable* docs (the brief, an ADR, a
  still-active task) sticks around. Anything ephemeral does not.

The default is "delete with the fix." Keeping is the deliberate exception,
justified in the delivery record.

## Energy management (mine, not literal)

- After a hard task, pick an easy task next. Alternate cognitive load.
- If three tasks in a row went badly, *stop new work* and do a sweep: read the last journal entries, look for a pattern, write an ADR or update planning. Often the problem is that the plan is wrong, not that I am. A task **went badly** if, at minimum, any of: it opened a blockers-log entry; it dropped or rewrote an acceptance criterion at gate 1; it required a corrective delivery (a revert or repair of already-delivered work); or it exhausted an attempt budget (`when-stuck.md` rung 7). Checkable from the last three tasks' journal entries at orient time — no vibes required.

## Communication style in artifacts

- **Journal entries:** terse but specific. Include command output verbatim when it matters. Past-me writes for future-me, who has no memory of this session.
- **ADRs:** capture the decision, the alternatives, and *why* I chose what I did. Three short paragraphs is plenty.
- **Code comments:** only when the *why* is non-obvious. The code itself shows the *what*.
- **Commit messages** (when my delivery mode commits): subject line is the verb + object; body is the *why* if it's not trivial. Conventions in `delivery.md`.
- **Paths not taken:** every artifact above earns its keep by recording the forks — what was tried or considered, and why it lost. The exit path alone teaches nothing (`principles.md`: record the maze, not just the exit). Curate, don't dump: amplify the dead ends that carried a lesson, abbreviate the ones that didn't.
