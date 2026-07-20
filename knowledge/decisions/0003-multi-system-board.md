# 0003 — Multi-system board: aggregate + server-side filter

**Date:** 2026-07-20 · **Status:** accepted

## Context

v1 treated a monorepo **root manifest** (a `systems/<name>` index, no planning
area of its own) as a helpful error: it printed the per-system `bklg …` commands
and exited. The user wants instead to *see tickets across all projects and
filter to any one* (feedback #1; spec §13.1). This is design-heavy — how to
combine systems, how to filter, whether the CLI/detail/`/_v` change — so the
choices are recorded here (the "sane defaults + ADR" the user opted into for the
autonomous v2 batch).

## Decision

At a root manifest, **aggregate**: resolve + parse every `systems/<name>`
instance and serve **one** board combining all their cards. Specifics:

- Each `Card` is tagged with its `System`; the board and detail page show a
  **system chip**. Single-system mode is unchanged (no `System`, no chip/bar).
- **Filtering is server-side** via `?system=<name>` — a plain-link filter bar,
  no client JS, shareable URLs. The bar lists **every discovered system**
  (including zero-card ones) so any project is filterable and the bar matches the
  startup's system count.
- **Detail** (`/<id>`) does a **global lookup** across the aggregated cards.
- **`/_v`** is the max mtime across all systems' files.
- A system in the index that **fails to resolve** is skipped with a `/_diag`
  `read-error` warning — never a crash.
- `main` builds a multi-server on `RootManifestError` instead of exiting; the
  startup line reads `aggregate: N systems — …`.

## Alternatives considered

- **System switcher (one at a time).** Rejected: the point is to see *across*
  projects; aggregate-with-filter is a superset (you can still focus one).
- **Client-side JS filter.** Rejected: server-side `?system=` is simpler, has no
  client state, and gives shareable/bookmarkable URLs — consistent with the
  tool's "no frontend framework" stance (only the live-reload poll is JS).
- **Per-system detail routes (`/<system>/<id>`) to avoid id collisions.**
  Deferred: namespaces differ per system, so a global `/<id>` lookup is
  unambiguous in practice; on the rare cross-system collision (same namespace +
  number in two systems) the first match wins. Revisit if it ever bites.

## Consequences

- `board()` re-resolves + parses N systems per request. Fine for a localhost
  tool with a handful of systems; if it ever feels slow, an mtime-keyed cache is
  the lever (spec §7 already anticipated caching as a later optimization).
- The root-manifest behaviour **changed** from error→aggregate. The old
  per-system invocation hint is gone; the aggregate board + filter bar replace
  it (a single system is still directly viewable via `--dir systems/<name>/knowledge`).
- Cross-system id collisions on the detail route are first-match — documented,
  low-risk.
