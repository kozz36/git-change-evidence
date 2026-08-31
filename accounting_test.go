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
