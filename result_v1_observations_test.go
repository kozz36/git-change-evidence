package changeevidence

import (
	"math"
	"reflect"
	"testing"
)

func TestResultV1ObservationsUseClosedOrderAndUnavailableVariants(t *testing.T) {
	policy := observationPolicyView(t, []CategoryInput{
		{Name: "source", LineThreshold: 2}, {Name: "binary", LineThreshold: 1},
		{Name: "reference"}, {Name: "empty"},
	}, "source", []RatioInput{
		{Category: "source", Reference: "reference", Maximum: "0.5"},
		{Category: "binary", Reference: "source", Maximum: "0.5"},
		{Category: "source", Reference: "binary", Maximum: "0.5"},
		{Category: "source", Reference: "empty", Maximum: "0.5"},
		{Category: "empty", Reference: "reference", Maximum: "0.5"},
	})
	got, err := resultObservations(policy, []CategoryTotal{
		{Category: "source", Additions: 2, Deletions: 1}, {Category: "binary", NonCountable: 1},
		{Category: "reference", Additions: 4}, {Category: "empty"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []ResultObservationV1{
		{Kind: ResultLineThresholdObservationV1, Category: "source", LineThresholdLimit: 2, Available: true, Actual: 3, Exceeded: true},
		{Kind: ResultLineThresholdObservationV1, Category: "binary", LineThresholdLimit: 1},
		{Kind: ResultRatioObservationV1, Category: "source", Reference: "reference", RatioLimit: "0.5", Available: true, Numerator: 3, Denominator: 4, Exceeded: true},
		{Kind: ResultRatioObservationV1, Category: "binary", Reference: "source", RatioLimit: "0.5"},
		{Kind: ResultRatioObservationV1, Category: "source", Reference: "binary", RatioLimit: "0.5"},
		{Kind: ResultRatioObservationV1, Category: "source", Reference: "empty", RatioLimit: "0.5"},
		{Kind: ResultRatioObservationV1, Category: "empty", Reference: "reference", RatioLimit: "0.5", Available: true, Denominator: 4},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("observations = %#v, want %#v", got, want)
	}
}

func TestResultV1StrictGreaterThanExcludesEqualityBoundary(t *testing.T) {
	policy := observationPolicyView(t, []CategoryInput{{Name: "source", LineThreshold: 3}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}})
	got, err := resultObservations(policy, []CategoryTotal{{Category: "source", Additions: 2, Deletions: 1}, {Category: "reference", Additions: 6}})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Exceeded || got[1].Exceeded || !got[0].Available || !got[1].Available {
		t.Fatalf("equality observations = %#v", got)
	}
}

func TestResultV1RejectsRequiredLineTotalOverflow(t *testing.T) {
	for _, test := range []struct {
		name   string
		policy PolicyView
		totals []CategoryTotal
	}{
		{"threshold", observationPolicyView(t, []CategoryInput{{Name: "source", LineThreshold: 1}}, "source", nil), []CategoryTotal{{Category: "source", Additions: math.MaxUint64, Deletions: 1}}},
		{"ratio numerator", observationPolicyView(t, []CategoryInput{{Name: "source"}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}}), []CategoryTotal{{Category: "source", Additions: math.MaxUint64, Deletions: 1}, {Category: "reference", Additions: 1}}},
		{"ratio reference", observationPolicyView(t, []CategoryInput{{Name: "source"}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}}), []CategoryTotal{{Category: "source", Additions: 1}, {Category: "reference", Additions: math.MaxUint64, Deletions: 1}}},
		{"reference overflow despite non-countable numerator", observationPolicyView(t, []CategoryInput{{Name: "source"}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}}), []CategoryTotal{{Category: "source", NonCountable: 1}, {Category: "reference", Additions: math.MaxUint64, Deletions: 1}}},
		{"late threshold after available threshold", observationPolicyView(t, []CategoryInput{{Name: "source", LineThreshold: 1}, {Name: "reference", LineThreshold: 1}}, "source", nil), []CategoryTotal{{Category: "source", Additions: 2}, {Category: "reference", Additions: math.MaxUint64, Deletions: 1}}},
		{"ratio reference after available threshold", observationPolicyView(t, []CategoryInput{{Name: "source", LineThreshold: 1}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}}), []CategoryTotal{{Category: "source", Additions: 2}, {Category: "reference", Additions: math.MaxUint64, Deletions: 1}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := resultObservations(test.policy, test.totals)
			if err == nil || got != nil {
				t.Fatalf("observations, error = %#v, %v", got, err)
			}
		})
	}
}

func TestResultV1NonCountableTotalsRemainUnavailableBeforeAddition(t *testing.T) {
	policy := observationPolicyView(t, []CategoryInput{{Name: "source", LineThreshold: 1}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}})
	got, err := resultObservations(policy, []CategoryTotal{{Category: "source", Additions: math.MaxUint64, Deletions: 1, NonCountable: 1}, {Category: "reference", Additions: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Available || got[0].Exceeded || got[1].Available || got[1].Exceeded {
		t.Fatalf("observations = %#v", got)
	}
}

func observationPolicyView(t *testing.T, categories []CategoryInput, defaultCategory string, ratios []RatioInput) PolicyView {
	t.Helper()
	document, err := NewAccountingPolicyV1(PolicyInput{Categories: categories, Default: defaultCategory, Ratios: ratios})
	if err != nil {
		t.Fatal(err)
	}
	return document.View()
}
