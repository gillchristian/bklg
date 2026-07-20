# Verification gates

A task is **not done** until every applicable gate below passes. This is the antidote to the "looks right, ship it" failure mode that wrecks autonomous work.

## Writing acceptance criteria

Written at the loop's Plan step, before implementation. A criterion is valid only if it **names its decider** — the check a fresh evaluator could run without the author's context:

- a **command + its expected exit/output** ("the local-CI gate exits 0", "the smoke prints PASS"), or
- a **countable delta** ("unit tests 68→71", "zero new console warnings"), or
- a **named manual probe + its expected observation** ("open the edited page, toggle the control, the count reads 1").

Litmus: **if a fresh evaluator couldn't check it, rewrite or split it.** "Feels responsive" is not a criterion; "the reference input's measured time stays under the recorded budget" is.

## The gates

For every task, before marking it complete:

1. **Acceptance criteria met.** Re-read the acceptance criteria written in `CURRENT.md`. Tick each one off explicitly. If a criterion is no longer relevant, document *why* in the task entry — don't silently drop it.
2. **It runs.** The code executes without crashing on the happy path. For libraries: an example import + call works. For CLIs: invoke it. For services: start it and hit it.
3. **It does the thing.** Demonstrate the behavior end-to-end with a real input and observe a real output. Type-checking is not behavior verification. A passing unit test on a mocked dependency is not behavior verification. Run the actual thing. If the project records an executable verification procedure for the surface touched (the reference area's local-CI file says where), run it and quote its output; where none exists, drive the surface by hand and say so.
4. **Tests where they earn their keep.** If the logic has branches, edge cases, or is likely to regress, add a test. Don't write tests for trivial getters.
5. **Type/lint/format clean.** Run the project's checkers — the commands recorded in the reference area's `local-ci.md`, whichever of these the project actually has; that file is the authority, not this list. Fix what they flag; don't suppress without a reason logged in the journal entry. Where the reference area records a numeric budget (a perf baseline, a test-count floor), quote the measured number against it.
6. **No new TODOs left behind.** If a TODO is necessary, it goes in `BACKLOG.md` as a real task, not as a comment buried in code.
7. **Independently reviewed (where the delivery profile says so).** If the enabled profile defines a pre-delivery review step (see `delivery.md` — the `pr` profile does), it has run before delivering: a fresh-context reviewer — whatever second-agent facility the environment provides (a built-in code-review skill, a spawned subagent) — is given **only the diff and the task's acceptance criteria**, never the authoring transcript, and grades each criterion pass/fail. Every confirmed correctness finding is fixed, or explicitly rebutted where the enabled profile records the delivery (the `pr` profile: in the PR description). Bookkeeping deliveries that carry no acceptance criteria of their own (e.g. the `pr` profile's close PR) are exempt. If the environment provides no second-agent facility, note that once in the blockers log and proceed under the remaining gates.
8. **Delivered per the project's delivery mode.** Apply the enabled profile's delivery gates (D1, D2, …) in `delivery.md` — the mode is declared in the project manifest. This gate cannot be ticked without opening that profile's checklist.

## How to verify without me lying to myself

The trap: writing code, running it once, seeing no error, declaring victory. Antidotes:

- **Test the failure case too.** If I implemented "X errors when Y is missing," actually pass missing-Y and watch it error.
- **Quote the output.** In the journal entry, quote the actual command output, not a paraphrase. If I can't quote it, I didn't run it.
- **Distinguish "compiles" from "works."** Note explicitly which gates passed: did I just check types, or did I run the program?
- **If I cannot verify, that is a blocker.** Don't fake-verify. Log in `blockers.md` what verification is missing and why, then pick a task whose verification *is* possible.

## End-of-session sweep

Before stopping (or at long natural breaks):

- The enabled delivery profile's end-of-session sweep passes? (See `delivery.md` — e.g. tree clean under `pr`, every modified file accounted for under `none`.)
- The journal has an entry for every task touched this session?
- The planning area's `CURRENT.md` reflects reality?
- Any silent assumptions that should be ADRs in the decisions area?
