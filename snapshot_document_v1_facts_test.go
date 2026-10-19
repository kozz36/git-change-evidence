package changeevidence_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func fullSnapshotCountV1(value uint64) *uint64 { return &value }

func fullSnapshotRevisionsV1() evidence.RevisionIdentity {
	return evidence.RevisionIdentity{Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40)}
}

func fullSnapshotEntryV1(path string) evidence.SnapshotEntryV1 {
	return evidence.SnapshotEntryV1{
		Status: "A",
		New: &evidence.SnapshotSideV1{Path: []byte(path), Mode: "100644", Kind: evidence.GitFile,
			ObjectID: evidence.GitObjectID(strings.Repeat("c", 40)),
			Content:  &evidence.SnapshotContentV1{ByteLength: 2, LogicalLines: 1}},
		Numstat: evidence.SnapshotNumstatV1{Kind: "text", AddedRaw: "1", DeletedRaw: "0",
			Added: fullSnapshotCountV1(1), Deleted: fullSnapshotCountV1(0)},
	}
}

func fullSnapshotPairV1(status, oldPath, newPath string) evidence.SnapshotEntryV1 {
	entry := fullSnapshotEntryV1(newPath)
	entry.Status = status
	entry.Old = fullSnapshotEntryV1(oldPath).New
	entry.Old.ObjectID = evidence.GitObjectID(strings.Repeat("d", 40))
	entry.New.Content = &evidence.SnapshotContentV1{ByteLength: 4, LogicalLines: 2}
	return entry
}

