# 0001 — Widen the input contract beyond the framework skeleton

**Date:** 2026-07-20 · **Status:** accepted

## Context

The v1 parser was keyed to the exact conventions in `bklg-spec.md` §2, which
mirror the framework `SETUP.md` **skeleton**: DONE entries as
`- <ID> — <title> — <date> — <summary> — <delivery>. See journal …` bullets, and
bare ids (`TRACK-000`). Running bklg against a real instance
(`/Users/bb8/dev/trail`, `systems/track`) showed the skeleton and reality
diverge: DONE.md uses `### <ID> — <title>` headings with a `**Completed:** … ·
**PR:** … · **Journal:** …` line + prose body, and backlog/parking ids are often
wrapped in emphasis (`**WI-8 — …**`). Consequences on trail: all 15 shipped
tasks fired the false "shipped item missing from DONE.md" warning (0 DONE entries
parsed), Done cards fell back to the huge `[x]` backlog lines, and decorated ids
produced id-less (unlinkable) cards. (Full findings:
`../whiteboard/trail-instance-findings.md`.)

## Decision

Extend the line-scanner to accept **both** DONE shapes (a `## Completed` section
with any `### ` line is parsed as heading entries; otherwise as bullets) and to
**strip leading markdown decoration** (`*` `` ` `` `_`) before matching an id.
Keep the zero-dependency line-scanner; do **not** adopt goldmark yet.

## Alternatives considered

- **Stay strict to the skeleton** (spec §2 as-written; divergence → best-effort +
  `/_diag`). Rejected: bklg's whole value is viewing *real* instances, and the
  skeleton is only one dialect. Faithfulness to real files beats faithfulness to
  the example.
- **Switch to a goldmark-backed parser now** (the D4 seam). Deferred: the
  line-scanner still handles the observed dialects cleanly, and "zero Go module
  deps" remains valuable (`go install`, no supply chain). D4 stays the escape
  hatch — revisit if the dialects keep multiplying or the metadata parsing gets
  gnarly. (Markdown *rendering* of field text is a separate concern — TASK-010.)

## Consequences

- The accepted input is now broader than spec §2 documents — a divergence from
  the written spec that must be kept in mind when reading it. The spec should get
  a v2 note; the parser stays behind the `Parser` interface (D4) so a swap is
  still cheap.
- More formats to keep in sync inside the line-scanner. If a third DONE dialect
  appears, that is the signal to reconsider goldmark rather than add a third
  branch.
