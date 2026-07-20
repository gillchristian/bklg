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

### TASK-003 — Planning parser → model
**Source:** BACKLOG (spec §15.3; contract §2, model §4, reconciliation §5, fixtures §11)
**Acceptance criteria:**
- [ ] AC1 — Cards & columns: parsing the `testdata/knowledge` instance yields the expected set of cards; each card's `{id, column, len(Acceptance), ParkingLot}` matches an expected table (asserted in `parse_test.go` and quoted in the journal). (Decider: table-equality test on the built board.)
- [ ] AC2 — Dedup, most-advanced-state wins (`Done > InProgress > Backlog`, spec §5): an id in both `BACKLOG` (unchecked) and `CURRENT` renders **once** as In-Progress; an id `[x]` in `BACKLOG` **and** present in `DONE` renders **once** as Done. (Decider: test asserts the deduped id appears exactly once and in the furthest column.)
- [ ] AC3 — The three reconciliation warnings fire on seeded inconsistencies (spec §5): (1) `[x]` in `BACKLOG` with no matching `DONE` entry → `shipped item missing from DONE.md`; (2) a `DONE` id not `[x]` in `BACKLOG` → `DONE item not ticked in BACKLOG`; (3) >1 task under `CURRENT ## Active` → `CURRENT holds >1 active task`. (Decider: test asserts each warning is present and references the seeded id; and that exactly these three fire — no unexpected warnings.)
- [ ] AC4 — Field extraction (spec §2 table): CURRENT tasks capture `Source`, `Acceptance` (each `Criterion` with its checked state), `Notes`, `DeliveryOverride`; DONE entries split into `{Date, Summary, DeliveryRecord, JournalPointer}`; a parking-lot item with no id → `ID==nil`, `ParkingLot==true`. `Card.Raw` keeps the source block. (Decider: test asserts these fields on specific fixture cards.)
- [ ] AC5 — Parse defensively (spec §2): a malformed entry never panics or blanks the board — captured best-effort, skipped otherwise, and recorded as a warning where applicable. (Decider: a unit test feeds malformed lines and asserts no panic + good entries still parse.)

**Notes:**
- New files: `internal/backlog/model.go` (types from §4 — `ID`, `Column`, `Criterion`, `DoneRecord`, `Card`, and `Board{Cards, Warnings, Meta}`; `Blocker` + `Board.Blockers` land in TASK-004) and `internal/backlog/parse.go` (line-oriented parser per §2, behind a small `Parser` interface per decision **D4** so a goldmark impl can replace it later — honored deliberately despite the usual "no premature abstraction" rule because the spec calls for the seam).
- Parser is keyed to the exact conventions in §2: IDs `[A-Z]+-\d+`, dates `\d{4}-\d{2}-\d{2}`, em-dash `—` (U+2014) as the DONE field separator, `- [ ]`/`- [x]` checkboxes, `### <ID> — <title>` for CURRENT, `- [ ] <ID> — <title>` for BACKLOG.
- Build the rich `testdata/knowledge/` fixture (§11) seeding: unchecked backlog item w/ id; a freeform parking-lot item w/o id; a `[x]` backlog item that **is** in DONE and one that **isn't** (warning 1); a DONE entry **not** `[x]` in backlog (warning 2); two CURRENT tasks (warning 3) — one with AC + delivery override, one with **no** AC (`no-ac`); and a `blockers.md` with one open + one resolved blocker (consumed by TASK-004). This instance doubles as the live demo (`bklg internal/backlog/testdata`).
- Wire the parser into `main` so `/` (still the placeholder page) can later render the board — but TASK-003 stops at model + `/_diag`-ready warnings; template rendering is TASK-005. Expose the board build so TASK-004/005 consume it.
- `Meta.LatestMTime` (for `/_v`) can be stubbed here and finished in TASK-007.
