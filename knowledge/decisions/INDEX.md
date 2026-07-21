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

- [0001](0001-widen-input-contract.md) — Widen the input contract beyond the
  framework skeleton (heading-style DONE entries + emphasized ids); keep the
  zero-dep line-scanner, defer goldmark (D4). **accepted**, 2026-07-20.
- [0002](0002-safe-inline-markdown-stdlib.md) — Hand-rolled safe inline markdown
  (escape-first + tag whitelist + href scheme-check), stdlib-only; keep the
  zero-dep constraint rather than pull in goldmark + a sanitizer. **accepted**,
  2026-07-20.
- [0003](0003-multi-system-board.md) — Multi-system board: a root manifest
  aggregates all `systems/<name>` into one board with per-card system chips and
  a server-side `?system=` filter (lists every system); global detail lookup;
  unresolvable systems skipped with a warning. **accepted**, 2026-07-20.
- [0004](0004-dashboard-adapter.md) — Dashboard adapter: a second, opt-in input
  convention (a `dashboard:` Locations key / `--dashboard` flag selects it) for
  single-file Active/Backlog/Done dashboards that identify work by inline Linear
  ids (0..N per card) rather than one id per card — for Pinata-shape KBs; pairs
  with the [`../reference/specs/dashboard-format.md`](../reference/specs/dashboard-format.md)
  contract. Built as backlog v3 (TASK-013…016). **accepted**, 2026-07-21.
