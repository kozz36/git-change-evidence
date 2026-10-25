package changeevidence_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func documentAddedV1(path []byte) evidence.SnapshotEntryV1 {
	added, deleted := uint64(1), uint64(0)
	return evidence.SnapshotEntryV1{
		Status: "A",
		New: &evidence.SnapshotSideV1{
			Path: path, Mode: "100644", Kind: evidence.GitFile,
			ObjectID: evidence.GitObjectID(strings.Repeat("c", 40)),
			Content:  &evidence.SnapshotContentV1{ByteLength: 2, LogicalLines: 1},
		},
		Numstat: evidence.SnapshotNumstatV1{Kind: "text", AddedRaw: "1", DeletedRaw: "0", Added: &added, Deleted: &deleted},
	}
}

func documentFactsV1(t *testing.T, entries []evidence.SnapshotEntryV1) evidence.SnapshotFactsV1 {
	t.Helper()
	facts, err := evidence.NewSnapshotFactsV1(evidence.RevisionIdentity{Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40)}, entries)
	if err != nil {
		t.Fatal(err)
	}
	return facts
}

func expectedDocumentV1(entries string) []byte {
	return []byte(`{"schema":"git-change-evidence.committed-snapshot/v1","revisions":{"base":"` + strings.Repeat("a", 40) +
		`","head":"` + strings.Repeat("b", 40) + `"},"entries":` + entries + "}\n")
}

func assertDocumentBytesV1(t *testing.T, doc evidence.SnapshotDocumentV1, expected []byte) {
	t.Helper()
	if !bytes.Equal(doc.CanonicalBytes(), expected) {
		t.Fatalf("canonical bytes = %q, want %q", doc.CanonicalBytes(), expected)
	}
	sum := sha256.Sum256(expected)
	if string(doc.Digest()) != hex.EncodeToString(sum[:]) {
		t.Fatalf("digest = %q", doc.Digest())
	}
}

func TestSnapshotDocumentV1ApprovedFileGolden(t *testing.T) {
	if evidence.CommittedSnapshotContractV1 != "git-change-evidence.committed-snapshot/v1" {
		t.Fatal("wrong independent schema")
	}
	facts := documentFactsV1(t, []evidence.SnapshotEntryV1{documentAddedV1([]byte("a.txt"))})
	doc, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	entry := `[{"status":"A","old":null,"new":{"path_b64":"YS50eHQ=","mode":"100644","kind":"file","object_id":"` + strings.Repeat("c", 40) +
		`","content":{"byte_length":2,"contains_nul":false,"logical_lines":1}},"numstat":{"kind":"text","added_raw":"1","deleted_raw":"0","added":1,"deleted":0}}]`
	assertDocumentBytesV1(t, doc, expectedDocumentV1(entry))
	if doc.Revisions() != facts.Revisions() || len(doc.Entries()) != 1 {
		t.Fatal("document views lost facts")
	}
}

func TestSnapshotDocumentV1EmptyAndZeroFacts(t *testing.T) {
	doc, err := evidence.NewSnapshotV1(documentFactsV1(t, nil))
	if err != nil {
		t.Fatal(err)
	}
	assertDocumentBytesV1(t, doc, expectedDocumentV1("[]"))
	if len(doc.Entries()) != 0 {
		t.Fatal("nonempty empty range")
	}
	zero, err := evidence.NewSnapshotV1(evidence.SnapshotFactsV1{})
	var contract *evidence.ContractError
	if !errors.As(err, &contract) || contract.Field != "revisions.base" || contract.Code != "invalid_revision" {
		t.Fatalf("zero-facts error = %v", err)
	}
	if len(zero.CanonicalBytes()) != 0 || zero.Digest() != "" || len(zero.Entries()) != 0 || zero.Revisions() != (evidence.RevisionIdentity{}) {
		t.Fatal("rejection returned partial document")
	}
}

func TestSnapshotDocumentV1CopyViews(t *testing.T) {
	entry := documentAddedV1([]byte("a"))
	old := *entry.New
	old.Path = []byte("a")
	old.ObjectID = evidence.GitObjectID(strings.Repeat("d", 40))
	entry.Status, entry.Old = "M", &old
	facts := documentFactsV1(t, []evidence.SnapshotEntryV1{entry})
	doc, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	before, digest := doc.CanonicalBytes(), doc.Digest()
	entry.New.Path[0] = 'x'
	entry.Old.Content.ByteLength = 99
	view := doc.Entries()
	view[0].New.Path[0], view[0].Old.Path[0] = 'y', 'z'
	view[0].New.Content.ByteLength, view[0].Old.Content.LogicalLines = 88, 77
	*view[0].Numstat.Added, *view[0].Numstat.Deleted = 66, 55
	canonical := doc.CanonicalBytes()
	canonical[0] = '!'
	after := doc.Entries()[0]
	if !bytes.Equal(doc.CanonicalBytes(), before) || doc.Digest() != digest ||
		string(after.New.Path) != "a" || string(after.Old.Path) != "a" ||
		after.New.Content.ByteLength != 2 || after.Old.Content.LogicalLines != 1 ||
		*after.Numstat.Added != 1 || *after.Numstat.Deleted != 0 || facts.Entries()[0].New.Path[0] != 'a' {
		t.Fatal("document/caller/getter alias")
	}
}

func TestSnapshotDocumentV1DeterministicOrdering(t *testing.T) {
	a, z := documentAddedV1([]byte("a")), documentAddedV1([]byte("z"))
	first, err := evidence.NewSnapshotV1(documentFactsV1(t, []evidence.SnapshotEntryV1{z, a}))
	if err != nil {
		t.Fatal(err)
	}
	second, err := evidence.NewSnapshotV1(documentFactsV1(t, []evidence.SnapshotEntryV1{a, z}))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first.CanonicalBytes(), second.CanonicalBytes()) || first.Digest() != second.Digest() {
		t.Fatal("equivalent facts changed byte identity")
	}
}
