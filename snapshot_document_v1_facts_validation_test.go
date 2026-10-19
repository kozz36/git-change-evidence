package changeevidence_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestSnapshotFactsV1RejectsInvalidFacts(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*evidence.RevisionIdentity, *[]evidence.SnapshotEntryV1)
		field  string
		code   string
	}{
		{"base", func(r *evidence.RevisionIdentity, _ *[]evidence.SnapshotEntryV1) { r.Base = "bad" }, "revisions.base", "invalid_revision"},
		{"head_format", func(r *evidence.RevisionIdentity, _ *[]evidence.SnapshotEntryV1) { r.Head = strings.Repeat("b", 64) }, "revisions.head", "invalid_revision"},
		{"same_revision_nonempty", func(r *evidence.RevisionIdentity, _ *[]evidence.SnapshotEntryV1) { r.Head = r.Base }, "entries", "inconsistent_range"},
		{"status", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Status = "Q" }, "entries.status", "invalid_status"},
		{"added_old", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Old = (*es)[0].New }, "entries.side", "invalid_side"},
		{"missing_new", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New = nil }, "entries.side", "invalid_side"},
		{"both_required", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Status = "M" }, "entries.side", "invalid_side"},
		{"path_empty", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.Path = nil }, "entries.path", "invalid_path"},
		{"path_absolute", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.Path = []byte("/a") }, "entries.path", "invalid_path"},
		{"path_nul", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.Path = []byte{'a', 0} }, "entries.path", "invalid_path"},
		{"path_dot", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Path = []byte("a/../b")
		}, "entries.path", "invalid_path"},
		{"path_trailing", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.Path = []byte("a/") }, "entries.path", "invalid_path"},
		{"mode", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.Mode = "040000" }, "entries.side", "invalid_side"},
		{"kind", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Kind = evidence.GitOther
		}, "entries.side", "invalid_side"},
		{"file_content_required", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.Content = nil }, "entries.side", "invalid_side"},
		{"link_content_forbidden", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Mode, (*es)[0].New.Kind = "120000", evidence.GitSymlink
		}, "entries.side", "invalid_side"},
		{"object_width", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].New.ObjectID = "bad" }, "entries.object_id", "invalid_object"},
		{"object_zero", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.ObjectID = evidence.GitObjectID(strings.Repeat("0", 40))
		}, "entries.object_id", "invalid_object"},
		{"object_upper", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.ObjectID = evidence.GitObjectID(strings.Repeat("C", 40))
		}, "entries.object_id", "invalid_object"},
		{"empty_metrics", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Content.ByteLength = 0
		}, "entries.content", "invalid_content"},
		{"no_lines", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Content.LogicalLines = 0
		}, "entries.content", "invalid_content"},
		{"too_many_lines", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Content.LogicalLines = 3
		}, "entries.content", "invalid_content"},
		{"numstat_kind", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Numstat.Kind = "unknown" }, "entries.numstat", "invalid_numstat"},
		{"decimal", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Numstat.AddedRaw = "01" }, "entries.numstat", "invalid_numstat"},
		{"count", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { *(*es)[0].Numstat.Added = 2 }, "entries.numstat", "invalid_numstat"},
		{"text_null", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Numstat.Deleted = nil }, "entries.numstat", "invalid_numstat"},
		{"binary_counts", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Numstat.Kind = "binary" }, "entries.numstat", "invalid_numstat"},
		{"duplicate_effective_path", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			*es = append(*es, fullSnapshotEntryV1("a"))
		}, "entries.path", "duplicate_path"},
		{"deleted_wrong_side", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) { (*es)[0].Status = "D" }, "entries.side", "invalid_side"},
		{"binary_nonnull_count", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].Numstat = evidence.SnapshotNumstatV1{Kind: "binary", AddedRaw: "-", DeletedRaw: "-", Added: fullSnapshotCountV1(0)}
		}, "entries.numstat", "invalid_numstat"},
		{"path_single_dot", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Path = []byte("a/./b")
		}, "entries.path", "invalid_path"},
		{"empty_with_nul", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0].New.Content = &evidence.SnapshotContentV1{ContainsNUL: true}
		}, "entries.content", "invalid_content"},
		{"base_uppercase", func(r *evidence.RevisionIdentity, _ *[]evidence.SnapshotEntryV1) { r.Base = strings.Repeat("A", 40) }, "revisions.base", "invalid_revision"},
		{"mode_only_nonzero", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			entry := fullSnapshotPairV1("M", "a", "a")
			entry.New.Mode, entry.New.ObjectID = "100755", entry.Old.ObjectID
			(*es)[0] = entry
		}, "entries.numstat", "invalid_numstat"},
		{"modified_kind_change", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			entry := fullSnapshotPairV1("M", "a", "a")
			entry.New.Mode, entry.New.Kind, entry.New.Content = "120000", evidence.GitSymlink, nil
			(*es)[0] = entry
		}, "entries.side", "invalid_side"},
		{"type_without_kind_change", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0] = fullSnapshotPairV1("T", "a", "a")
		}, "entries.side", "invalid_side"},
		{"rename_same_path", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0] = fullSnapshotPairV1("R", "a", "a")
		}, "entries.path_pair", "invalid_path_pair"},
		{"modified_changed_path", func(_ *evidence.RevisionIdentity, es *[]evidence.SnapshotEntryV1) {
			(*es)[0] = fullSnapshotPairV1("M", "a", "b")
		}, "entries.path_pair", "invalid_path_pair"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revisions, entries := fullSnapshotRevisionsV1(), []evidence.SnapshotEntryV1{fullSnapshotEntryV1("a")}
			tt.mutate(&revisions, &entries)
			before := snapshotEntriesJSONV1(t, entries)
			got, err := evidence.NewSnapshotFactsV1(revisions, entries)
			var contract *evidence.ContractError
			if !errors.As(err, &contract) || contract.Field != tt.field || contract.Code != tt.code {
				t.Fatalf("error = %v, want ContractError{%q, %q}", err, tt.field, tt.code)
			}
			if got.Revisions() != (evidence.RevisionIdentity{}) || len(got.Entries()) != 0 || !bytes.Equal(before, snapshotEntriesJSONV1(t, entries)) {
				t.Fatal("rejection exposed partial result or mutated input")
			}
		})
	}
}
