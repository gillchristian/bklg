# Glossary

"**Term** — definition." Add entries the moment a term first appears with a
specific meaning. Surprising semantics first.

- **Card** — one task as rendered on the board; derived by joining a task id
  across the parsed files. Parking-lot items without an id are cards too, but
  get no detail page.
- **Column** — one of Backlog / In Progress / Done. Sourced, respectively, from
  `BACKLOG.md` (`## Active` unchecked + all `## Parking lot`), `CURRENT.md`
  (`## Active`), `DONE.md` (`## Completed`).
- **Dedup rule (most-advanced-state wins)** — a task id may appear in several
  files; it renders as exactly one card, in its furthest column:
  `Done > InProgress > Backlog`.
- **Badge** — cross-cutting state chip on a card: `blocked` (id is the
  `Task affected` of an *open* blocker — the load-bearing one), `parking`,
  `override` (card has a `**Delivery override:**`), `no-ac` (in-progress card
  with zero acceptance criteria), and the namespace chip (the id's namespace,
  e.g. `MONO`, `TRAIL`).
- **ID** — `{Namespace, Number, Raw}`, e.g. `{"MONO", 6, "MONO-006"}`. Matches
  `[A-Z]+-\d+`. Route key is the **full** id (`/MONO-006`), case-insensitive —
  bare numbers collide across namespaces (D2).
- **Planning area / progress area** — the two instance areas bklg reads,
  resolved through the manifest's `## Locations` block (repo-root-relative),
  falling back to `<dir>/planning` + `<dir>/progress`.
- **Locations block** — the `## Locations` section of a knowledge `README.md`
  manifest mapping roles (`planning`, `progress`, …) → repo-root-relative paths.
- **Root manifest / system index** — a monorepo-root `README.md` with no
  planning area of its own but a table whose rows name `systems/<name>`; bklg
  can't board it directly (v1) and instead lists the per-system invocations.
- **`/_v`** — endpoint returning the max mtime across parsed files as a bare
  integer; the page polls it (~3s) and reloads on change.
- **`/_diag`** — endpoint listing parse warnings verbatim, one per line.
- **Warning** — a non-fatal parse/reconciliation diagnostic collected during a
  parse and surfaced in a subtle banner linking to `/_diag`.
