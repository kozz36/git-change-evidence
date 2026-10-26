package changeevidence_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestCheckSnapshotResultRangeV1Mismatches(t *testing.T) {
	revisions := documentFactsV1(t, nil).Revisions()
	entry := documentAddedV1([]byte("a"))
	snapshot := snapshotRangeFactsFixtureV1(t, revisions, []evidence.SnapshotEntryV1{entry})
	result := snapshotRangeResultFixtureV1(t, revisions, "a")
	base, head := revisions, revisions
	base.Base, head.Head = strings.Repeat("c", 40), strings.Repeat("c", 40)
	cases := []struct {
		name     string
		snapshot evidence.SnapshotDocumentV1
		result   evidence.ResultDocumentV1
		field    string
		code     string
	}{
		{"zero_snapshot", evidence.SnapshotDocumentV1{}, result, "snapshot", "invalid_document"},
		{"zero_result", snapshot, evidence.ResultDocumentV1{}, "result", "invalid_document"},
		{"different_base", snapshotRangeFactsFixtureV1(t, base, []evidence.SnapshotEntryV1{entry}), result, "revisions", "mismatch"},
		{"different_head", snapshotRangeFactsFixtureV1(t, head, []evidence.SnapshotEntryV1{entry}), result, "revisions", "mismatch"},
		{"missing_path", snapshotRangeFactsFixtureV1(t, revisions, nil), result, "entries", "mismatch"},
		{"extra_path", snapshotRangeFactsFixtureV1(t, revisions, []evidence.SnapshotEntryV1{entry, documentAddedV1([]byte("b"))}), result, "entries", "mismatch"},
		{"same_count_different_path", snapshot, snapshotRangeResultFixtureV1(t, revisions, "b"), "entries", "mismatch"},
		{"case_is_raw", snapshot, snapshotRangeResultFixtureV1(t, revisions, "A"), "entries", "mismatch"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			beforeSnapshot, beforeResult := tc.snapshot.CanonicalBytes(), tc.result.CanonicalBytes()
			snapshotDigest, resultDigest := tc.snapshot.Digest(), tc.result.Digest()
			err := evidence.CheckSnapshotResultRangeV1(tc.snapshot, tc.result)
			var contract *evidence.ContractError
			if !errors.As(err, &contract) || contract.Field != tc.field || contract.Code != tc.code {
				t.Fatalf("error = %v, want ContractError{%q, %q}", err, tc.field, tc.code)
			}
			if !bytes.Equal(tc.snapshot.CanonicalBytes(), beforeSnapshot) || !bytes.Equal(tc.result.CanonicalBytes(), beforeResult) || tc.snapshot.Digest() != snapshotDigest || tc.result.Digest() != resultDigest {
				t.Fatal("mismatch rewrote an antecedent")
			}
		})
	}
}
