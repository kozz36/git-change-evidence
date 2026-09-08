package changeevidence_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestW2FrozenFixtureUsesNeutralPublicChain(t *testing.T) {
	vector, err := w2Translate(w2FixtureBytes(t))
	if err != nil {
		t.Fatal(err)
	}
	wantCategories := []string{"artifacts", "tests", "mechanical", "production", "production"}
	if !reflect.DeepEqual(vector.categories, wantCategories) || vector.totals != [4]uint64{5, 3, 2, 4} {
		t.Fatalf("frozen expected = %#v/%#v", vector.categories, vector.totals)
	}
	result := w2Result(t, vector)
	if result.Revisions() != (evidence.RevisionIdentity{Base: w2Base, Head: w2Head}) {
		t.Fatalf("synthetic revisions = %#v", result.Revisions())
	}
	byPath := make(map[string]string, len(vector.entries))
	for i, entry := range vector.entries {
		byPath[entry.path] = vector.categories[i]
	}
	var gotPaths []string
	for _, entry := range result.Entries() {
		path := string(entry.Path)
		gotPaths = append(gotPaths, path)
		if entry.Category != byPath[path] {
			t.Fatalf("category for raw path %q = %q, want %q", path, entry.Category, byPath[path])
		}
	}
	if want := []string{"binary.bin", "generated.json", "openspec/changes/vector.md", "src/core.py", "tests/vector.py"}; !reflect.DeepEqual(gotPaths, want) {
		t.Fatalf("result raw-path order = %#v, want %#v", gotPaths, want)
	}
	if entries := result.Entries(); entries[0].Measurement.Kind != evidence.ResultEntryNonCountableV1 {
		t.Fatalf("binary entry measurement = %#v, want non-countable", entries[0].Measurement)
	}
	wantTotals := []evidence.ResultCategoryTotalV1{
		{Category: "artifacts", Additions: 5}, {Category: "tests", Additions: 3},
		{Category: "mechanical", Additions: 2}, {Category: "production", Additions: 4, NonCountable: 1},
	}
	if !reflect.DeepEqual(result.Totals(), wantTotals) {
		t.Fatalf("neutral totals = %#v, want %#v", result.Totals(), wantTotals)
	}
	wantObservations := []evidence.ResultObservationV1{{
		Kind: evidence.ResultLineThresholdObservationV1, Category: "production", LineThresholdLimit: 400,
	}}
	if !reflect.DeepEqual(result.Observations(), wantObservations) {
		t.Fatalf("binary production observation = %#v, want unavailable %#v", result.Observations(), wantObservations)
	}
	if known := w2KnownProduction(vector); known != 4 || known > 400 {
		t.Fatalf("legacy known production threshold = %d/%t, want 4/false", known, known > 400)
	}
}

func TestW2SyntheticCopySeparatesLegacyThresholdFromNeutralAvailability(t *testing.T) {
	vector := w2Vector{mechanical: "generated.json", categories: []string{"artifacts", "tests", "mechanical", "production", "production"}, entries: []w2Entry{
		{path: "openspec/changes/vector.md", additions: 5, countable: true},
		{path: "tests/vector.py", additions: 3, countable: true},
		{path: "generated.json", additions: 2, countable: true},
		{path: "src/core.py", additions: 401, countable: true}, {path: "binary.bin"},
	}}
	result := w2Result(t, vector)
	wantTotals := []evidence.ResultCategoryTotalV1{
		{Category: "artifacts", Additions: 5}, {Category: "tests", Additions: 3},
		{Category: "mechanical", Additions: 2}, {Category: "production", Additions: 401, NonCountable: 1},
	}
	if !reflect.DeepEqual(result.Totals(), wantTotals) {
		t.Fatalf("synthetic totals = %#v, want %#v", result.Totals(), wantTotals)
	}
	wantObservation := evidence.ResultObservationV1{Kind: evidence.ResultLineThresholdObservationV1, Category: "production", LineThresholdLimit: 400}
	if got := result.Observations(); len(got) != 1 || got[0] != wantObservation || w2KnownProduction(vector) != 401 || !(w2KnownProduction(vector) > 400) {
		t.Fatalf("synthetic threshold separation = %#v/%d", got, w2KnownProduction(vector))
	}
}

func w2KnownProduction(vector w2Vector) uint64 {
	var total uint64
	for i, entry := range vector.entries {
		if vector.categories[i] == "production" && entry.countable {
			total += entry.additions + entry.deletions
		}
	}
	return total
}

func w2FixtureBytes(t *testing.T) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "compatibility", "v1", "accounting-v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
