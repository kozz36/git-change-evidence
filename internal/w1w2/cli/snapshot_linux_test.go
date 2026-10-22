//go:build linux

package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func snapshotArguments() []string {
	return []string{"snapshot", "--repo", "/fixture", "--base", "--literal-base", "--head", "--literal-head"}
}

func snapshotFacts(t *testing.T) evidence.SnapshotFactsV1 {
	t.Helper()
	facts, err := evidence.NewSnapshotFactsV1(evidence.RevisionIdentity{Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return facts
}

func TestSnapshotCLIIndependentCanonicalOutput(t *testing.T) {
	facts, calls := snapshotFacts(t), 0
	source := func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error) {
		calls++
		return facts, nil
	}
	var out, diagnostic bytes.Buffer
	code := runWithSources(context.Background(), snapshotArguments(), unusedInput{}, &out, &diagnostic, nil, source)
	want := "{\"schema\":\"git-change-evidence.committed-snapshot/v1\",\"revisions\":{\"base\":\"" + strings.Repeat("a", 40) + "\",\"head\":\"" + strings.Repeat("b", 40) + "\"},\"entries\":[]}\n"
	if code != 0 || out.String() != want || diagnostic.Len() != 0 || calls != 1 {
		t.Fatalf("code%d out%q diagnostic%q calls%d", code, out.Bytes(), diagnostic.String(), calls)
	}
}

func TestSnapshotCLINativeCommittedFactsIgnoreDirtyIndexAndWorktree(t *testing.T) {
	root := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		command := exec.Command("git", args...)
		command.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0"}
		raw, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("fixture Git: %v %s", err, raw)
		}
		return strings.TrimSpace(string(raw))
	}
	git("init", root)
	commit := func(message string) {
		git("-C", root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "commit", "--allow-empty", "-m", message)
	}
	commit("base")
	base := git("-C", root, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(root, "a"), []byte("x\n"), 0600); err != nil {
		t.Fatal(err)
	}
	git("-C", root, "add", "a")
	commit("head")
	head, object := git("-C", root, "rev-parse", "HEAD"), git("-C", root, "rev-parse", "HEAD:a")
	if err := os.WriteFile(filepath.Join(root, "a"), []byte("not committed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	git("-C", root, "add", "a")
	one, zero := uint64(1), uint64(0)
	facts, err := evidence.NewSnapshotFactsV1(evidence.RevisionIdentity{Base: base, Head: head}, []evidence.SnapshotEntryV1{{Status: "A", New: &evidence.SnapshotSideV1{Path: []byte("a"), Mode: "100644", Kind: "file", ObjectID: evidence.GitObjectID(object), Content: &evidence.SnapshotContentV1{ByteLength: 2, LogicalLines: 1}}, Numstat: evidence.SnapshotNumstatV1{Kind: "text", AddedRaw: "1", DeletedRaw: "0", Added: &one, Deleted: &zero}}})
	if err != nil {
		t.Fatal(err)
	}
	want, err := evidence.NewSnapshotV1(facts)
	if err != nil {
		t.Fatal(err)
	}
	code, out, diagnostic := invoke(t, []string{"snapshot", "--repo", root, "--base", base, "--head", head})
	if code != 0 || !bytes.Equal(out, want.CanonicalBytes()) || diagnostic != "" {
		t.Fatalf("immutable facts: %d %q %q", code, out, diagnostic)
	}
}

func TestSnapshotCLINativeUnavailableAndForeignFailures(t *testing.T) {
	root := t.TempDir()
	requireFailure(t, []string{"snapshot", "--repo", root, "--base", "base", "--head", "head"}, 2, "w1w2 snapshot: invalid_arguments\n")
	command := exec.Command("git", "init", root)
	command.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0"}
	if raw, err := command.CombinedOutput(); err != nil {
		t.Fatalf("fixture Git: %v %s", err, raw)
	}
	requireFailure(t, []string{"snapshot", "--repo", root, "--base", "missing-base", "--head", "missing-head"}, 3, "w1w2 snapshot: git_unavailable\n")
	t.Setenv("PATH", "")
	requireFailure(t, []string{"snapshot", "--repo", root, "--base", "base", "--head", "head"}, 3, "w1w2 snapshot: git_unavailable\n")
}
