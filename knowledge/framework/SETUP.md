# Adopting the framework in a new project

> **Guard: is it already installed?** If a `CURRENT.md` already exists for this
> instance's planning area — for a standalone project that's
> `knowledge/planning/CURRENT.md`; in general, wherever the manifest's Locations
> block points — the framework is installed: **stop**. This guide is for
> first-time adoption only; running it against a live project clobbers real
> state. Every copy step below is copy-if-absent: never overwrite an existing
> file.

## Steps

1. **Copy `framework/` wholesale** into `<project>/knowledge/framework/` —
   all files, unmodified. Don't edit them afterward (see "Changing the
   framework" in `README.md`).
2. **Create the instance areas** from the skeletons below:
   `knowledge/{planning,progress,decisions,reference,whiteboard}/`.
3. **Write the manifest** (`knowledge/README.md`) from its skeleton. The
   decisions that matter:
   - **Delivery mode** (`pr` / `commits` / `none` — see `delivery.md`).
     If unsure, or if the user hasn't said the agent may touch version
     control, choose `none`. Record the mode *and* its operative one-liner.
   - **Identity/attribution rules.** Whose name goes on commits/PRs, and
     whether agent attribution is allowed. If unstated: no agent attribution
     (the framework default).
   - **Locations.** The role → path map (`framework`, `planning`, `progress`,
     `decisions`, `reference`, `whiteboard`) the framework dereferences. For a
     standalone project these are the default `knowledge/…` paths created in
     step 2; declare them explicitly so the framework resolves areas by role
     rather than assuming a layout.
4. **Write the project's `CLAUDE.md`** (or extend an existing one) from its
   skeleton: the explicit reading chain plus inlined non-negotiables,
   including the delivery-mode-specific line.
