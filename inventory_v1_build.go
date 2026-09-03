package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"sort"
)

type inventoryWire struct {
	Schema                 ContractVersion      `json:"schema"`
	AccountingPolicySHA256 DocumentDigest       `json:"accounting_policy_sha256"`
	Entries                []inventoryWireEntry `json:"entries"`
}
type inventoryWireEntry struct {
	PathB64       string `json:"path_b64"`
	ContentSHA256 Digest `json:"content_sha256"`
	ByteLength    uint64 `json:"byte_length"`
}

func buildInventoryV1(policy PolicyDocumentV1, input []InventoryEntryInput) (InventoryDocumentV1, error) {
	policyDigest, err := validatedPolicyDigest(policy)
	if err != nil {
		return InventoryDocumentV1{}, err
	}
	entries := make([]InventoryEntryV1, len(input))
	for n, entry := range input {
		derived, err := deriveInventoryEntry(entry)
		if err != nil {
			return InventoryDocumentV1{}, err
		}
		entries[n] = derived
	}
	sort.Slice(entries, func(left, right int) bool { return bytes.Compare(entries[left].Path, entries[right].Path) < 0 })
	for n := 1; n < len(entries); n++ {
		if bytes.Equal(entries[n-1].Path, entries[n].Path) {
			return InventoryDocumentV1{}, &ContractError{"entries.path", "duplicate_path"}
		}
	}
	return inventoryDocument(policyDigest, entries)
}

func deriveInventoryEntry(input InventoryEntryInput) (InventoryEntryV1, error) {
	path, content := copyInventoryPath(input.Path), append([]byte(nil), input.Content...)
	if !validInventoryPath(path) {
		return InventoryEntryV1{}, &ContractError{"entries.path", "invalid_path"}
	}
	sum := sha256.Sum256(content)
	return InventoryEntryV1{path, Digest(hex.EncodeToString(sum[:])), uint64(len(content))}, nil
}

func validInventoryPath(value []byte) bool {
	if len(value) == 0 || value[0] == '/' || bytes.IndexByte(value, 0) >= 0 {
		return false
	}
	for _, component := range bytes.Split(value, []byte("/")) {
		if len(component) == 0 || dotInventoryComponent(component) {
			return false
		}
	}
	return true
}

func dotInventoryComponent(value []byte) bool {
	return len(value) == 1 && value[0] == '.' || len(value) == 2 && value[0] == '.' && value[1] == '.'
}

func validatedPolicyDigest(policy PolicyDocumentV1) (DocumentDigest, error) {
	raw := policy.CanonicalBytes()
	decoded, err := DecodeAccountingPolicyV1(raw)
	if len(raw) == 0 || err != nil || !bytes.Equal(decoded.CanonicalBytes(), raw) {
		return "", &ContractError{"accounting_policy", "invalid_document"}
	}
	sum := sha256.Sum256(raw)
	digest := DocumentDigest(hex.EncodeToString(sum[:]))
	if policy.Digest() != digest || decoded.Digest() != digest {
		return "", &ContractError{"accounting_policy", "invalid_document"}
	}
	return digest, nil
}

func inventoryDocument(policyDigest DocumentDigest, entries []InventoryEntryV1) (InventoryDocumentV1, error) {
	wire := inventoryWire{InventoryContractV1, policyDigest, make([]inventoryWireEntry, len(entries))}
	for n, entry := range entries {
		wire.Entries[n] = inventoryWireEntry{base64.StdEncoding.EncodeToString(entry.Path), entry.ContentSHA256, entry.ByteLength}
	}
	canonical, err := json.Marshal(wire)
	if err != nil {
		return InventoryDocumentV1{}, err
	}
	canonical = append(canonical, '\n')
	sum := sha256.Sum256(canonical)
	return InventoryDocumentV1{append([]byte(nil), canonical...), DocumentDigest(hex.EncodeToString(sum[:])), policyDigest, cloneInventoryEntries(entries)}, nil
}

func cloneInventoryEntries(entries []InventoryEntryV1) []InventoryEntryV1 {
	result := make([]InventoryEntryV1, len(entries))
	for n, entry := range entries {
		result[n] = InventoryEntryV1{copyInventoryPath(entry.Path), entry.ContentSHA256, entry.ByteLength}
	}
	return result
}

func copyInventoryPath(value []byte) []byte { return append([]byte(nil), value...) }
