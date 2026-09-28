# Design first Python syntax census (planning only)

Issue: https://github.com/kozz36/git-change-evidence/issues/141
Branch/worktree: `docs/python-census-design` at `/data/Projects/git-change-evidence-worktrees/python-census-design`; base `1abe22dbcf51bae7453087b23d64a15c6b1b725e`.

## Objective and fixed choices
Produce a reviewable design, not implementation, for an evidence-only Python syntax census. Owner selected direct `Name.Attribute(...)` calls, explicitly selected `.py` and `.pyi` files, atomic failure on invalid/unsupported syntax, and a standalone Go binary with neither installed CPython nor cgo. A pure-Go parser is a candidate only after independent feasibility proof of grammar/stubs, physical byte spans, error behavior and resource limits. No release, CNSIC pilot, semantic resolution or delivery authority.

## Tasks
- [x] D1 — Map Go census invariants and post-preview scope boundary; confirm choices with owner, create planning issue #141 and isolate clean worktree. Findings: `docs/census.md`, public `census.QueryV1` and Go AST spec; current post-preview neutral-capabilities exploration explicitly excludes multilingual census. External parser README claims are unverified product evidence.
- [x] D2 — Wrote `docs/python-census-design.md` with selected match predicate/negative matrix, proposed complete inventory/API-CLI/byte-evidence contract and parser proof gate. Corrected two post-writer readback defects: Go-derived invariants were mislabelled OWNER-DECIDED rather than PROPOSED, and the candidate parser URL pointed to the wrong owner; line/column remains provisional while byte offsets require proof. No dependency chosen.
- [x] D3 — Independent read-only review found two low issues, corrected: Go limit wording and stale task next action. Parent read back both files; staged diff check passed; five relative Markdown links resolve. An initial naive regex falsely treated `os.open[0]("x")` inside inline code as a link; a code-aware check corrected that verification error. Planning-only work-unit commit `2dd2c7c6492ed01565f5259691b158e40a52e007` created. PR remains gated by issue #141 approval and owner delivery decision; no push or merge.

## Checks and rollback
`git diff --check`; resolve added relative Markdown links; inspect diff for implementation/dependency changes (must be none). Runtime harness N/A because this is a planning document with no executable boundary. Rollback: design doc and this task record only. Delivery strategy ask-on-risk; forecast <400 authored changed lines.

Memory mirror pending: this Pi session belongs to CNSIC and cannot write the GCE-project Engram task mirror. Local file is the recovery copy.
Next action: owner reviews/approves issue #141 on GitHub, then confirm label before any push/PR. This task closure is a separate follow-up commit. No parser implementation is authorized by the design.
