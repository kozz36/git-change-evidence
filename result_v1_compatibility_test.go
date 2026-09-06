package changeevidence

import (
	"reflect"
	"strconv"
	"testing"
)

func TestResultV1AndAccountAgreeOnSharedDomainVectors(t *testing.T) {
	ordinary := PolicyInput{Categories: []CategoryInput{{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 2}, {Name: "other"}}, Default: "other", Ratios: []RatioInput{{Category: "source", Reference: "other", Maximum: "0.5"}}}
	for _, test := range []struct {
		name    string
		input   PolicyInput
		entries []CommittedChange
		want    []ResultEntryV1
		ratio   *ResultObservationV1
		fatal   bool
	}{
		{
			name: "overlap default and non-countable", input: resultUnit5Policy(),
			entries: []CommittedChange{{Path: string([]byte{0xff})}, {Path: "a0", Lines: CommittedLineCounts{Additions: 5, Countable: true}}, {Path: string(append([]byte("a/"), 0xff)), Lines: CommittedLineCounts{Additions: 4, Countable: true}}, {Path: "a", Lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}}},
			want:    []ResultEntryV1{{Path: []byte("a"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 2, Deletions: 1}}, {Path: append([]byte("a/"), 0xff), Category: "raw", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 4}}, {Path: []byte("a0"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 5}}, {Path: []byte{0xff}, Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryNonCountableV1}}},
		},
		{
			name: "ordinary threshold and ratio", input: ordinary,
			entries: []CommittedChange{{Path: "src/main.go", Lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}}, {Path: "note", Lines: CommittedLineCounts{Additions: 4, Countable: true}}},
			want:    []ResultEntryV1{{Path: []byte("note"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 4}}, {Path: []byte("src/main.go"), Category: "source", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 2, Deletions: 1}}},
			ratio:   &ResultObservationV1{RatioLimit: "0.5", Numerator: 3, Denominator: 4},
		},
		{
			name: "accumulation overflow", input: resultUnit5Policy(),
			entries: []CommittedChange{{Path: "a/x", Lines: CommittedLineCounts{Additions: ^uint64(0), Countable: true}}, {Path: "a/y", Lines: CommittedLineCounts{Additions: 1, Countable: true}}}, fatal: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			policy, inventory := resultUnit5Antecedents(t, test.input)
			result, resultErr := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot(test.entries))
			legacy, legacyErr := Account(resultUnit5Snapshot(test.entries), resultUnit5LegacyPolicy(t, test.input))
			if test.fatal {
				if resultErr == nil || legacyErr == nil {
					t.Fatalf("errors = (%v, %v)", resultErr, legacyErr)
				}
				assertZeroResultDocument(t, result)
				return
			}
			if resultErr != nil || legacyErr != nil || !reflect.DeepEqual(result.Entries(), test.want) {
				t.Fatalf("result, legacy = (%#v, %v), (%#v, %v)", result, resultErr, legacy, legacyErr)
			}
			assertResultMatchesLegacy(t, result, legacy)
			if test.ratio != nil {
				gotRatio := result.Observations()[len(result.Observations())-1]
				if gotRatio.RatioLimit != test.ratio.RatioLimit || gotRatio.Numerator != test.ratio.Numerator || gotRatio.Denominator != test.ratio.Denominator {
					t.Fatalf("ratio fields = %#v, want %#v", gotRatio, test.ratio)
				}
			}
		})
	}
}

func resultUnit5LegacyPolicy(t *testing.T, input PolicyInput) AccountingPolicy {
	t.Helper()
	policy := AccountingPolicy{Categories: make([]AccountingCategory, len(input.Categories)), Default: input.Default, Ratios: make([]CategoryRatio, len(input.Ratios))}
	for i, category := range input.Categories {
		globs := make([]string, len(category.PathGlobs))
		for j, glob := range category.PathGlobs {
			globs[j] = string(glob)
		}
		policy.Categories[i] = AccountingCategory{Name: category.Name, PathGlobs: globs, LineThreshold: int64(category.LineThreshold)}
	}
	for i, ratio := range input.Ratios {
		maximum, err := strconv.ParseFloat(ratio.Maximum, 64)
		if err != nil {
			t.Fatal(err)
		}
		policy.Ratios[i] = CategoryRatio{Category: ratio.Category, Reference: ratio.Reference, Maximum: maximum}
	}
	return policy
}

func assertResultMatchesLegacy(t *testing.T, result ResultDocumentV1, legacy AccountingResult) {
	t.Helper()
	if len(result.Totals()) != len(legacy.Totals) || len(result.Observations()) != len(legacy.Observations) {
		t.Fatalf("totals or observations differ: %#v, %#v", result, legacy)
	}
	for i, total := range result.Totals() {
		if total.Category != legacy.Totals[i].Category || total.Additions != legacy.Totals[i].Additions || total.Deletions != legacy.Totals[i].Deletions || total.NonCountable != legacy.Totals[i].NonCountable {
			t.Fatalf("total %d = %#v, want %#v", i, total, legacy.Totals[i])
		}
	}
	for i, observation := range result.Observations() {
		legacyObservation := legacy.Observations[i]
		if observation.Category != legacyObservation.Category || observation.Reference != legacyObservation.Reference || observation.Available != legacyObservation.Available || observation.Exceeded != legacyObservation.Exceeded {
			t.Fatalf("observation %d = %#v, want %#v", i, observation, legacyObservation)
		}
		if observation.Kind == ResultLineThresholdObservationV1 && (float64(observation.LineThresholdLimit) != legacyObservation.Limit || observation.Available && float64(observation.Actual) != legacyObservation.Actual) {
			t.Fatalf("threshold %d = %#v, want %#v", i, observation, legacyObservation)
		}
		if observation.Kind == ResultRatioObservationV1 && observation.Available && float64(observation.Numerator)/float64(observation.Denominator) != legacyObservation.Actual {
			t.Fatalf("ratio %d = %#v, want %#v", i, observation, legacyObservation)
		}
	}
}
