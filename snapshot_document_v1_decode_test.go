package changeevidence_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func snapshotDecodeFixtureV1(t *testing.T, entries []evidence.SnapshotEntryV1) evidence.SnapshotDocumentV1 {
	t.Helper()
	doc, err := evidence.NewSnapshotV1(documentFactsV1(t, entries))
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestDecodeSnapshotV1ValidShapes(t *testing.T) {
	binary := documentAddedV1([]byte("binary"))
	binary.New.Content.ContainsNUL = true
	binary.Numstat = evidence.SnapshotNumstatV1{Kind: "binary", AddedRaw: "-", DeletedRaw: "-"}
	link, gitlink := documentAddedV1([]byte("link")), documentAddedV1([]byte("gitlink"))
	link.New.Mode, link.New.Kind, link.New.Content = "120000", evidence.GitSymlink, nil
	gitlink.New.Mode, gitlink.New.Kind, gitlink.New.Content = "160000", evidence.GitGitlink, nil
	deleted := documentAddedV1([]byte("deleted"))
	deleted.Status, deleted.Old, deleted.New = "D", deleted.New, nil
	deleted.Numstat.AddedRaw, *deleted.Numstat.Added = "0", 0
	typed := fullSnapshotPairV1("T", "a", "a")
	typed.New.Mode, typed.New.Kind, typed.New.Content = "120000", evidence.GitSymlink, nil
	mode := fullSnapshotPairV1("M", "a", "a")
	mode.New.Mode, mode.New.ObjectID = "100755", mode.Old.ObjectID
	mode.Numstat.AddedRaw, *mode.Numstat.Added = "0", 0
	cases := []struct {
		name    string
		entries []evidence.SnapshotEntryV1
	}{
		{"empty", nil},
		{"raw_file", []evidence.SnapshotEntryV1{documentAddedV1([]byte{0xff, '\n'})}},
		{"binary", []evidence.SnapshotEntryV1{binary}},
		{"symlink", []evidence.SnapshotEntryV1{link}},
		{"gitlink", []evidence.SnapshotEntryV1{gitlink}},
		{"deleted", []evidence.SnapshotEntryV1{deleted}},
		{"modified", []evidence.SnapshotEntryV1{fullSnapshotPairV1("M", "a", "a")}},
		{"renamed", []evidence.SnapshotEntryV1{fullSnapshotPairV1("R", "a", "b")}},
		{"copied", []evidence.SnapshotEntryV1{fullSnapshotPairV1("C", "a", "b")}},
		{"type", []evidence.SnapshotEntryV1{typed}},
		{"mode", []evidence.SnapshotEntryV1{mode}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := snapshotDecodeFixtureV1(t, tc.entries)
			raw := expected.CanonicalBytes()
			before := append([]byte(nil), raw...)
			got, err := evidence.DecodeSnapshotV1(raw)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(raw, before) || !bytes.Equal(got.CanonicalBytes(), raw) || got.Digest() != expected.Digest() || got.Revisions() != expected.Revisions() || !reflect.DeepEqual(got.Entries(), expected.Entries()) {
				t.Fatal("decoded identity, facts or input changed")
			}
		})
	}
}

func TestDecodeSnapshotV1CopyIsolation(t *testing.T) {
	expected := snapshotDecodeFixtureV1(t, []evidence.SnapshotEntryV1{fullSnapshotPairV1("M", "a", "a")})
	raw := expected.CanonicalBytes()
	got, err := evidence.DecodeSnapshotV1(raw)
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = '!'
	canonical, entries := got.CanonicalBytes(), got.Entries()
	canonical[0] = '?'
	entries[0].New.Path[0], entries[0].Old.Path[0] = 'x', 'y'
	entries[0].New.Content.ByteLength, entries[0].Old.Content.LogicalLines = 99, 88
	*entries[0].Numstat.Added, *entries[0].Numstat.Deleted = 77, 66
	if !bytes.Equal(got.CanonicalBytes(), expected.CanonicalBytes()) || got.Digest() != expected.Digest() || !reflect.DeepEqual(got.Entries(), expected.Entries()) {
		t.Fatal("decoded input or getter alias")
	}
}

func TestDecodeSnapshotV1Native64AndUint64Maximum(t *testing.T) {
	maximum := ^uint64(0)
	entry := documentAddedV1([]byte("a"))
	entry.New.ObjectID = evidence.GitObjectID(strings.Repeat("c", 64))
	entry.New.Content = &evidence.SnapshotContentV1{ByteLength: maximum, LogicalLines: maximum}
	entry.Numstat.AddedRaw, *entry.Numstat.Added = "18446744073709551615", maximum
	facts, err := evidence.NewSnapshotFactsV1(evidence.RevisionIdentity{Base: strings.Repeat("a", 64), Head: strings.Repeat("b", 64)}, []evidence.SnapshotEntryV1{entry})
	if err != nil {
		t.Fatal(err)
	}
	expected, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	got, err := evidence.DecodeSnapshotV1(expected.CanonicalBytes())
	if err != nil || !bytes.Equal(got.CanonicalBytes(), expected.CanonicalBytes()) || got.Digest() != expected.Digest() || !reflect.DeepEqual(got.Entries(), expected.Entries()) {
		t.Fatalf("native format or uint64 boundary: %v", err)
	}
}
