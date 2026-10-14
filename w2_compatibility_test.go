package changeevidence_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/cnsicprofile"
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

func TestCNSICProfileW2FrozenChainAndDistinctNeutrality(t *testing.T) {
	vector, err := w2Translate(w2FixtureBytes(t))
	if err != nil {
		t.Fatal(err)
	}
	policy, err := evidence.NewAccountingPolicyV1(w2Policy(vector.mechanical))
	if err != nil {
		t.Fatal(err)
	}
	changes := make([]evidence.CommittedChange, len(vector.entries))
	for i, entry := range vector.entries {
		changes[i] = evidence.CommittedChange{Status: "A", Path: entry.path, Lines: evidence.CommittedLineCounts{Additions: entry.additions, Deletions: entry.deletions, Countable: entry.countable}}
	}
	build := func(p evidence.PolicyDocumentV1) evidence.ResultDocumentV1 {
		t.Helper()
		inventory, err := evidence.NewInventoryV1(p, nil)
		if err != nil {
			t.Fatal(err)
		}
		result, err := evidence.NewAccountingResultV1(p, inventory, evidence.NewCommittedSnapshot(w2Base, w2Head, changes))
		if err != nil {
			t.Fatal(err)
		}
		if result.AccountingPolicyDigest() != p.Digest() || result.InventoryDigest() != inventory.Digest() {
			t.Fatal("policy binding")
		}
		return result
	}
	want := w2Result(t, vector)
	for pass := 0; pass < 2; pass++ {
		admitted, err := cnsicprofile.DecodePolicyV1(policy.CanonicalBytes())
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(admitted.CanonicalBytes(), policy.CanonicalBytes()) || admitted.Digest() != policy.Digest() || !reflect.DeepEqual(admitted.View(), policy.View()) {
			t.Fatal("policy parity")
		}
		got := build(admitted)
		if len(got.Entries()) != 5 || len(got.Observations()) != 1 || !bytes.Equal(got.CanonicalBytes(), want.CanonicalBytes()) || !reflect.DeepEqual(got.Entries(), want.Entries()) || !reflect.DeepEqual(got.Totals(), want.Totals()) || !reflect.DeepEqual(got.Observations(), want.Observations()) {
			t.Fatal("frozen chain parity")
		}
		if pass == 1 {
			break
		}
		other, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{{Name: "note", PathGlobs: [][]byte{[]byte("tests/**")}}, {Name: "engine", LineThreshold: 7}}, Default: "engine"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cnsicprofile.DecodePolicyV1(other.CanonicalBytes()); err == nil {
			t.Fatal("other profile admitted")
		}
		result := build(other)
		if totals := result.Totals(); !reflect.DeepEqual(totals, []evidence.ResultCategoryTotalV1{{Category: "note", Additions: 3}, {Category: "engine", Additions: 11, NonCountable: 1}}) {
			t.Fatalf("distinct totals: %#v", totals)
		}
		for _, entry := range result.Entries() {
			category := "engine"
			if string(entry.Path) == "tests/vector.py" {
				category = "note"
			}
			if entry.Category != category {
				t.Fatal("distinct classification")
			}
		}
		if observations := result.Observations(); len(observations) != 1 || observations[0] != (evidence.ResultObservationV1{Kind: evidence.ResultLineThresholdObservationV1, Category: "engine", LineThresholdLimit: 7}) {
			t.Fatalf("distinct threshold: %#v", observations)
		}
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
