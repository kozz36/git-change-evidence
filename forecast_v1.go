package changeevidence

// ForecastMetricKindV1 is closed to the existing category-total axes.
type ForecastMetricKindV1 string

const (
	ForecastAdditionsV1    ForecastMetricKindV1 = "additions"     // lines
	ForecastDeletionsV1    ForecastMetricKindV1 = "deletions"     // lines
	ForecastNonCountableV1 ForecastMetricKindV1 = "non_countable" // entries
)

type ForecastMetricKeyV1 struct {
	Category string
	Kind     ForecastMetricKindV1
}

// ForecastMetricV1 retains nil as explicit unknown; a pointer to zero is known.
type ForecastMetricV1 struct {
	Category string
	Kind     ForecastMetricKindV1
	Value    *uint64
}

type ForecastInputV1 struct {
	AccountingPolicyDigest DocumentDigest
	InventoryDigest        DocumentDigest
	Revisions              RevisionIdentity
	Metrics                []ForecastMetricV1
}

// ForecastV1 owns validated supplied input, without completing missing metrics.
type ForecastV1 struct{ input ForecastInputV1 }

func (f ForecastV1) AccountingPolicyDigest() DocumentDigest { return f.input.AccountingPolicyDigest }
func (f ForecastV1) InventoryDigest() DocumentDigest        { return f.input.InventoryDigest }
func (f ForecastV1) Revisions() RevisionIdentity            { return f.input.Revisions }
func (f ForecastV1) Metrics() []ForecastMetricV1            { return cloneForecastMetricsV1(f.input.Metrics) }

func cloneForecastMetricsV1(metrics []ForecastMetricV1) []ForecastMetricV1 {
	if metrics == nil {
		return nil
	}
	cloned := make([]ForecastMetricV1, len(metrics))
	for i, metric := range metrics {
		cloned[i] = metric
		if metric.Value != nil {
			value := *metric.Value
			cloned[i].Value = &value
		}
	}
	return cloned
}

func IngestForecastV1(input ForecastInputV1) (ForecastV1, error) {
	for _, binding := range []struct {
		field  string
		digest DocumentDigest
	}{
		{"accounting_policy_sha256", input.AccountingPolicyDigest},
		{"inventory_sha256", input.InventoryDigest},
	} {
		if !isDigest(string(binding.digest)) {
			return ForecastV1{}, &ContractError{binding.field, "invalid_digest"}
		}
	}
	if !isRevision(input.Revisions.Base) {
		return ForecastV1{}, &ContractError{"revisions.base", "invalid_revision"}
	}
	if !isRevision(input.Revisions.Head) || len(input.Revisions.Base) != len(input.Revisions.Head) {
		return ForecastV1{}, &ContractError{"revisions.head", "invalid_revision"}
	}
	seen := make(map[ForecastMetricKeyV1]bool, len(input.Metrics))
	for _, metric := range input.Metrics {
		if !validPolicyText(metric.Category) {
			return ForecastV1{}, &ContractError{"metrics.category", "invalid"}
		}
		switch metric.Kind {
		case ForecastAdditionsV1, ForecastDeletionsV1, ForecastNonCountableV1:
		default:
			return ForecastV1{}, &ContractError{"metrics.kind", "invalid"}
		}
		key := ForecastMetricKeyV1{metric.Category, metric.Kind}
		if seen[key] {
			return ForecastV1{}, &ContractError{"metrics", "duplicate"}
		}
		seen[key] = true
	}
	input.Metrics = cloneForecastMetricsV1(input.Metrics)
	return ForecastV1{input: input}, nil
}
