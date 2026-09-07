package changeevidence

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestResultV1DecodeValuesMaterializesOwnedCandidate(t *testing.T) {
	raw := resultDecodeValuesDocument(`{"path_b64":"YQ==","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1}},{"path_b64":"/wph","category":"source","measurement":{"kind":"non_countable"}}`, resultUnit5Base, resultUnit5Head)
	candidate, err := materializeResultDecodeCandidate(raw)
	if err != nil {
		t.Fatal(err)
	}
	entries := candidate.snapshot.Entries()
	if candidate.policyDigest != DocumentDigest(strings.Repeat("1", 64)) || candidate.inventoryDigest != DocumentDigest(strings.Repeat("2", 64)) || candidate.snapshot.Base != resultUnit5Base || candidate.snapshot.Head != resultUnit5Head || len(entries) != 2 || entries[0].Lines != (CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}) || entries[1].Lines.Countable || !bytes.Equal([]byte(entries[1].Path), []byte{0xff, '\n', 'a'}) {
		t.Fatalf("candidate = %#v, entries = %#v", candidate, entries)
	}
	alternate := resultDecodeValuesDocument(`{"path_b64":"Yg==","category":"source","measurement":{"kind":"countable","additions":3,"deletions":4}},{"path_b64":"/wpp","category":"source","measurement":{"kind":"non_countable"}}`, strings.Repeat("5", len(resultUnit5Base)), strings.Repeat("6", len(resultUnit5Head)))
	alternate = bytes.ReplaceAll(alternate, []byte(strings.Repeat("1", 64)), []byte(strings.Repeat("3", 64)))
	alternate = bytes.ReplaceAll(alternate, []byte(strings.Repeat("2", 64)), []byte(strings.Repeat("4", 64)))
	if len(alternate) != len(raw) || bytes.Equal(raw, alternate) {
		t.Fatal("alternate raw fixture must be distinct and equal length")
	}
	if copied := copy(raw, alternate); copied != len(raw) || !bytes.Equal(raw, alternate) {
		t.Fatalf("alternate raw fixture copied = %d, want %d", copied, len(raw))
	}
	control, err := materializeResultDecodeCandidate(raw)
	if err != nil {
		t.Fatalf("mutated raw control: %v", err)
	}
	controlEntries := control.snapshot.Entries()
	if control.policyDigest != DocumentDigest(strings.Repeat("3", 64)) || control.inventoryDigest != DocumentDigest(strings.Repeat("4", 64)) || control.snapshot.Base != GitObjectID(strings.Repeat("5", len(resultUnit5Base))) || control.snapshot.Head != GitObjectID(strings.Repeat("6", len(resultUnit5Head))) || len(controlEntries) != 2 || controlEntries[0].Path != "b" || controlEntries[0].Lines != (CommittedLineCounts{Additions: 3, Deletions: 4, Countable: true}) || controlEntries[1].Lines.Countable || !bytes.Equal([]byte(controlEntries[1].Path), []byte{0xff, '\n', 'i'}) {
		t.Fatalf("mutated raw control = %#v, entries = %#v", control, controlEntries)
	}
	entries[0].Path, entries[0].Lines.Additions = "changed", 99
	later := candidate.snapshot.Entries()
	if candidate.policyDigest != DocumentDigest(strings.Repeat("1", 64)) || candidate.inventoryDigest != DocumentDigest(strings.Repeat("2", 64)) || candidate.snapshot.Base != resultUnit5Base || candidate.snapshot.Head != resultUnit5Head || len(later) != 2 || later[0].Path != "a" || later[0].Lines != (CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}) || later[1].Lines.Countable || !bytes.Equal([]byte(later[1].Path), []byte{0xff, '\n', 'a'}) {
		t.Fatalf("candidate retained mutable input or returned entry: %#v", candidate)
	}
}

