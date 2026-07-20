# Running bklg against the real `trail` monorepo — findings

**Status:** open — informs the v2 feedback tasks TASK-008…012 in `BACKLOG.md`.
**Date:** 2026-07-20.

## Context

After the v1 MVP shipped, the user ran `bklg` against their real repo
`/Users/bb8/dev/trail` (a monorepo of five systems: trail, cadence, gateway,
track, reflect) and against `systems/track` specifically. The clean in-repo
fixture parses perfectly; the real instance surfaced several gaps. This entry
records what's actually in the trail files so the tasks are grounded in fact,
not guesses.

## What the real instance looks like (systems/track/knowledge/planning)

1. **DONE.md uses `### <ID> — <title>` headings, not `- <ID> — … — …` bullets.**
   Each completed entry is:
   ```
   ### TRACK-000 — Swift/iOS toolchain bootstrap + orientation
   **Completed:** 2026-06-25 · **PR:** #161 (squash-merged) · **Journal:** … .
   <multi-line prose summary>
   ```
   Our parser only reads `- ` bullets under `## Completed`, so it parses **zero**
   DONE entries here. Consequence: `doneIDs` is empty → all 15 `[x]` backlog
   items fire "shipped item missing from DONE.md" (the 15 `/_diag` warnings the
   user saw), and each Done-column card is built from the bare `[x]` BACKLOG line
   instead of the (short-titled) DONE entry. → **TASK-009.**

2. **BACKLOG `[x]` / parking lines carry the full record inline (huge).** e.g.
   `- [x] TRACK-000 — **Swift/iOS toolchain bootstrap …** ✓ PR #161 <several
   sentences>`. Since the card title is "everything after the first ` — `", the
   card renders a wall of text. The card should show a trimmed title; the full
   text belongs on the detail page. → **TASK-008.**

3. **Ids are sometimes wrapped in `**…**`.** Parking items read
   `- [ ] **WI-8 — `.trace` export …**`. `parseID` anchors at string start
   (`^[A-Z]+-\d+`), so `**WI-8` yields no id → the card is id-less (no detail
   link). Recognizing emphasized/decorated ids is part of parser robustness. →
   **TASK-009** (same "handle real instances" slice).

4. **Markdown is shown literally.** `**bold**`, `` `code` ``, `*italic*`,
   `→`-style arrows appear raw on cards and detail. A safe markdown subset would
   make these readable. → **TASK-010.**

5. **`/_diag` is noisy and not actionable.** ~~almost all the DONE-format false
   positive~~ **Corrected during TASK-009:** trail's DONE.md contains **only**
   TRACK-000; TRACK-001…014 live solely in their rich `[x]` BACKLOG lines, never
   duplicated into DONE.md. So only *one* of the 15 warnings was a format
   parse-failure (fixed by TASK-009); the other 14 are the "shipped item missing
   from DONE.md" check **correctly** firing against a repo that diverges from the
   skeleton convention ("the full record lives in DONE.md"). Reducing that noise
   is a **policy** decision — suppress/downgrade when the `[x]` line already
   carries a full record, and/or make each warning link to the file + say what to
   do. → **TASK-011.** (TASK-009 only removes the genuine parse-failure + gives
   real warnings; it is *not* the bulk fix I first assumed.)

6. **Root manifest errors instead of aggregating.** `bklg` at the trail root
   prints the per-system invocation list (v1 behaviour, spec §13.1). The user
   wants an aggregate board across systems with a per-project filter. →
   **TASK-012.**

## Where we landed

Turn the six observations into the five tasks TASK-008…012 (2+3 merge into the
one "parser robustness" slice). Ordering favours the quick fixes that repair the
visibly-broken trail experience (trim, DONE-headings, markdown) before the big
multi-system feature. The parser changes **widen the input contract** beyond the
spec §2 skeleton — that is deliberate v2 scope (bklg's job is real instances, not
only the skeleton), and should be recorded in an ADR when TASK-009 is built.

## Follow-ups

- When building TASK-009, decide (ADR) whether to keep the line-scanner or move
  to the goldmark-backed parser the spec's D4 seam anticipated — real instances
  are messier than the skeleton, which is exactly the trigger D4 named.
- TASK-010 (markdown) collides with the stdlib-only constraint: a safe renderer
  usually wants a library + sanitizer. Decide (ADR) whether to relax "zero Go
  module deps" or hand-roll a strict inline subset.
