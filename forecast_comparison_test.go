package changeevidence

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func forecastComparisonCheck(t *testing.T, valid bool, message string) {
	t.Helper()
	if !valid {
		t.Fatal(message)
	}
}

func forecastComparisonFixture(t *testing.T) (ResultDocumentV1, ForecastInputV1) {
	t.Helper()
	policy, inventory := resultUnit5Antecedents(t, resultUnit5Policy())
	result, err := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot([]CommittedChange{
		{Path: "a/file", Lines: CommittedLineCounts{Additions: 4, Deletions: 2, Countable: true}},
		{Path: "binary", Lines: CommittedLineCounts{}},
	}))
	forecastComparisonCheck(t, err == nil, "public result constructor failed")
	return result, ForecastInputV1{AccountingPolicyDigest: result.AccountingPolicyDigest(), InventoryDigest: result.InventoryDigest(), Revisions: result.Revisions()}
}

func forecastComparisonRun(t *testing.T, r ResultDocumentV1, input ForecastInputV1, thresholds []ForecastThresholdV1) ForecastComparisonV1 {
	t.Helper()
	f, err := IngestForecastV1(input)
	forecastComparisonCheck(t, err == nil, "forecast ingest failed")
	c, err := CompareForecastV1(r, &f, thresholds)
	forecastComparisonCheck(t, err == nil, "valid comparison returned an error")
	return c
}

func TestForecastComparisonUnionAvailability(t *testing.T) {
	r, input := forecastComparisonFixture(t)
	zero, max := uint64(0), uint64(math.MaxUint64)
	input.Metrics = []ForecastMetricV1{{"unmeasured", ForecastAdditionsV1, &max}, {"raw", ForecastDeletionsV1, nil}, {"source", ForecastAdditionsV1, &zero}, {"other", ForecastAdditionsV1, &zero}, {"other", ForecastNonCountableV1, &zero}}
	c := forecastComparisonRun(t, r, input, nil)
	forecastComparisonCheck(t, c.State() == ForecastCompatibleV1 && reflect.DeepEqual(c.Forecast().Metrics(), input.Metrics), "binding or retained forecast changed")
	observations := c.Observations()
	forecastComparisonCheck(t, len(observations) == 10, "wrong union length")
	for i, o := range observations {
		if i > 0 {
			previous := observations[i-1].Key
			forecastComparisonCheck(t, previous.Category < o.Key.Category || previous.Category == o.Key.Category && previous.Kind < o.Key.Kind, "nonunique or unordered union")
		}
		t.Run(o.Key.Category+"/"+string(o.Key.Kind), func(t *testing.T) {
			forecastComparisonCheck(t, o.AbsoluteLimit == nil && o.Exceeded == nil, "invented threshold")
			switch {
			case o.Key.Category == "unmeasured":
				forecastComparisonCheck(t, o.Measured == nil && o.Forecast != nil && *o.Forecast == max && o.Delta == nil, "missing measurement became zero")
			case o.Key.Category == "other" && o.Key.Kind != ForecastNonCountableV1:
				forecastComparisonCheck(t, o.Measured == nil && o.Delta == nil, "partial lines compared")
			case o.Key.Category == "raw" && o.Key.Kind == ForecastDeletionsV1:
				forecastComparisonCheck(t, o.ForecastPresent && o.Forecast == nil && o.Delta == nil && o.Measured != nil && *o.Measured == 2, "unknown became missing or zero")
			case o.ForecastPresent:
				want := ForecastDeltaV1{"equal", 0}
				measured := uint64(0)
				if o.Key.Category == "other" {
					want, measured = ForecastDeltaV1{"higher", 1}, 1
				}
				forecastComparisonCheck(t, o.Measured != nil && *o.Measured == measured && o.Delta != nil && *o.Delta == want, "known zero or noncountable count changed")
			default:
				forecastComparisonCheck(t, o.Forecast == nil && o.Delta == nil, "omission filled")
			}
		})
	}
	thresholds := []ForecastThresholdV1{}
	for _, key := range []ForecastMetricKeyV1{{"unmeasured", ForecastAdditionsV1}, {"raw", ForecastDeletionsV1}, {"raw", ForecastAdditionsV1}, {"other", ForecastAdditionsV1}} {
		thresholds = append(thresholds, ForecastThresholdV1{key, 0})
	}
	c = forecastComparisonRun(t, r, input, thresholds)
	found := 0
	for _, o := range c.Observations() {
		if o.AbsoluteLimit != nil {
			found++
			forecastComparisonCheck(t, *o.AbsoluteLimit == 0 && o.Delta == nil && o.Exceeded == nil, "unavailable threshold computed a verdict")
		}
	}
	forecastComparisonCheck(t, found == 4, "unavailable thresholds missing")
}

