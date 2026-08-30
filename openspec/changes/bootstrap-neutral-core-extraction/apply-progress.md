# Apply Progress: Bootstrap Neutral Core Extraction

## PR 1 Historical Slice
**Completed:** PR 1 / Phase 1 — Bootstrap Decision (`auto-chain`, `stacked-to-main`).

This planning/configuration slice establishes Go-first boundaries only. It adds no Go source, tests, `go.mod`, compatibility vectors, CNSIC changes, commits, pushes, or pull requests.
## Completed Tasks

- [x] **1.1** Recorded the single root `AGENTS.md`, Go 1.25.10 target, root public API/`cmd/git-change-evidence/`/`internal/` layout, co-located tests, Go gates, and immutable typed configuration boundary.
- [x] **1.2** Replaced the pre-bootstrap skip rubric with the authoritative Go strict-TDD rubric, exact commands, layout thresholds, and a zero-file initial crossing baseline.

## Files Changed

- `AGENTS.md`
- `README.md`
- `docs/product-charter.md`
- `openspec/config.yaml`
- `openspec/changes/bootstrap-neutral-core-extraction/exploration.md`
- `openspec/changes/bootstrap-neutral-core-extraction/proposal.md`
- `openspec/changes/bootstrap-neutral-core-extraction/design.md`
- `openspec/changes/bootstrap-neutral-core-extraction/specs/neutral-evidence-contracts/spec.md`
- `openspec/changes/bootstrap-neutral-core-extraction/tasks.md`
- `openspec/changes/bootstrap-neutral-core-extraction/apply-progress.md`

## Verification Evidence

| Check | Result |
|---|---|
| `git diff --check` | PASS |
| `python3` with PyYAML: parse `openspec/config.yaml` | PASS |
| `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 \| xargs -0 -r gofmt -l)"` | PASS; no Go files exist and the check performed no writes |
| Search for active legacy-language library requirement in `README.md`, `docs/`, and `openspec/` | PASS; none remains |
| Root `AGENTS.md` inventory | PASS; exactly one file at `./AGENTS.md` |
| `go test ./...` | N/A for PR 1; the approved no-module boundary means no `go.mod` or Go package exists yet |
| `go test -race ./...` | N/A for PR 1; no concurrency-bearing Go surface exists |
| Runtime harness | N/A; this slice has no runtime boundary |

The executor host reports `go version go1.27.0-X:nodwarf5 linux/amd64`, not the selected Go 1.25.10 baseline. The host toolchain was not used as evidence for the target baseline.

## TDD Cycle Evidence

| Task | RED | GREEN | TRIANGULATE | REFACTOR | Result |
|---|---|---|---|---|---|
| 1.1 | N/A — planning documentation only; no production behavior | N/A | N/A | N/A | Complete |
| 1.2 | N/A — test-policy configuration only; no production behavior | N/A | N/A | N/A | Complete |

Strict TDD is active for later Go production work. The required support file `.pi/gentle-ai/support/strict-tdd.md` is absent, so later implementation must follow the RED → GREEN → TRIANGULATE → REFACTOR contract directly and retain this risk until the support file is supplied.

## Design Conformance and Deviations

- The runtime is Go-first; all active Python-library requirements were replaced with a neutral versioned Go package API and canonical-document contracts.
- The Functional Core accepts immutable typed Go structs/interfaces. Configuration parsing remains adapter-only; TOML, YAML, JSON, and other syntax choices remain deferred.
- The public/internal layout and thresholds are development guidance only, never core evidence or governance policy.
- **Intentional deviation from future implementation:** no `go.mod` or Go files were added. Module/publication identity is still unresolved and outside PR 1.

## Workload and PR Boundary

```text
main
└── PR 1 — Bootstrap Decision 📍
    └── PR 2 — Neutral W1/W2 implementation
```

- **Strategy:** `auto-chain` with `stacked-to-main`.
- **Current boundary and review budget:** only tasks 1.1 and 1.2; 285 additions + 114 deletions = 399 changed lines.
- **Dependency:** immutable CNSIC W1/W2 predecessor oracle at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806`.
- **Rollback boundary:** remove the listed PR 1 decision/configuration artifacts or restore their pre-PR-1 wording; no production or consumer behavior is affected.
- **Follow-up:** PR 2 may begin task 2.1 only after its owner-authorized module identity is available. Do not begin any later task in this PR.

## Remaining Tasks

- [ ] 2.1 through 2.3 — Neutral W1/W2 implementation.
- [ ] 3.1 through 3.4 — Staged carveout, report, CLI/publication, and forecast capabilities.
- [ ] 4.1 through 4.3 — CNSIC profile, separate cutover handoff, and distribution decision.

## Unresolved Decisions and Risks

- Module/publication identity, packaging, CI, license, release, distribution, and concrete configuration-file syntax remain unresolved by design.
- The executor's installed Go toolchain is not Go 1.25.10.
- `.pi/gentle-ai/support/strict-tdd.md` is absent. This did not affect planning-only work, but it is a strict-TDD risk for the first Go implementation slice.

## PR 2 / Task 2.1 — Neutral W1/W2 contracts
- [x] **2.1** Closed/versioned contracts, canonical bytes/SHA-256/provenance, authority rejection, and the corrected predecessor corpus pass all acceptance criteria.
- **Initial record (corrected):** the prior “frozen W1/W2 vectors” claim was overbroad: both vectors were successor-authored, not predecessor bytes.
- **Independent rejection / correction:** replaced them with exact predecessor `contracts-v1.json`, `accounting-v1.json`, and `manifest.sha256`; deleted both misleading vectors.
- **Files:** `go.mod`, `contract.go`, `decode.go`, `contract_test.go`, `decode_test.go`, `compatibility_test.go`, the four corpus/manifest files, `openspec/config.yaml`, `tasks.md`, and this record.
- **Verification:** focused corpus RED/GREEN and decoder characterization, `go test .`, `go test ./...`, `go vet ./...`, and `git diff --check origin/main` → PASS.
- **Format:** check-only `gofmt -l` → PASS; **Race:** N/A (no concurrency).
### TDD Cycle Evidence
| Task | RED | GREEN | TRIANGULATE | REFACTOR |
|---|---|---|---|---|
| 2.1 initial | Public-symbol compile failure | Constructor/canonical pass | Earlier rejection proof, now narrowed | Decoder simplification pass |
| 2.1 correction | Corpus test failed on false manifest | Exact predecessor imports pass | Negative decoder cases already passed (characterization) | `slices.Contains` pass |
### Design Conformance and Deviations
- Functional Core only; the standalone manifest is successor-authored metadata binding predecessor identities, not a predecessor artifact.
### Workload / PR Boundary
- **Budget:** 393 additions + 4 deletions = 397 changed lines relative to `origin/main` (hard cap: 400).
- **Boundary:** stacked-to-main PR 2 task 2.1 only; rollback removes PR 2 contracts, tests, and corpus without affecting PR 1 planning.
### Remaining Tasks / Risks
- [ ] 2.2 onward remain outside this slice; strict-TDD support was read from the predecessor checkout because it is absent here.
