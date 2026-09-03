# Sync Report: Inventory V1

## Status

**`disabled/unmanaged` ordinary sync.** The active change remains in place; no archive or native status action occurred.

## Canonical domains

| Canonical domain | Active source | Copied requirement content | Byte-identity proof |
|---|---|---|---|
| `inventory-v1-content-identity` | `openspec/changes/inventory-v1/specs/inventory-v1-content-identity/spec.md` | All four requirements copied unchanged. | `cmp -s` PASS |
| `inventory-v1-contract` | `openspec/changes/inventory-v1/specs/inventory-v1-contract/spec.md` | All five requirements copied unchanged. | `cmp -s` PASS |

## Same-domain active-change census

The only active-change source specs for these two Inventory V1 domains are the two paths listed above; each has one canonical counterpart. The byte comparison was performed source-to-counterpart with `cmp -s` and both pairs are identical.

## Budget and next boundary

Projected archive count after exact historical restoration: 318 additions + 26 deletions = **344/400** changed lines. A byte-preserving archive move adds no authored lines, retaining **344/400**; any archive report must keep the complete closure within 400. Archive only after this report, `verify-report.md`, and the zero-unchecked-task census are revalidated.
