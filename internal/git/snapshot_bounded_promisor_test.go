package gitadapter

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestBoundedSnapshotQualifiedLocalNoFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("uses local-only partial clones")
	}
	origin := repository(t, false)
	write(t, origin, "file", "base content\n")
	base := commit(t, origin)
	blob := git(t, origin, "rev-parse", base+":file")
	write(t, origin, "file", "head content\n")
	head := commit(t, origin)
	git(t, origin, "config", "uploadpack.allowFilter", "true")
	clones, markers := []string{}, []string{}
	for _, name := range []string{"control", "candidate"} {
		root := t.TempDir()
		repo := filepath.Join(root, name)
		git(t, root, "clone", "--filter=blob:none", "--no-checkout", "--", "file://"+origin, repo)
		for _, id := range []string{base, head} {
			if got := git(t, repo, "cat-file", "-t", id); got != "commit" {
				t.Fatalf("commit unavailable: %q", got)
			}
		}
		// --missing=print traverses without materializing missing promised blobs.
		missing := git(t, repo, "rev-list", "--objects", "--missing=print", base, head)
		if !strings.Contains("\n"+missing+"\n", "\n?"+blob+"\n") {
			t.Fatalf("known blob is not absent: %q", missing)
		}
		trap := filepath.Join(root, "upload-pack")
		if err := os.WriteFile(trap, []byte("#!/bin/sh\n: > \"$0.marker\"\nexit 1\n"), 0700); err != nil {
			t.Fatal(err)
		}
		git(t, repo, "config", "remote.origin.uploadpack", trap)
		clones, markers = append(clones, repo), append(markers, trap+".marker")
	}
	// The unsuppressed control is allowed to attempt local transport and mutate
	// its own state; this is not part of the candidate no-write assertion.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	control := exec.CommandContext(ctx, "git", "-C", clones[0], "diff", "--numstat", base, head, "--")
	control.Env = []string{"PATH=/usr/bin:/bin", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_COUNT=0", "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0", "LANG=C", "LC_ALL=C"}
	output, controlErr := control.CombinedOutput()
	if controlErr == nil || ctx.Err() != nil {
		t.Fatalf("control must fail through local trap, not timeout: %v %q", controlErr, output)
	}
	t.Logf("unsuppressed local control: exit=%d combined=%q", control.ProcessState.ExitCode(), output)
	if _, err := os.Stat(markers[0]); err != nil {
		t.Fatalf("positive lazy-fetch attempt control missing: %v", err)
	}
	before := boundedFixtureState(t, clones[1])
	got, err := AcquireBounded(ctx, Request{clones[1], base, head}, snapshotTestBounds())
	// Commits resolve, but numstat cannot read promised content: the existing
	// neutral acquisition contract reports SnapshotRacing, with no partial data.
	assertFailure(t, got, err, evidence.SnapshotRacing)
	if _, err := os.Stat(markers[1]); !os.IsNotExist(err) {
		t.Fatalf("candidate attempted lazy fetch: %v", err)
	}
	if after := boundedFixtureState(t, clones[1]); !reflect.DeepEqual(before, after) {
		t.Fatal("candidate repository changed")
	}
	t.Log("resolved-range missing blob: SnapshotRacing, empty snapshot, candidate trap absent, repository hashes unchanged")
}
