package gitadapter_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func metadataRawFixture(status, path, previous string, width int) []byte {
	oldMode, newMode := "100644", "100644"
	oldObject, newObject := strings.Repeat("d", width), strings.Repeat("c", width)
	switch status[0] {
	case 'A':
		oldMode, oldObject = "000000", strings.Repeat("0", width)
	case 'D':
		newMode, newObject = "000000", strings.Repeat("0", width)
	case 'T':
		newMode = "120000"
	}
	header := fmt.Sprintf(":%s %s %s %s %s\x00", oldMode, newMode, oldObject, newObject, status)
	if status[0] == 'R' || status[0] == 'C' {
		return []byte(header + previous + "\x00" + path + "\x00")
	}
	return []byte(header + path + "\x00")
}

func metadataNumstatFixture(added, deleted, path, previous string) []byte {
	if previous != "" {
		return []byte(added + "\t" + deleted + "\t\x00" + previous + "\x00" + path + "\x00")
	}
	return []byte(added + "\t" + deleted + "\t" + path + "\x00")
}

func TestParseSnapshotMetadataPreservesFieldsAndTokens(t *testing.T) {
	cases := []struct {
		status, path, previous string
		binary                 bool
	}{
		{"A", "a", "", false}, {"D", "a", "", false}, {"M", "a", "", false},
		{"R100", "new", "old", false}, {"C075", "new", "old", false}, {"T", "a", "", false},
		{"M", "binary", "", true}, {"R090", "new\t\n\xff", "old\t\n\xff", false},
	}
	for _, tc := range cases {
		t.Run(tc.status+tc.path, func(t *testing.T) {
			added, deleted := "2", "1"
			if tc.binary {
				added, deleted = "-", "-"
			}
			raw, numstat := metadataRawFixture(tc.status, tc.path, tc.previous, 40), metadataNumstatFixture(added, deleted, tc.path, tc.previous)
			rawBefore, numstatBefore := append([]byte(nil), raw...), append([]byte(nil), numstat...)
			got, err := gitadapter.ParseSnapshotMetadata(raw, numstat)
			if err != nil || len(got) != 1 {
				t.Fatalf("metadata = %v, error = %v", got, err)
			}
			want := evidence.CommittedChange{Status: tc.status[:1], Path: tc.path, PreviousPath: tc.previous,
				OldMode: "100644", NewMode: "100644", OldKind: evidence.GitFile, NewKind: evidence.GitFile,
				OldObject: evidence.GitObjectID(strings.Repeat("d", 40)), NewObject: evidence.GitObjectID(strings.Repeat("c", 40)),
				Binary: tc.binary, Lines: evidence.CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}}
			if tc.status == "A" {
				want.OldMode, want.OldKind, want.OldObject = "000000", evidence.GitAbsent, evidence.GitObjectID(strings.Repeat("0", 40))
			}
			if tc.status == "D" {
				want.NewMode, want.NewKind, want.NewObject = "000000", evidence.GitAbsent, evidence.GitObjectID(strings.Repeat("0", 40))
			}
			if tc.status == "T" {
				want.NewMode, want.NewKind = "120000", evidence.GitSymlink
			}
			if tc.binary {
				want.Lines = evidence.CommittedLineCounts{}
			}
			if got[0].Change != want || got[0].Numstat.AddedRaw != added || got[0].Numstat.DeletedRaw != deleted {
				t.Fatalf("fields or tokens lost: %#v", got[0])
			}
			if tc.binary {
				if got[0].Numstat.Kind != "binary" || got[0].Numstat.Added != nil || got[0].Numstat.Deleted != nil {
					t.Fatal("binary tokens became numeric zero")
				}
			} else if got[0].Numstat.Kind != "text" || got[0].Numstat.Added == nil || got[0].Numstat.Deleted == nil || *got[0].Numstat.Added != 2 || *got[0].Numstat.Deleted != 1 {
				t.Fatal("numeric token mismatch")
			}
			if !bytes.Equal(raw, rawBefore) || !bytes.Equal(numstat, numstatBefore) {
				t.Fatal("parser mutated input")
			}
			raw[0], numstat[0] = '!', '!'
			if got[0].Change != want || got[0].Numstat.AddedRaw != added {
				t.Fatal("parsed metadata aliases caller bytes")
			}
		})
	}
}

