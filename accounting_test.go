package changeevidence

import (
	"errors"
	"math"
	"testing"
)

func TestAccountUsesFirstMatchingRawByteGlobAndDefault(t *testing.T) {
	raw := string([]byte{0xff}) + "file.go"
	snapshot := NewCommittedSnapshot("base", "head", []CommittedChange{
		{Path: raw, Lines: CommittedLineCounts{Additions: 4, Countable: true}},
		{Path: "src/main.go", Lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}},
		{Path: "assets/blob", Lines: CommittedLineCounts{}},
		{Path: "note.txt", Lines: CommittedLineCounts{Deletions: 5, Countable: true}},
	})
	got, err := Account(snapshot, AccountingPolicy{Categories: []AccountingCategory{
		{Name: "raw", PathGlobs: []string{string([]byte{0xff}) + "*.go"}, LineThreshold: 3},
		{Name: "first", PathGlobs: []string{"**/*.go"}, LineThreshold: 2},
		{Name: "second", PathGlobs: []string{"src/**"}},
		{Name: "binary", PathGlobs: []string{"assets/**"}, LineThreshold: 1},
		{Name: "other"},
	}, Default: "other", Ratios: []CategoryRatio{{"raw", "other", .5}, {"binary", "raw", .5}}})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []CategoryTotal{
		{"raw", 4, 0, 0}, {"first", 2, 1, 0}, {"second", 0, 0, 0}, {"binary", 0, 0, 1}, {"other", 0, 5, 0},
	} {
		if got := categoryTotal(got, want.Category); got != want {
			t.Fatalf("total %q = %#v, want %#v", want.Category, got, want)
		}
	}
	if got := accountingObservation(got, "raw", ""); !got.Available || !got.Exceeded || got.Actual != 4 {
		t.Fatalf("raw threshold = %#v", got)
	}
	if got := accountingObservation(got, "binary", ""); got.Available || got.Exceeded {
		t.Fatalf("binary threshold = %#v", got)
	}
	if got := accountingObservation(got, "raw", "other"); !got.Available || !got.Exceeded {
		t.Fatalf("ratio = %#v", got)
	}
	if got := accountingObservation(got, "binary", "raw"); got.Available || got.Exceeded {
		t.Fatalf("non-countable ratio = %#v", got)
	}
}

func TestAccountRejectsInvalidPolicyBeforeAccounting(t *testing.T) {
	base := AccountingPolicy{Categories: []AccountingCategory{{Name: "one"}, {Name: "other"}}, Default: "other"}
	for _, test := range []struct {
		name   string
		policy AccountingPolicy
	}{
		{"categories", AccountingPolicy{}},
		{"duplicate category", AccountingPolicy{Categories: []AccountingCategory{{Name: "one"}, {Name: "one"}}, Default: "one"}},
		{"default required", AccountingPolicy{Categories: base.Categories}},
		{"default", AccountingPolicy{Categories: base.Categories, Default: "missing"}},
		{"pattern", AccountingPolicy{Categories: []AccountingCategory{{Name: "one", PathGlobs: []string{"["}}, {Name: "other"}}, Default: "other"}},
		{"reference", AccountingPolicy{Categories: base.Categories, Default: "other", Ratios: []CategoryRatio{{"one", "missing", .5}}}},
		{"threshold", AccountingPolicy{Categories: []AccountingCategory{{Name: "one", LineThreshold: -1}, {Name: "other"}}, Default: "other"}},
		{"ratio", AccountingPolicy{Categories: base.Categories, Default: "other", Ratios: []CategoryRatio{{"one", "other", math.NaN()}}}},
		{"duplicate ratio", AccountingPolicy{Categories: base.Categories, Default: "other", Ratios: []CategoryRatio{{"one", "other", .5}, {"one", "other", .5}}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := Account(NewCommittedSnapshot("base", "head", []CommittedChange{{Path: "x", Lines: CommittedLineCounts{Additions: math.MaxUint64, Countable: true}}, {Path: "x", Lines: CommittedLineCounts{Additions: 1, Countable: true}}}), test.policy)
			var accountingError *AccountingError
			if !errors.As(err, &accountingError) || accountingError.Code != AccountingInvalidPolicy || len(got.Totals) != 0 || len(got.Observations) != 0 {
				t.Fatalf("result, error = %#v, %v", got, err)
			}
		})
	}
}

