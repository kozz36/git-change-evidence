package changeevidence_test

import (
	"bytes"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func snapshotRangeResultFixtureV1(t *testing.T, revisions evidence.RevisionIdentity, paths ...string) evidence.ResultDocumentV1 {
	t.Helper()
	policy, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{{Name: "other"}}, Default: "other"})
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := evidence.NewInventoryV1(policy, nil)
	if err != nil {
		t.Fatal(err)
	}
	changes := make([]evidence.CommittedChange, len(paths))
	for i, path := range paths {
		changes[i] = evidence.CommittedChange{Path: path, Lines: evidence.CommittedLineCounts{Additions: 1, Countable: true}}
	}
	committed := evidence.NewCommittedSnapshot(evidence.GitObjectID(revisions.Base), evidence.GitObjectID(revisions.Head), changes)
	result, err := evidence.NewAccountingResultV1(policy, inventory, committed)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func snapshotRangeFactsFixtureV1(t *testing.T, revisions evidence.RevisionIdentity, entries []evidence.SnapshotEntryV1) evidence.SnapshotDocumentV1 {
	t.Helper()
	facts, err := evidence.NewSnapshotFactsV1(revisions, entries)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestCheckSnapshotResultRangeV1MatchesWithoutCountParity(t *testing.T) {
	revisions := documentFactsV1(t, nil).Revisions()
	deleted := documentAddedV1([]byte("a"))
	deleted.Status, deleted.Old, deleted.New = "D", deleted.New, nil
	deleted.Numstat.AddedRaw, *deleted.Numstat.Added = "0", 0
	counted := documentAddedV1([]byte("a"))
	counted.Numstat.AddedRaw, *counted.Numstat.Added = "99", 99
	cases := []struct {
		name    string
		entries []evidence.SnapshotEntryV1
		paths   []string
	}{
		{"empty", nil, nil},
		{"addition", []evidence.SnapshotEntryV1{documentAddedV1([]byte("a"))}, []string{"a"}},
		{"raw", []evidence.SnapshotEntryV1{documentAddedV1([]byte{0xff})}, []string{string([]byte{0xff})}},
		{"deleted_old_path", []evidence.SnapshotEntryV1{deleted}, []string{"a"}},
		{"renamed_new_path", []evidence.SnapshotEntryV1{fullSnapshotPairV1("R", "old", "new")}, []string{"new"}},
		{"different_counts", []evidence.SnapshotEntryV1{counted}, []string{"a"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			snapshot := snapshotRangeFactsFixtureV1(t, revisions, tc.entries)
			result := snapshotRangeResultFixtureV1(t, revisions, tc.paths...)
			beforeSnapshot, beforeResult := snapshot.CanonicalBytes(), result.CanonicalBytes()
			snapshotDigest, resultDigest := snapshot.Digest(), result.Digest()
			if tc.name == "different_counts" && *snapshot.Entries()[0].Numstat.Added == result.Entries()[0].Measurement.Additions {
				t.Fatal("fixture must diverge in numeric observations")
			}
			if err := evidence.CheckSnapshotResultRangeV1(snapshot, result); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(snapshot.CanonicalBytes(), beforeSnapshot) || !bytes.Equal(result.CanonicalBytes(), beforeResult) || snapshot.Digest() != snapshotDigest || result.Digest() != resultDigest {
				t.Fatal("association rewrote an antecedent")
			}
		})
	}
}

func TestCheckSnapshotResultRangeV1NativeAndEffectiveBoundaries(t *testing.T) {
	revisions := documentFactsV1(t, nil).Revisions()
	same := revisions
	same.Head = same.Base
	wide := evidence.RevisionIdentity{Base: strings.Repeat("a", 64), Head: strings.Repeat("b", 64)}
	wideEntry := documentAddedV1([]byte("a"))
	wideEntry.New.ObjectID = evidence.GitObjectID(strings.Repeat("c", 64))
	cases := []struct {
		name      string
		revisions evidence.RevisionIdentity
		entries   []evidence.SnapshotEntryV1
		paths     []string
	}{
		{"shared_old_copy_paths", revisions, []evidence.SnapshotEntryV1{fullSnapshotPairV1("C", "old", "a"), fullSnapshotPairV1("C", "old", "b")}, []string{"b", "a"}},
		{"native64", wide, []evidence.SnapshotEntryV1{wideEntry}, []string{"a"}},
		{"same_revision_empty", same, nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			snapshot := snapshotRangeFactsFixtureV1(t, tc.revisions, tc.entries)
			result := snapshotRangeResultFixtureV1(t, tc.revisions, tc.paths...)
			beforeSnapshot, beforeResult := snapshot.CanonicalBytes(), result.CanonicalBytes()
			snapshotDigest, resultDigest := snapshot.Digest(), result.Digest()
			if err := evidence.CheckSnapshotResultRangeV1(snapshot, result); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(snapshot.CanonicalBytes(), beforeSnapshot) || !bytes.Equal(result.CanonicalBytes(), beforeResult) || snapshot.Digest() != snapshotDigest || result.Digest() != resultDigest {
				t.Fatal("boundary association rewrote evidence")
			}
		})
	}
}