func TestForecastComparisonDeltaThreshold(t *testing.T) {
	r, input := forecastComparisonFixture(t)
	key := ForecastMetricKeyV1{"raw", ForecastAdditionsV1}
	for _, tt := range []struct {
		name                    string
		value, limit, magnitude uint64
		direction               string
		exceeded                bool
	}{
		{"higher below", 1, 2, 3, "higher", true}, {"higher equal", 1, 3, 3, "higher", false}, {"higher above", 1, 4, 3, "higher", false},
		{"lower below", 7, 2, 3, "lower", true}, {"lower equal", 7, 3, 3, "lower", false}, {"lower above", 7, 4, 3, "lower", false},
		{"equal zero", 4, 0, 0, "equal", false}, {"higher zero", 0, 0, 4, "higher", true}, {"maximum supplied", math.MaxUint64, 0, math.MaxUint64 - 4, "lower", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input.Metrics = []ForecastMetricV1{{key.Category, key.Kind, &tt.value}}
			c := forecastComparisonRun(t, r, input, []ForecastThresholdV1{{key, tt.limit}})
			matches := []ForecastComparisonObservationV1{}
			for _, o := range c.Observations() {
				if o.Key == key {
					matches = append(matches, o)
				}
			}
			forecastComparisonCheck(t, len(matches) == 1, "target key must occur exactly once")
			o := matches[0]
			forecastComparisonCheck(t, o.Delta != nil && *o.Delta == (ForecastDeltaV1{tt.direction, tt.magnitude}) && o.Exceeded != nil && *o.Exceeded == tt.exceeded && o.AbsoluteLimit != nil && *o.AbsoluteLimit == tt.limit, "wrong delta or threshold")
		})
	}
	for _, pair := range [][2]uint64{{0, math.MaxUint64}, {math.MaxUint64, 0}} {
		direction := "lower"
		if pair[0] > pair[1] {
			direction = "higher"
		}
		forecastComparisonCheck(t, forecastComparisonDelta(pair[0], pair[1]) == (ForecastDeltaV1{direction, math.MaxUint64}), "pure maximum delta failed")
	}
}

