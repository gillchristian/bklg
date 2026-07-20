# When stuck

The user is not available. I cannot ask. This document is what I do instead.

## Diagnostic ladder

Walk down this ladder when blocked. Don't skip rungs.

1. **Re-read the task.** Open the planning area's `CURRENT.md`. Did I drift from the acceptance criteria? Often the answer is "yes, I'm solving a problem the task didn't ask me to solve."
2. **Re-read recent journal entries.** The journal (last ~5 entries). Did past-me leave a hint? Did I already try this and learn something?
3. **Read the code.** Not the code I'm writing — the code around it. Most "I don't know how X works" questions are answered by reading 50 more lines of existing code.
4. **Read the docs.** For external libraries/APIs, fetch fresh documentation with whatever docs tooling the environment provides (e.g. a context7 MCP server, or the published docs themselves). Training data is stale; published docs are not.
5. **Reduce the problem.** Build the smallest possible reproduction. If I can't reproduce the bug in isolation, my mental model is wrong somewhere.
6. **Try the boring option.** When choosing between a clever path and an obvious path, pick obvious. The boring option is reversible; the clever one usually isn't.
7. **Three attempts and pivot.** If three distinct attempts at the same obstacle have produced no new information — not just no fix, no new *fact* — stop. Number the attempts as I go; the blocker entry lists them. Anchor any time reasoning to clocks I can actually read (`date`, journal-entry timestamps), never to felt time. Write down (a) what each attempt tried, (b) what I observed, (c) what I currently believe, (d) what would unblock me. Move that to the blockers log (the progress area's `blockers.md`) and pick a different task from the backlog.

## Classes of "stuck" and the right response

| Symptom | Response |
|---|---|
| Ambiguous requirement | Make a reasonable interpretation, log it as an ADR in the decisions area, continue. Note in the blockers log only if the interpretation might be wrong in a way that wastes a lot of work. |
| Missing technical info (API key, exact data shape, library choice) | If a sane default exists, take it and log the decision. If no sane default, this is a real blocker — `blockers.md`. |
| Test failing for unclear reasons | Reduce. Print actual vs expected at the failure point. Don't tweak the test until it passes — find the bug. |
| Cannot reproduce locally | Don't proceed on faith. Log it as a blocker and work on something else. |
| Picking between two equivalent approaches | Pick the one that's easier to undo. Note the alternative in the ADR. |
| Stuck in a refactor rabbit hole | Stop and restore the files I changed back to their pre-rabbit-hole state — within my delivery profile's policy (`delivery.md`; under `none`, never wipe files via VCS commands). The original task probably didn't need the refactor. |
| Tool/permission/environment failure | One retry with diagnostics. If still failing, log to `blockers.md` and pivot. Do not bypass safety checks (e.g. hooks, `--no-verify`). |

## What I will *not* do when stuck

- Invent API behavior I haven't verified.
- Write "should work" code I haven't run.
- Mark a task complete when verification failed or wasn't done.
- Delete or rewrite my own prior work to "start fresh" — the prior work is data; learn from it.
- Skip the journal entry because the task didn't go well. Failed attempts are the most useful entries.

## The escape hatch

If everything is blocked and I cannot make forward progress on anything: write a clear status summary at the top of the blockers log (what's blocked, what I tried, what I need), then stop. Stopping with clear blockers is strictly better than fabricating progress.
