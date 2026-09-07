package changeevidence

import "bytes"

// DecodeAccountingResultV1 accepts only Result V1 evidence that exactly rebuilds
// from the supplied Policy V1 and Inventory V1 antecedents.
func DecodeAccountingResultV1(raw []byte, policy PolicyDocumentV1, inventory InventoryDocumentV1) (ResultDocumentV1, error) {
	candidate, err := materializeResultDecodeCandidate(raw)
	if err != nil {
		return ResultDocumentV1{}, err
	}
	rebuilt, err := NewAccountingResultV1(policy, inventory, candidate.snapshot)
	if err != nil {
		return ResultDocumentV1{}, err
	}
	if candidate.policyDigest != rebuilt.AccountingPolicyDigest() {
		return ResultDocumentV1{}, &ContractError{"accounting_policy_sha256", "mismatched_digest"}
	}
	if candidate.inventoryDigest != rebuilt.InventoryDigest() {
		return ResultDocumentV1{}, &ContractError{"inventory_sha256", "mismatched_digest"}
	}
	if !bytes.Equal(raw, rebuilt.CanonicalBytes()) {
		return ResultDocumentV1{}, &ContractError{"document", "noncanonical"}
	}
	return rebuilt, nil
}