func snapshotEntriesJSONV1(t *testing.T, value []evidence.SnapshotEntryV1) []byte {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestSnapshotFactsV1CopiesAndSorts(t *testing.T) {
	input := []evidence.SnapshotEntryV1{fullSnapshotPairV1("R", "source", "z"), fullSnapshotPairV1("C", "source", "a")}
	before := snapshotEntriesJSONV1(t, input)
	facts, err := evidence.NewSnapshotFactsV1(fullSnapshotRevisionsV1(), input)
	if err != nil {
		t.Fatal(err)
	}
	if facts.Revisions() != fullSnapshotRevisionsV1() || !bytes.Equal(before, snapshotEntriesJSONV1(t, input)) {
		t.Fatal("revision mismatch or constructor modified input")
	}
	want := snapshotEntriesJSONV1(t, []evidence.SnapshotEntryV1{input[1], input[0]})
	if got := snapshotEntriesJSONV1(t, facts.Entries()); !bytes.Equal(got, want) {
		t.Fatalf("entries = %s, want %s", got, want)
	}
	input[0].Old.Path[0] = 'X'
	input[0].New.Content.ByteLength = 999
	*input[0].Numstat.Added = 999
	input[1] = evidence.SnapshotEntryV1{}
	view := facts.Entries()
	view[0].Old.Path[0] = 'Y'
	view[0].New.Path[0] = 'Y'
	view[0].Old.Content.LogicalLines = 99
	view[0].New.Content.ContainsNUL = true
	*view[0].Numstat.Deleted = 99
	view[1] = evidence.SnapshotEntryV1{}
	if got := snapshotEntriesJSONV1(t, facts.Entries()); !bytes.Equal(got, want) {
		t.Fatalf("mutable alias changed facts: %s, want %s", got, want)
	}
}

func TestSnapshotFactsV1EmptyAndNativeIDs(t *testing.T) {
	for _, width := range []int{40, 64} {
		revisions := evidence.RevisionIdentity{Base: strings.Repeat("a", width), Head: strings.Repeat("a", width)}
		empty, err := evidence.NewSnapshotFactsV1(revisions, nil)
		if err != nil || empty.Revisions() != revisions || len(empty.Entries()) != 0 {
			t.Fatalf("empty range = (%#v, %v)", empty, err)
		}
		revisions.Head = strings.Repeat("b", width)
		entry := fullSnapshotEntryV1(string([]byte{0xff, '\\', 'z'}))
		entry.New.ObjectID = evidence.GitObjectID(strings.Repeat("c", width))
		facts, err := evidence.NewSnapshotFactsV1(revisions, []evidence.SnapshotEntryV1{entry})
		if err != nil || facts.Revisions() != revisions || !bytes.Equal(facts.Entries()[0].New.Path, entry.New.Path) {
			t.Fatalf("native/raw identity width%d = (%#v, %v)", width, facts, err)
		}
	}
}

func TestSnapshotFactsV1ValidShapes(t *testing.T) {
	deleted := fullSnapshotEntryV1("deleted")
	deleted.Status, deleted.Old, deleted.New = "D", deleted.New, nil
	deleted.Numstat = evidence.SnapshotNumstatV1{Kind: "text", AddedRaw: "0", DeletedRaw: "1",
		Added: fullSnapshotCountV1(0), Deleted: fullSnapshotCountV1(1)}
	changedType := fullSnapshotPairV1("T", "type", "type")
	changedType.New.Mode, changedType.New.Kind, changedType.New.Content = "120000", evidence.GitSymlink, nil
	mode := fullSnapshotPairV1("M", "mode", "mode")
	mode.New.Mode, mode.New.ObjectID = "100755", mode.Old.ObjectID
	mode.New.Content = &evidence.SnapshotContentV1{ByteLength: 2, LogicalLines: 1}
	mode.Numstat.AddedRaw, *mode.Numstat.Added = "0", 0
	binary := fullSnapshotEntryV1("binary")
	binary.New.Content.ContainsNUL = true
	binary.Numstat = evidence.SnapshotNumstatV1{Kind: "binary", AddedRaw: "-", DeletedRaw: "-"}
	symlink, gitlink := fullSnapshotEntryV1("symlink"), fullSnapshotEntryV1("gitlink")
	symlink.New.Mode, symlink.New.Kind, symlink.New.Content = "120000", evidence.GitSymlink, nil
	gitlink.New.Mode, gitlink.New.Kind, gitlink.New.Content = "160000", evidence.GitGitlink, nil
	nulText := fullSnapshotEntryV1("nul-text")
	nulText.New.Content.ContainsNUL = true
	empty := fullSnapshotEntryV1("empty")
	empty.New.Content = &evidence.SnapshotContentV1{}
	empty.Numstat.AddedRaw, *empty.Numstat.Added = "0", 0
	for _, entry := range []evidence.SnapshotEntryV1{
		fullSnapshotEntryV1("added"), deleted, fullSnapshotPairV1("M", "edit", "edit"),
		fullSnapshotPairV1("R", "old", "renamed"), fullSnapshotPairV1("C", "source", "copy"),
		changedType, mode, binary, symlink, gitlink, nulText, empty,
	} {
		t.Run(entry.Status+"_"+string(effectiveTestPathV1(entry)), func(t *testing.T) {
			facts, err := evidence.NewSnapshotFactsV1(fullSnapshotRevisionsV1(), []evidence.SnapshotEntryV1{entry})
			if err != nil || !bytes.Equal(snapshotEntriesJSONV1(t, facts.Entries()), snapshotEntriesJSONV1(t, []evidence.SnapshotEntryV1{entry})) {
				t.Fatalf("valid shape = (%#v, %v)", facts, err)
			}
		})
	}
}

func TestSnapshotFactsV1RawByteOrder(t *testing.T) {
	input := []evidence.SnapshotEntryV1{fullSnapshotEntryV1(string([]byte{0xff})), fullSnapshotEntryV1("a")}
	facts, err := evidence.NewSnapshotFactsV1(fullSnapshotRevisionsV1(), input)
	if err != nil {
		t.Fatal(err)
	}
	view := facts.Entries()
	if len(view) != 2 || !bytes.Equal(view[0].New.Path, []byte("a")) || !bytes.Equal(view[1].New.Path, []byte{0xff}) {
		t.Fatalf("raw-byte order = %#v", view)
	}
}

func effectiveTestPathV1(entry evidence.SnapshotEntryV1) []byte {
	if entry.New != nil {
		return entry.New.Path
	}
	return entry.Old.Path
}
