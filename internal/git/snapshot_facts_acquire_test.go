package gitadapter_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func metadataAcquireScript(width int, raw, numstat []byte) ([][]byte, [][]string) {
	base, head := strings.Repeat("a", width), strings.Repeat("b", width)
	format := "sha1\n"
	if width == 64 {
		format = "sha256\n"
	}
	outputs := [][]byte{[]byte(format), []byte(base + "\n"), nil, []byte(head + "\n"), nil, raw, numstat, nil, nil}
	commands := [][]string{
		{"rev-parse", "--show-object-format"},
		{"rev-parse", "--verify", "--end-of-options", "base^{commit}"}, {"cat-file", "-e", base + "^{commit}"},
		{"rev-parse", "--verify", "--end-of-options", "head^{commit}"}, {"cat-file", "-e", head + "^{commit}"},
		{"diff", "--raw", "-z", "-M", "--no-abbrev", "--no-ext-diff", "--no-textconv", "--no-relative", base, head, "--"},
		{"diff", "--numstat", "-z", "-M", "--no-ext-diff", "--no-textconv", "--no-relative", base, head, "--"},
		{"cat-file", "-e", base + "^{commit}"}, {"cat-file", "-e", head + "^{commit}"},
	}
	return outputs, commands
}

func metadataAcquireBounds() gitadapter.SnapshotBounds {
	return gitadapter.SnapshotBounds{MaxResponseBytes: 16 << 20, MaxTotalBytes: 64 << 20, MaxEntries: 10000}
}

func TestAcquireSnapshotMetadataBoundedSequenceAndCopies(t *testing.T) {
	for _, width := range []int{40, 64} {
		t.Run(strings.Repeat("x", width), func(t *testing.T) {
			outputs, commands := metadataAcquireScript(width, metadataRawFixture("M", "a", "", width), metadataNumstatFixture("2", "1", "a", ""))
			baseRef := "base"
			if width == 64 {
				baseRef = "--literal-ref"
			}
			commands[1][3] = baseRef + "^{commit}"
			bounds, calls := metadataAcquireBounds(), 0
			remaining := bounds.MaxTotalBytes
			ctx := context.Background()
			runner := func(gotCtx context.Context, repository string, cap int, args ...string) ([]byte, error) {
				if gotCtx != ctx || repository != "/fixture" || cap != min(bounds.MaxResponseBytes, remaining) || calls >= len(commands) || !reflect.DeepEqual(args, commands[calls]) {
					t.Fatalf("command/context/budget changed at %d: %q cap%d", calls, args, cap)
				}
				out := append([]byte(nil), outputs[calls]...)
				calls++
				remaining -= len(out)
				return out, nil
			}
			got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(ctx, runner, gitadapter.Request{Repository: "/fixture", Base: baseRef, Head: "head"}, bounds)
			if err != nil || got == nil || calls != len(commands) || got.Revisions() != (evidence.RevisionIdentity{Base: strings.Repeat("a", width), Head: strings.Repeat("b", width)}) {
				t.Fatalf("incomplete acquisition: %v calls%d", err, calls)
			}
			view := got.Records()
			if len(view) != 1 || view[0].Change.Path != "a" || view[0].Numstat.AddedRaw != "2" || *view[0].Numstat.Added != 2 {
				t.Fatal("metadata or exact tokens lost")
			}
			view[0].Change.Path, view[0].Numstat.AddedRaw = "changed", "99"
			*view[0].Numstat.Added, *view[0].Numstat.Deleted = 99, 88
			after := got.Records()[0]
			if after.Change.Path != "a" || after.Numstat.AddedRaw != "2" || *after.Numstat.Added != 2 || *after.Numstat.Deleted != 1 {
				t.Fatal("session getter alias")
			}
		})
	}
	var zero *gitadapter.SnapshotMetadataAcquisition
	if zero.Revisions() != (evidence.RevisionIdentity{}) || zero.Records() != nil {
		t.Fatal("nil session has populated views")
	}
}

func TestAcquireSnapshotMetadataBoundedNativeGit(t *testing.T) {
	root := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		command := exec.Command("git", args...)
		command.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0"}
		out, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("fixture Git: %v %s", err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run("init", root)
	commit := func(message string) {
		run("-C", root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "commit", "--allow-empty", "-m", message)
	}
	commit("base")
	base := run("-C", root, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(root, "a"), []byte("x\n"), 0600); err != nil {
		t.Fatal(err)
	}
	run("-C", root, "add", "a")
	commit("head")
	got, err := gitadapter.AcquireSnapshotMetadataBounded(context.Background(), gitadapter.Request{Repository: root, Base: base, Head: "HEAD"}, metadataAcquireBounds())
	if err != nil || got == nil {
		t.Fatalf("native metadata: %v", err)
	}
	view := got.Records()
	if got.Revisions().Base != base || got.Revisions().Head != run("-C", root, "rev-parse", "HEAD") || len(view) != 1 || view[0].Change.Status != "A" || view[0].Change.Path != "a" || view[0].Numstat.AddedRaw != "1" || view[0].Numstat.DeletedRaw != "0" {
		t.Fatal("native range or numstat evidence lost")
	}
}

func TestAcquireSnapshotMetadataBoundedInconsistentResponses(t *testing.T) {
	raw, numstat := metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("2", "1", "a", "")
	cases := []struct {
		name   string
		at     int
		output []byte
		code   evidence.SnapshotErrorCode
		calls  int
	}{
		{"unknown_format", 0, []byte("unknown\n"), evidence.SnapshotForeign, 1},
		{"invalid_resolved_id", 1, []byte("private invalid id\n"), evidence.SnapshotUnresolved, 2},
		{"malformed_raw", 5, []byte("broken"), evidence.SnapshotRacing, 6},
		{"duplicate_numstat", 6, append(append([]byte(nil), numstat...), numstat...), evidence.SnapshotRacing, 7},
		{"missing_numstat", 6, nil, evidence.SnapshotRacing, 7},
		{"wrong_blob_width", 5, metadataRawFixture("M", "a", "", 64), evidence.SnapshotRacing, 7},
		{"nonempty_identical_revisions", 3, []byte(strings.Repeat("a", 40) + "\n"), evidence.SnapshotRacing, 6},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outputs, _ := metadataAcquireScript(40, raw, numstat)
			outputs[tc.at] = tc.output
			calls := 0
			runner := func(context.Context, string, int, ...string) ([]byte, error) {
				out := outputs[calls]
				calls++
				return out, nil
			}
			got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, metadataAcquireBounds())
			var snapshotError *evidence.SnapshotError
			if !errors.As(err, &snapshotError) || snapshotError.Code != tc.code || got != nil || calls != tc.calls || strings.Contains(err.Error(), "private") {
				t.Fatalf("inconsistent response: %v %v calls%d", got, err, calls)
			}
		})
	}
}
