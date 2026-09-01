package gitadapter

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestAcquirePreservesCommittedEvidence(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	old, renamed := "old-\xff", "new-\xff\npath"
	write(t, repo, old, "rename")
	write(t, repo, "binary", "before\x00")
	write(t, repo, "text", "before\n")
	write(t, repo, "executable", "run")
	base := commit(t, repo)
	git(t, repo, "mv", old, renamed)
	write(t, repo, "binary", "after\x00")
	write(t, repo, "text", "after\nmore\n")
	if err := os.Chmod(filepath.Join(repo, "executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("binary", filepath.Join(repo, "link")); err != nil {
		t.Fatal(err)
	}
	sub := repository(t, false)
	write(t, sub, "file", "submodule")
	subhead := commit(t, sub)
	git(t, repo, "add", "-A")
	git(t, repo, "update-index", "--add", "--cacheinfo", "160000,"+subhead+",module")
	git(t, repo, "commit", "-qm", "snapshot")
	head := git(t, repo, "rev-parse", "HEAD")

	got, err := Acquire(Request{Repository: repo, Base: base, Head: head})
	if err != nil {
		t.Fatal(err)
	}
	if got.Base != evidence.GitObjectID(base) || got.Head != evidence.GitObjectID(head) {
		t.Fatalf("revisions = %#v", got)
	}
	entries := map[string]evidence.CommittedChange{}
	for _, entry := range got.Entries() {
		entries[entry.Path] = entry
	}
	copy := got.Entries()
	copy[0].Path = "mutated"
	if got.Entries()[0].Path == "mutated" {
		t.Fatal("snapshot entries are mutable")
	}
	assertChange(t, entries[renamed], "R", old, "100644", evidence.GitFile, false)
	assertChange(t, entries["binary"], "M", "", "100644", evidence.GitFile, true)
	assertChange(t, entries["text"], "M", "", "100644", evidence.GitFile, false)
	assertChange(t, entries["executable"], "M", "", "100755", evidence.GitFile, false)
	assertChange(t, entries["link"], "A", "", "120000", evidence.GitSymlink, false)
	assertChange(t, entries["module"], "A", "", "160000", evidence.GitGitlink, false)
	assertLineCounts(t, entries[renamed], 0, 0, true)
	assertLineCounts(t, entries["binary"], 0, 0, false)
	assertLineCounts(t, entries["text"], 2, 1, true)
	assertLineCounts(t, entries["executable"], 0, 0, true)
	assertLineCounts(t, entries["link"], 1, 0, true)
	assertLineCounts(t, entries["module"], 1, 0, true)

	write(t, repo, "binary", "dirty")
	git(t, repo, "add", "binary")
	write(t, repo, "untracked", "ignored")
	git(t, repo, "replace", head, base)
	t.Setenv("GIT_DIR", "/missing")
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "diff.external")
	t.Setenv("GIT_CONFIG_VALUE_0", "/bin/false")
	isolated, err := Acquire(Request{Repository: repo, Base: base, Head: head})
	if err != nil || !reflect.DeepEqual(got, isolated) {
		t.Fatalf("isolated acquisition = (%#v, %v), want %#v", isolated, err, got)
	}
}

func TestAcquireReportsTypedFailuresWithoutSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	write(t, repo, "file", "base")
	base := commit(t, repo)
	write(t, repo, "file", "head")
	head := commit(t, repo)
	for _, test := range []struct {
		revision string
		code     evidence.SnapshotErrorCode
	}{
		{"unknown", evidence.SnapshotUnresolved}, {strings.Repeat("0", 40), evidence.SnapshotMissing},
	} {
		snapshot, err := Acquire(Request{Repository: repo, Base: base, Head: test.revision})
		assertFailure(t, snapshot, err, test.code)
	}
	snapshot, err := Acquire(Request{Repository: t.TempDir(), Base: base, Head: head})
	assertFailure(t, snapshot, err, evidence.SnapshotForeign)
	git(t, repo, "update-ref", "HEAD", base)
	snapshot, err = mustRace(t, repo, base, head)
	assertFailure(t, snapshot, err, evidence.SnapshotRacing)
}

func TestParseNumstatKeepsBinaryAndOverflowExplicitlyNonCountable(t *testing.T) {
	got, ok := parseNumstat([]byte("18446744073709551616\t0\toverflow\x00-\t-\tbinary\x00"))
	if !ok || got["overflow"].Lines.Countable || got["overflow"].Binary || got["binary"].Lines.Countable || !got["binary"].Binary {
		t.Fatalf("measurements = %#v, %v", got, ok)
	}
}

func TestAcquirePreservesSHA256Identity(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, true)
	write(t, repo, "file", "base")
	base := commit(t, repo)
	write(t, repo, "file", "head")
	head := commit(t, repo)
	snapshot, err := Acquire(Request{Repository: repo, Base: base, Head: head})
	if err != nil || len(snapshot.Base) != 64 || len(snapshot.Head) != 64 {
		t.Fatalf("snapshot = %#v, %v", snapshot, err)
	}
}

func TestAcquireBlobCancelsRunningGitCommand(t *testing.T) {
	for _, test := range []struct {
		name  string
		cause error
		want  BlobErrorCode
	}{
		{"custom canceled", errors.New("custom cancellation"), BlobCanceled}, {"deadline", context.DeadlineExceeded, BlobDeadline},
	} {
		t.Run(test.name, func(t *testing.T) { assertAcquireBlobStops(t, test.cause, test.want) })
	}
}

