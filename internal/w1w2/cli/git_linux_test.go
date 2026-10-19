package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func fixtureGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "git", append([]string{"-C", repo}, args...)...)
	command.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0"}
	out, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("fixture Git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestResultUsesEntireCommittedRangeNotInventoryOrWorktree(t *testing.T) {
	repo := t.TempDir()
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.name", "Fixture"}, {"config", "user.email", "fixture@example.invalid"}} {
		fixtureGit(t, repo, args...)
	}
	for _, content := range []string{"old\n", "new\nline\n"} {
		for _, name := range []string{"a", "b"} {
			if err := os.WriteFile(repo+"/"+name, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
		}
		fixtureGit(t, repo, "add", "--", "a", "b")
		fixtureGit(t, repo, "commit", "-qm", content)
	}
	head := fixtureGit(t, repo, "rev-parse", "HEAD")
	base := fixtureGit(t, repo, "rev-parse", "HEAD^")
	policy, _, _ := documents(t)
	inventory, err := evidence.NewInventoryV1(policy, []evidence.InventoryEntryInput{{Path: []byte("a"), Content: []byte("captured view")}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repo+"/b", []byte("dirty worktree"), 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"result", "--repo", repo, "--base", base, "--head", head, "--policy", documentFile(t, policy.CanonicalBytes()), "--inventory", documentFile(t, inventory.CanonicalBytes())}
	code, raw, diagnostic := invoke(t, args)
	snapshot, err := gitadapter.AcquireBounded(context.Background(), gitadapter.Request{Repository: repo, Base: base, Head: head}, gitadapter.SnapshotBounds{MaxResponseBytes: contract.MaxGitResponseBytes, MaxTotalBytes: contract.MaxGitBytes, MaxEntries: contract.MaxDiffEntries})
	if err != nil {
		t.Fatal(err)
	}
	want, err := evidence.NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil || code != 0 || diagnostic != "" || !bytes.Equal(raw, want.CanonicalBytes()) || len(want.Entries()) != 2 {
		t.Fatalf("code=%d stderr=%q entries=%d error=%v", code, diagnostic, len(want.Entries()), err)
	}
	if data, err := os.ReadFile(repo + "/b"); err != nil || string(data) != "dirty worktree" {
		t.Fatal("worktree changed", err)
	}
}
