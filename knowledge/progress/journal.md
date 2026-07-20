# Journal

Append-only. Newest at the bottom. Each entry is a snapshot for future-me
with no memory of this session.

## Entry format

---
## YYYY-MM-DD HH:MM — <short heading>
   (HH:MM may be dropped only if the heading alone makes the entry findable
   from a DONE.md pointer)

**Task:** TASK-NNN (or "scaffolding" / "exploration" / "blocker triage")
**What I did:** 1–3 sentences.
**What I verified:** which gates ran, with literal command output worth keeping.
**What changed:** files touched + the delivery record per my profile.
**What I learned:** anything that would surprise future-me. Non-obvious only.
**Next:** the very next thing to do on resume.

## Entries

---
## 2026-07-20 13:43 — Framework init + spec → backlog + bootstrap

**Task:** scaffolding
**What I did:** Ran SETUP.md adoption against this project. Created the five
instance areas and their files from the skeletons; wrote the manifest
(`knowledge/README.md`) declaring **delivery: pr**, the Locations block (default
`knowledge/…` paths), attribution = user identity only (no agent attribution —
user decision), and the session envelope = ship the MVP (TASK-001…007). Wrote
the project `CLAUDE.md`, `.gitignore`, root `README.md`, MIT `LICENSE`. Filled
the brief and glossary from the spec; recorded the real Go local-CI commands in
`reference/local-ci.md`. Relocated the spec to
`reference/specs/bklg-spec.md` (its own closing note says it belongs there).
Processed the spec's §15 into a 7-task `BACKLOG.md`; pulled **TASK-001** into
`CURRENT.md` with acceptance criteria (each naming its decider). Initialized git
and made the one allowed bootstrap commit on `main`.

**What I verified:** No code yet — nothing to run. Framework guard checked
before starting: no pre-existing `planning/CURRENT.md`, so this was a genuine
first-time adoption, not a re-run against live state. Bootstrap commit + push
output quoted in the git-bootstrap step below.

**What changed:** New: `CLAUDE.md`, `.gitignore`, `README.md`, `LICENSE`,
`knowledge/README.md`, `knowledge/planning/{CURRENT,BACKLOG,DONE}.md`,
`knowledge/progress/{journal,blockers}.md`, `knowledge/decisions/INDEX.md`,
`knowledge/reference/{project-brief,glossary,local-ci}.md`,
`knowledge/reference/specs/bklg-spec.md` (moved), `knowledge/whiteboard/README.md`.
Delivery record: **bootstrap commit on `main`** (pre-framework, covered by the
manifest's Rule-1 exception — the only direct-to-main commit this project gets).

**What I learned:** The GitHub repo `gillchristian/bklg` was pre-created by the
user — **empty, public**, default branch not yet set — so the bootstrap push
establishes `main` as the default branch, after which every change goes via PR.
`gh` is authed as `gillchristian` over SSH. The framework's upstream repo/path
is not recorded (it was pre-placed at `knowledge/framework/`); logged as a
non-blocking open question in the brief. No remote CI exists, so delivery gate
D3 is vacuous until GitHub Actions is added (parked).

**Next:** Start TASK-001 (CLI + server skeleton): branch `feat/task-001-…`,
`go mod init github.com/gillchristian/bklg`, implement `cmd/bklg/main.go`
(`splitArgs`, flags, 127.0.0.1 bind, `GET /{$}`→200, startup block), local CI,
PR, fresh-context review, squash-merge, close PR.
