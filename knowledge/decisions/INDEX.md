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
