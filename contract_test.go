package changeevidence

import (
	"bytes"
	"testing"
)

func validReportInput() ReportInput {
	return ReportInput{Subject: "frozen-evidence", Provenance: Provenance{
		Revisions:         RevisionIdentity{Base: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Head: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
		InputPolicyDigest: Digest("cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"),
		InventoryDigest:   Digest("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"),
		AccountingDigest:  Digest("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
	}}
}

func TestNewReportV1CanonicalBytesAndDigest(t *testing.T) {
	evidence, err := NewReportV1(validReportInput())
	if err != nil {
		t.Fatalf("NewReportV1() error = %v", err)
	}
	want := []byte("{\"kind\":\"report\",\"provenance\":{\"accounting_sha256\":\"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee\",\"input_policy_sha256\":\"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\",\"inventory_sha256\":\"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd\",\"revisions\":{\"base\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"head\":\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"}},\"schema\":\"git-change-evidence.evidence/v1\",\"subject\":\"frozen-evidence\"}\n")
	if got := evidence.CanonicalBytes(); !bytes.Equal(got, want) {
		t.Fatalf("CanonicalBytes() = %q, want %q", got, want)
	}
	if got, want := evidence.Digest(), Digest("b7145afd882b9cb401fb7d61873b6ea748d816fdbc385eea1f732bece0fc43d0"); got != want {
		t.Fatalf("Digest() = %q, want %q", got, want)
	}
}
