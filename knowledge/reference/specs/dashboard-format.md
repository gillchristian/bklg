# Dashboard format — bklg's second input contract

**Status:** proposed — paired with [ADR-0004](../../decisions/0004-dashboard-adapter.md);
implemented by backlog v3 (TASK-013…016).

This is the input contract for bklg's **dashboard adapter**: the mode that reads
a single-file Active / Backlog / Done dashboard instead of the framework
`planning/` + `progress/` skeleton. The main spec's §2 is the framework-mode
contract; this is its sibling for KBs (like the real Pinata KB) that track work
as one dashboard file keyed by external ticket ids.

As in §2, **the patterns below *are* the contract.** The adapter parses
defensively — a malformed row degrades to a `/_diag` warning, never a crash or a
blank board — but the board is only as useful as the file's regularity, so a
target KB should keep to these rules. The two calls this contract makes
deliberately (both recorded in ADR-0004): mode is selected by an explicit
`dashboard:` config key, not auto-detected; and "blocked" is a **leading** `⛔`
only, because dashboards already use `⛔` mid-prose as a decision marker.

The version handed to a target KB's maintainers as standalone instructions is
reproduced at the repo root as `pinata.md`; keep the two in sync.

---

## Selection & resolution

- The manifest lookup tries `README.md` then `index.md` under the knowledge dir.
- A `## Locations` block key `dashboard: <path>` (repo-root-relative) points at
  the single dashboard file and selects the adapter. A `--dashboard <file>` flag
  is the zero-config escape hatch.
- In dashboard mode the `planning/` directory requirement is lifted; the
  resolved target is the one dashboard file. Framework mode is unchanged.

## The dashboard file

### Sections — exactly three headings

`## Active`, `## Backlog`, `## Done` (case-insensitive). Other `##` sections are
ignored.

### `## Active` and `## Done` — pipe tables

GitHub-style pipe tables with these exact header columns:

```
## Active
| Work | Material | Status / next step |
|---|---|---|
| **Linkable form pages** — CMS pages that render a report form | [linkable_form_pages/](x) | ⛔ Blocked on PINATA-602. Planned 2026-07-21; 4 issues. |

## Done
| When | What | Record |
|---|---|---|
| 2026-06-12 | **Same-worksite clone (PINATA-478)** — clone pages in place | [same_worksite_clone.md](x) |
```

- **One row = one physical line** (cells may be long, freeform prose).
- **Literal pipes inside a cell are escaped `\|`** — an unescaped `|` starts a
  new column.
- **Done's first column is a date** — `YYYY-MM-DD` preferred; partial
  (`2026-05`, `~2026-03`) tolerated.

### `## Backlog` — bullet groups

Optional bold group headers (a bold line ending in a colon), then `-` bullets
each leading with a bold title:

```
## Backlog

**Product / code:**

- **CMS Page Types (PINATA-599, PINATA-601)** — new CMS page types. PINATA-601 remains.

**Knowledge base:**

- **Retrofit frontmatter** — cms/ done; chat_nav/, ai_nav/ still unswept.
```

The group name becomes a chip on the card; ungrouped bullets still parse (no
chip).

### Card titles

Every Active row, Done row, and Backlog bullet **begins with a `**bold**`
phrase** — the card title. Text after an em-dash (` — `, U+2014) is a subtitle.

### Linear tickets

Bare ids matching `[A-Z]+-\d+` (e.g. `PINATA-602`) anywhere in a row become
clickable chips (0..N per card). The Linear base URL is configurable
(`linear:` Locations key / `--linear-base`, default
`https://linear.app/gopinata/issue/`). GitHub PR refs like `#11531` are ignored.

### Blocked

A card is blocked iff its status cell (Active) or bullet prose (Backlog)
**starts** with `⛔` or `**Blocked**`. A `⛔` used mid-prose (a decision/attention
marker) is ignored.

## Model & routes (dashboard mode)

- A `Card` carries 0..N tickets, an optional group label, a `blocked` flag, and
  the captured Material link; no acceptance criteria, no cross-file dedup.
- Detail route key is a slug of the title (collisions disambiguated); the detail
  page shows title, column, linked tickets, Material link, and the raw row.

## What a target KB does *not* need

Freeform prose in cells, no per-row ids, no `planning/`/`progress/` dirs, no
`blockers.md`; Linear stays the source of truth for status — this file is the
map. Divergence from the rules above shows up in `/_diag`, not as silent
corruption.
