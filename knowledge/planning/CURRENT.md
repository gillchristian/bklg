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

### TASK-010 — Render a safe markdown subset
**Source:** BACKLOG (v2 feedback; `../whiteboard/trail-instance-findings.md` #4)
**Acceptance criteria:**
- [ ] AC1 — Inline markdown renders: `**bold**`/`__b__` → `<strong>`, `*italic*`/`_i_` → `<em>`, `` `code` `` → `<code>`, `[text](url)` → `<a>` — in card titles and the detail fields (title, summary, notes, acceptance text, blocker body). (Decider: a render/unit test asserts the tags for a markdown input.)
- [ ] AC2 — **Escape-first safety (no injection):** the renderer HTML-escapes the input **before** applying any markdown transform, so raw repo HTML never survives — `<script>` → `&lt;script&gt;`; a `[x](javascript:alert(1))` link is neutralized (unsafe scheme dropped); only whitelisted tags (`strong`/`em`/`code`/`a`) with escaped attributes are emitted. (Decider: unit test with hostile inputs asserts no live `<script>`, no `javascript:` href; `go test` + the reviewer's probe.)
- [ ] AC3 — Applied on real fields: on `trail --dir systems/track/knowledge`, `**Swift/iOS toolchain…**` shows as **bold**, `` `code` `` as code — not literal `**`/backticks. (Decider: binary smoke — card/detail HTML contains `<strong>`, no literal `**` in the rendered field.)
- [ ] AC4 — Trim + markdown compose safely: the board still trims (`truncate` then render), and an unmatched marker left by truncation stays **literal** (no unclosed `<strong>`/broken tag). (Decider: render test on a long bold title asserts no unclosed tag; matched-pairs-only rule.)
- [ ] AC5 — ADR-0002 records the decision: **hand-rolled strict inline subset, stdlib-only** (keep the brief's zero-dep hard constraint) vs. pulling in goldmark + a sanitizer. (Decider: `decisions/0002-*.md` exists + indexed.)

**Notes:**
- **Zero-dep fork (the consequential one):** default = hand-roll a small, safe inline renderer, keeping "zero Go module deps" (a brief hard constraint). Safe because it **escapes first** (`html.EscapeString`) then applies a whitelist of inline patterns on the escaped text, returning `template.HTML`. Do NOT relax zero-dep without a strong reason — record the call in ADR-0002.
- `renderMarkdown(s) template.HTML` in `render.go`; register as a `md` template func. Order: escape → protect code spans → links (validate scheme: allow http/https/mailto/relative `#`,`/`; escape href) → bold → italic. Unmatched markers stay literal (matched-pairs-only) → AC4.
- Compose with trim in `board.html`: `{{md (truncate .Title 140)}}`. Detail uses `{{md .Card.Title}}` etc. Introducing `template.HTML` bypasses auto-escaping, so the escape-first discipline is load-bearing — this is the one place to be paranoid; expect a hard security review.
- Scope v1 to **inline** markdown (bold/italic/code/links). Block-level lists/headings are a follow-up (note in parking if wanted). Line breaks already preserved via `whitespace-pre-wrap` where used.
