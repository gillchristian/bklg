# 0004 — Dashboard adapter: a second input convention

**Date:** 2026-07-21 · **Status:** proposed

## Context

bklg's parser is built around the framework skeleton: a `planning/` area of
`CURRENT.md`/`BACKLOG.md`/`DONE.md`, a `progress/blockers.md`, and one
structured id (`[A-Z]+-\d+`) per card that joins across those files. Pointed at
the real Pinata KB (`~/dev/Pinata-dev/Pinata/knowledge`, investigated
2026-07-21) it exits `no planning area at knowledge/planning`. That KB is not
malformed and does not lack a work lifecycle — it runs a full Active/Backlog/Done
board, but in a different shape: one file (`work/index.md`), Active/Done as pipe
**tables**, Backlog as bullet groups under bold subheads, identity via inline
Linear ids (`PINATA-\d+`, 0..N per row) rather than one id per card, and
"blocked" expressed in prose (`⛔`), with Linear as the source of truth for
status. It is a second, legitimate way a knowledge base tracks work — "Path B"
in the investigation (Path A, conforming the KB to the skeleton, is the rejected
alternative below).

## Decision

Teach bklg a **second input convention** — a "dashboard adapter" — that reads a
single-file Active/Backlog/Done dashboard, alongside (not replacing) the
framework parser:

- **Opt-in, config-driven.** A `dashboard:` key in the manifest's `## Locations`
  block (or a `--dashboard <file>` flag) points at the one dashboard file and
  selects the adapter; without it, nothing changes. No auto-detection in v1
  (parked — explicit config is safer than guessing).
- **Manifest lookup widens** to try `README.md` then `index.md`, because
  dashboard-shaped KBs (Pinata) name their manifest `index.md`.
- **The adapter sits behind the same parser seam** as the line-scanner (spec
  §10 / D4): resolution selects it, and the rest of the pipeline (model, badges,
  board, detail, `/_v`) is shared.
- **Identity is relaxed:** a card carries 0..N Linear tickets (chips/links), not
  one join key; the route key is a title slug; there is no cross-file dedup (one
  file) and no acceptance-criteria or blocker-file join.
- **The convention is a two-sided contract.** The adapter is defensive (degrade
  to `/_diag`, never crash), but a useful board needs the dashboard to stay
  regular; the shape the target KB follows is
  [`../reference/specs/dashboard-format.md`](../reference/specs/dashboard-format.md).

Implemented as backlog batch **v3 (TASK-013…016)**.

## Alternatives considered

- **Conform the KB to the framework skeleton (Path A).** Add `planning/` +
  `progress/` to Pinata and mirror `work/index.md` into the skeleton. Rejected:
  it means hand-maintaining a rigid second work tracker in parallel with
  `work/index.md` *and* Linear — exactly the duplication that KB avoids by
  delegating status to Linear. The cost lands on the KB (ongoing) instead of on
  the tool (once).
- **Auto-detect dashboard mode** (no `planning/` but a `work/index.md` present).
  Deferred to a follow-up: explicit config is less surprising and won't misfire
  on a half-migrated repo.
- **`contains(⛔)` ⇒ blocked.** Rejected: the Pinata KB already uses `⛔`
  mid-paragraph as a decision/attention marker, so "contains" false-positives
  (the C4 row has several and is *not* blocked). Only a **leading** `⛔` /
  `**Blocked**` marks a card blocked.
- **goldmark for the tables/markdown.** Still rejected (as in D4 / ADR-0001):
  the shapes are regular (pipe tables, bullet groups), the domain interpretation
  is needed either way, and the zero-dep story is worth keeping.

## Consequences

- A second convention to keep working — but additive and behind the existing
  seam, so framework-mode behaviour is untouched.
- The `Card` model gains an optional multi-ticket slice + a dashboard flag;
  dashboard-mode detail pages show the raw row + links out rather than id-joined
  fields.
- bklg now has a stake in an external KB's format. The contract doc is that
  stake made explicit and is what the target KB agrees to; divergence surfaces
  in `/_diag`, not as silent board corruption.
- Linear stays authoritative for status; the board reflects the file, which may
  lag Linear. A later "Linear status sync on read" is parked (needs network +
  auth; likely stays a non-goal for a zero-dep localhost tool).
- **Proposed, not accepted:** this records the plan; building it is the v3
  backlog batch, greenlit and pulled through `CURRENT.md` per the normal flow.
  Flip to **accepted** when TASK-013 is promoted.
