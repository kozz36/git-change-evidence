package changeevidence

import (
	"bytes"
	"encoding/json"
	"slices"
	"strings"
)

func DecodeCanonical(value []byte) (Evidence, error) {
	object, err := decodeObject(value, "document")
	if err != nil {
		return Evidence{}, err
	}
	if err := exactFields(object, "document", "kind", "provenance", "schema", "subject"); err != nil {
		return Evidence{}, err
	}
	provenance, err := decodeObject(object["provenance"], "provenance")
	if err != nil {
		return Evidence{}, err
	}
	if err := exactFields(provenance, "provenance", "accounting_sha256", "input_policy_sha256", "inventory_sha256", "revisions"); err != nil {
		return Evidence{}, err
	}
	revisions, err := decodeObject(provenance["revisions"], "provenance.revisions")
	if err != nil {
		return Evidence{}, err
	}
	if err := exactFields(revisions, "provenance.revisions", "base", "head"); err != nil {
		return Evidence{}, err
	}
	var wire canonicalReport
	if err := json.Unmarshal(value, &wire); err != nil {
		return Evidence{}, &ContractError{"document", "invalid_value"}
	}
	if wire.Schema != EvidenceContractV1 {
		return Evidence{}, &ContractError{"schema", "unsupported_version"}
	}
	if wire.Kind != ReportDocument {
		return Evidence{}, &ContractError{"kind", "unsupported_variant"}
	}
	evidence, err := NewReportV1(ReportInput{Subject: wire.Subject, Provenance: wire.Provenance})
	if err != nil {
		return Evidence{}, err
	}
	if !bytes.Equal(evidence.CanonicalBytes(), value) {
		return Evidence{}, &ContractError{"document", "noncanonical"}
	}
	return evidence, nil
}

func decodeObject(value []byte, field string) (map[string]json.RawMessage, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(value, &object); err != nil || object == nil {
		return nil, &ContractError{field, "invalid_object"}
	}
	return object, nil
}

func exactFields(object map[string]json.RawMessage, field string, allowed ...string) error {
	for _, name := range allowed {
		if _, found := object[name]; !found {
			return &ContractError{field + "." + name, "required"}
		}
	}
	for name := range object {
		if !slices.Contains(allowed, name) {
			code := "unknown_field"
			if isAuthorityField(name) {
				code = "authority_field"
			}
			return &ContractError{field + "." + name, code}
		}
	}
	return nil
}

func isAuthorityField(field string) bool {
	normalized := strings.NewReplacer("_", "", "-", "", ".", "").Replace(strings.ToLower(field))
	for _, term := range []string{"approval", "approve", "denial", "deny", "gate", "block", "merge", "deploy", "release", "deliveryauthority"} {
		if strings.Contains(normalized, term) {
			return true
		}
	}
	return false
}
