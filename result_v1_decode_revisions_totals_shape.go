package changeevidence

import (
	"encoding/json"
	"fmt"
)

func validateResultRevisionsShape(raw json.RawMessage, field string) error {
	revisions, err := resultShapeObject(raw, field, "base", "head")
	if err != nil {
		return err
	}
	if err := inventoryRequired(revisions, field, "base", "head"); err != nil {
		return err
	}
	for _, name := range []string{"base", "head"} {
		if !policyString(revisions[name]) {
			return &ContractError{field + "." + name, "invalid"}
		}
	}
	return nil
}

func validateResultTotalsShapes(raw []byte) error {
	totals, ok := inventoryArray(raw)
	if !ok {
		return &ContractError{"totals", "invalid_array"}
	}
	for index, total := range totals {
		if err := validateResultTotalShape(total, resultTotalField(index)); err != nil {
			return err
		}
	}
	return nil
}

func validateResultTotalShape(raw json.RawMessage, field string) error {
	total, err := resultShapeObject(raw, field, "category", "additions", "deletions", "non_countable")
	if err != nil {
		return err
	}
	if err := inventoryRequired(total, field, "category", "additions", "deletions", "non_countable"); err != nil {
		return err
	}
	if !policyString(total["category"]) {
		return &ContractError{field + ".category", "invalid"}
	}
	for _, name := range []string{"additions", "deletions", "non_countable"} {
		if !policyUint(total[name]) {
			return &ContractError{field + "." + name, "invalid_uint"}
		}
	}
	return nil
}

func resultTotalField(index int) string { return fmt.Sprintf("totals[%d]", index) }
