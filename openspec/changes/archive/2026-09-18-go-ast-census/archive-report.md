# Archive report: go-ast-census

## Status

**PASS — archived.** Native `gentle-ai.sdd-status` v2 was refreshed immediately before archive work: selected change `go-ast-census`, store `openspec`, archive dependency `ready`, tasks `9/9` complete, next recommended `archive`, no blocked reasons, and repo-local action context confined to `/data/Projects/git-change-evidence-worktrees/go-ast-census-archive`.

## Artifacts read and preservation

Read `proposal.md`, `specs/go-ast-census/spec.md`, `design.md`, `tasks.md`, `apply-progress.md`, `verify-report.md`, and `openspec/config.yaml`. The optional current verification report is canonical **PASS** with no unresolved FAIL, BLOCKED, CRITICAL, or verification blockers. The final persisted task read found no `- [ ]` implementation markers.

The current PASS verification report and all three historical failed reports were retained byte-for-byte in the archived change:

| Artifact | SHA-256 before archive |
|---|---|
| `verify-report.md` | `e54814e0bc8222e04a493bc5710dd7bfe66912a3b62f94226c08574756e87368` |
| `verify-report.failed-initial.md` | `f088b195421dc4ba682c6fa48b2d016de86715c64cee9d0a8db9493d80dba7e9` |
| `verify-report.failed-c3.md` | `67973625eb60816fa3489713054e359914e3391dcb083a39e2b072914d7a8a97` |
| `verify-report.failed-docs.md` | `e9938fc780edb3cd06adf7c7299c045aa5abccd0be313bddc5691ea76918551f` |

## Canonical specification composition

Domain `go-ast-census` had no existing canonical target. The validated change spec is a full domain specification, so it was mechanically copied to `openspec/specs/go-ast-census/spec.md`; the native existing-canonical composition provider was not applicable. The empty `diff -r` readback recorded below confirms byte identity.

- ADDED / MODIFIED / REMOVED delta operations: none; the source is a full new-domain specification.
- Full-domain requirement content created: GAC-001, GAC-002, GAC-003, GAC-004, GAC-005, GAC-006, GAC-007, GAC-008, GAC-009, and GAC-010.
- Unsupported RENAMED operations: none.
- Same-domain active-change collision: none, per refreshed native status and active-change scan.
- Destructive merge: none; no REMOVED or MODIFIED operation was applied, so no destructive approval was required.

## Archive move

Destination collision check passed for `openspec/changes/archive/2026-09-18-go-ast-census`. The archive report was written before the mechanical move. The source tree was snapshotted with `cp -R`, moved with filesystem `mv` to honor the explicit no-staging instruction, and compared recursively against the snapshot. The empty `diff -r` readback confirms the archive preserved every artifact, including this report.

- Archived path: `openspec/changes/archive/2026-09-18-go-ast-census/`
- Active source after move: absent.
- Non-critical partial archive approval: none.
- Stale-checkbox reconciliation: none; no unchecked implementation task existed.
