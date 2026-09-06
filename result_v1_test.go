package changeevidence

import (
	"bytes"
	"reflect"
	"testing"
)

func TestResultV1ModelZeroValue(t *testing.T) {
	document := ResultDocumentV1{}
	for _, test := range []struct {
		name string
		got  any
		want any
	}{
		{"canonical bytes", document.CanonicalBytes(), []byte(nil)},
		{"digest", document.Digest(), DocumentDigest("")},
		{"policy digest", document.AccountingPolicyDigest(), DocumentDigest("")},
		{"inventory digest", document.InventoryDigest(), DocumentDigest("")},
		{"revisions", document.Revisions(), RevisionIdentity{}},
		{"entries", document.Entries(), []ResultEntryV1(nil)},
		{"totals", document.Totals(), []ResultCategoryTotalV1(nil)},
		{"observations", document.Observations(), []ResultObservationV1(nil)},
	} {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.got, test.want) {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestResultV1ModelPopulatedAccessors(t *testing.T) {
	document := modelResultDocument()
	for _, test := range []struct {
		name string
		got  any
		want any
	}{
		{"canonical bytes", document.CanonicalBytes(), []byte("canonical\n")},
		{"digest", document.Digest(), DocumentDigest("result")},
		{"policy digest", document.AccountingPolicyDigest(), DocumentDigest("policy")},
		{"inventory digest", document.InventoryDigest(), DocumentDigest("inventory")},
		{"revisions", document.Revisions(), RevisionIdentity{Base: "base", Head: "head"}},
		{"entries", document.Entries(), modelResultEntries()},
		{"totals", document.Totals(), modelResultTotals()},
		{"observations", document.Observations(), modelResultObservations()},
	} {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.got, test.want) {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestResultV1ModelDefensiveCopiesPreserveState(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(ResultDocumentV1)
	}{
		{"canonical bytes", func(document ResultDocumentV1) { document.CanonicalBytes()[0] = 'x' }},
		{"entry slice", func(document ResultDocumentV1) { document.Entries()[0] = ResultEntryV1{} }},
		{"entry path", func(document ResultDocumentV1) { document.Entries()[0].Path[0] = 'x' }},
		{"totals slice", func(document ResultDocumentV1) { document.Totals()[0].Category = "changed" }},
		{"observations slice", func(document ResultDocumentV1) { document.Observations()[0].Category = "changed" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			document := modelResultDocument()
			test.mutate(document)
			if got := document.CanonicalBytes(); !bytes.Equal(got, []byte("canonical\n")) {
				t.Fatalf("canonical bytes changed: %q", got)
			}
			if got := document.Entries(); !reflect.DeepEqual(got, modelResultEntries()) {
				t.Fatalf("entries changed: %#v", got)
			}
			if got := document.Totals(); !reflect.DeepEqual(got, modelResultTotals()) {
				t.Fatalf("totals changed: %#v", got)
			}
			if got := document.Observations(); !reflect.DeepEqual(got, modelResultObservations()) {
				t.Fatalf("observations changed: %#v", got)
			}
		})
	}
}

func TestResultV1ModelReturnsEmptyViewsForPopulatedDocument(t *testing.T) {
	document := ResultDocumentV1{canonical: []byte("{}\n")}
	for _, test := range []struct {
		name string
		got  any
		want any
	}{
		{"entries", document.Entries(), []ResultEntryV1{}},
		{"totals", document.Totals(), []ResultCategoryTotalV1{}},
		{"observations", document.Observations(), []ResultObservationV1{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.got, test.want) {
				t.Fatalf("got %#v, want %#v", test.got, test.want)
			}
		})
	}
}

func TestResultV1ModelDiscriminators(t *testing.T) {
	for _, test := range []struct{ name, got, want string }{
		{"countable measurement", string(ResultEntryCountableV1), "countable"},
		{"non-countable measurement", string(ResultEntryNonCountableV1), "non_countable"},
		{"line threshold observation", string(ResultLineThresholdObservationV1), "line_threshold"},
		{"ratio observation", string(ResultRatioObservationV1), "ratio"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %q, want %q", test.got, test.want)
			}
		})
	}
}

func modelResultDocument() ResultDocumentV1 {
	return ResultDocumentV1{
		canonical:              []byte("canonical\n"),
		digest:                 "result",
		accountingPolicyDigest: "policy",
		inventoryDigest:        "inventory",
		revisions:              RevisionIdentity{Base: "base", Head: "head"},
		entries:                modelResultEntries(),
		totals:                 modelResultTotals(),
		observations:           modelResultObservations(),
	}
}

func modelResultEntries() []ResultEntryV1 {
	return []ResultEntryV1{
		{Path: []byte{0xff, '/', 'a'}, Category: "source", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 2, Deletions: 1}},
		{Path: []byte("asset"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryNonCountableV1}},
	}
}

func modelResultTotals() []ResultCategoryTotalV1 {
	return []ResultCategoryTotalV1{{Category: "source", Additions: 2, Deletions: 1}, {Category: "other", NonCountable: 1}}
}

func modelResultObservations() []ResultObservationV1 {
	return []ResultObservationV1{
		{Kind: ResultLineThresholdObservationV1, Category: "source", LineThresholdLimit: 3, Available: true, Actual: 3},
		{Kind: ResultRatioObservationV1, Category: "source", Reference: "other", RatioLimit: "0.5", Available: true, Numerator: 1, Denominator: 2},
	}
}