func TestAccountReturnsEmptyResultOnOverflow(t *testing.T) {
	got, err := Account(NewCommittedSnapshot("base", "head", []CommittedChange{{Path: "x", Lines: CommittedLineCounts{Additions: math.MaxUint64, Countable: true}}, {Path: "x", Lines: CommittedLineCounts{Additions: 1, Countable: true}}}), AccountingPolicy{Categories: []AccountingCategory{{Name: "other"}}, Default: "other"})
	var accountingError *AccountingError
	if !errors.As(err, &accountingError) || accountingError.Code != AccountingOverflow || len(got.Totals) != 0 || len(got.Observations) != 0 {
		t.Fatalf("result, error = %#v, %v", got, err)
	}
}

func TestPathGlobGrammar(t *testing.T) {
	for _, test := range []struct {
		pattern, path string
		matches       bool
	}{
		{"*.go", "dir/file.go", false}, {"**/*.go", "file.go", true}, {"a?.go", "ab.go", true}, {"file-[0-9].txt", "file-7.txt", true},
	} {
		if got := matchPathGlob([]byte(test.pattern), []byte(test.path)); got != test.matches {
			t.Fatalf("%q %q = %v", test.pattern, test.path, got)
		}
	}
}

func TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable(t *testing.T) {
	policy := AccountingPolicy{Categories: []AccountingCategory{
		{Name: "raw", PathGlobs: []string{string([]byte{0xff}) + "*.go"}},
		{Name: "source", PathGlobs: []string{"**/*.go"}},
		{Name: "second", PathGlobs: []string{"src/**"}},
		{Name: "binary", PathGlobs: []string{"assets/**"}},
		{Name: "other"},
	}, Default: "other"}
	index, err := validateAccountingPolicy(policy)
	if err != nil {
		t.Fatal(err)
	}
	domain := legacyAccountingDomain(policy, index)
	classifications, totals, ok := domain.classifyAndTotal([]accountingEntry{
		{path: append([]byte{0xff}, []byte("file.go")...), lines: CommittedLineCounts{Additions: 4, Countable: true}},
		{path: []byte("src/main.go"), lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}},
		{path: []byte("assets/blob"), lines: CommittedLineCounts{Additions: 99, Deletions: 99}},
		{path: []byte("note.txt"), lines: CommittedLineCounts{Deletions: 5, Countable: true}},
	})
	if !ok {
		t.Fatal("classification failed")
	}
	for i, want := range []int{0, 1, 3, 4} {
		if got := classifications[i].category; got != want {
			t.Fatalf("classification %d category = %d, want %d", i, got, want)
		}
	}
	for i, want := range []CategoryTotal{
		{"raw", 4, 0, 0}, {"source", 2, 1, 0}, {"second", 0, 0, 0}, {"binary", 0, 0, 1}, {"other", 0, 5, 0},
	} {
		if got := totals[i]; got != want {
			t.Fatalf("total %d = %#v, want %#v", i, got, want)
		}
	}
}

func TestAccountingDomainRejectsCheckedAccumulationOverflow(t *testing.T) {
	policy := AccountingPolicy{Categories: []AccountingCategory{{Name: "other"}}, Default: "other"}
	index, err := validateAccountingPolicy(policy)
	if err != nil {
		t.Fatal(err)
	}
	_, _, ok := legacyAccountingDomain(policy, index).classifyAndTotal([]accountingEntry{
		{path: []byte("first"), lines: CommittedLineCounts{Additions: math.MaxUint64, Countable: true}},
		{path: []byte("second"), lines: CommittedLineCounts{Additions: 1, Countable: true}},
	})
	if ok {
		t.Fatal("classification succeeded after accumulation overflow")
	}
}

