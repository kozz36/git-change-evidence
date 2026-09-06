package changeevidence

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths(t *testing.T) {
	sha1, sha256 := strings.Repeat("a", 40), strings.Repeat("b", 64)
	metadata := CommittedChange{Path: string([]byte{0xff, '\n', '\r', 'a'}), PreviousPath: "/old", Status: "?", OldObject: "invalid", NewObject: "invalid", OldMode: "invalid", NewMode: "invalid", OldKind: "invalid", NewKind: "invalid", Binary: true, Lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}}
	identity, entries, err := validateResultSnapshot(NewCommittedSnapshot(GitObjectID(sha256), GitObjectID(strings.Repeat("c", 64)), []CommittedChange{metadata}))
	if err != nil || identity.Base != sha256 || !bytes.Equal(entries[0].path, []byte{0xff, '\n', '\r', 'a'}) || entries[0].lines != metadata.Lines {
		t.Fatalf("valid snapshot = (%#v, %#v, %v)", identity, entries, err)
	}
	for _, test := range []struct {
		name, base, head, path, field, code string
	}{
		{"sha1", sha1, strings.Repeat("d", 40), "a", "", ""},
		{"empty equal revisions", sha1, sha1, "", "", ""},
		{"empty base", "", sha1, "a", "snapshot.base", "invalid_revision"},
		{"abbreviated base", sha1[:39], sha1, "a", "snapshot.base", "invalid_revision"},
		{"symbolic head", sha1, "main", "a", "snapshot.head", "invalid_revision"},
		{"nonhex head", sha1, strings.Repeat("g", 40), "a", "snapshot.head", "invalid_revision"},
		{"uppercase head", sha1, strings.ToUpper(sha1), "a", "snapshot.head", "invalid_revision"},
		{"mixed case head", sha1, strings.Repeat("d", 39) + "E", "a", "snapshot.head", "invalid_revision"},
		{"mixed widths", sha1, sha256, "a", "snapshot.head", "invalid_revision"},
		{"equal revisions with entry", sha1, sha1, "a", "snapshot.head", "invalid_revision"},
		{"empty path", sha1, strings.Repeat("d", 40), "", "entries.path_b64", "invalid_path"},
		{"absolute path", sha1, strings.Repeat("d", 40), "/a", "entries.path_b64", "invalid_path"},
		{"NUL path", sha1, strings.Repeat("d", 40), "a\x00b", "entries.path_b64", "invalid_path"},
		{"empty component", sha1, strings.Repeat("d", 40), "a//b", "entries.path_b64", "invalid_path"},
		{"dot component", sha1, strings.Repeat("d", 40), "a/./b", "entries.path_b64", "invalid_path"},
		{"dotdot component", sha1, strings.Repeat("d", 40), "a/../b", "entries.path_b64", "invalid_path"},
		{"trailing component", sha1, strings.Repeat("d", 40), "a/", "entries.path_b64", "invalid_path"},
	} {
		t.Run(test.name, func(t *testing.T) {
			changes := []CommittedChange(nil)
			if test.path != "" || test.field != "" {
				changes = []CommittedChange{{Path: test.path, Lines: CommittedLineCounts{Countable: true}}}
			}
			identity, entries, err := validateResultSnapshot(NewCommittedSnapshot(GitObjectID(test.base), GitObjectID(test.head), changes))
			if test.field == "" {
				if err != nil || identity.Base != test.base || len(entries) != len(changes) {
					t.Fatalf("accepted snapshot = (%#v, %#v, %v)", identity, entries, err)
				}
				return
			}
			assertResultSnapshotValidationError(t, err, test.field, test.code)
			if identity != (RevisionIdentity{}) || entries != nil {
				t.Fatalf("rejected snapshot = (%#v, %#v)", identity, entries)
			}
		})
	}
	duplicate := NewCommittedSnapshot(GitObjectID(sha1), GitObjectID(strings.Repeat("d", 40)), []CommittedChange{{Path: "same"}, {Path: "same"}})
	_, _, err = validateResultSnapshot(duplicate)
	assertResultSnapshotValidationError(t, err, "entries.path_b64", "duplicate_path")
	atomic := NewCommittedSnapshot(GitObjectID(sha1), GitObjectID(strings.Repeat("d", 40)), []CommittedChange{{Path: "valid"}, {Path: "a/../b"}})
	identity, entries, err = validateResultSnapshot(atomic)
	assertResultSnapshotValidationError(t, err, "entries.path_b64", "invalid_path")
	if identity != (RevisionIdentity{}) || entries != nil {
		t.Fatalf("late failure returned partial validation: (%#v, %#v)", identity, entries)
	}
	owned := NewCommittedSnapshot(GitObjectID(sha1), GitObjectID(strings.Repeat("d", 40)), []CommittedChange{{Path: "kept", Lines: CommittedLineCounts{Additions: 1, Countable: true}}})
	_, entries, err = validateResultSnapshot(owned)
	entries[0].path[0], entries[0].lines.Additions = 'x', 9
	identity, entries, err = validateResultSnapshot(owned)
	if err != nil || identity.Base != sha1 || string(entries[0].path) != "kept" || entries[0].lines.Additions != 1 {
		t.Fatalf("returned entry mutation affected later validation: (%#v, %#v, %v)", identity, entries, err)
	}
}

func assertResultSnapshotValidationError(t *testing.T, err error, field, code string) {
	t.Helper()
	var contract *ContractError
	if !errors.As(err, &contract) || contract.Field != field || contract.Code != code {
		t.Fatalf("error = %#v, want %s/%s", err, field, code)
	}
}
