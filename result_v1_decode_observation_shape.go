package changeevidence

import (
	"encoding/json"
	"fmt"
)

func validateResultObservationShapes(raw []byte) error {
	observations, ok := inventoryArray(raw)
	if !ok {
		return &ContractError{"observations", "invalid_array"}
	}
	for index, observation := range observations {
		if err := validateResultObservationShape(observation, resultObservationField(index)); err != nil {
			return err
		}
	}
	return nil
}

func resultObservationField(index int) string { return fmt.Sprintf("observations[%d]", index) }

func validateResultObservationShape(raw json.RawMessage, field string) error {
	observation, err := resultShapeObject(raw, field, "kind", "category", "reference", "limit", "available", "actual", "numerator", "denominator", "exceeded")
	if err != nil {
		return err
	}
	if err := inventoryRequired(observation, field, "kind", "available"); err != nil {
		return err
	}
	if !policyString(observation["kind"]) {
		return &ContractError{field + ".kind", "invalid"}
	}
	if !resultShapeBool(observation["available"]) {
		return &ContractError{field + ".available", "invalid"}
	}
	var kind string
	var available bool
	_ = json.Unmarshal(observation["kind"], &kind)
	_ = json.Unmarshal(observation["available"], &available)
	switch ResultObservationKind(kind) {
	case ResultLineThresholdObservationV1:
		return validateResultLineThresholdObservation(observation, field, available)
	case ResultRatioObservationV1:
		return validateResultRatioObservation(observation, field, available)
	default:
		if resultShapeAuthority(kind) {
			return &ContractError{field + ".kind", "authority_field"}
		}
		return &ContractError{field + ".kind", "unsupported_kind"}
	}
}

func validateResultLineThresholdObservation(observation map[string]json.RawMessage, field string, available bool) error {
	if available {
		if err := validateResultObservationFields(observation, field, "kind", "category", "limit", "available", "actual", "exceeded"); err != nil {
			return err
		}
	} else if err := validateResultObservationFields(observation, field, "kind", "category", "limit", "available", "exceeded"); err != nil {
		return err
	}
	if !policyString(observation["category"]) {
		return &ContractError{field + ".category", "invalid"}
	}
	if !policyUint(observation["limit"]) {
		return &ContractError{field + ".limit", "invalid_uint"}
	}
	if available && !policyUint(observation["actual"]) {
		return &ContractError{field + ".actual", "invalid_uint"}
	}
	if !resultShapeBool(observation["exceeded"]) {
		return &ContractError{field + ".exceeded", "invalid"}
	}
	return nil
}

func validateResultRatioObservation(observation map[string]json.RawMessage, field string, available bool) error {
	if available {
		if err := validateResultObservationFields(observation, field, "kind", "category", "reference", "limit", "available", "numerator", "denominator", "exceeded"); err != nil {
			return err
		}
	} else if err := validateResultObservationFields(observation, field, "kind", "category", "reference", "limit", "available", "exceeded"); err != nil {
		return err
	}
	for _, name := range []string{"category", "reference", "limit"} {
		if !policyString(observation[name]) {
			return &ContractError{field + "." + name, "invalid"}
		}
	}
	for _, name := range []string{"numerator", "denominator"} {
		if available && !policyUint(observation[name]) {
			return &ContractError{field + "." + name, "invalid_uint"}
		}
	}
	if !resultShapeBool(observation["exceeded"]) {
		return &ContractError{field + ".exceeded", "invalid"}
	}
	return nil
}

func validateResultObservationFields(observation map[string]json.RawMessage, field string, names ...string) error {
	if err := inventoryRequired(observation, field, names...); err != nil {
		return err
	}
	for name := range observation {
		allowed := false
		for _, candidate := range names {
			allowed = allowed || name == candidate
		}
		if !allowed {
			return &ContractError{field + "." + name, "unknown_field"}
		}
	}
	return nil
}

func resultShapeBool(raw json.RawMessage) bool {
	var value bool
	return len(raw) > 0 && (raw[0] == 't' || raw[0] == 'f') && json.Unmarshal(raw, &value) == nil
}