func TestForecastComparisonAbsenceBindingOwnership(t *testing.T) {
	r, input := forecastComparisonFixture(t)
	c, err := CompareForecastV1(r, nil, nil)
	forecastComparisonCheck(t, err == nil && c.State() == ForecastAbsentV1 && c.Forecast() == nil, "absent forecast changed")
	for _, metrics := range [][]ForecastMetricV1{nil, {}} {
		input.Metrics = metrics
		c = forecastComparisonRun(t, r, input, nil)
		forecastComparisonCheck(t, c.State() == ForecastCompatibleV1 && reflect.DeepEqual(c.Forecast().Metrics(), metrics), "empty supplied forecast lost")
	}
	for _, field := range []string{"policy", "inventory", "base", "head"} {
		t.Run(field, func(t *testing.T) {
			changed := input
			switch field {
			case "policy":
				changed.AccountingPolicyDigest = DocumentDigest(strings.Repeat("c", 64))
			case "inventory":
				changed.InventoryDigest = DocumentDigest(strings.Repeat("c", 64))
			case "base":
				changed.Revisions.Base = strings.Repeat("c", 40)
			case "head":
				changed.Revisions.Head = strings.Repeat("c", 40)
			}
			value := uint64(0)
			changed.Metrics = []ForecastMetricV1{{"raw", ForecastAdditionsV1, &value}}
			c := forecastComparisonRun(t, r, changed, []ForecastThresholdV1{{ForecastMetricKeyV1{"raw", ForecastAdditionsV1}, 0}})
			forecastComparisonCheck(t, c.State() == ForecastIncompatibleV1 && reflect.DeepEqual(c.Forecast().Metrics(), changed.Metrics), "mismatch not retained")
			for _, o := range c.Observations() {
				forecastComparisonCheck(t, o.Delta == nil && o.Exceeded == nil, "incompatible comparison computed")
			}
		})
	}
	value := uint64(1)
	input.Metrics = []ForecastMetricV1{{"raw", ForecastAdditionsV1, &value}}
	f, err := IngestForecastV1(input)
	forecastComparisonCheck(t, err == nil, "ingest failed")
	thresholds := []ForecastThresholdV1{{ForecastMetricKeyV1{"raw", ForecastAdditionsV1}, 0}}
	c, err = CompareForecastV1(r, &f, thresholds)
	forecastComparisonCheck(t, err == nil, "comparison failed")
	before := c.Observations()
	f, value, thresholds[0] = ForecastV1{}, 99, ForecastThresholdV1{ForecastMetricKeyV1{"changed", ForecastDeletionsV1}, 99}
	returned := c.Forecast()
	*returned.input.Metrics[0].Value, returned.input.Metrics[0].Category = 99, "changed"
	for i, o := range c.Observations() {
		if o.Measured != nil {
			*o.Measured = 99
		}
		if o.Forecast != nil {
			*o.Forecast = 99
		}
		if o.Delta != nil {
			o.Delta.Magnitude = 99
		}
		if o.AbsoluteLimit != nil {
			*o.AbsoluteLimit = 99
		}
		if o.Exceeded != nil {
			*o.Exceeded = !*o.Exceeded
		}
		returnedObservations := c.Observations()
		returnedObservations[i] = ForecastComparisonObservationV1{}
	}
	forecastComparisonCheck(t, reflect.DeepEqual(c.Observations(), before) && *c.Forecast().Metrics()[0].Value == 1, "comparison borrowed mutable state")
}

func TestForecastComparisonInvalidContracts(t *testing.T) {
	r, input := forecastComparisonFixture(t)
	f, err := IngestForecastV1(input)
	forecastComparisonCheck(t, err == nil, "ingest failed")
	key := ForecastMetricKeyV1{"raw", ForecastAdditionsV1}
	for _, tt := range []struct {
		name        string
		result      ResultDocumentV1
		forecast    *ForecastV1
		thresholds  []ForecastThresholdV1
		field, code string
	}{
		{"result", ResultDocumentV1{}, &f, nil, "result", "invalid_document"},
		{"forecast", r, &ForecastV1{}, nil, "accounting_policy_sha256", "invalid_digest"},
		{"category", r, &f, []ForecastThresholdV1{{ForecastMetricKeyV1{"", ForecastAdditionsV1}, 0}}, "thresholds.category", "invalid"},
		{"kind", r, &f, []ForecastThresholdV1{{ForecastMetricKeyV1{"raw", "actual"}, 0}}, "thresholds.kind", "invalid"},
		{"duplicate", r, &f, []ForecastThresholdV1{{key, 0}, {key, 1}}, "thresholds", "duplicate"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c, err := CompareForecastV1(tt.result, tt.forecast, tt.thresholds)
			forecastComparisonCheck(t, reflect.DeepEqual(err, &ContractError{tt.field, tt.code}) && reflect.DeepEqual(c, ForecastComparisonV1{}), "wrong technical error or nonzero invalid return")
		})
	}
}