func TestParseSnapshotMetadataEmptyAndByteExactPairing(t *testing.T) {
	empty, err := gitadapter.ParseSnapshotMetadata(nil, nil)
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty metadata: %v %v", empty, err)
	}
	raw := append(metadataRawFixture("M", "a", "", 40), metadataRawFixture("R100", "b", "old", 40)...)
	numstat := append(metadataNumstatFixture("3", "0", "b", "old"), metadataNumstatFixture("2", "1", "a", "")...)
	got, err := gitadapter.ParseSnapshotMetadata(raw, numstat)
	if err != nil || len(got) != 2 || got[0].Change.Path != "a" || got[0].Numstat.AddedRaw != "2" || got[1].Change.Path != "b" || got[1].Change.PreviousPath != "old" || got[1].Numstat.AddedRaw != "3" {
		t.Fatalf("raw order or exact path-pair association lost: %v %v", got, err)
	}
	*got[0].Numstat.Added = 99
	if *got[1].Numstat.Added != 3 {
		t.Fatal("records share numeric storage")
	}
}

func TestParseSnapshotMetadataNative64AndUint64Maximum(t *testing.T) {
	raw := metadataRawFixture("M100", "a", "", 64)
	numstat := metadataNumstatFixture("18446744073709551615", "0", "a", "")
	beforeRaw, beforeNumstat := append([]byte(nil), raw...), append([]byte(nil), numstat...)
	got, err := gitadapter.ParseSnapshotMetadata(raw, numstat)
	if err != nil || len(got) != 1 || got[0].Change.Status != "M" ||
		string(got[0].Change.OldObject) != strings.Repeat("d", 64) || string(got[0].Change.NewObject) != strings.Repeat("c", 64) ||
		got[0].Numstat.AddedRaw != "18446744073709551615" || got[0].Numstat.DeletedRaw != "0" ||
		got[0].Numstat.Added == nil || *got[0].Numstat.Added != ^uint64(0) || got[0].Numstat.Deleted == nil || *got[0].Numstat.Deleted != 0 ||
		got[0].Change.Lines.Additions != ^uint64(0) || !got[0].Change.Lines.Countable {
		t.Fatalf("native identity/count boundary: %v %v", got, err)
	}
	if !bytes.Equal(raw, beforeRaw) || !bytes.Equal(numstat, beforeNumstat) {
		t.Fatal("boundary input changed")
	}
}

func TestParseSnapshotMetadataSharedOldPaths(t *testing.T) {
	raw := append(metadataRawFixture("C100", "a", "old", 40), metadataRawFixture("C100", "b", "old", 40)...)
	numstat := append(metadataNumstatFixture("3", "0", "b", "old"), metadataNumstatFixture("2", "0", "a", "old")...)
	beforeRaw, beforeNumstat := append([]byte(nil), raw...), append([]byte(nil), numstat...)
	got, err := gitadapter.ParseSnapshotMetadata(raw, numstat)
	if err != nil || len(got) != 2 || got[0].Change.PreviousPath != "old" || got[1].Change.PreviousPath != "old" || got[0].Change.Path != "a" || got[1].Change.Path != "b" || got[0].Numstat.AddedRaw != "2" || got[1].Numstat.AddedRaw != "3" {
		t.Fatalf("shared old path is not duplicate identity: %v %v", got, err)
	}
	if !bytes.Equal(raw, beforeRaw) || !bytes.Equal(numstat, beforeNumstat) {
		t.Fatal("shared-path input changed")
	}
}
