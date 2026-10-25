package changeevidence_test

import (
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestSnapshotDocumentV1RawOrderAndNulls(t *testing.T) {
	raw := documentAddedV1([]byte{0xff})
	raw.New.Content.ContainsNUL = true
	link := documentAddedV1([]byte("a"))
	link.Status, link.Old, link.New = "D", link.New, nil
	link.Old.Mode, link.Old.Kind, link.Old.Content = "120000", evidence.GitSymlink, nil
	link.Numstat = evidence.SnapshotNumstatV1{Kind: "binary", AddedRaw: "-", DeletedRaw: "-"}
	doc, err := evidence.NewSnapshotV1(documentFactsV1(t, []evidence.SnapshotEntryV1{raw, link}))
	if err != nil {
		t.Fatal(err)
	}
	linkJSON := `{"status":"D","old":{"path_b64":"YQ==","mode":"120000","kind":"symlink","object_id":"` + strings.Repeat("c", 40) +
		`","content":null},"new":null,"numstat":{"kind":"binary","added_raw":"-","deleted_raw":"-","added":null,"deleted":null}}`
	rawJSON := `{"status":"A","old":null,"new":{"path_b64":"/w==","mode":"100644","kind":"file","object_id":"` + strings.Repeat("c", 40) +
		`","content":{"byte_length":2,"contains_nul":true,"logical_lines":1}},"numstat":{"kind":"text","added_raw":"1","deleted_raw":"0","added":1,"deleted":0}}`
	assertDocumentBytesV1(t, doc, expectedDocumentV1("["+linkJSON+","+rawJSON+"]"))
	view := doc.Entries()
	if view[0].Old.Content != nil || view[0].New != nil || view[0].Numstat.Added != nil || view[0].Numstat.Deleted != nil || view[1].New.Path[0] != 0xff {
		t.Fatal("null state or raw identity lost")
	}
}

func TestSnapshotDocumentV1Native64AndUint64Maximum(t *testing.T) {
	maximum := ^uint64(0)
	entry := documentAddedV1([]byte("a"))
	entry.New.ObjectID = evidence.GitObjectID(strings.Repeat("c", 64))
	entry.New.Content = &evidence.SnapshotContentV1{ByteLength: maximum, LogicalLines: maximum}
	entry.Numstat.AddedRaw, *entry.Numstat.Added = "18446744073709551615", maximum
	facts, err := evidence.NewSnapshotFactsV1(evidence.RevisionIdentity{Base: strings.Repeat("a", 64), Head: strings.Repeat("b", 64)}, []evidence.SnapshotEntryV1{entry})
	if err != nil {
		t.Fatal(err)
	}
	doc, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"schema":"git-change-evidence.committed-snapshot/v1","revisions":{"base":"` + strings.Repeat("a", 64) +
		`","head":"` + strings.Repeat("b", 64) + `"},"entries":[{"status":"A","old":null,"new":{"path_b64":"YQ==","mode":"100644","kind":"file","object_id":"` + strings.Repeat("c", 64) +
		`","content":{"byte_length":18446744073709551615,"contains_nul":false,"logical_lines":18446744073709551615}},"numstat":{"kind":"text","added_raw":"18446744073709551615","deleted_raw":"0","added":18446744073709551615,"deleted":0}}]}` + "\n"
	assertDocumentBytesV1(t, doc, []byte(expected))
}

func TestSnapshotDocumentV1GitlinkContentNull(t *testing.T) {
	entry := documentAddedV1([]byte("a"))
	entry.New.Mode, entry.New.Kind, entry.New.Content = "160000", evidence.GitGitlink, nil
	doc, err := evidence.NewSnapshotV1(documentFactsV1(t, []evidence.SnapshotEntryV1{entry}))
	if err != nil {
		t.Fatal(err)
	}
	entries := `[{"status":"A","old":null,"new":{"path_b64":"YQ==","mode":"160000","kind":"gitlink","object_id":"` + strings.Repeat("c", 40) +
		`","content":null},"numstat":{"kind":"text","added_raw":"1","deleted_raw":"0","added":1,"deleted":0}}]`
	assertDocumentBytesV1(t, doc, expectedDocumentV1(entries))
}

func TestSnapshotDocumentV1EmptyFileHasExplicitZeros(t *testing.T) {
	entry := documentAddedV1([]byte("a"))
	entry.New.Content = &evidence.SnapshotContentV1{}
	entry.Numstat.AddedRaw, *entry.Numstat.Added = "0", 0
	doc, err := evidence.NewSnapshotV1(documentFactsV1(t, []evidence.SnapshotEntryV1{entry}))
	if err != nil {
		t.Fatal(err)
	}
	entries := `[{"status":"A","old":null,"new":{"path_b64":"YQ==","mode":"100644","kind":"file","object_id":"` + strings.Repeat("c", 40) +
		`","content":{"byte_length":0,"contains_nul":false,"logical_lines":0}},"numstat":{"kind":"text","added_raw":"0","deleted_raw":"0","added":0,"deleted":0}}]`
	assertDocumentBytesV1(t, doc, expectedDocumentV1(entries))
}
