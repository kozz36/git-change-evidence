package changeevidence_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestDecodeSnapshotV1InvalidDocuments(t *testing.T) {
	text := string(snapshotDecodeFixtureV1(t, []evidence.SnapshotEntryV1{documentAddedV1([]byte("a"))}).CanonicalBytes())
	empty := string(snapshotDecodeFixtureV1(t, nil).CanonicalBytes())
	change := func(from, to string) string {
		t.Helper()
		if strings.Count(text, from) != 1 {
			t.Fatalf("ambiguous fixture replacement %q", from)
		}
		return strings.Replace(text, from, to, 1)
	}
	cases := []struct{ name, raw, field, code string }{
		{"empty_input", "", "document", "invalid"},
		{"syntax", "{", "document", "invalid"},
		{"invalid_utf8", string([]byte{0xff}), "document", "invalid"},
		{"unknown_root", change(`"schema":`, `"sensitive_payload":0,"schema":`), "document", "invalid"},
		{"unknown_nested", change(`"base":`, `"extra":0,"base":`), "document", "invalid"},
		{"wrong_schema", change("committed-snapshot/v1", "committed-snapshot/v2"), "schema", "unsupported_version"},
		{"duplicate_root", change(`"schema":`, `"schema":"git-change-evidence.committed-snapshot/v1","schema":`), "document", "noncanonical"},
		{"duplicate_nested", change(`"byte_length":2`, `"byte_length":2,"byte_length":2`), "document", "noncanonical"},
		{"missing_old", change(`"old":null,`, ``), "document", "noncanonical"},
		{"missing_false", change(`"contains_nul":false,`, ``), "document", "noncanonical"},
		{"null_boolean", change(`"contains_nul":false`, `"contains_nul":null`), "document", "noncanonical"},
		{"null_entries", strings.Replace(empty, `"entries":[]`, `"entries":null`, 1), "document", "noncanonical"},
		{"missing_lf", strings.TrimSuffix(text, "\n"), "document", "noncanonical"},
		{"extra_lf", text + "\n", "document", "noncanonical"},
		{"leading_space", " " + text, "document", "noncanonical"},
		{"trailing_json", text + "{}", "document", "invalid"},
		{"field_case", change(`"schema":`, `"Schema":`), "document", "noncanonical"},
		{"base64_unpadded", change("YQ==", "YQ"), "entries.path_b64", "invalid"},
		{"base64_pad_bits", change("YQ==", "YR=="), "entries.path_b64", "invalid"},
		{"base64_linebreak", change("YQ==", `YQ==\r\n`), "entries.path_b64", "invalid"},
		{"escaped_schema", change("git-change-evidence", `git\u002dchange-evidence`), "document", "noncanonical"},
		{"base64_unsafe_path", change("YQ==", "Li4vYQ=="), "entries.path", "invalid_path"},
		{"base64_empty_path", change("YQ==", ""), "entries.path", "invalid_path"},
		{"boolean_uint", change(`"byte_length":2`, `"byte_length":true`), "document", "invalid"},
		{"fraction_uint", change(`"byte_length":2`, `"byte_length":2.0`), "document", "invalid"},
		{"exponent_uint", change(`"byte_length":2`, `"byte_length":2e0`), "document", "invalid"},
		{"negative_uint", change(`"byte_length":2`, `"byte_length":-1`), "document", "invalid"},
		{"overflow_uint", change(`"byte_length":2`, `"byte_length":18446744073709551616`), "document", "invalid"},
		{"null_metric", change(`"byte_length":2`, `"byte_length":null`), "entries.content", "invalid_content"},
		{"null_text_count", change(`"added":1`, `"added":null`), "entries.numstat", "invalid_numstat"},
		{"noncanonical_token", change(`"added_raw":"1"`, `"added_raw":"01"`), "entries.numstat", "invalid_numstat"},
		{"mismatched_count", change(`"added":1`, `"added":2`), "entries.numstat", "invalid_numstat"},
		{"wrong_absence", change(`"status":"A"`, `"status":"D"`), "entries.side", "invalid_side"},
		{"unsupported_kind", change(`"kind":"file"`, `"kind":"other"`), "entries.side", "invalid_side"},
		{"missing_measurement", change(`"content":{"byte_length":2,"contains_nul":false,"logical_lines":1}`, `"content":null`), "entries.side", "invalid_side"},
		{"inconsistent_metric", change(`"logical_lines":1`, `"logical_lines":3`), "entries.content", "invalid_content"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := []byte(tc.raw)
			before := append([]byte(nil), raw...)
			got, err := evidence.DecodeSnapshotV1(raw)
			var contract *evidence.ContractError
			if !errors.As(err, &contract) || contract.Field != tc.field || contract.Code != tc.code || strings.Contains(err.Error(), "sensitive_payload") {
				t.Fatalf("error = %v, want ContractError{%q, %q}", err, tc.field, tc.code)
			}
			if !bytes.Equal(raw, before) || len(got.CanonicalBytes()) != 0 || got.Digest() != "" || got.Revisions() != (evidence.RevisionIdentity{}) || len(got.Entries()) != 0 {
				t.Fatal("rejected input changed or returned partial document")
			}
		})
	}
}

func TestDecodeSnapshotV1MembershipAndOrder(t *testing.T) {
	text := string(snapshotDecodeFixtureV1(t, []evidence.SnapshotEntryV1{documentAddedV1([]byte("a")), documentAddedV1([]byte("b"))}).CanonicalBytes())
	cases := []struct{ name, raw, field, code string }{
		{"duplicate_effective_path", strings.Replace(text, "Yg==", "YQ==", 1), "entries.path", "duplicate_path"},
		{"unsorted", strings.NewReplacer("YQ==", "Yg==", "Yg==", "YQ==").Replace(text), "document", "noncanonical"},
		{"same_revision_nonempty", strings.Replace(text, strings.Repeat("b", 40), strings.Repeat("a", 40), 1), "entries", "inconsistent_range"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := []byte(tc.raw)
			before := append([]byte(nil), raw...)
			if bytes.Equal(raw, []byte(text)) {
				t.Fatal("fixture mutation was not applied")
			}
			got, err := evidence.DecodeSnapshotV1(raw)
			var contract *evidence.ContractError
			if !errors.As(err, &contract) || contract.Field != tc.field || contract.Code != tc.code {
				t.Fatalf("error = %v, want ContractError{%q, %q}", err, tc.field, tc.code)
			}
			if !bytes.Equal(raw, before) || len(got.CanonicalBytes()) != 0 || got.Digest() != "" || got.Revisions() != (evidence.RevisionIdentity{}) || len(got.Entries()) != 0 {
				t.Fatal("rejected range changed input or returned partial evidence")
			}
		})
	}
}
