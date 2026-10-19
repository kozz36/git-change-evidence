package gitadapter

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestBoundedSnapshotRefMovementAfterResolution(t *testing.T) {
	if testing.Short() {
		t.Skip("uses fixture Git")
	}
	repo := repository(t, false)
	write(t, repo, "file", "before\n")
	base := commit(t, repo)
	write(t, repo, "file", "after\n")
	head := commit(t, repo)
	request := Request{repo, base, "HEAD"}
	want, err := Acquire(request)
	if err != nil {
		t.Fatal(err)
	}
	moved := false
	runner := func(ctx context.Context, root string, cap int, args ...string) ([]byte, error) {
		if args[0] == "diff" && !moved {
			git(t, repo, "update-ref", "HEAD", base)
			moved = true
		}
		return snapshotCommand(ctx, root, cap, args...)
	}
	got, err := AcquireBoundedWithRunner(context.Background(), runner, request, snapshotTestBounds())
	if err != nil || !moved || got.Head != evidence.GitObjectID(head) || !reflect.DeepEqual(got, want) {
		t.Fatalf("%#v %v moved=%v", got, err, moved)
	}
}

func TestBoundedSnapshotDoesNotFetchMissingPromisorObject(t *testing.T) {
	if testing.Short() {
		t.Skip("uses fixture Git")
	}
	repo := repository(t, false)
	write(t, repo, "file", "present\n")
	base := commit(t, repo)
	// A local-only transport trap records attempted lazy fetch, without a network
	// endpoint. Its upload-pack replacement can only write in this owned fixture.
	trapRoot := t.TempDir()
	marker := filepath.Join(trapRoot, "fetch-attempt")
	trap := filepath.Join(trapRoot, "upload-pack")
	script := "#!/bin/sh\nprintf attempted > '" + marker + "'\nexit 1\n"
	if err := os.WriteFile(trap, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "config", "remote.origin.url", repo)
	git(t, repo, "config", "remote.origin.promisor", "true")
	git(t, repo, "config", "remote.origin.partialclonefilter", "blob:none")
	git(t, repo, "config", "remote.origin.uploadpack", trap)
	got, err := AcquireBounded(context.Background(), Request{repo, base, strings.Repeat("0", 40)}, snapshotTestBounds())
	assertFailure(t, got, err, evidence.SnapshotMissing)
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("implicit fetch trap: %v", err)
	}
}
