package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

func validateResultAntecedents(policy PolicyDocumentV1, inventory InventoryDocumentV1) (PolicyView, DocumentDigest, DocumentDigest, error) {
	policyRaw := policy.CanonicalBytes()
	decodedPolicy, err := DecodeAccountingPolicyV1(policyRaw)
	if len(policyRaw) == 0 || err != nil || !bytes.Equal(decodedPolicy.CanonicalBytes(), policyRaw) {
		return PolicyView{}, "", "", resultValidationFailure("accounting_policy")
	}
	policyDigest := resultDocumentDigest(policyRaw)
	if policy.Digest() != policyDigest || decodedPolicy.Digest() != policyDigest {
		return PolicyView{}, "", "", resultValidationFailure("accounting_policy")
	}

	inventoryRaw := inventory.CanonicalBytes()
	decodedInventory, err := DecodeInventoryV1(inventoryRaw, decodedPolicy)
	if len(inventoryRaw) == 0 || err != nil || !bytes.Equal(decodedInventory.CanonicalBytes(), inventoryRaw) {
		return PolicyView{}, "", "", resultValidationFailure("inventory")
	}
	inventoryDigest := resultDocumentDigest(inventoryRaw)
	if inventory.Digest() != inventoryDigest || decodedInventory.Digest() != inventoryDigest || inventory.AccountingPolicyDigest() != policyDigest || decodedInventory.AccountingPolicyDigest() != policyDigest {
		return PolicyView{}, "", "", resultValidationFailure("inventory")
	}
	return decodedPolicy.View(), policyDigest, inventoryDigest, nil
}

func validateResultSnapshot(snapshot CommittedSnapshot) (RevisionIdentity, []accountingEntry, error) {
	base, head := string(snapshot.Base), string(snapshot.Head)
	if !isRevision(base) {
		return RevisionIdentity{}, nil, &ContractError{"snapshot.base", "invalid_revision"}
	}
	if !isRevision(head) || len(base) != len(head) || base == head && len(snapshot.Entries()) != 0 {
		return RevisionIdentity{}, nil, &ContractError{"snapshot.head", "invalid_revision"}
	}

	changes := snapshot.Entries()
	entries := make([]accountingEntry, len(changes))
	paths := make(map[string]struct{}, len(changes))
	for index, change := range changes {
		path := []byte(change.Path)
		if !validInventoryPath(path) {
			return RevisionIdentity{}, nil, &ContractError{"entries.path_b64", "invalid_path"}
		}
		key := string(path)
		if _, duplicate := paths[key]; duplicate {
			return RevisionIdentity{}, nil, &ContractError{"entries.path_b64", "duplicate_path"}
		}
		paths[key] = struct{}{}
		entries[index] = accountingEntry{path: append([]byte(nil), path...), lines: change.Lines}
	}
	return RevisionIdentity{Base: base, Head: head}, entries, nil
}

func resultDocumentDigest(raw []byte) DocumentDigest {
	sum := sha256.Sum256(raw)
	return DocumentDigest(hex.EncodeToString(sum[:]))
}

func resultValidationFailure(field string) error { return &ContractError{field, "invalid_document"} }
