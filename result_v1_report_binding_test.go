package changeevidence

import (
	"bytes"
	"testing"
)

const (
	reportBindingAccounting = "fa1bd02a06e213818cc916c2a26f7f96128e3c1e68b4ac1f87a68e07d1ff9f92"
	reportBindingPolicy     = "4b851200a84eff95ff8e68d3c3caeaa8a4e802054890f555cad44131f2b19764"
	reportBindingInventory  = "79df4a4347656656736557d09fdfbd9e482a190c7903500bab1be2cacf70d66e"
	reportBindingBase       = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	reportBindingHead       = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

func TestResultV1ReportBindingDerivesAllProvenance(t *testing.T) {
	binding, alternate := resultReportBindings(t)
	report, err := NewReportV1FromResult("result-v1-bound", binding)
	if err != nil {
		t.Fatal(err)
	}
	want := Provenance{
		AccountingDigest:  Digest(reportBindingAccounting),
		InputPolicyDigest: Digest(reportBindingPolicy),
		InventoryDigest:   Digest(reportBindingInventory),
		Revisions:         RevisionIdentity{Base: reportBindingBase, Head: reportBindingHead},
	}
	if got := report.Provenance(); got != want {
		t.Fatalf("derived provenance = %#v, want %#v", got, want)
	}
	legacy, err := NewReportV1(ReportInput{Subject: "result-v1-bound", Provenance: want})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(report.CanonicalBytes(), legacy.CanonicalBytes()) || report.Digest() != legacy.Digest() {
		t.Fatalf("bound report differs from unchanged NewReportV1: %#v, %#v", report, legacy)
	}
	if err := ValidateReportV1ResultProvenance(report, binding); err != nil {
		t.Fatal(err)
	}
	cached := binding
	cached.Result.digest = alternate.Result.Digest()
	cached.Result.accountingPolicyDigest = alternate.Result.AccountingPolicyDigest()
	cached.Result.inventoryDigest = alternate.Result.InventoryDigest()
	cached.Result.revisions = alternate.Result.Revisions()
	fromCanonical, err := NewReportV1FromResult("result-v1-bound", cached)
	if err != nil || !bytes.Equal(fromCanonical.CanonicalBytes(), report.CanonicalBytes()) {
		t.Fatalf("result cached fields displaced canonical proof: %#v, %v", fromCanonical, err)
	}
}

func TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces(t *testing.T) {
	binding, alternate := resultReportBindings(t)
	report, err := NewReportV1FromResult("result-v1-bound", binding)
	if err != nil {
		t.Fatal(err)
	}
	alternateReport, err := NewReportV1FromResult("result-v1-alternate", alternate)
	if err != nil || ValidateReportV1ResultProvenance(alternateReport, alternate) != nil {
		t.Fatalf("independent alternate chain/report = %#v, %v", alternateReport, err)
	}
	alternateProvenance := alternateReport.Provenance()
	for _, test := range []struct {
		name, field, code string
		change            func(*Provenance)
	}{
		{"accounting", "provenance.accounting_sha256", "mismatched_digest", func(p *Provenance) { p.AccountingDigest = alternateProvenance.AccountingDigest }},
		{"policy", "provenance.input_policy_sha256", "mismatched_digest", func(p *Provenance) { p.InputPolicyDigest = alternateProvenance.InputPolicyDigest }},
		{"inventory", "provenance.inventory_sha256", "mismatched_digest", func(p *Provenance) { p.InventoryDigest = alternateProvenance.InventoryDigest }},
		{"base", "provenance.revisions.base", "mismatched_revision", func(p *Provenance) { p.Revisions.Base = alternateProvenance.Revisions.Base }},
		{"head", "provenance.revisions.head", "mismatched_revision", func(p *Provenance) { p.Revisions.Head = alternateProvenance.Revisions.Head }},
	} {
		t.Run(test.name, func(t *testing.T) {
			provenance := report.Provenance()
			test.change(&provenance)
			candidate, err := NewReportV1(ReportInput{Subject: "result-v1-bound", Provenance: provenance})
			if err != nil {
				t.Fatal(err)
			}
			assertContractErrorValue(t, ValidateReportV1ResultProvenance(candidate, binding), test.field, test.code)
		})
	}
	cached := report
	cached.digest, cached.provenance = "not-the-canonical-digest", alternateProvenance
	if err := ValidateReportV1ResultProvenance(cached, binding); err != nil {
		t.Fatalf("canonical report bytes were not authoritative: %v", err)
	}
	assertContractErrorValue(t, ValidateReportV1ResultProvenance(Evidence{canonical: []byte("null")}, ResultV1ReportBinding{}), "document", "invalid_object")
}

func TestNewReportV1FromResultReturnsZeroEvidenceForInvalidChains(t *testing.T) {
	binding, alternate := resultReportBindings(t)
	noncanonical := binding
	noncanonical.Result.canonical = append(noncanonical.Result.CanonicalBytes(), '\n')
	for _, test := range []struct {
		name    string
		binding ResultV1ReportBinding
	}{
		{"zero result", ResultV1ReportBinding{Policy: binding.Policy, Inventory: binding.Inventory}},
		{"noncanonical result", noncanonical},
		{"valid wrong antecedents", ResultV1ReportBinding{Policy: alternate.Policy, Inventory: alternate.Inventory, Result: binding.Result}},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewReportV1FromResult("result-v1-bound", test.binding)
			if err == nil {
				t.Fatal("invalid chain was accepted")
			}
			assertZeroEvidence(t, got)
		})
	}
}

func resultReportBindings(t *testing.T) (ResultV1ReportBinding, ResultV1ReportBinding) {
	t.Helper()
	policy, inventory, snapshot := resultDecodeFixture(t)
	result, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	alternatePolicy, alternateInventory := resultUnit5Antecedents(t, PolicyInput{
		Categories: []CategoryInput{{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 3}, {Name: "docs", PathGlobs: [][]byte{[]byte("docs/**")}}, {Name: "other", LineThreshold: 2}},
		Default:    "other",
		Ratios:     []RatioInput{{Category: "source", Reference: "docs", Maximum: "0.5"}},
	})
	alternateResult, err := NewAccountingResultV1(alternatePolicy, alternateInventory, NewCommittedSnapshot(
		"cccccccccccccccccccccccccccccccccccccccc", "dddddddddddddddddddddddddddddddddddddddd", snapshot.Entries()))
	if err != nil {
		t.Fatal(err)
	}
	return ResultV1ReportBinding{Policy: policy, Inventory: inventory, Result: result}, ResultV1ReportBinding{Policy: alternatePolicy, Inventory: alternateInventory, Result: alternateResult}
}

func assertZeroEvidence(t *testing.T, got Evidence) {
	t.Helper()
	if got.CanonicalBytes() != nil || got.Digest() != "" || got.Provenance() != (Provenance{}) {
		t.Fatalf("nonzero evidence = %#v", got)
	}
}
