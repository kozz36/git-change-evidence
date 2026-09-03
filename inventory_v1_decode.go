package changeevidence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"unicode/utf8"
)

func DecodeInventoryV1(raw []byte, policy PolicyDocumentV1) (InventoryDocumentV1, error) {
	policyDigest, err := validatedPolicyDigest(policy)
	if err != nil {
		return InventoryDocumentV1{}, err
	}
	if !utf8.Valid(raw) {
		return InventoryDocumentV1{}, &ContractError{"document", "invalid"}
	}
	wire, err := decodeInventoryWire(raw)
	if err != nil {
		return InventoryDocumentV1{}, err
	}
	if wire.Schema != InventoryContractV1 {
		return InventoryDocumentV1{}, &ContractError{"schema", "unsupported_version"}
	}
	if !isDigest(string(wire.AccountingPolicySHA256)) {
		return InventoryDocumentV1{}, &ContractError{"accounting_policy_sha256", "invalid_digest"}
	}
	if wire.AccountingPolicySHA256 != policyDigest {
		return InventoryDocumentV1{}, &ContractError{"accounting_policy_sha256", "mismatch"}
	}
	entries, paths := make([]InventoryEntryV1, len(wire.Entries)), make(map[string]bool, len(wire.Entries))
	for n, entry := range wire.Entries {
		path, _ := base64.StdEncoding.DecodeString(entry.PathB64)
		if paths[string(path)] {
			return InventoryDocumentV1{}, &ContractError{"entries.path_b64", "duplicate_path"}
		}
		paths[string(path)] = true
		entries[n] = InventoryEntryV1{path, entry.ContentSHA256, entry.ByteLength}
		if n > 0 && bytes.Compare(entries[n-1].Path, path) > 0 {
			return InventoryDocumentV1{}, &ContractError{"entries", "noncanonical_order"}
		}
	}
	document, err := inventoryDocument(wire.AccountingPolicySHA256, entries)
	if err != nil {
		return InventoryDocumentV1{}, err
	}
	if !bytes.Equal(raw, document.CanonicalBytes()) {
		return InventoryDocumentV1{}, &ContractError{"document", "noncanonical"}
	}
	return document, nil
}

func decodeInventoryWire(raw []byte) (inventoryWire, error) {
	root, err := inventoryObject(raw, "document", "schema", "accounting_policy_sha256", "entries")
	if err != nil {
		return inventoryWire{}, err
	}
	if err := inventoryRequired(root, "document", "schema", "accounting_policy_sha256", "entries"); err != nil {
		return inventoryWire{}, err
	}
	if !policyString(root["accounting_policy_sha256"]) {
		return inventoryWire{}, &ContractError{"accounting_policy_sha256", "invalid_digest"}
	}
	entries, ok := inventoryArray(root["entries"])
	if !ok {
		return inventoryWire{}, &ContractError{"entries", "invalid_array"}
	}
	for _, rawEntry := range entries {
		entry, err := inventoryObject(rawEntry, "entries", "path_b64", "content_sha256", "byte_length")
		if err != nil {
			return inventoryWire{}, err
		}
		if err := inventoryRequired(entry, "entries", "path_b64", "content_sha256", "byte_length"); err != nil {
			return inventoryWire{}, err
		}
		if err := validateInventoryEntry(entry); err != nil {
			return inventoryWire{}, err
		}
	}
	var wire inventoryWire
	if err := json.Unmarshal(raw, &wire); err != nil {
		return inventoryWire{}, &ContractError{"document", "invalid"}
	}
	return wire, nil
}

func validateInventoryEntry(entry map[string]json.RawMessage) error {
	var encoded string
	if !policyString(entry["path_b64"]) || json.Unmarshal(entry["path_b64"], &encoded) != nil {
		return &ContractError{"entries.path_b64", "invalid"}
	}
	path, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || base64.StdEncoding.EncodeToString(path) != encoded {
		return &ContractError{"entries.path_b64", "invalid"}
	}
	if !validInventoryPath(path) {
		return &ContractError{"entries.path_b64", "invalid_path"}
	}
	var digest string
	if !policyString(entry["content_sha256"]) || json.Unmarshal(entry["content_sha256"], &digest) != nil || !isDigest(digest) {
		return &ContractError{"entries.content_sha256", "invalid_digest"}
	}
	if !inventoryUint(entry["byte_length"]) {
		return &ContractError{"entries.byte_length", "invalid_uint"}
	}
	return nil
}

func inventoryObject(raw []byte, field string, allowed ...string) (map[string]json.RawMessage, error) {
	object, err := decodePolicyObject(raw, field, allowed...)
	if contract, ok := err.(*ContractError); ok && contract.Code == "unknown_field" && inventoryAuthority(contract.Field) {
		return nil, &ContractError{contract.Field, "authority_field"}
	}
	return object, err
}

func inventoryAuthority(field string) bool {
	normalized := strings.NewReplacer("_", "", "-", "", ".", "").Replace(strings.ToLower(field))
	for _, term := range []string{"approval", "approve", "rejection", "reject", "blocking", "block", "gating", "gate", "merge", "deploy", "release", "delivery", "deliver"} {
		if strings.Contains(normalized, term) {
			return true
		}
	}
	return false
}

func inventoryRequired(object map[string]json.RawMessage, field string, names ...string) error {
	for _, name := range names {
		if _, ok := object[name]; !ok {
			return &ContractError{field + "." + name, "required"}
		}
	}
	return nil
}

func inventoryArray(raw json.RawMessage) ([]json.RawMessage, bool) {
	var values []json.RawMessage
	ok := len(raw) > 0 && raw[0] == '[' && json.Unmarshal(raw, &values) == nil
	return values, ok
}

func inventoryUint(raw json.RawMessage) bool {
	var value uint64
	return len(raw) > 0 && raw[0] >= '0' && raw[0] <= '9' && json.Unmarshal(raw, &value) == nil
}
