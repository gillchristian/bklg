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

### TASK-009 — Parser robustness for real instances
**Source:** BACKLOG (v2 feedback; `../whiteboard/trail-instance-findings.md` #1, #3, #5)
**Acceptance criteria:**
- [ ] AC1 — Heading-style DONE entries parse: a `### <ID> — <title>` entry under `## Completed` (with a `**Completed:** <date> · **PR:** #N · **Journal:** …` line + prose body) yields a Done card with id, title, and a `DoneRecord` (Date, DeliveryRecord, JournalPointer, Summary) extracted from that body. (Decider: unit test on an inline heading-format `## Completed` string asserts each field.)
- [ ] AC2 — Both DONE formats supported (no regression): the skeleton's `- <ID> — … — … — …. See journal …` bullet format still parses. (Decider: existing `TestParseCardTable`/`TestParseFields`/`TestParseDedup` stay green.)
- [ ] AC3 — Emphasized ids recognized: `parseID` strips leading markdown decoration (`*`/`` ` ``/`_`/space), so `**WI-8 — …` → id `WI-8`. (Decider: unit test `parseID("**WI-8 — x")` → `{WI,8,"WI-8"}`; and a `- [ ] **WI-8 — …**` backlog line yields an id-bearing card.)
- [ ] AC4 — Real-repo effect: against `trail --dir systems/track/knowledge`, a `### <ID>` DONE entry that **exists** now parses and matches its `[x]` backlog id — TRACK-000's warning is gone and `/TRACK-000` shows the parsed Date/PR/Journal/Summary from the heading body; its Done card uses the short heading title. **Revised from the original "15 → 0":** trail's DONE.md contains **only** TRACK-000 — TRACK-001…014 are kept solely in their rich `[x]` BACKLOG lines, so the other 14 "shipped item missing from DONE.md" warnings are the check *correctly* flagging that divergence. Reducing / making that noise actionable is **TASK-011's** job, not this parse task. (Decider: `/_diag` no longer mentions TRACK-000 (15 → 14); `/TRACK-000` shows `2026-06-25` + `PR #161`.)
- [ ] AC5 — ADR recorded: an ADR in `../decisions/` documents (a) widening the §2 input contract to the heading DONE format + emphasized ids, and (b) the line-scanner-vs-goldmark call (spec D4). (Decider: `decisions/0001-*.md` exists + indexed.)

**Notes:**
- `parseDone`: a `## Completed` section is heading-style if it contains any `### ` line → parse `### <ID>` blocks (ignore body `- ` bullets); else bullet-style (existing). New `parseDoneHeading(block)` + helpers `fieldAfter(body, marker)` (text after a `**Marker:**` up to the next ` · `/newline) and `doneSummary(body)` (prose minus the metadata line). `prPrefixed` normalizes `#161 …` → `PR #161 …`.
- `parseID`: `strings.TrimLeft(TrimSpace(s), "*`_ ")` before matching `idRe`. Affects all id extraction uniformly; normal ids unaffected. (Trailing `**` on titles is left for TASK-010 markdown to clean.)
- This **widens the input contract** deliberately (v2) — the tool's job is real instances, not just the skeleton. Record ADR-0001. If the heading/prose handling gets gnarly, that's the D4 signal to consider goldmark (note in the ADR; don't switch this task).
- Test the parser with inline strings (no change to `testdata/knowledge` counts); verify the trail effect via binary smoke.
