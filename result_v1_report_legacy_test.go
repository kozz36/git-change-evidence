package changeevidence

import (
	"bytes"
	"testing"
)

func TestResultV1ReportBindingPreservesLegacyStandaloneReports(t *testing.T) {
	legacy, err := NewReportV1(validReportInput())
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeCanonical(legacy.CanonicalBytes())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded.CanonicalBytes(), legacy.CanonicalBytes()) || decoded.Digest() != legacy.Digest() || decoded.Provenance() != legacy.Provenance() {
		t.Fatalf("legacy standalone report changed: %#v, %#v", decoded, legacy)
	}
	binding, _ := resultReportBindings(t)
	assertContractErrorValue(t, ValidateReportV1ResultProvenance(decoded, binding), "provenance.accounting_sha256", "mismatched_digest")
}
