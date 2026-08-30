package changeevidence

import (
	"bytes"
	"errors"
	"testing"
)

func TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments(t *testing.T) {
	valid := canonicalTestInput(t)
	cases := []struct{ name, old, new, field, code string }{
		{"missing document field", `,"subject":"frozen-evidence"`, "", "document.subject", "required"},
		{"missing provenance field", `"accounting_sha256":"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",`, "", "provenance.accounting_sha256", "required"},
		{"missing revisions field", `"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",`, "", "provenance.revisions.base", "required"},
		{"malformed document value", `"kind":"report"`, `"kind":false`, "document", "invalid_value"},
		{"malformed provenance value", `"provenance":{"accounting_sha256":"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","input_policy_sha256":"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc","inventory_sha256":"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}},"schema`, `"provenance":false,"schema`, "provenance", "invalid_object"},
		{"malformed revisions value", `"revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}},"schema`, `"revisions":false},"schema`, "provenance.revisions", "invalid_object"},
		{"unsupported schema", `"schema":"git-change-evidence.evidence/v1"`, `"schema":"git-change-evidence.evidence/v2"`, "schema", "unsupported_version"},
		{"unsupported kind", `"kind":"report"`, `"kind":"other"`, "kind", "unsupported_variant"},
		{"duplicate document field", `,"subject":"frozen-evidence"}`, `,"subject":"frozen-evidence","subject":"frozen-evidence"}`, "document", "noncanonical"},
		{"duplicate provenance field", `"accounting_sha256":"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","input`, `"accounting_sha256":"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","accounting_sha256":"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","input`, "document", "noncanonical"},
		{"duplicate revisions field", `"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head`, `"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head`, "document", "noncanonical"},
		{"nested provenance authority", `"provenance":{"accounting`, `"provenance":{"approval":true,"accounting`, "provenance.approval", "authority_field"},
		{"nested revisions authority", `"revisions":{"base`, `"revisions":{"release":true,"base`, "provenance.revisions.release", "authority_field"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) { assertContractError(t, replace(t, valid, tt.old, tt.new), tt.field, tt.code) })
	}
	for _, field := range []string{"approval", "denial", "gate", "block", "merge", "deploy", "release", "delivery_authority"} {
		t.Run("document "+field, func(t *testing.T) {
			assertContractError(t, replace(t, valid, `{"kind`, `{"`+field+`":true,"kind`), "document."+field, "authority_field")
		})
	}
	assertContractError(t, replace(t, valid, `{"kind`, `{"unknown":true,"kind`), "document.unknown", "unknown_field")
	assertContractError(t, append([]byte(" "), valid...), "document", "noncanonical")
}

func TestNewReportV1RejectsNonSHA256ProvenanceDigest(t *testing.T) {
	input := validReportInput()
	input.Provenance.InventoryDigest = Digest("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	_, err := NewReportV1(input)
	assertContractErrorValue(t, err, "provenance.inventory_sha256", "invalid_digest")
}

func canonicalTestInput(t *testing.T) []byte {
	t.Helper()
	evidence, err := NewReportV1(validReportInput())
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeCanonical(evidence.CanonicalBytes())
	if err != nil || decoded.Digest() != evidence.Digest() {
		t.Fatalf("DecodeCanonical() = (%v, %v), want matching evidence", decoded, err)
	}
	return evidence.CanonicalBytes()
}

func replace(t *testing.T, input []byte, old, new string) []byte {
	t.Helper()
	output := bytes.Replace(input, []byte(old), []byte(new), 1)
	if bytes.Equal(input, output) {
		t.Fatalf("missing %q", old)
	}
	return output
}

func assertContractError(t *testing.T, input []byte, field, code string) {
	t.Helper()
	_, err := DecodeCanonical(input)
	assertContractErrorValue(t, err, field, code)
}

func assertContractErrorValue(t *testing.T, err error, field, code string) {
	t.Helper()
	var contractError *ContractError
	if !errors.As(err, &contractError) || contractError.Field != field || contractError.Code != code {
		t.Fatalf("error = %v, want %s: %s", err, field, code)
	}
}