func TestAccountSharedDomainTriangulatesLegacyBehavior(t *testing.T) {
	for _, test := range []struct {
		name     string
		snapshot CommittedSnapshot
		policy   AccountingPolicy
		want     AccountingResult
		errCode  string
	}{
		{
			name: "first raw byte match wins",
			snapshot: NewCommittedSnapshot("base", "head", []CommittedChange{
				{Path: string([]byte{0xff}) + "file.go", Lines: CommittedLineCounts{Additions: 4, Countable: true}},
			}),
			policy: AccountingPolicy{Categories: []AccountingCategory{
				{Name: "raw", PathGlobs: []string{string([]byte{0xff}) + "*.go"}, LineThreshold: 3},
				{Name: "source", PathGlobs: []string{"**/*.go"}, LineThreshold: 1}, {Name: "other"},
			}, Default: "other"},
			want: AccountingResult{Totals: []CategoryTotal{{"raw", 4, 0, 0}, {"source", 0, 0, 0}, {"other", 0, 0, 0}}, Observations: []AccountingObservation{
				{Category: "raw", Limit: 3, Available: true, Actual: 4, Exceeded: true},
				{Category: "source", Limit: 1, Available: true},
			}},
		},
		{
			name:     "default retains all zero category and threshold order",
			snapshot: NewCommittedSnapshot("base", "head", []CommittedChange{{Path: "note.txt", Lines: CommittedLineCounts{Deletions: 3, Countable: true}}}),
			policy: AccountingPolicy{Categories: []AccountingCategory{
				{Name: "source", PathGlobs: []string{"**/*.go"}, LineThreshold: 1}, {Name: "other", LineThreshold: 2},
			}, Default: "other"},
			want: AccountingResult{Totals: []CategoryTotal{{"source", 0, 0, 0}, {"other", 0, 3, 0}}, Observations: []AccountingObservation{
				{Category: "source", Limit: 1, Available: true}, {Category: "other", Limit: 2, Available: true, Actual: 3, Exceeded: true},
			}},
		},
		{
			name:     "non countable ignores supplied lines",
			snapshot: NewCommittedSnapshot("base", "head", []CommittedChange{{Path: "assets/blob", Lines: CommittedLineCounts{Additions: 12, Deletions: 7}}}),
			policy: AccountingPolicy{Categories: []AccountingCategory{
				{Name: "binary", PathGlobs: []string{"assets/**"}, LineThreshold: 1}, {Name: "other"},
			}, Default: "other", Ratios: []CategoryRatio{{Category: "binary", Reference: "other", Maximum: .5}}},
			want: AccountingResult{Totals: []CategoryTotal{{"binary", 0, 0, 1}, {"other", 0, 0, 0}}, Observations: []AccountingObservation{
				{Category: "binary", Limit: 1}, {Category: "binary", Reference: "other", Limit: .5},
			}},
		},
		{
			name: "max uint plus one fails",
			snapshot: NewCommittedSnapshot("base", "head", []CommittedChange{
				{Path: "first", Lines: CommittedLineCounts{Additions: math.MaxUint64, Countable: true}},
				{Path: "second", Lines: CommittedLineCounts{Additions: 1, Countable: true}},
			}),
			policy:  AccountingPolicy{Categories: []AccountingCategory{{Name: "other"}}, Default: "other"},
			errCode: AccountingOverflow,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := Account(test.snapshot, test.policy)
			if test.errCode != "" {
				var accountingError *AccountingError
				if !errors.As(err, &accountingError) || accountingError.Code != test.errCode || len(got.Totals) != 0 || len(got.Observations) != 0 {
					t.Fatalf("result, error = %#v, %v", got, err)
				}
				return
			}
			if err != nil || len(got.Totals) != len(test.want.Totals) || len(got.Observations) != len(test.want.Observations) {
				t.Fatalf("result, error = %#v, %v", got, err)
			}
			for i := range test.want.Totals {
				if got.Totals[i] != test.want.Totals[i] {
					t.Fatalf("total %d = %#v, want %#v", i, got.Totals[i], test.want.Totals[i])
				}
			}
			for i := range test.want.Observations {
				if got.Observations[i] != test.want.Observations[i] {
					t.Fatalf("observation %d = %#v, want %#v", i, got.Observations[i], test.want.Observations[i])
				}
			}
		})
	}
}

func TestResultAccountingDomainProjectsPolicyViewLosslessly(t *testing.T) {
	view := PolicyView{Categories: []CategoryInput{
		{Name: "source", PathGlobs: [][]byte{{0xff, '*'}}},
		{Name: "other"},
	}, Default: "other"}
	domain := resultAccountingDomain(view)
	view.Categories[0].PathGlobs[0][0] = 'x'
	classifications, totals, ok := domain.classifyAndTotal([]accountingEntry{
		{path: []byte{0xff, 'f'}, lines: CommittedLineCounts{Additions: 2, Countable: true}},
		{path: []byte("note"), lines: CommittedLineCounts{Deletions: 1, Countable: true}},
	})
	if !ok || len(classifications) != 2 || classifications[0].category != 0 || classifications[1].category != 1 {
		t.Fatalf("classifications, success = %#v, %v", classifications, ok)
	}
	if totals[0] != (CategoryTotal{Category: "source", Additions: 2}) || totals[1] != (CategoryTotal{Category: "other", Deletions: 1}) {
		t.Fatalf("totals = %#v", totals)
	}
}

func categoryTotal(result AccountingResult, category string) CategoryTotal {
	for _, total := range result.Totals {
		if total.Category == category {
			return total
		}
	}
	return CategoryTotal{}
}
func accountingObservation(result AccountingResult, category, reference string) AccountingObservation {
	for _, observation := range result.Observations {
		if observation.Category == category && observation.Reference == reference {
			return observation
		}
	}
	return AccountingObservation{}
}
