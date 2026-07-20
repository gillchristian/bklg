# Delivery

How verified work leaves my working tree. This is the framework's pluggable
piece: every project declares exactly **one delivery mode** in its manifest,
and only that profile's rules below apply to me.

> **If the manifest declares no delivery mode, the mode is `none`.** I never
> assume version-control liberties that were not declared. If the declaration
> is missing, I flag it in the blockers log and work under `none`.

## Reading rules

- I follow **only** the profile the manifest enables. The other profiles are
  not background reading — a `none`-mode agent that skims the `pr` profile
  "for context" is how a forbidden PR gets opened.
- Wherever a profile says *commit* or *PR*, I first apply the project's
  identity/attribution rules **from the manifest — re-reading that section at
  the moment of delivery**, not from memory of having read it hours ago.
- Framework default when a manifest says nothing about attribution: I add
  **no attribution of my own** — no `Co-Authored-By` trailers, no
  "generated with" footers. Commits and PRs carry the user's identity only.

## Commit conventions (any profile that commits)

- **Subject line** ≤ 72 chars, imperative mood ("add", "fix", "rename"), no trailing period.
- **Body** explains *why* if it's not obvious. Skip the body for trivial commits.
- One logical change per commit. Fixing a bug *and* renaming a function is two commits.
- No emojis unless the existing history uses them.
- Attribution per the manifest (see reading rules above).

## Per-task overrides

The user can grant delivery liberties beyond the enabled mode for a single
task ("you may make several commits for this one"). The contract:

1. **Written down before acted on.** An override is in effect only once it is
   recorded in the task's entry in the planning area's `CURRENT.md`:
   `**Delivery override:** <granted behavior> — user, YYYY-MM-DD`.
   A grant that lives only in conversation does not survive a session restart
   or context compaction. If it isn't written down, it isn't in effect.
2. **Scoped to that task.** The override expires the moment the task leaves
   `CURRENT.md`. It sets no precedent — my own recent commits are *not*
   evidence of standing permission, in either direction: I neither extend the
   grant to the next task nor "correct" recorded-override commits I find when
   resuming.
3. **Named, not extrapolated.** A grant covers exactly what the user named —
   typically the `commits` profile's semantics. **Pushing, force operations,
   committing to protected branches, and opening or merging PRs are never
   implied**; each requires being individually named in the grant.
4. **Bounded by the manifest.** An override relaxes *this profile's*
   defaults; it does not defeat a rule the manifest marks as a hard
   constraint. If the user's grant collides with one (say, a merge strategy
   the manifest forbids), the manifest is what changes: update its rule as
   part of recording the grant, quoting the user's words — so the file future
   sessions actually read stays true.

---

## Profile: `pr` — branch → PR → self-merge

I own the full cycle: branch, commit, open the PR, merge it.

### Rules

1. **The default branch is sacrosanct.** No direct pushes — every change
   lands via a PR I open and merge myself. (Whatever bootstrap commits predate
   the framework are covered by the project manifest, not by this rule.)
2. **One PR = one logical unit of work.** Something I'd call out in a
   changelog — usually one `CURRENT.md` task.
3. **Branch off the latest default branch.** Fetch/pull before branching.
4. **Iterate freely on the branch.** Commit early and often. Local CI (the
   project's recorded commands — see the reference area's `local-ci.md`) must
   pass before the PR opens.
5. **Merge my own PRs** once CI is green and the description is complete —
   this profile assumes a repo where I'm authorized to (solo repos, usually).
   Merge strategy: whatever the manifest declares; default `--squash`.

### Branch naming

`<kind>/<task-id>-<slug>`, kind one of `feat`, `fix`, `chore`, `refactor`, `docs`, `test`.

Examples: `feat/task-007-user-auth`, `fix/task-012-trailing-slash-redirect`.

### The full cycle for one PR

1. **Pull the task** into the planning area's `CURRENT.md`. Write acceptance criteria.
2. **Branch** from the up-to-date default branch.
3. **Implement, committing as I go.** Each commit leaves the branch in a sane state.
4. **Run local CI** (the commands in the reference area's `local-ci.md`), plus a
   manual smoke test where applicable (the `verification.md` gates still apply).
5. **If CI fails, fix on the branch.** Don't open a red PR.
6. **Open the PR** (e.g. `gh pr create`): title imperative ≤ 72 chars
   mirroring the task; body from the template below; identity/attribution per
   the manifest.
7. **Obtain the fresh-context review** (verification gate 7): hand the
   reviewer only the diff and the task's acceptance criteria — never the
   authoring transcript — and have it grade each criterion pass/fail. Fix
   each confirmed correctness finding, or rebut it explicitly in the PR
   description, before merging.
8. **Merge** per the manifest's declared strategy (default
   `gh pr merge --squash --delete-branch`).
9. **Sync the local default branch** (`pull --ff-only`).
10. **Run the remote check, if the project records one** (gate D3 below) —
    resolve it before the close PR merges.
11. **Close the task** with a small follow-up PR (see below).
12. **Pick the next task.**

### The close PR

Post-merge bookkeeping can't land on the default branch directly (Rule 1 has
no exceptions), so it ships as its own tiny PR immediately after the task PR:

- Branch: `docs/task-NNN-close`. Title: `docs: close TASK-NNN (<short>, merged <sha>)`.
- Contents: move the task from `CURRENT.md` to `DONE.md` (with PR number +
  merge sha), tick it off in `BACKLOG.md` if it was queued there, append the
  journal entry quoting the PR URL and merge commit.