type deadlineContext struct{ context.Context }

func (ctx deadlineContext) Err() error {
	if ctx.Context.Err() != nil {
		return context.DeadlineExceeded
	}
	return nil
}
func assertAcquireBlobStops(t *testing.T, cause error, want BlobErrorCode) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	write(t, repo, "file", "x")
	head := commit(t, repo)
	readyRead, readyWrite, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	releaseRead, releaseWrite, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		readyRead.Close()
		readyWrite.Close()
		releaseRead.Close()
		releaseWrite.Close()
	})
	child := make(chan *exec.Cmd, 1)
	runner := Runner(func(ctx context.Context, repository string, byteCap int, args ...string) ([]byte, error) {
		if len(args) > 1 && args[0] == "cat-file" && args[1] == "blob" {
			return controlledOutput(ctx, os.Args[0], []string{"-test.run=^TestAcquireBlobHelperProcess$", "--", "blob-helper"}, nil, byteCap, func(command *exec.Cmd) {
				command.Stdin, command.ExtraFiles = releaseRead, []*os.File{readyWrite}
				child <- command
			})
		}
		return gitOutputBounded(ctx, repository, byteCap, args...)
	})
	ctx, cancel := context.WithCancelCause(context.Background())
	if cause == context.DeadlineExceeded {
		inner, end := context.WithCancel(context.Background())
		ctx, cancel = deadlineContext{inner}, func(error) { end() }
	}
	defer cancel(cause)
	result := make(chan struct {
		blob Blob
		err  error
	}, 1)
	go func() {
		blob, err := AcquireBlob(ctx, runner, BlobRequest{Repository: repo, Revision: head, Path: "file"}, BlobBounds{ByteCap: 1})
		result <- struct {
			blob Blob
			err  error
		}{blob, err}
	}()
	ready := make(chan error, 1)
	go func() {
		_, err := io.ReadFull(readyRead, make([]byte, 1))
		ready <- err
	}()
	select {
	case err := <-ready:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("helper did not report readiness")
	}
	command := <-child
	cancel(cause)
	select {
	case outcome := <-result:
		assertBlobFailure(t, outcome.blob, outcome.err, want)
	case <-time.After(time.Second):
		releaseWrite.Close()
		select {
		case <-result:
			t.Fatal("canceled Git command did not return before helper release")
		case <-time.After(time.Second):
			t.Fatal("helper did not exit after release")
		}
	}
	if command.ProcessState == nil {
		t.Fatal("direct child was not reaped")
	}
}

func TestAcquireBlobHelperProcess(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "blob-helper" {
		return
	}
	ready := os.NewFile(3, "ready")
	if _, err := ready.Write([]byte{1}); err != nil {
		os.Exit(2)
	}
	_, _ = io.Copy(io.Discard, os.Stdin)
	os.Exit(0)
}

func assertChange(t *testing.T, entry evidence.CommittedChange, status, previous, mode string, kind evidence.GitEntryKind, binary bool) {
	t.Helper()
	if entry.Status != status || entry.PreviousPath != previous || entry.NewMode != mode || entry.NewKind != kind || entry.Binary != binary || len(entry.NewObject) != 40 {
		t.Fatalf("entry = %#v", entry)
	}
}
func assertLineCounts(t *testing.T, entry evidence.CommittedChange, additions, deletions uint64, countable bool) {
	t.Helper()
	if lines := entry.Lines; lines.Additions != additions || lines.Deletions != deletions || lines.Countable != countable {
		t.Fatalf("lines = %#v", lines)
	}
}
func assertFailure(t *testing.T, snapshot evidence.CommittedSnapshot, err error, want evidence.SnapshotErrorCode) {
	t.Helper()
	var failure *evidence.SnapshotError
	if snapshot.Base != "" || snapshot.Head != "" || len(snapshot.Entries()) != 0 || !errors.As(err, &failure) || failure.Code != want {
		t.Fatalf("snapshot, error = %#v, %v; want %s", snapshot, err, want)
	}
}
func mustRace(t *testing.T, repo, base, head string) (evidence.CommittedSnapshot, error) {
	t.Helper()
	return acquire(Request{Repository: repo, Base: base, Head: head}, func() {
		git(t, repo, "reflog", "expire", "--expire=now", "--all")
		git(t, repo, "prune", "--expire=now")
	})
}
func repository(t *testing.T, sha256 bool) string {
	t.Helper()
	repo := t.TempDir()
	args := []string{"init", "-q"}
	if sha256 {
		args = append(args, "--object-format=sha256")
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	if output, err := cmd.CombinedOutput(); err != nil {
		if sha256 {
			t.Skipf("SHA-256 Git unavailable: %s", output)
		}
		t.Fatal(string(output))
	}
	git(t, repo, "config", "user.email", "test@example.com")
	git(t, repo, "config", "user.name", "Test")
	return repo
}
func write(t *testing.T, repo, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
func commit(t *testing.T, repo string) string {
	t.Helper()
	git(t, repo, "add", "-A")
	git(t, repo, "commit", "-qm", "snapshot")
	return git(t, repo, "rev-parse", "HEAD")
}
func git(t *testing.T, repo string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s", args, output)
	}
	return strings.TrimSpace(string(output))
}
