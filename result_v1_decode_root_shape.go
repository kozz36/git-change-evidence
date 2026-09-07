package changeevidence

func validateResultRootShape(raw []byte) error {
	root, err := resultShapeObject(raw, "document", "schema", "accounting_policy_sha256", "inventory_sha256", "revisions", "entries", "totals", "observations")
	if err != nil {
		return err
	}
	if err := inventoryRequired(root, "document", "schema", "accounting_policy_sha256", "inventory_sha256", "revisions", "entries", "totals", "observations"); err != nil {
		return err
	}
	for _, name := range []string{"schema", "accounting_policy_sha256", "inventory_sha256"} {
		if !policyString(root[name]) {
			return &ContractError{name, "invalid"}
		}
	}
	if err := validateResultRevisionsShape(root["revisions"], "revisions"); err != nil {
		return err
	}
	if err := validateResultEntryShapes(root["entries"]); err != nil {
		return err
	}
	if err := validateResultTotalsShapes(root["totals"]); err != nil {
		return err
	}
	return validateResultObservationShapes(root["observations"])
}