- This is the one PR class that carries no acceptance criteria of its own —
  it documents the task that just merged. Merge it the same way (it is exempt
  from the fresh-context review, per verification gate 7).
- Where the project records a remote check (gate D3), the close PR also notes
  the task PR's remote-check outcome.
- The close PR's **mechanical shell** — branch, title, file moves, open, merge,
  sync — may be scripted if the project records one in the reference area
  (`working-style.md`, scripts over re-derivation). The journal entry and
  `DONE.md` prose remain **authored input** to the script; no script ever
  performs a task PR's merge decision.
- When the next task is already known, the close PR may also pull it into
  `CURRENT.md` — closing and orienting in one step.

### Delivery record

How "delivered" is written in `DONE.md` and the journal: ``PR #N, merged `sha` ``.

### Delivery gates (the gate-8 checklist for this profile)

- **D1 — Committed.** One commit per logical change, message explains *why*,
  attribution per the manifest, working tree clean before the next task.
- **D2 — PR merged.** Local CI green before opening; the fresh-context review
  passed (verification gate 7 — findings fixed or rebutted in the PR
  description); template body; merged per the manifest's strategy; branch
  deleted; default branch synced.
- **D3 — Remote check resolved.** If the reference area's local-CI file
  records a remote-check command (the CI/deploy status a default-branch merge
  triggers), run it after the merge and before merging the close PR — close-PR
  prep absorbs the remote latency. The command must resolve to **green or
  red**; "in progress" is neither — poll, or use the command's wait mode,
  until it resolves or times out. **Red:** the fix (per the Hotfix exception)
  becomes the next task, jumping the backlog, and the close PR notes it.
  **Unresolvable** (outage, auth): log it to the blockers log — this gate is
  never ticked on an unresolved state. If the local-CI file records no remote
  check, this gate is vacuous — and recording one the first time a remote
  check exists is part of keeping that file true.

### End-of-session sweep (this profile)

- `git status` clean? (Or every dirty file accounted for in `CURRENT.md`?)
- The last merge's remote check resolved green (or its unresolved state is
  logged)?

### PR description template

```markdown
## What
One-paragraph summary of the change.

## Why
The motivation — typically pulled from the task in `CURRENT.md`/`DONE.md`.

## How
Bullets describing the approach, anything unusual, anything explicitly out of scope.

## Verification
- [ ] Acceptance criterion 1 (with how it was verified)
- [ ] Acceptance criterion 2
- ... command outputs or specific manual checks as evidence.

## Notes
Anything reviewers (future-me, the user) should know: follow-ups, known limitations, ADRs added.
```

### When *not* to merge

- Local CI is failing.
- A `verification.md` gate failed and was not addressed.
- The fresh-context review confirmed a correctness finding that is neither
  fixed nor rebutted in the PR description.
- I'm uncertain about behavior and would normally want a second pair of eyes —
  open the PR as draft, log the uncertainty in `blockers.md`, pivot to another
  task. Better one held PR than a bad merge.

### Hotfix exception

If a merged PR broke the default branch, the right move is a small `fix/...`
PR that reverts or repairs — never a force-push or direct commit. Same
workflow, just fast.

---

## Profile: `commits` — commit, but don't manage branches or PRs

I commit my work on the branch I was given; branching, pushing, PRs, and
merging belong to the user.

- Commit conventions above apply; attribution per the manifest.
- I never push, never create/switch/delete branches, never open or merge PRs —
  unless individually named in a per-task override.
- **Delivery record:** `commits <shaA>..<shaB> on <branch>`.
- **Delivery gates:** D1 — committed, one logical change per commit, *why* in
  the message. D2 — working tree clean, or every dirty file accounted for in
  `CURRENT.md`.
- **End-of-session sweep:** `git status` shows only what the user expects to
  find: my commits on their branch, nothing pushed, no surprise dirty files.

This profile is also the usual target of a per-task override granted inside a
`none`-mode project.

---

## Profile: `none` — I never mutate version control

The user owns version control entirely. "Deliver" means: the verified change
sits in the working tree, and the planning and journal files say exactly what
changed and how it was verified — the user can commit without asking me what
happened.

### Version-control command policy (exact)

- **Allowed (read-only):** `status`, `diff`, `log`, `show`, `blame`.
- **Never:** `commit`, `push`, `pull`, `checkout`/`switch`/`restore`, `reset`,
  `stash`, `clean`, `rebase`, `merge`, `add`, branch/tag creation or deletion,
  anything `--force`. Not even when it looks like cleanup.
- **Undoing my own work:** per-file only, only for files *this session*
  changed, and by re-writing their content from what I know — never by
  `checkout`/`restore`/`clean`. The tree may hold the user's own uncommitted
  work; wiping it is the worst failure available to me. If I can't cleanly
  reconstruct a file's prior state, I say so in the journal and leave the file
  as it stands.

### Delivery record

In the journal entry: `handed off YYYY-MM-DD — verified <how>; files touched: <list>`.
`DONE.md` points at that entry instead of a PR number.

### Delivery gates (the gate-8 checklist for this profile)

- **D1 — Accounted for.** Every file I modified is listed in the journal entry
  (and in `CURRENT.md` while the task stays open).
- **D2 — Hand-off ready.** Verification output quoted; the user can act on the
  tree without asking me what changed or whether it works.

### End-of-session sweep (this profile)

- Every modified file accounted for in `CURRENT.md` or the journal?
  (Read-only `status`/`diff` is the checklist source.)
