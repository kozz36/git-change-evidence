package changeevidence

import "sort"

// ForecastThresholdV1 is a profile input, not part of Policy V1.
type ForecastThresholdV1 struct {
	Key           ForecastMetricKeyV1
	AbsoluteLimit uint64
}

type ForecastComparisonStateV1 string

const (
	ForecastAbsentV1       ForecastComparisonStateV1 = "absent"
	ForecastCompatibleV1   ForecastComparisonStateV1 = "compatible"
	ForecastIncompatibleV1 ForecastComparisonStateV1 = "incompatible"
)

// ForecastDeltaV1 represents measured minus supplied without signed conversion.
type ForecastDeltaV1 struct {
	Direction string // lower, equal, or higher
	Magnitude uint64
}

// Nil values are unavailable/unknown; ForecastPresent distinguishes omission.
// Delta and Exceeded exist only when a coherent numeric comparison is possible.
type ForecastComparisonObservationV1 struct {
	Key                ForecastMetricKeyV1
	ForecastPresent    bool
	Measured, Forecast *uint64
	Delta              *ForecastDeltaV1
	AbsoluteLimit      *uint64
	Exceeded           *bool
}

// ForecastComparisonV1 owns technical evidence only; it conveys no authority.
type ForecastComparisonV1 struct {
	state        ForecastComparisonStateV1
	forecast     *ForecastV1
	observations []ForecastComparisonObservationV1
}

func (c ForecastComparisonV1) State() ForecastComparisonStateV1 { return c.state }

func (c ForecastComparisonV1) Forecast() *ForecastV1 {
	if c.forecast == nil {
		return nil
	}
	copy := *c.forecast
	copy.input.Metrics = cloneForecastMetricsV1(copy.input.Metrics)
	return &copy
}

func (c ForecastComparisonV1) Observations() []ForecastComparisonObservationV1 {
	if c.observations == nil {
		return nil
	}
	copy := append([]ForecastComparisonObservationV1{}, c.observations...)
	for i := range copy {
		copy[i].Measured = forecastComparisonCopy(copy[i].Measured)
		copy[i].Forecast = forecastComparisonCopy(copy[i].Forecast)
		copy[i].Delta = forecastComparisonCopy(copy[i].Delta)
		copy[i].AbsoluteLimit = forecastComparisonCopy(copy[i].AbsoluteLimit)
		copy[i].Exceeded = forecastComparisonCopy(copy[i].Exceeded)
	}
	return copy
}

func forecastComparisonCopy[T any](value *T) *T {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func CompareForecastV1(result ResultDocumentV1, forecast *ForecastV1, thresholds []ForecastThresholdV1) (ForecastComparisonV1, error) {
	// Public Result V1 values are immutable constructor output; reject the zero document.
	if raw := result.CanonicalBytes(); len(raw) == 0 || resultDocumentDigest(raw) != result.Digest() {
		return ForecastComparisonV1{}, resultValidationFailure("result")
	}
	comparison := ForecastComparisonV1{state: ForecastAbsentV1}
	if forecast != nil {
		retained, err := IngestForecastV1(forecast.input)
		if err != nil {
			return ForecastComparisonV1{}, err
		}
		comparison.state, comparison.forecast = ForecastCompatibleV1, &retained
		if retained.AccountingPolicyDigest() != result.AccountingPolicyDigest() || retained.InventoryDigest() != result.InventoryDigest() || retained.Revisions() != result.Revisions() {
			comparison.state = ForecastIncompatibleV1
		}
	}
	limits := make(map[ForecastMetricKeyV1]uint64, len(thresholds))
	for _, threshold := range thresholds {
		if !validPolicyText(threshold.Key.Category) {
			return ForecastComparisonV1{}, &ContractError{"thresholds.category", "invalid"}
		}
		if threshold.Key.Kind != ForecastAdditionsV1 && threshold.Key.Kind != ForecastDeletionsV1 && threshold.Key.Kind != ForecastNonCountableV1 {
			return ForecastComparisonV1{}, &ContractError{"thresholds.kind", "invalid"}
		}
		if _, duplicate := limits[threshold.Key]; duplicate {
			return ForecastComparisonV1{}, &ContractError{"thresholds", "duplicate"}
		}
		limits[threshold.Key] = threshold.AbsoluteLimit
	}
	byKey := make(map[ForecastMetricKeyV1]ForecastComparisonObservationV1)
	for _, total := range result.Totals() {
		for _, metric := range []ForecastMetricV1{{total.Category, ForecastAdditionsV1, &total.Additions}, {total.Category, ForecastDeletionsV1, &total.Deletions}, {total.Category, ForecastNonCountableV1, &total.NonCountable}} {
			key := ForecastMetricKeyV1{metric.Category, metric.Kind}
			observation := ForecastComparisonObservationV1{Key: key}
			if metric.Kind == ForecastNonCountableV1 || total.NonCountable == 0 {
				observation.Measured = metric.Value
			}
			byKey[key] = observation
		}
	}
	if comparison.forecast != nil {
		for _, metric := range comparison.forecast.Metrics() {
			key := ForecastMetricKeyV1{metric.Category, metric.Kind}
			observation := byKey[key]
			observation.Key, observation.ForecastPresent, observation.Forecast = key, true, metric.Value
			byKey[key] = observation
		}
	}
	comparison.observations = make([]ForecastComparisonObservationV1, 0, len(byKey))
	for key, observation := range byKey {
		if limit, present := limits[key]; present {
			observation.AbsoluteLimit = &limit
		}
		if comparison.state == ForecastCompatibleV1 && observation.Measured != nil && observation.Forecast != nil {
			delta := forecastComparisonDelta(*observation.Measured, *observation.Forecast)
			observation.Delta = &delta
			if observation.AbsoluteLimit != nil {
				exceeded := delta.Magnitude > *observation.AbsoluteLimit
				observation.Exceeded = &exceeded
			}
		}
		comparison.observations = append(comparison.observations, observation)
	}
	sort.Slice(comparison.observations, func(i, j int) bool {
		a, b := comparison.observations[i].Key, comparison.observations[j].Key
		return a.Category < b.Category || a.Category == b.Category && a.Kind < b.Kind
	})
	return comparison, nil
}

func forecastComparisonDelta(measured, forecast uint64) ForecastDeltaV1 {
	if measured < forecast {
		return ForecastDeltaV1{"lower", forecast - measured}
	}
	if measured > forecast {
		return ForecastDeltaV1{"higher", measured - forecast}
	}
	return ForecastDeltaV1{"equal", 0}
}
