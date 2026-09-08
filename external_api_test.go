package changeevidence_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

const (
	externalBase = "0123456789abcdef0123456789abcdef01234567"
	externalHead = "89abcdef0123456789abcdef0123456789abcdef"
)

func TestExternalPublicResultV1ChainIsReproducibleAndBound(t *testing.T) {
	policyA, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{
		{Name: "docs", PathGlobs: [][]byte{[]byte("docs/**")}, LineThreshold: 2},
		{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 5},
		{Name: "other", LineThreshold: 3},
	}, Default: "other", Ratios: []evidence.RatioInput{{Category: "source", Reference: "docs", Maximum: "0.5"}}})
	if err != nil {
		t.Fatal(err)
	}
	inventoryA, err := evidence.NewInventoryV1(policyA, []evidence.InventoryEntryInput{{Path: []byte("docs/guide.md"), Content: []byte("guide\n")}, {Path: []byte("src/new.go"), Content: []byte("package source\n")}})
	if err != nil {
		t.Fatal(err)
	}
	snapshotA := evidence.NewCommittedSnapshot(externalBase, externalHead, []evidence.CommittedChange{
		{Status: "M", Path: "docs/guide.md", Lines: evidence.CommittedLineCounts{Additions: 3, Deletions: 1, Countable: true}},
		{Status: "R", PreviousPath: "docs/old.go", Path: "src/new.go", Lines: evidence.CommittedLineCounts{Additions: 4, Countable: true}},
		{Status: "M", Path: "assets/logo.bin"},
	})
	resultA, err := evidence.NewAccountingResultV1(policyA, inventoryA, snapshotA)
	if err != nil {
		t.Fatal(err)
	}
	reportA, err := evidence.NewReportV1FromResult("external consumer", evidence.ResultV1ReportBinding{Policy: policyA, Inventory: inventoryA, Result: resultA})
	if err != nil {
		t.Fatal(err)
	}

	policyB, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{
		{Name: "docs", PathGlobs: [][]byte{[]byte("docs/**")}, LineThreshold: 2},
		{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 5},
		{Name: "other", LineThreshold: 3},
	}, Default: "other", Ratios: []evidence.RatioInput{{Category: "source", Reference: "docs", Maximum: "0.5"}}})
	if err != nil {
		t.Fatal(err)
	}
	inventoryB, err := evidence.NewInventoryV1(policyB, []evidence.InventoryEntryInput{{Path: []byte("docs/guide.md"), Content: []byte("guide\n")}, {Path: []byte("src/new.go"), Content: []byte("package source\n")}})
	if err != nil {
		t.Fatal(err)
	}
	resultB, err := evidence.NewAccountingResultV1(policyB, inventoryB, evidence.NewCommittedSnapshot(externalBase, externalHead, []evidence.CommittedChange{
		{Status: "M", Path: "docs/guide.md", Lines: evidence.CommittedLineCounts{Additions: 3, Deletions: 1, Countable: true}},
		{Status: "R", PreviousPath: "docs/old.go", Path: "src/new.go", Lines: evidence.CommittedLineCounts{Additions: 4, Countable: true}},
		{Status: "M", Path: "assets/logo.bin"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	reportB, err := evidence.NewReportV1FromResult("external consumer", evidence.ResultV1ReportBinding{Policy: policyB, Inventory: inventoryB, Result: resultB})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resultA.CanonicalBytes(), resultB.CanonicalBytes()) || resultA.Digest() != resultB.Digest() || !bytes.Equal(reportA.CanonicalBytes(), reportB.CanonicalBytes()) || reportA.Digest() != reportB.Digest() {
		t.Fatal("equivalent public typed chains are not canonical and reproducible")
	}
	wantEntries := []evidence.ResultEntryV1{{Path: []byte("assets/logo.bin"), Category: "other", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryNonCountableV1}}, {Path: []byte("docs/guide.md"), Category: "docs", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryCountableV1, Additions: 3, Deletions: 1}}, {Path: []byte("src/new.go"), Category: "source", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryCountableV1, Additions: 4}}}
	wantTotals := []evidence.ResultCategoryTotalV1{{Category: "docs", Additions: 3, Deletions: 1}, {Category: "source", Additions: 4}, {Category: "other", NonCountable: 1}}
	wantObservations := []evidence.ResultObservationV1{{Kind: evidence.ResultLineThresholdObservationV1, Category: "docs", LineThresholdLimit: 2, Available: true, Actual: 4, Exceeded: true}, {Kind: evidence.ResultLineThresholdObservationV1, Category: "source", LineThresholdLimit: 5, Available: true, Actual: 4}, {Kind: evidence.ResultLineThresholdObservationV1, Category: "other", LineThresholdLimit: 3}, {Kind: evidence.ResultRatioObservationV1, Category: "source", Reference: "docs", RatioLimit: "0.5", Available: true, Numerator: 4, Denominator: 4, Exceeded: true}}
	if !reflect.DeepEqual(resultA.Entries(), wantEntries) || !reflect.DeepEqual(resultA.Totals(), wantTotals) || !reflect.DeepEqual(resultA.Observations(), wantObservations) || resultA.Revisions() != (evidence.RevisionIdentity{Base: externalBase, Head: externalHead}) || evidence.ValidateReportV1ResultProvenance(reportA, evidence.ResultV1ReportBinding{Policy: policyA, Inventory: inventoryA, Result: resultA}) != nil {
		t.Fatalf("public result/report oracle = (%#v, %#v, %#v, %#v)", resultA.Entries(), resultA.Totals(), resultA.Observations(), reportA.Provenance())
	}

	otherPolicy, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{{Name: "other"}}, Default: "other"})
	if err != nil || otherPolicy.Digest() == policyA.Digest() {
		t.Fatalf("different nonempty policy fixture = (%q, %v)", otherPolicy.Digest(), err)
	}
	otherInventory, err := evidence.NewInventoryV1(otherPolicy, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := evidence.NewAccountingResultV1(policyA, otherInventory, snapshotA); err == nil || len(got.CanonicalBytes()) != 0 || got.Digest() != "" {
		t.Fatalf("wrong antecedent = (%#v, %v)", got, err)
	}
	if got, err := evidence.NewReportV1FromResult("external consumer", evidence.ResultV1ReportBinding{Policy: otherPolicy, Inventory: otherInventory, Result: resultA}); err == nil || len(got.CanonicalBytes()) != 0 || got.Digest() != "" {
		t.Fatalf("invalid binding = (%#v, %v)", got, err)
	}
	for _, test := range []struct {
		name, field, code string
		call              func() error
	}{
		{"empty policy", "categories", "required", func() error { _, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{}); return err }},
		{"invalid inventory path", "entries.path", "invalid_path", func() error {
			_, err := evidence.NewInventoryV1(policyA, []evidence.InventoryEntryInput{{Path: []byte("../invalid")}})
			return err
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var got *evidence.ContractError
			if err := test.call(); !errors.As(err, &got) || got.Field != test.field || got.Code != test.code {
				t.Fatalf("error = %#v, want %s/%s", err, test.field, test.code)
			}
		})
	}
}
