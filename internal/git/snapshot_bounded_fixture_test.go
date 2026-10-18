package gitadapter

import (
	"context"
	"crypto/sha256"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestBoundedSnapshotFixtureIsolationAndEntryBudgets(t *testing.T) {
	if testing.Short() {
		t.Skip("uses fixture Git")
	}
	repo := repository(t, false)
	write(t, repo, "old", "rename\n")
	write(t, repo, "binary", "before\x00")
	base := commit(t, repo)
	git(t, repo, "mv", "old", "new")
	write(t, repo, "binary", "after\x00")
	head := commit(t, repo)
	request := Request{repo, base, head}
	want, err := Acquire(request)
	if err != nil {
		t.Fatal(err)
	}
	git(t, repo, "update-ref", "refs/heads/literal;$ref", head)
	request.Head = "literal;$ref"
	// Inherited controls must not reach the new subprocess route.
	t.Setenv("GIT_DIR", "/missing")
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "diff.external")
	t.Setenv("GIT_CONFIG_VALUE_0", "/bin/false")
	before := boundedFixtureState(t, repo)
	for _, limit := range []int{1, 2, 3} {
		bounds := snapshotTestBounds()
		bounds.MaxEntries = limit
		got, err := AcquireBounded(context.Background(), request, bounds)
		if limit == 1 {
			if !errors.Is(err, ErrSnapshotResourceLimit) || !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
				t.Fatalf("limit %d: %#v %v", limit, got, err)
			}
		} else if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("limit %d: %#v %v", limit, got, err)
		}
	}
	if after := boundedFixtureState(t, repo); !reflect.DeepEqual(before, after) {
		t.Fatal("fixture repository/source/index changed")
	}
	for _, tc := range []struct {
		repo, head string
		code       evidence.SnapshotErrorCode
	}{
		{t.TempDir(), head, evidence.SnapshotForeign},
		{repo, "unknown", evidence.SnapshotUnresolved},
		{repo, strings.Repeat("0", 40), evidence.SnapshotMissing},
		{"relative", head, evidence.SnapshotForeign},
	} {
		got, err := AcquireBounded(context.Background(), Request{tc.repo, base, tc.head}, snapshotTestBounds())
		assertFailure(t, got, err, tc.code)
	}
	// A real stdout cap failure is classified, not truncated or echoed.
	got, err := AcquireBounded(context.Background(), Request{repo, base, head}, SnapshotBounds{1, 64, 2})
	if !errors.Is(err, ErrSnapshotResourceLimit) || !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
		t.Fatalf("real cap: %#v %v", got, err)
	}
}

func boundedFixtureState(t *testing.T, root string) map[string][32]byte {
	t.Helper()
	state := map[string][32]byte{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		state[path] = sha256.Sum256(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func TestBoundedSnapshotPreservesOverflowAndBinaryMeasurements(t *testing.T) {
	id := strings.Repeat("c", 40)
	raw := ":000000 100644 " + strings.Repeat("0", 40) + " " + id + " A\x00overflow\x00"
	calls := 0
	runner := func(context.Context, string, int, ...string) ([]byte, error) {
		calls++
		switch calls {
		case 1:
			return []byte("sha1\n"), nil
		case 2, 4:
			return []byte(id + "\n"), nil
		case 6:
			return []byte(raw), nil
		case 7:
			return []byte("18446744073709551616\t0\toverflow\x00"), nil
		}
		return nil, nil
	}
	got, err := AcquireBoundedWithRunner(context.Background(), runner, Request{"/fixture", "base", "head"}, snapshotTestBounds())
	if err != nil || len(got.Entries()) != 1 || got.Entries()[0].Lines.Countable || got.Entries()[0].Binary {
		t.Fatalf("%#v %v", got, err)
	}
}
