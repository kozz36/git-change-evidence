package changeevidence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"unicode/utf8"
)

func DecodeAccountingPolicyV1(raw []byte) (PolicyDocumentV1, error) {
	wire, err := decodePolicyWire(raw)
	if err != nil {
		return PolicyDocumentV1{}, err
	}
	if wire.Schema != AccountingPolicyContractV1 {
		return PolicyDocumentV1{}, policyError("schema", "unsupported_version")
	}
	input := PolicyInput{Categories: make([]CategoryInput, len(wire.Categories)), Default: wire.Default, Ratios: make([]RatioInput, len(wire.Ratios))}
	for i, category := range wire.Categories {
		globs := make([][]byte, len(category.PathGlobs))
		for j, encoded := range category.PathGlobs {
			glob, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil || base64.StdEncoding.EncodeToString(glob) != encoded {
				return PolicyDocumentV1{}, policyError("categories.path_globs_b64", "invalid")
			}
			globs[j] = glob
		}
		input.Categories[i] = CategoryInput{category.Name, globs, category.LineThreshold}
	}
	for i, ratio := range wire.Ratios {
		input.Ratios[i] = RatioInput{ratio.Category, ratio.Reference, ratio.Maximum}
	}
	document, err := makePolicy(input, false)
	if err != nil {
		return PolicyDocumentV1{}, err
	}
	if !bytes.Equal(raw, document.canonical) {
		return PolicyDocumentV1{}, policyError("document", "noncanonical")
	}
	return document, nil
}

func decodePolicyWire(raw []byte) (policyWire, error) {
	if !utf8.Valid(raw) {
		return policyWire{}, policyError("document", "invalid")
	}
	root, err := policyObject(raw, "document", "schema", "categories", "default", "ratios")
	if err != nil {
		return policyWire{}, err
	}
	if len(root) != 4 || !policyString(root["schema"]) || !policyString(root["default"]) {
		return policyWire{}, policyError("document", "invalid")
	}
	categories, ok := policyArray(root["categories"])
	if !ok {
		return policyWire{}, policyError("categories", "invalid_array")
	}
	for _, rawCategory := range categories {
		category, err := policyObject(rawCategory, "categories", "name", "path_globs_b64", "line_threshold")
		if err != nil {
			return policyWire{}, err
		}
		globs, ok := policyArray(category["path_globs_b64"])
		if len(category) != 3 || !ok || !policyString(category["name"]) || !policyUint(category["line_threshold"]) {
			return policyWire{}, policyError("categories", "invalid")
		}
		for _, glob := range globs {
			if !policyString(glob) {
				return policyWire{}, policyError("categories.path_globs_b64", "invalid")
			}
		}
	}
	ratios, ok := policyArray(root["ratios"])
	if !ok {
		return policyWire{}, policyError("ratios", "invalid_array")
	}
	for _, rawRatio := range ratios {
		ratio, err := policyObject(rawRatio, "ratios", "category", "reference", "maximum")
		if err != nil {
			return policyWire{}, err
		}
		if len(ratio) != 3 || !policyString(ratio["category"]) || !policyString(ratio["reference"]) || !policyString(ratio["maximum"]) {
			return policyWire{}, policyError("ratios", "invalid")
		}
	}
	var wire policyWire
	if json.Unmarshal(raw, &wire) != nil {
		return policyWire{}, policyError("document", "invalid")
	}
	return wire, nil
}

func policyObject(raw []byte, field string, allowed ...string) (map[string]json.RawMessage, error) {
	object, err := decodePolicyObject(raw, field, allowed...)
	if contract, ok := err.(*ContractError); ok && contract.Code == "unknown_field" && policyAuthority(contract.Field) {
		return nil, policyError(contract.Field, "authority_field")
	}
	return object, err
}
func policyAuthority(field string) bool {
	normalized := strings.NewReplacer("_", "", "-", "", ".", "").Replace(strings.ToLower(field))
	return isAuthorityField(field) || strings.Contains(normalized, "reject") || strings.Contains(normalized, "gating") || strings.Contains(normalized, "deliver")
}
func policyArray(raw []byte) ([]json.RawMessage, bool) {
	var values []json.RawMessage
	ok := len(raw) != 0 && raw[0] == '[' && json.Unmarshal(raw, &values) == nil
	return values, ok
}
func policyString(raw []byte) bool {
	var value string
	return len(raw) != 0 && raw[0] == '"' && json.Unmarshal(raw, &value) == nil
}
func policyUint(raw []byte) bool {
	var value uint64
	return len(raw) != 0 && raw[0] >= '0' && raw[0] <= '9' && json.Unmarshal(raw, &value) == nil
}