5. **Fill the brief** (the reference area's `project-brief.md`) from the user's
   description, derive an initial `BACKLOG.md`, pull the first task into
   `CURRENT.md`. If the
   user hasn't described the product yet, record what *is* known plus an
   explicit Unknowns list in the brief, log a blocker, and stop after setup —
   don't invent a backlog.
6. **Record the local CI commands** in the reference area's `local-ci.md` as
   soon as they exist; until then leave the stub saying they're undefined.
7. In the manifest, **record the framework version and upstream**: the
   version is the `Framework vN` line atop `framework/README.md`; the
   upstream is the repo + path this directory was copied from, named so a
   teammate on another machine can find it.

One more placement decision for `none`-mode projects: if the agent may not
commit, the `knowledge/` files themselves can't be committed by the agent
either. Either the user commits them as part of their own flow, or
`knowledge/` lives outside the repo (a sibling directory) — pick one in the
manifest so future sessions don't wonder.

## Skeletons

### `knowledge/README.md` (the manifest)

```markdown
# knowledge/ — project manifest (<project>)

The working system lives in [`framework/`](framework/README.md); this file is
what makes it THIS project. Reading order: this file → `framework/README.md`
→ the enabled profile in `framework/delivery.md`.

## Delivery mode

delivery: <pr | commits | none>

**Operative meaning:** <one sentence: what the agent may and may not do with
version control — e.g. "never run VCS-mutating commands; the user handles
git" for `none`>.

## Locations

The role → path map the framework dereferences. It refers to areas by role
("the planning area", "the journal"); this block says where they live. Paths
are repo-root-relative. For a standalone project these are the defaults below;
when one framework copy serves several instances (e.g. a monorepo), each
instance's manifest points its roles at that instance's own areas while
`framework` points at the shared copy.

framework:  knowledge/framework
planning:   knowledge/planning
progress:   knowledge/progress
decisions:  knowledge/decisions
reference:  knowledge/reference
whiteboard: knowledge/whiteboard

## Project rules

- **Identity/attribution:** <whose name on commits/PRs; agent attribution
  allowed or not>.
- **Session envelope (optional):** <when an unattended session stops — e.g.
  "run until BACKLOG's active list is empty or the first hard blocker". Omit
  the line to use the framework default, which is *not* that example: backlog's
  active list empty, or the escape hatch fires; see `working-style.md`.>
- <other standing rules: merge strategy, protected branches, bootstrap
  exceptions, where knowledge/ lives if not in-repo, …>
- **Framework copy:** v<N>, copied from <upstream repo + path>. Don't edit
  `framework/` here; improvements go upstream.

## The loop, instantiated

<the 8 steps from framework/README.md with the §-steps made concrete for
this project's delivery mode — keep it crisp enough to follow verbatim>

## Layout

- framework/ — the reusable system. planning/ — CURRENT (one task), BACKLOG,
  DONE. progress/ — journal (append-only), blockers. decisions/ — ADRs +
  INDEX. reference/ — brief (wins conflicts), glossary, local-ci, specs.
  whiteboard/ — discussions in flight.
```

### `CLAUDE.md` (project root)

```markdown
# CLAUDE.md — <project>

**Read these, in order, before doing anything else:** `knowledge/README.md`
(the manifest: delivery mode + Locations + project rules),
`knowledge/framework/README.md` (the working system), then the enabled profile
in `knowledge/framework/delivery.md`. The summary below is just the headline
rules so you can't accidentally violate them while still loading the rest.

## Non-negotiables

1. **One task at a time.** Pull from the planning area's `CURRENT.md` (the
   manifest's Locations block says where that is); if empty, promote the top
   unchecked item of `BACKLOG.md`. Acceptance criteria before code.
2. **Delivery: <mode>.** <the mode's operative one-liner, verbatim from the
   manifest — e.g. "NEVER run VCS-mutating commands; the user handles git.">
3. **Attribution:** <the identity rule — e.g. "commits/PRs carry the user's
   name only; no Co-Authored-By, no generated-with footers.">
4. **Verify before declaring done.** Gates in
   `knowledge/framework/verification.md`. Run the program, quote actual
   output, don't confuse "compiles" with "works."
5. **Journal everything.** Append to the journal after every task.
6. **When stuck, follow `framework/when-stuck.md`.** Don't ask the user; log
   real blockers to the blockers log, then pivot.
```

### `planning/CURRENT.md`

```markdown
# Current task

> One task at a time. When this file is empty, pull the next item from `BACKLOG.md`.

## Entry template

### TASK-NNN — <title>
**Source:** BACKLOG / parking lot / user request
**Acceptance criteria:**
- [ ] criterion (how it will be verified)
**Notes:** scope cuts, links, anything decided while planning.
(Add `**Delivery override:** … — user, YYYY-MM-DD` only when the user grants
one; see framework/delivery.md.)

## Active

_(none yet)_
```

### `planning/BACKLOG.md`

```markdown
# Backlog

Ordered. The top *unchecked* item is next. Promote into `CURRENT.md` when started.

Conventions:
- Completed tasks stay here, checked `[x]` with their delivery record — the
  list doubles as a one-line shipping index. The full record lives in `DONE.md`.
- TASK ids are global and monotonically increasing across all of `planning/`:
  next id = highest id appearing anywhere + 1.

## Active

_(derive from the project brief)_

## Parking lot

_(deferred items, known bugs not yet scheduled)_
```

### `planning/DONE.md`

```markdown
# Done

Tasks that passed all verification gates. Newest at top.

Each entry: id, title, completion date, a summary (a sentence for trivial
tasks, a paragraph when the detail earns its keep — the journal holds the
full story), the delivery record (see your profile in framework/delivery.md),
and a journal pointer.

- TASK-NNN — title — YYYY-MM-DD — summary — <delivery record>. See journal YYYY-MM-DD [HH:MM]

(Under `none`, the delivery record is itself a journal pointer — write it once.)

## Completed

_(none yet)_
```

### `progress/journal.md`

```markdown
# Journal

Append-only. Newest at the bottom. Each entry is a snapshot for future-me
with no memory of this session.

## Entry format

---
## YYYY-MM-DD HH:MM — <short heading>
   (HH:MM may be dropped only if the heading alone makes the entry findable
   from a DONE.md pointer)

**Task:** TASK-NNN (or "scaffolding" / "exploration" / "blocker triage")
**What I did:** 1–3 sentences.
**What I verified:** which gates ran, with literal command output worth keeping.
**What changed:** files touched + the delivery record per my profile.
**What I learned:** anything that would surprise future-me. Non-obvious only.
**Next:** the very next thing to do on resume.

## Entries
```

### `progress/blockers.md`

```markdown
# Blockers

Things that genuinely require user input. If this file has entries, **the
user needs to see them** at the top of any hand-off.

## Format

## BLOCKER-NNN — short title — opened YYYY-MM-DD HH:MM
**Task affected:** TASK-NNN
**What I tried / observed / currently believe / need:** facts first,
speculation flagged as such, then the specific input or decision needed.
**Workaround in place:** what happened instead, if anything.

## Open

_(none)_

## Resolved

_(none)_
```

### `decisions/INDEX.md`

```markdown
# Decisions (ADRs)

One file per non-trivial decision: `NNNN-short-title.md`.

## Template

# NNNN — <decision title>
**Date:** YYYY-MM-DD · **Status:** proposed | accepted | superseded by NNNN
## Context — the situation and forces, 2–4 sentences.
## Decision — what was decided, specifically.
## Alternatives considered — bullets, each with why not.
## Consequences — what this makes easy/hard; what to revisit.

## What deserves an ADR

Language/framework/major-library picks; choosing between viable
architectures; interpreting an ambiguous requirement; anything a teammate
should be able to understand later without asking.

## Index

_(none yet)_
```

### `reference/` stubs

```markdown
# project-brief.md — What it is / Why it exists / Hard constraints /
# Out of scope / Stack / Success criteria / Open questions / Raw notes (user words)

# glossary.md — "**Term** — definition." Add entries the moment a term first
# appears with a specific meaning. Surprising semantics first.

# local-ci.md — the commands behind "run local CI" (verification gate 5 and
# the delivery gates). Until defined: "_(no CI commands recorded yet — record
# them the first time they exist; do not claim a green CI before then.)_"
```

### `whiteboard/README.md`

```markdown
# whiteboard/

Discussions in flight — explored but not committed to action. Not a backlog:
entries record *what we considered* so future sessions don't reopen settled
ground. One file per topic: `<topic-slug>.md` with Status line (open /
resolved / tabled), Context, Options considered, Where we landed, Follow-ups.
When a topic resolves: promote to an ADR, a backlog task, or a brief update —
and update the Status line; keep the entry (the recorded thinking is the value).

## Index of current entries

_(none yet)_
```
