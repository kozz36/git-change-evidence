package changeevidence

import "encoding/json"

func resultShapeObject(raw []byte, field string, allowed ...string) (map[string]json.RawMessage, error) {
	object, err := decodePolicyObject(raw, field, allowed...)
	if contract, ok := err.(*ContractError); ok && contract.Code == "unknown_field" && resultShapeAuthority(contract.Field) {
		return nil, &ContractError{contract.Field, "authority_field"}
	}
	return object, err
}

func resultShapeAuthority(field string) bool {
	return inventoryAuthority(field) || isAuthorityField(field)
}
