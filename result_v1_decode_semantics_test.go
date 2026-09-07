package changeevidence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeAccountingResultV1RoundTripsOwnedSelfConsistentEvidence(t *testing.T) {
	policy, inventory, snapshot := resultDecodeFixture(t)
	want, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	raw := want.CanonicalBytes()
	got, err := DecodeAccountingResultV1(raw, policy, inventory)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.CanonicalBytes(), want.CanonicalBytes()) || got.Digest() != want.Digest() || got.AccountingPolicyDigest() != want.AccountingPolicyDigest() || got.InventoryDigest() != want.InventoryDigest() || got.Revisions() != want.Revisions() || !reflect.DeepEqual(got.Entries(), want.Entries()) || !reflect.DeepEqual(got.Totals(), want.Totals()) || !reflect.DeepEqual(got.Observations(), want.Observations()) {
		t.Fatalf("decoded result differs: %#v", got)
	}

	alternate, err := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot([]CommittedChange{{Path: "src/a", Lines: CommittedLineCounts{Additions: 3, Deletions: 1, Countable: true}}, {Path: "docs/a", Lines: CommittedLineCounts{Additions: 4, Countable: true}}, {Path: "assets/a", Lines: CommittedLineCounts{Countable: false}}}))
	if err != nil || len(alternate.CanonicalBytes()) != len(raw) {
		t.Fatalf("alternate = %#v, %v", alternate, err)
	}
	if copy(raw, alternate.CanonicalBytes()) != len(raw) {
		t.Fatal("alternate overwrite was incomplete")
	}
	control, err := DecodeAccountingResultV1(raw, policy, inventory)
	if err != nil || !bytes.Equal(control.CanonicalBytes(), alternate.CanonicalBytes()) || bytes.Equal(control.CanonicalBytes(), got.CanonicalBytes()) || control.Entries()[2].Measurement.Additions != 3 || got.Entries()[2].Measurement.Additions != 2 {
		t.Fatalf("effective mutated-input control = %#v, %v", control, err)
	}
	canonical, entries := got.CanonicalBytes(), got.Entries()
	canonical[0] ^= 1
	entries[0].Category, entries[0].Path[0], entries[0].Measurement.Additions = "changed", 'X', 99
	if !bytes.Equal(got.CanonicalBytes(), want.CanonicalBytes()) || !reflect.DeepEqual(got.Entries(), want.Entries()) || !reflect.DeepEqual(got.Totals(), want.Totals()) || !reflect.DeepEqual(got.Observations(), want.Observations()) {
		t.Fatal("decoded result retained caller input or returned views")
	}
}

