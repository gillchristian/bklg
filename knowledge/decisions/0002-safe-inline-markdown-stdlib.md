# 0002 — Hand-rolled safe inline markdown, stdlib-only

**Date:** 2026-07-20 · **Status:** accepted

## Context

TASK-010 renders a markdown subset in card/detail text (real instances write
`**bold**`, `` `code` ``, links; v1 showed them literally). The textbook safe
approach pairs a markdown parser (goldmark) with an HTML sanitizer
(bluemonday). But the project brief lists **"zero Go module dependencies"** as a
hard constraint (pure-stdlib `go install`, no supply chain). So there is a real
fork: relax that hard constraint, or hand-roll.

## Decision

Hand-roll a small **inline** renderer (`markdown.go`), stdlib-only, keeping
zero-dep. It is safe by construction: (1) HTML-escape the input **first**
(`html.EscapeString`) so no repo-authored markup can survive; (2) then apply a
fixed whitelist of inline patterns (`**bold**`, `*italic*`, `` `code` ``,
`[text](url)`), emitting only `strong`/`em`/`code`/`a`; (3) scheme-check link
hrefs (`safeURL`) so `javascript:`/`data:` never become live links. Returns
`template.HTML` — the one and only place bklg emits raw HTML.

Scope: **inline only** (no block lists/headings). **Asterisk emphasis only** —
underscore emphasis is dropped so `snake_case` and `__dunder__` identifiers
aren't mangled.

## Alternatives considered

- **goldmark + a sanitizer (bluemonday).** Rejected: two dependencies, defeats
  the brief's hard constraint, and heavier than an inline subset needs. The
  escape-first design is safe *without* a sanitizer precisely because we never
  emit repo HTML — the sanitizer's job (scrub untrusted HTML) doesn't arise.
- **Keep rendering literal (v1 behaviour).** Rejected: the user explicitly asked
  for markdown.
- **Full CommonMark by hand.** Rejected: emphasis flanking rules, nested
  constructs, and block parsing are exactly what a library is for.

## Consequences

- Safe by construction, zero-dep preserved. The escape-first invariant is
  load-bearing — any future edit to `renderMarkdown` must keep "escape before
  transform, whitelist tags, scheme-check hrefs."
- Known limits: no block markdown; simple regexes have flanking quirks (e.g.
  `a*b*c` can italicize `b`; a glob like `**/*.go` can render oddly). If real
  content needs block constructs or the quirks bite, that is the trigger to
  adopt goldmark **with** a sanitizer — the D4 seam / ADR-0001 escalation path.
