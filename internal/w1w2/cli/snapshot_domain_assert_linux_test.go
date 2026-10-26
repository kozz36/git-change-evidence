package cli

import (
	"bytes"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func domainContent(data string) *evidence.SnapshotContentV1 {
	// Literal oracle: deliberately not the production measurement function.
	metrics := map[string]evidence.SnapshotContentV1{
		"deletion distinct payload\n":        {ByteLength: 26, LogicalLines: 1},
		"old\n":                              {ByteLength: 4, LogicalLines: 1},
		"new\ntail":                          {ByteLength: 8, LogicalLines: 2},
		"uniquely retained rename payload\n": {ByteLength: 33, LogicalLines: 1},
		"mode-only\n":                        {ByteLength: 10, LogicalLines: 1},
		"regular before link\n":              {ByteLength: 20, LogicalLines: 1},
		"regular after link\n":               {ByteLength: 19, LogicalLines: 1},
		"old\x00\n":                          {ByteLength: 5, ContainsNUL: true, LogicalLines: 1},
		"new\x00tail":                        {ByteLength: 8, ContainsNUL: true, LogicalLines: 1},
		"":                                   {}, "\n": {ByteLength: 1, LogicalLines: 1},
		"tail":  {ByteLength: 4, LogicalLines: 1},
		"a\r\n": {ByteLength: 3, LogicalLines: 1},
	}
	m, ok := metrics[data]
	if !ok {
		panic("missing literal metric oracle")
	}
	return &m
}

func domainSide(t *testing.T, f domainFixture, rev, path string, files []domainFile) *evidence.SnapshotSideV1 {
	t.Helper()
	for _, file := range files {
		if file.path != path {
			continue
		}
		side := &evidence.SnapshotSideV1{Path: []byte(path), Mode: file.mode, Kind: evidence.GitEntryKind("file"), ObjectID: evidence.GitObjectID(domainGit(t, f.repo, "rev-parse", rev+":"+path))}
		if file.mode == "120000" {
			side.Kind = evidence.GitEntryKind("symlink")
		} else {
			side.Content = domainContent(file.data)
		}
		return side
	}
	t.Fatalf("missing fixture path %q", path)
	return nil
}

func domainNumstat(added, deleted uint64) evidence.SnapshotNumstatV1 {
	tokens := []string{"0", "1", "2"}
	return evidence.SnapshotNumstatV1{Kind: "text", AddedRaw: tokens[added], DeletedRaw: tokens[deleted], Added: &added, Deleted: &deleted}
}

func domainExpected(t *testing.T, f domainFixture) []evidence.SnapshotEntryV1 {
	t.Helper()
	type vector struct {
		status, old, new string
		added, deleted   uint64
	}
	vectors := []vector{
		{"D", "deleted", "", 0, 1}, {"M", "modified", "modified", 2, 1},
		{"R", "rename-old", "rename-new", 0, 0}, {"M", "mode", "mode", 0, 0},
		{"T", "to-link", "to-link", 1, 1}, {"T", "to-file", "to-file", 1, 1},
		{"M", "binary", "binary", 0, 0}, {"A", "", "empty", 0, 0},
		{"A", "", "tab\tname", 1, 0}, {"A", "", "line\nname", 1, 0},
		{"A", "", "raw\xff", 1, 0},
	}
	var entries []evidence.SnapshotEntryV1
	for _, v := range vectors {
		e := evidence.SnapshotEntryV1{Status: v.status, Numstat: domainNumstat(v.added, v.deleted)}
		if v.old != "" {
			e.Old = domainSide(t, f, f.base, v.old, f.before)
		}
		if v.new != "" {
			e.New = domainSide(t, f, f.head, v.new, f.after)
		}
		if v.new == "binary" {
			e.Numstat = evidence.SnapshotNumstatV1{Kind: "binary", AddedRaw: "-", DeletedRaw: "-"}
		}
		entries = append(entries, e)
	}
	return append(entries, evidence.SnapshotEntryV1{Status: "A", New: &evidence.SnapshotSideV1{Path: []byte("gitlink"), Mode: "160000", Kind: evidence.GitEntryKind("gitlink"), ObjectID: evidence.GitObjectID(f.base)}, Numstat: domainNumstat(1, 0)})
}

func domainAssert(t *testing.T, f domainFixture, base, head string, entries []evidence.SnapshotEntryV1) []byte {
	t.Helper()
	facts, err := evidence.NewSnapshotFactsV1(evidence.RevisionIdentity{Base: base, Head: head}, entries)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	code, out, diagnostic := invoke(t, []string{"snapshot", "--repo", f.repo, "--base", base, "--head", head})
	if code != 0 || diagnostic != "" {
		t.Fatalf("code=%d stderr=%q stdout=%q", code, diagnostic, out)
	}
	if !bytes.Equal(out, expected.CanonicalBytes()) {
		t.Fatalf("complete canonical bytes mismatch\ngot %s\nwant %s", out, expected.CanonicalBytes())
	}
	decoded, err := evidence.DecodeSnapshotV1(out)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded.Entries(), facts.Entries()) || decoded.Revisions() != facts.Revisions() || !bytes.Equal(decoded.CanonicalBytes(), out) {
		t.Fatal("public decoder roundtrip mismatch")
	}
	return out
}
