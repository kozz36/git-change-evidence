# Apply Progress: Inventory V1

## Work Unit 2A — Complete

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
