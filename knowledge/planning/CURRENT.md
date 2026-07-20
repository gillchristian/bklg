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

### TASK-008 — Trim card text on the board; full text on the detail page
**Source:** BACKLOG (v2 feedback; `../whiteboard/trail-instance-findings.md` #2)
**Acceptance criteria:**
- [ ] AC1 — Long titles are trimmed on the board: a card whose title exceeds the limit renders a truncated title ending in `…` on `/`. (Decider: `truncate` unit test + a render test on a crafted long-title card asserts the board output contains the truncated prefix + `…` and **not** the full string.)
- [ ] AC2 — Full text preserved on the detail page: `/<id>` shows the **full**, untrimmed title. (Decider: render the task template for the same long-title card; assert the full string is present.)
- [ ] AC3 — Short titles unchanged: a title at/under the limit renders verbatim, no `…`. (Decider: `truncate` unit test + the existing fixture cards (short) still render their full titles — `TestBoardRender` stays green.)
- [ ] AC4 — Rune-safe: truncation never splits a multibyte rune. (Decider: `truncate` unit test on a multibyte string.)

**Notes:**
- Add a `truncate(s string, max int) string` template helper in `render.go` (rune-count based, light word-boundary backoff, appends `…`); use it in `board.html` for the card title only. `task.html` keeps `{{.Card.Title}}` (full).
- Limit ≈140 runes (≈2–3 lines). Trim operates on the raw title now; TASK-010 (markdown) will refine so a trim can't cut mid-`**` once titles render as markdown.
- Don't touch the fixture card table (keep `TestParseCardTable` stable); test trimming via `truncate` unit cases + crafted view models so the board/detail split is asserted without rippling counts.
- This is display-only — the model/`Card.Title` is unchanged; only the board's rendering trims.