func TestResultV1DecodeValuesRejectsLexicalAndStructuralCandidates(t *testing.T) {
	entries := `{"path_b64":"YQ==","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1}}`
	raw := resultDecodeValuesDocument(entries, resultUnit5Base, resultUnit5Head)
	invalidUTF8 := append([]byte(nil), raw...)
	invalidUTF8[0] = 0xff
	for _, test := range []struct{ name, old, new, field, code string }{
		{"invalid UTF-8", "", "", "document", "invalid"},
		{"unsupported schema", `"git-change-evidence.accounting-result/v1"`, `"other"`, "schema", "unsupported_version"},
		{"uppercase policy digest", strings.Repeat("1", 64), strings.Repeat("A", 64), "accounting_policy_sha256", "invalid_digest"},
		{"short inventory digest", strings.Repeat("2", 64), strings.Repeat("2", 63), "inventory_sha256", "invalid_digest"},
		{"URL Base64", "YQ==", "_Q==", "entries[0].path_b64", "invalid_base64"},
		{"unpadded Base64", "YQ==", "YQ", "entries[0].path_b64", "invalid_base64"},
		{"CRLF Base64", "YQ==", `YQ\r\n==`, "entries[0].path_b64", "invalid_base64"},
		{"nonzero Base64 pad bits", "YQ==", "YR==", "entries[0].path_b64", "invalid_base64"},
		{"duplicate entry key", `"path_b64":"YQ=="`, `"path_b64":"YQ==","path_b64":"Yg=="`, "entries[0].path_b64", "duplicate_field"},
		{"authority entry key", `"deletions":1}}`, `"deletions":1},"approval":true}`, "entries[0].approval", "authority_field"},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidateRaw := raw
			if test.old == "" {
				candidateRaw = invalidUTF8
			} else {
				candidateRaw = resultDecodeValuesReplace(t, raw, test.old, test.new)
			}
			candidate, err := materializeResultDecodeCandidate(candidateRaw)
			assertResultDecodeValuesError(t, err, test.field, test.code)
			if candidate.policyDigest != "" || candidate.inventoryDigest != "" || candidate.snapshot.Base != "" || candidate.snapshot.Head != "" || candidate.snapshot.Entries() != nil {
				t.Fatalf("failure returned candidate = %#v", candidate)
			}
		})
	}
}

func TestResultV1DecodeValuesAcceptsLexicallyCanonicalEmptyPath(t *testing.T) {
	candidate, err := materializeResultDecodeCandidate(resultDecodeValuesDocument(`{"path_b64":"","category":"source","measurement":{"kind":"non_countable"}}`, resultUnit5Base, resultUnit5Head))
	if err != nil || candidate.snapshot.Entries()[0].Path != "" {
		t.Fatalf("empty Base64 candidate = %#v, %v", candidate, err)
	}
}

func TestResultV1DecodeValuesDefersSnapshotSemanticsToConstructor(t *testing.T) {
	policy, inventory := resultUnit5Antecedents(t, resultUnit5Policy())
	valid, err := materializeResultDecodeCandidate(resultDecodeValuesDocument(`{"path_b64":"YQ==","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1}}`, resultUnit5Base, resultUnit5Head))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewAccountingResultV1(policy, inventory, valid.snapshot); err != nil {
		t.Fatalf("valid composed construction: %v", err)
	}
	for _, test := range []struct{ name, entries, base, head, field, code string }{
		{"invalid decoded path", `{"path_b64":"L2Ficw==","category":"source","measurement":{"kind":"non_countable"}}`, resultUnit5Base, resultUnit5Head, "entries.path_b64", "invalid_path"},
		{"decoded duplicate paths", `{"path_b64":"YQ==","category":"source","measurement":{"kind":"non_countable"}},{"path_b64":"YQ==","category":"source","measurement":{"kind":"non_countable"}}`, resultUnit5Base, resultUnit5Head, "entries.path_b64", "duplicate_path"},
		{"symbolic base", `{"path_b64":"YQ==","category":"source","measurement":{"kind":"non_countable"}}`, "main", resultUnit5Head, "snapshot.base", "invalid_revision"},
		{"mixed revision widths", `{"path_b64":"YQ==","category":"source","measurement":{"kind":"non_countable"}}`, resultUnit5Base, strings.Repeat("b", 64), "snapshot.head", "invalid_revision"},
		{"equal revisions with entries", `{"path_b64":"YQ==","category":"source","measurement":{"kind":"non_countable"}}`, resultUnit5Base, resultUnit5Base, "snapshot.head", "invalid_revision"},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate, err := materializeResultDecodeCandidate(resultDecodeValuesDocument(test.entries, test.base, test.head))
			if err != nil {
				t.Fatalf("materialization rejected delegated semantics: %v", err)
			}
			got, err := NewAccountingResultV1(policy, inventory, candidate.snapshot)
			assertResultSnapshotValidationError(t, err, test.field, test.code)
			assertZeroResultDocument(t, got)
		})
	}
}

func resultDecodeValuesDocument(entries, base, head string) []byte {
	return []byte(`{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"` + strings.Repeat("1", 64) + `","inventory_sha256":"` + strings.Repeat("2", 64) + `","revisions":{"base":"` + base + `","head":"` + head + `"},"entries":[` + entries + `],"totals":[],"observations":[]}`)
}

func resultDecodeValuesReplace(t *testing.T, raw []byte, old, new string) []byte {
	t.Helper()
	updated := strings.Replace(string(raw), old, new, 1)
	if updated == string(raw) {
		t.Fatalf("fixture mutation %q was not applied", old)
	}
	return []byte(updated)
}

func assertResultDecodeValuesError(t *testing.T, err error, field, code string) {
	t.Helper()
	var contract *ContractError
	if !errors.As(err, &contract) || contract.Field != field || contract.Code != code {
		t.Fatalf("error = %#v, want ContractError{%q, %q}", err, field, code)
	}
}