func TestDecodeAccountingResultV1RejectsSemanticAndOrderMutations(t *testing.T) {
	policy, inventory, snapshot := resultDecodeFixture(t)
	want, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ name, old, new string }{
		{"entry category", `"path_b64":"YXNzZXRzL2E=","category":"other"`, `"path_b64":"YXNzZXRzL2E=","category":"docs"`},
		{"measurement kind", `"measurement":{"kind":"countable","additions":2,"deletions":1}`, `"measurement":{"kind":"non_countable"}`},
		{"entry additions", `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2`, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":3`},
		{"entry deletions", `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1`, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2,"deletions":2`},
		{"total", `"category":"source","additions":2,"deletions":1,"non_countable":0`, `"category":"source","additions":3,"deletions":1,"non_countable":0`},
		{"threshold limit", `"category":"source","limit":2,"available":true,"actual":3`, `"category":"source","limit":3,"available":true,"actual":3`},
		{"threshold availability", `{"kind":"line_threshold","category":"source","limit":2,"available":true,"actual":3,"exceeded":true}`, `{"kind":"line_threshold","category":"source","limit":2,"available":false,"exceeded":false}`},
		{"threshold actual", `"limit":2,"available":true,"actual":3`, `"limit":2,"available":true,"actual":2`},
		{"threshold exceeded", `"actual":3,"exceeded":true`, `"actual":3,"exceeded":false`},
		{"ratio limit", `"reference":"docs","limit":"0.5","available":true`, `"reference":"docs","limit":"0.6","available":true`},
		{"ratio availability", `{"kind":"ratio","category":"source","reference":"docs","limit":"0.5","available":true,"numerator":3,"denominator":4,"exceeded":true}`, `{"kind":"ratio","category":"source","reference":"docs","limit":"0.5","available":false,"exceeded":false}`},
		{"ratio numerator", `"reference":"docs","limit":"0.5","available":true,"numerator":3`, `"reference":"docs","limit":"0.5","available":true,"numerator":2`},
		{"ratio denominator", `"numerator":3,"denominator":4,"exceeded":true`, `"numerator":3,"denominator":3,"exceeded":true`},
		{"ratio exceeded", `"denominator":4,"exceeded":true`, `"denominator":4,"exceeded":false`},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DecodeAccountingResultV1(resultDecodeReplaceOne(t, want.CanonicalBytes(), test.old, test.new), policy, inventory)
			if err == nil {
				t.Fatal("semantic mutation was accepted")
			}
			assertZeroResultDocument(t, got)
		})
	}
	for _, test := range []struct{ name, first, second string }{
		{"entries", `{"path_b64":"YXNzZXRzL2E=","category":"other","measurement":{"kind":"non_countable"}}`, `{"path_b64":"ZG9jcy9h","category":"docs","measurement":{"kind":"countable","additions":4,"deletions":0}}`},
		{"totals", `{"category":"source","additions":2,"deletions":1,"non_countable":0}`, `{"category":"docs","additions":4,"deletions":0,"non_countable":0}`},
		{"observations", `{"kind":"line_threshold","category":"source","limit":2,"available":true,"actual":3,"exceeded":true}`, `{"kind":"line_threshold","category":"other","limit":1,"available":false,"exceeded":false}`},
	} {
		t.Run(test.name+" order", func(t *testing.T) {
			got, err := DecodeAccountingResultV1(resultDecodeSwapOne(t, want.CanonicalBytes(), test.first, test.second), policy, inventory)
			if err == nil {
				t.Fatal("reordered array was accepted")
			}
			assertZeroResultDocument(t, got)
		})
	}
}

func TestDecodeAccountingResultV1ReturnsZeroOnFatalBoundaries(t *testing.T) {
	policy, inventory, snapshot := resultDecodeFixture(t)
	valid, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	overflow := resultDecodeReplaceOne(t, valid.CanonicalBytes(), `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2`, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":18446744073709551615`)
	overflow = resultDecodeReplaceOne(t, overflow, `"path_b64":"ZG9jcy9h","category":"docs"`, `"path_b64":"c3JjL2I=","category":"source"`)
	invalidSnapshot := resultDecodeReplaceOne(t, valid.CanonicalBytes(), `"base":"`+resultUnit5Base+`"`, `"base":"main"`)
	for _, test := range []struct {
		name, field, code string
		raw               []byte
		p                 PolicyDocumentV1
		i                 InventoryDocumentV1
	}{
		{"materializer", "document", "invalid", []byte{0xff}, policy, inventory},
		{"policy antecedent", "accounting_policy", "invalid_document", valid.CanonicalBytes(), PolicyDocumentV1{}, inventory},
		{"inventory antecedent", "inventory", "invalid_document", valid.CanonicalBytes(), policy, InventoryDocumentV1{}},
		{"snapshot", "snapshot.base", "invalid_revision", invalidSnapshot, policy, inventory},
		{"overflow", "totals", "overflow", overflow, policy, inventory},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DecodeAccountingResultV1(test.raw, test.p, test.i)
			assertResultSnapshotValidationError(t, err, test.field, test.code)
			assertZeroResultDocument(t, got)
		})
	}
}

func resultDecodeFixture(t *testing.T) (PolicyDocumentV1, InventoryDocumentV1, CommittedSnapshot) {
	t.Helper()
	policy, inventory := resultUnit5Antecedents(t, PolicyInput{
		Categories: []CategoryInput{{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 2}, {Name: "docs", PathGlobs: [][]byte{[]byte("docs/**")}}, {Name: "other", LineThreshold: 1}},
		Default:    "other",
		Ratios:     []RatioInput{{Category: "source", Reference: "docs", Maximum: "0.5"}},
	})
	return policy, inventory, resultUnit5Snapshot([]CommittedChange{{Path: "src/a", Lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}}, {Path: "docs/a", Lines: CommittedLineCounts{Additions: 4, Countable: true}}, {Path: "assets/a", Lines: CommittedLineCounts{Countable: false}}})
}

func resultDecodeReplaceOne(t *testing.T, raw []byte, old, new string) []byte {
	t.Helper()
	if count := strings.Count(string(raw), old); count != 1 {
		t.Fatalf("replacement %q occurred %d times, want 1", old, count)
	}
	return []byte(strings.Replace(string(raw), old, new, 1))
}

func resultDecodeSwapOne(t *testing.T, raw []byte, first, second string) []byte {
	t.Helper()
	raw = resultDecodeReplaceOne(t, raw, first, "__RESULT_DECODE_SWAP__")
	raw = resultDecodeReplaceOne(t, raw, second, first)
	return resultDecodeReplaceOne(t, raw, "__RESULT_DECODE_SWAP__", second)
}

func resultDecodeBase64CRLFPremise(t *testing.T, raw []byte) {
	t.Helper()
	var decoded struct {
		Entries []struct {
			Path string `json:"path_b64"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil || len(decoded.Entries) != 1 || decoded.Entries[0].Path != "/w\r\n==" {
		t.Fatalf("Base64 CRLF JSON premise = %#v, %v", decoded, err)
	}
	if value, err := base64.StdEncoding.DecodeString(decoded.Entries[0].Path); err != nil || !bytes.Equal(value, []byte{0xff}) {
		t.Fatalf("standard Base64 CRLF = %x, %v", value, err)
	}
}
