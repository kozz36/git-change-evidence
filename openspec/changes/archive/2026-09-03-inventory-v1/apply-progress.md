# Apply Progress: Inventory V1

## Work Unit 2A — Complete
The entire 2A section below is historical state at 2A close; the later 2B section is current and authoritative.

Implemented the pure additive Inventory V1 construction contract only. `NewInventoryV1` revalidates the exact Policy V1 canonical antecedent, derives its link, copies and derives content evidence, preserves raw path bytes, orders by `bytes.Compare`, rejects invalid/duplicate paths, and exposes defensive views. No decoder or 2B behavior was added.

### TDD Cycle Evidence

| Task | Phase | Command | Result |
|---|---|---|---|
| 2A.1 | RED | `go test . -run '^TestNewInventoryV1CanonicalContract$'` | exit 1: `InventoryEntryInput` and `NewInventoryV1` undefined. |
| 2A.2 | GREEN | same focused command | exit 0; exact `/wph` canonical fixture, policy link, and Inventory digest passed. |
| 2A.3 | TRIANGULATE | same focused command | exit 0; byte-identical policy reproduction, byte-distinct policy divergence, and no digest/length input fields passed. |
| 2A.4 | REFACTOR | same focused command | exit 0 after private entry-derivation helper extraction. |
| 2A.5 | RED | `go test . -run '^TestNewInventoryV1ExactContentAndPaths$'` | exit 1: zero-length content rejected as `entries.content: required`. |
| 2A.6 | GREEN | same focused command | exit 0; exact binary, empty, LF, and CRLF content evidence passed. |
| 2A.7 | TRIANGULATE | same focused command | exit 0 after raw-byte ordering plus invalid/duplicate path vectors passed. |
| 2A.8 | REFACTOR | same focused command | exit 0 after byte-only path-component helper extraction. |
| 2A.9 | RED | `go test . -run '^TestInventoryV1OwnershipEmptyAndLegacyCoexistence$'` | exit 1: empty inventory rejected as `entries: required`. |
| 2A.10 | GREEN | same focused command | exit 0; verified-empty, zero/failed value, ownership, and legacy smoke assertions passed. |
| 2A.11 | TRIANGULATE | same focused command | exit 0 after caller/egress slice, bytes, nested-path, and reslice mutation assertions. |
| 2A.12 | REFACTOR | closure commands below | exit 0 after named high-bit fixture and private path-copy helper. |

### Verification

- `go test . -run '^(TestNewInventoryV1CanonicalContract|TestNewInventoryV1ExactContentAndPaths|TestInventoryV1OwnershipEmptyAndLegacyCoexistence)$'` — exit 0.
- `go test ./...` — exit 0 (root, CLI, internal Git, inventory, and publication packages).
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` — exit 0.
- Runtime harness: N/A. This is a pure root-package value constructor with no executable, filesystem, network, Git, CLI, or adapter boundary.
- Race suite: N/A. The implementation creates no goroutines and has no synchronization or concurrent-execution dependency.

### Files and Scope

- Added: `inventory_v1.go`, `inventory_v1_build.go`, `inventory_v1_contract_test.go`, `inventory_v1_ownership_test.go`.
- Updated: `tasks.md` (2A.1–2A.12 only) and this progress record.
- Legacy/policy/acquisition surfaces remain unmodified; no decoder (`DecodeInventoryV1`) exists.
- Design deviations: none. The strict-TDD support file `.pi/gentle-ai/support/strict-tdd.md` was absent; the configured RED → GREEN → TRIANGULATE → REFACTOR contract was followed directly.

### Workload and Rollback

- PR boundary: approved stacked-to-main work unit `inventory-v1-2a-canonical-construction` / issue #38 only. Work unit 2B remains unchecked.
- Budget at closure: production 37 + 106 = 143; tests 119 + 45 = 164; artifacts 12 checkbox replacements (12 additions + 12 deletions) plus 47 progress additions = 71; total changed-line accounting = 378, under 400. Each Go file is at most 119 changed lines and four root Go files changed.
- Rollback boundary: delete only the four new `inventory_v1*` files. This removes 2A without changing Policy V1, legacy inventory/acquisition, or consumer integration.

### Remaining

- All Work Unit 2B decoder tasks remain unchecked and out of scope for this apply.

## Work Unit 2B — Complete (ordinary disabled/unmanaged sync)

This section is authoritative for the current Inventory V1 state; the preceding 2A-close sections preserve their historical chronology.

The strict decoder is merged. `DecodeInventoryV1` accepts only constructor-equivalent canonical bytes linked to the supplied exact Policy V1 document; it rejects malformed, ambiguous, authority-bearing, noncanonical, wrong-link, duplicate-path, and invalid-value inputs without returning a partial document. 2A behavior and legacy/acquisition surfaces remain unchanged.

### Clean-redo/correction TDD chronology

- Original 2B worktree commit `2865a2f39671b6c868a8d38e9e75e7fb7e390e43` was clean-redone as merge commit `9f9eef30ac9b6e983468109c9910ebab45706e38` (PR #45); `git log --cherry` reports the final patch as equivalent.
- The merged commit does not preserve independently rerunnable historical pre-implementation RED output. This ordinary documentation sync makes no Go behavior change, so it does not claim a recreated RED run. The merged decoder tests cover the named RED/GREEN/TRIANGULATE/REFACTOR task classes, and fresh focused verification below is the final GREEN evidence.

### Merged evidence and final verification

- PR #42 / `f61e2680b4567f33587a8fcbc573ae570c3de253`: proposal and design.
- PR #43 / `6d3795065c23a9ec54f7871208d96e6c6a85adda`: active domain specs and tasks.
- PR #44 / `4f033d7fe426800fec946908bad7e54408477473`: 2A construction contract.
- PR #45 / `9f9eef30ac9b6e983468109c9910ebab45706e38`: 2B decoder; final budget 365/400 additions + deletions against its parent, three root Go files, and maximum 146/160 changed lines per Go file.
- Focused 2A+2B: `go test -count=1 . -run '^(TestNewInventoryV1CanonicalContract|TestNewInventoryV1ExactContentAndPaths|TestInventoryV1OwnershipEmptyAndLegacyCoexistence|TestDecodeInventoryV1CanonicalAndClosedShape|TestDecodeInventoryV1ValueValidation|TestDecodeInventoryV1OrderLinkAndCanonicality)$'` — PASS.
- Policy: `go test -count=1 . -run '^(TestAccountingPolicyV1Contract|TestPolicyDecimal|TestPolicyDecimalCompare|TestDecodeAccountingPolicyV1FailsClosed|TestPolicyGlob|TestPolicyGlobByteIdentity|TestPolicyObject)$'` — PASS.
- Full: `go test -count=1 ./...` — PASS.
- Gofmt: `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` — PASS. Whitespace: `git diff --check` — PASS.
- Runtime harness: N/A; Inventory V1 is a pure value/decoder contract with no executable, filesystem, network, Git, CLI, or adapter boundary. Race suite: N/A; it creates no goroutines and has no synchronization or concurrent-execution dependency.

### Rollback and status

- 2B rollback boundary: remove only `inventory_v1_decode.go`, `inventory_v1_decode_test.go`, and `inventory_v1_decode_strict_test.go`; 2A remains a valid construction-only contract.
- Status: `disabled/unmanaged`. The native successor attempt for #4040 remains unresolved. This record claims no native settle, native verification, or native approval.
