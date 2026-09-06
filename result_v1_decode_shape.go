package changeevidence

import (
	"encoding/json"
	"fmt"
)

func validateResultEntryShapes(raw []byte) error {
	entries, ok := inventoryArray(raw)
	if !ok {
		return &ContractError{"entries", "invalid_array"}
	}
	for index, rawEntry := range entries {
		if err := validateResultEntryShape(rawEntry, resultEntryField(index)); err != nil {
			return err
		}
	}
	return nil
}

func validateResultEntryShape(raw json.RawMessage, field string) error {
	entry, err := resultShapeObject(raw, field, "path_b64", "category", "measurement")
	if err != nil {
		return err
	}
	if err := inventoryRequired(entry, field, "path_b64", "category", "measurement"); err != nil {
		return err
	}
	if !policyString(entry["path_b64"]) {
		return &ContractError{field + ".path_b64", "invalid"}
	}
	if !policyString(entry["category"]) {
		return &ContractError{field + ".category", "invalid"}
	}
	return validateResultMeasurementShape(entry["measurement"], field+".measurement")
}

func validateResultMeasurementShape(raw json.RawMessage, field string) error {
	measurement, err := resultShapeObject(raw, field, "kind", "additions", "deletions")
	if err != nil {
		return err
	}
	if err := inventoryRequired(measurement, field, "kind"); err != nil {
		return err
	}
	if !policyString(measurement["kind"]) {
		return &ContractError{field + ".kind", "invalid"}
	}
	var kind string
	if err := json.Unmarshal(measurement["kind"], &kind); err != nil {
		return &ContractError{field + ".kind", "invalid"}
	}
	switch ResultEntryMeasurementKind(kind) {
	case ResultEntryCountableV1:
		return validateResultCountableMeasurement(measurement, field)
	case ResultEntryNonCountableV1:
		return validateResultNonCountableMeasurement(measurement, field)
	default:
		if resultShapeAuthority(kind) {
			return &ContractError{field + ".kind", "authority_field"}
		}
		return &ContractError{field + ".kind", "unsupported_kind"}
	}
}

func validateResultCountableMeasurement(measurement map[string]json.RawMessage, field string) error {
	if err := inventoryRequired(measurement, field, "additions", "deletions"); err != nil {
		return err
	}
	if !policyUint(measurement["additions"]) {
		return &ContractError{field + ".additions", "invalid_uint"}
	}
	if !policyUint(measurement["deletions"]) {
		return &ContractError{field + ".deletions", "invalid_uint"}
	}
	return nil
}

func validateResultNonCountableMeasurement(measurement map[string]json.RawMessage, field string) error {
	for _, name := range []string{"additions", "deletions"} {
		if _, found := measurement[name]; found {
			return &ContractError{field + "." + name, "unknown_field"}
		}
	}
	return nil
}

func resultEntryField(index int) string { return fmt.Sprintf("entries[%d]", index) }
