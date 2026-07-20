# The knowledge framework

**Framework v5 (2026-07-08).** A file-based working system for autonomous
agent sessions: pick one task, verify it honestly, deliver it per the
project's rules, write down what happened — so a future session with zero
memory can pick up exactly where this one stopped.

> **Changelog.** **v5 (2026-07-08):** decidable contracts and countable caps —
> acceptance criteria must name their decider (`verification.md`, new "writing
> acceptance criteria" rule; numeric budgets quoted against measurements; run
> recorded executable verification procedures and quote their output);
> when-stuck rung 7 pivots on **three fact-less attempts**, not felt minutes;
> "went badly" is a checkable floor and task size is "one verifiable slice"
> (`working-style.md`); deterministic rituals done from prose three times earn
> a script (never the merge decision). **v4 (2026-07-08):** the
> unattended-delivery safety net —
> a fresh-context review before task deliveries (verification gate 7 + the
> `pr` profile's pre-merge review step; the old gate 7 is now gate 8), a
> **session envelope** saying when the *session* stops, not just the task
> (`working-style.md`; checked at the loop's Advance step), and the `pr`
> profile's post-merge remote check (gate D3). **v3 (2026-06-24):** path
> indirection — the framework names instance areas by *role* ("the planning
> area", "the journal", "the manifest") and each project's manifest maps roles
> → paths via a **Locations block**, so a single framework copy can serve
> instance areas that live anywhere (e.g. a monorepo's per-system knowledge
> trees). **v2 (2026-06-09):** framework extracted from the project into this
> self-contained directory.

This directory is project-agnostic and travels as a unit. The copy you are
reading may be a downstream copy; the **project manifest** records which
project this is, the enabled delivery mode, the project's own rules, where the
upstream copy lives, and — in its **Locations block** — where each instance
area actually sits. The framework refers to those areas by *role* ("the
planning area", "the journal", "the manifest"); the manifest maps each role to
a path. The **only** hardcoded hop is `CLAUDE.md` → the manifest; everything
else dereferences through Locations — which is what lets one framework copy
serve many instance areas at once. If you landed here first, read the manifest
before acting on anything below.

New project? `SETUP.md` is the adoption guide. Already-installed project?
Never run SETUP.md's steps — and never edit the files in this directory to
say project-specific things; project specifics belong in the manifest.

## The pieces

The framework names six areas by **role** (the same keys the manifest's
Locations block maps to paths): `framework` (this directory) plus the five
instance areas below. Refer to them by role; let the manifest say where they
live.

- **This directory** (role `framework`) — the stable "how I work":
  `principles.md` (values, in priority order), `verification.md` (when a task
  counts as done), `when-stuck.md` (what to do instead of asking the user),
  `working-style.md` (cadence, scope discipline, anti-patterns), `delivery.md`
  (the pluggable part: how verified work leaves the tree — profiles `pr` /
  `commits` / `none`, per-task overrides).
- **The planning area** (role `planning`) — the active plan: `CURRENT.md`
  (exactly one task, with acceptance criteria), `BACKLOG.md` (ordered queue),
  `DONE.md` (archive).
- **The progress area** (role `progress`) — `journal.md` (the journal:
  append-only log, future-me's memory), `blockers.md` (the things only the user
  can resolve).
- **The decisions area** (role `decisions`) — one ADR per non-trivial decision,
  plus `INDEX.md` (which also holds the template and the "what deserves an ADR"
  bar).
- **The reference area** (role `reference`) — project facts: the brief (product
  intent — it wins conflicts with planning), glossary, local-CI commands,
  specs.
- **The whiteboard** (role `whiteboard`) — discussions in flight that haven't
  earned an ADR or a backlog task; conventions and index in its README.

## The loop

Every working session follows the same shape. Steps marked § are defined by
the enabled delivery profile (`delivery.md`) — the manifest usually shows
them instantiated for the project.

1. **Orient** — Read the manifest, then the planning area's `CURRENT.md` for
   the active task. If empty, refill from `BACKLOG.md`.
2. **Plan** — Write the acceptance criteria into `CURRENT.md` before touching
   code — each naming its decider (the "writing acceptance criteria" rule in
   `verification.md`).
3. **Stage §** — Prepare per the delivery profile (`pr`: branch off the
   default branch; `commits` / `none`: nothing to do).
4. **Execute** — Implement, checkpointing per the profile as I go.
5. **Verify** — Run the gates in `verification.md`, including gates 7–8 (the
   profile's review step and delivery gates) and the project's local CI
   commands. If any fail, fix before moving on.
6. **Deliver §** — The profile's delivery steps (`pr`: open the PR, pass the
   fresh-context review, merge; `commits`: final commit; `none`: verified tree
   + hand-off record).
7. **Log** — Append a journal entry: timestamp, what was done, what was
   verified (with quoted output), the profile's delivery record, what's next.
8. **Advance** — Move the task to `DONE.md`, do the profile's post-delivery
   sync if it defines one (`pr`: sync the default branch, resolve the remote
   check where the project records one; `commits` / `none`: nothing), then
   check the **session envelope** (`working-style.md`) before pulling the next
   task into `CURRENT.md`. An empty backlog or a spent envelope is a terminal
   state, not an error: run the end-of-session sweep and stop.

If I ever feel stuck or unsure, the answer is in `when-stuck.md` — not in
asking the user.

## Changing the framework

Framework files stay instance-free: no project names, no stack names, no
user identities, no repo-specific rituals, and **no hardcoded instance-area
paths** — refer to areas by role and let the manifest's Locations block resolve
them; all of that lives in the manifest. Downstream copies don't edit these
files in place (a local edit forks the framework silently); improvements go to
the upstream copy recorded in the manifest, then re-copy the directory
wholesale and bump the version line above.
