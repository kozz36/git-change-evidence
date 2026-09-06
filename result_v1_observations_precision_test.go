package changeevidence

import "testing"

func TestResultV1ExactRatioPrecisionEdge(t *testing.T) {
	policy := observationPolicyView(t, []CategoryInput{{Name: "source"}, {Name: "reference"}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}})
	got, err := resultObservations(policy, []CategoryTotal{{Category: "source", Additions: 9007199254740993}, {Category: "reference", Additions: 18014398509481984}})
	if err != nil {
		t.Fatal(err)
	}
	want := ResultObservationV1{Kind: ResultRatioObservationV1, Category: "source", Reference: "reference", RatioLimit: "0.5", Available: true, Numerator: 9007199254740993, Denominator: 18014398509481984, Exceeded: true}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("observation = %#v, want %#v", got, want)
	}
}

func TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference(t *testing.T) {
	for _, test := range []struct {
		name                              string
		source, reference                 CommittedLineCounts
		wantAvailable, wantResult, legacy bool
	}{
		{"ordinary shared domain", CommittedLineCounts{Additions: 3, Countable: true}, CommittedLineCounts{Additions: 4, Countable: true}, true, true, true},
		{"ordinary equality boundary", CommittedLineCounts{Additions: 1, Countable: true}, CommittedLineCounts{Additions: 2, Countable: true}, true, false, false},
		{"exact precision edge", CommittedLineCounts{Additions: 9007199254740993, Countable: true}, CommittedLineCounts{Additions: 18014398509481984, Countable: true}, true, true, false},
		{"unavailable non-countable numerator", CommittedLineCounts{Additions: 3}, CommittedLineCounts{Additions: 4, Countable: true}, false, false, false},
		{"unavailable non-countable reference", CommittedLineCounts{Additions: 3, Countable: true}, CommittedLineCounts{Additions: 4}, false, false, false},
		{"unavailable zero denominator", CommittedLineCounts{Additions: 3, Countable: true}, CommittedLineCounts{Countable: true}, false, false, false},
		{"available zero numerator", CommittedLineCounts{Countable: true}, CommittedLineCounts{Additions: 4, Countable: true}, true, false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			policy := observationPolicyView(t, []CategoryInput{{Name: "source", PathGlobs: [][]byte{[]byte("source")}}, {Name: "reference", PathGlobs: [][]byte{[]byte("reference")}}}, "source", []RatioInput{{Category: "source", Reference: "reference", Maximum: "0.5"}})
			snapshot := NewCommittedSnapshot("base", "head", []CommittedChange{{Path: "source", Lines: test.source}, {Path: "reference", Lines: test.reference}})
			entries := snapshot.Entries()
			_, totals, ok := resultAccountingDomain(policy).classifyAndTotal([]accountingEntry{{path: []byte("source"), lines: entries[0].Lines}, {path: []byte("reference"), lines: entries[1].Lines}})
			if !ok {
				t.Fatal("result totals overflowed")
			}
			got, err := resultObservations(policy, totals)
			if err != nil {
				t.Fatal(err)
			}
			legacy, err := Account(snapshot, AccountingPolicy{Categories: []AccountingCategory{{Name: "source", PathGlobs: []string{"source"}}, {Name: "reference", PathGlobs: []string{"reference"}}}, Default: "source", Ratios: []CategoryRatio{{Category: "source", Reference: "reference", Maximum: .5}}})
			if err != nil {
				t.Fatal(err)
			}
			if got[0].Available != test.wantAvailable || legacy.Observations[0].Available != test.wantAvailable {
				t.Fatalf("availability result, legacy = %t, %t; want %t", got[0].Available, legacy.Observations[0].Available, test.wantAvailable)
			}
			if got[0].Exceeded != test.wantResult || legacy.Observations[0].Exceeded != test.legacy {
				t.Fatalf("exceeded result, legacy = %t, %t; want %t, %t", got[0].Exceeded, legacy.Observations[0].Exceeded, test.wantResult, test.legacy)
			}
			if totals[0] != legacy.Totals[0] || totals[1] != legacy.Totals[1] {
				t.Fatalf("result totals, legacy totals = %#v, %#v", totals, legacy.Totals)
			}
		})
	}
}
