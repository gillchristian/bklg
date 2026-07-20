# knowledge/ — project manifest (bklg)

The working system lives in [`framework/`](framework/README.md); this file is
what makes it THIS project. Reading order: this file → `framework/README.md`
→ the enabled profile in `framework/delivery.md`.

`bklg` is a read-only backlog-board viewer for framework knowledge instances
(the full product intent lives in [`reference/project-brief.md`](reference/project-brief.md)
and the spec at [`reference/specs/bklg-spec.md`](reference/specs/bklg-spec.md)).
It is built *with* this framework and, once working, can view *itself* — the
in-repo `knowledge/` instance doubles as a live dogfood target.

## Delivery mode

delivery: pr

**Operative meaning:** I own the full branch → PR → self-merge cycle. The
default branch (`main`) is sacrosanct: never pushed to or committed on
directly — every change lands via a PR I open, review (fresh-context), and
merge myself with `--squash --delete-branch`. Read-only git (`status`, `diff`,
`log`, `show`) is always fine. The single exception is the pre-framework
**bootstrap commit** on `main` that establishes the repo (this manifest, the
framework, the instance areas, scaffolding) — explicitly covered here, per the
`pr` profile's Rule 1 parenthetical.

## Locations

The role → path map the framework dereferences. Paths are repo-root-relative.
These are the standalone defaults.

framework:  knowledge/framework
planning:   knowledge/planning
progress:   knowledge/progress
decisions:  knowledge/decisions
reference:  knowledge/reference
whiteboard: knowledge/whiteboard

## Project rules

- **Identity/attribution:** commits and PRs carry the user's identity only —
  `gillchristian <gillchristiang@gmail.com>`. **No agent attribution:** no
  `Co-Authored-By` trailer, no "generated with" footer. (User decision,
  2026-07-20; this overrides the harness's default Co-Authored-By line.)
- **Session envelope:** ship the **MVP** — the seven executable slices of the
  spec's §15 (TASK-001…TASK-007) — then run the end-of-session sweep and stop.
  Stop early only if the escape hatch fires (`framework/when-stuck.md`:
  everything blocked). The spec's §13 non-goals are explicitly v2+ and out of
  this envelope.
- **Merge strategy:** `gh pr merge --squash --delete-branch`.
- **Remote check (gate D3):** none recorded — the repo has no remote CI
  (GitHub Actions) in v1, so D3 is vacuous. Adding CI is a parking-lot item;
  record the remote-check command in `reference/local-ci.md` the first time
  one exists.
- **PR rhythm:** each task ships as a feature PR (implementation) followed by a
  tiny `docs/task-NNN-close` PR (move to DONE + journal + tick BACKLOG, and —
  when known — orient the next task into CURRENT). Bootstrap is the one direct
  `main` commit.
- **Where knowledge/ lives:** in-repo, committed by the agent (allowed under
  `pr`). It records how the tool was built and serves as bklg's dogfood fixture.
- **Repo:** `github.com/gillchristian/bklg` (module path); remote `origin` =
  `git@github.com:gillchristian/bklg.git` (SSH); default branch `main`;
  visibility public (pre-created by the user).
- **Framework copy:** v5 (2026-07-08). Upstream not recorded by the provider —
  the `framework/` directory was pre-placed in this project. Don't edit
  `framework/` here; improvements go upstream (fill in the upstream when known).

## The loop, instantiated (pr mode)

1. **Orient** — Read this manifest, then `planning/CURRENT.md`. If empty,
   promote the top unchecked `BACKLOG.md` item.
2. **Plan** — Write acceptance criteria into `CURRENT.md`, each naming its
   decider (`framework/verification.md`). For task N>1 this was already done by
   task N−1's close PR ("orient in one step"); for TASK-001 it was done in the
   bootstrap.
3. **Stage** — `git switch main && git pull --ff-only`, then branch
   `<kind>/task-NNN-<slug>` (kind ∈ feat/fix/chore/refactor/docs/test).
4. **Execute** — Implement, committing as I go; each commit leaves the branch
   sane.
5. **Verify** — Run local CI (`reference/local-ci.md`: `go build ./...`,
   `go vet ./...`, `gofmt -l` clean, `go test ./...`) + a real smoke test of the
   touched surface; quote output. Gates in `framework/verification.md`.
6. **Deliver** — Push; `gh pr create` (imperative title ≤72, template body);
   obtain the **fresh-context review** (spawn a subagent given ONLY the diff +
   AC — never this transcript — grading each criterion pass/fail); fix or rebut
   confirmed findings; `gh pr merge --squash --delete-branch`;
   `git switch main && git pull --ff-only`.
7. **Log** — Append a `journal.md` entry: timestamp, what I did, what I verified
   (quoted output), delivery record (``PR #N, merged `sha` ``), next.
8. **Advance** — Ship the **close PR** (`docs/task-NNN-close`): move CURRENT→DONE
   with PR#+sha, tick BACKLOG, add the journal pointer, and pull the next task
   into CURRENT with its AC. Merge it (review-exempt). Then check the session
   envelope; if the MVP is shipped, run the end-of-session sweep and stop.

## Layout

- framework/ — the reusable system. planning/ — CURRENT (one task), BACKLOG,
  DONE. progress/ — journal (append-only), blockers. decisions/ — ADRs +
  INDEX. reference/ — brief (wins conflicts), glossary, local-ci, specs/.
  whiteboard/ — discussions in flight.
