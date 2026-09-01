package gitadapter

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestAcquireBlobUsesCommittedObjectIdentity(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	write(t, repo, "file", "before")
	base := commit(t, repo)
	content := string([]byte{0xff, 0, '\n', 0xfe})
	write(t, repo, "file", content)
	head := commit(t, repo)
	object := git(t, repo, "rev-parse", head+":file")
	git(t, repo, "replace", head, base)

	got, err := acquireBlob(BlobRequest{Repository: repo, Revision: "HEAD", Path: "file"}, func() {
		git(t, repo, "update-ref", "HEAD", base)
	})
	if err != nil {
		t.Fatal(err)
	}
	assertBlob(t, got, head, object, "100644", evidence.GitFile, "blob", content)
}

func TestAcquireBlobIgnoresIndexAndWorktree(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	write(t, repo, "file", "committed")
	head, object := commit(t, repo), ""
	object = git(t, repo, "rev-parse", head+":file")
	assertCommitted := func() {
		got, err := AcquireBlob(BlobRequest{Repository: repo, Revision: head, Path: "file"})
		if err != nil {
			t.Fatal(err)
		}
		assertBlob(t, got, head, object, "100644", evidence.GitFile, "blob", "committed")
	}
	write(t, repo, "file", "unstaged")
	assertCommitted()
	git(t, repo, "add", "--", "file")
	assertCommitted()
	if err := os.Remove(filepath.Join(repo, "file")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(repo, "file"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, repo, "file/child", "untracked conflict")
	assertCommitted()
	if err := os.Remove(filepath.Join(repo, "file/child")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(repo, "file")); err != nil {
		t.Fatal(err)
	}
	write(t, repo, "file", "worktree")
	other := filepath.Join(t.TempDir(), "other")
	git(t, repo, "worktree", "add", "-q", "-b", "divergent", other, head)
	write(t, other, "file", "divergent worktree")
	git(t, other, "add", "--", "file")
	t.Setenv("GIT_DIR", "/missing")
	t.Setenv("GIT_WORK_TREE", other)
	t.Setenv("GIT_INDEX_FILE", filepath.Join(t.TempDir(), "empty-index"))
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "core.attributesfile")
	t.Setenv("GIT_CONFIG_VALUE_0", "/missing")
	assertCommitted()
}

func TestAcquireBlobPreservesLiteralPathsAndSymlinks(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	name, content := "literal*\n-leading-dash", string([]byte{0xff, 0, '\n'})
	write(t, repo, name, content)
	write(t, repo, "literalx", "wrong wildcard")
	write(t, repo, "-leading-dash", "dash")
	if err := os.Symlink("literal*\n-leading-dash", filepath.Join(repo, "link")); err != nil {
		t.Fatal(err)
	}
	head := commit(t, repo)
	for _, test := range []struct {
		path, mode, typ, content string
		kind                     evidence.GitEntryKind
	}{{name, "100644", "blob", content, evidence.GitFile}, {"-leading-dash", "100644", "blob", "dash", evidence.GitFile}, {"link", "120000", "blob", name, evidence.GitSymlink}} {
		t.Run(test.path, func(t *testing.T) {
			got, err := AcquireBlob(BlobRequest{Repository: repo, Revision: head, Path: test.path})
			if err != nil {
				t.Fatal(err)
			}
			assertBlob(t, got, head, git(t, repo, "rev-parse", head+":"+test.path), test.mode, test.kind, test.typ, test.content)
		})
	}
}

func TestAcquireBlobReportsTypedFailuresWithoutContent(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, false)
	write(t, repo, "file", "content")
	if err := os.Mkdir(filepath.Join(repo, "dir"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, repo, "dir/file", "nested")
	base := commit(t, repo)
	sub := repository(t, false)
	write(t, sub, "file", "submodule")
	subhead := commit(t, sub)
	git(t, repo, "update-index", "--add", "--cacheinfo", "160000,"+subhead+",module")
	git(t, repo, "commit", "-qm", "gitlink")
	head := git(t, repo, "rev-parse", "HEAD")
	for _, test := range []struct {
		name string
		in   BlobRequest
		want BlobErrorCode
	}{
		{"hostile revision", BlobRequest{repo, "--help", "file"}, BlobUnresolved},
		{"missing revision", BlobRequest{repo, strings.Repeat("0", len(head)), "file"}, BlobMissing},
		{"missing path", BlobRequest{repo, head, "missing"}, BlobMissing},
		{"foreign repository", BlobRequest{t.TempDir(), head, "file"}, BlobForeign},
		{"tree", BlobRequest{repo, head, "dir"}, BlobNotBlob},
		{"gitlink", BlobRequest{repo, head, "module"}, BlobNotBlob},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := AcquireBlob(test.in)
			assertBlobFailure(t, got, err, test.want)
		})
	}
	for _, path := range []string{"", "/absolute", "../file", "dir/../file", "dir/./file", "dir//file", "nul\x00path"} {
		t.Run("invalid path "+path, func(t *testing.T) {
			got, err := AcquireBlob(BlobRequest{repo, head, path})
			assertBlobFailure(t, got, err, BlobInvalidPath)
		})
	}
	got, err := acquireBlob(BlobRequest{repo, head, "file"}, func() {
		git(t, repo, "update-ref", "HEAD", base)
		git(t, repo, "reflog", "expire", "--expire=now", "--all")
		git(t, repo, "prune", "--expire=now")
	})
	assertBlobFailure(t, got, err, BlobRacing)
}

func TestAcquireBlobPreservesSHA256Identity(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := repository(t, true)
	content := []byte{0xff, 0, '\n', 0xfe}
	write(t, repo, "file", string(content))
	head := commit(t, repo)
	object := hashBlob(t, repo, content)
	committedObject := git(t, repo, "rev-parse", head+":file")
	if len(head) != 64 || len(object) != 64 || object != committedObject {
		t.Fatalf("SHA-256 commit, content object, committed object = %q, %q, %q", head, object, committedObject)
	}

	got, err := AcquireBlob(BlobRequest{Repository: repo, Revision: head, Path: "file"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Commit != evidence.GitObjectID(head) || got.Object != evidence.GitObjectID(object) || got.Mode != "100644" || got.Type != "blob" || got.Kind != evidence.GitFile || !bytes.Equal([]byte(got.Content), content) {
		t.Fatalf("blob = %#v; want commit %q, object %q, mode 100644, type blob, kind %q, content %x", got, head, object, evidence.GitFile, content)
	}
}

func hashBlob(t *testing.T, repo string, content []byte) string {
	t.Helper()
	cmd := exec.Command("git", "hash-object", "-t", "blob", "--stdin")
	cmd.Dir = repo
	cmd.Stdin = bytes.NewReader(content)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git hash-object: %s", output)
	}
	return strings.TrimSpace(string(output))
}

func assertBlob(t *testing.T, got Blob, commit, object, mode string, kind evidence.GitEntryKind, typ, content string) {
	t.Helper()
	if got.Commit != evidence.GitObjectID(commit) || got.Object != evidence.GitObjectID(object) || got.Mode != mode || got.Kind != kind || got.Type != typ || got.Content != content {
		t.Fatalf("blob = %#v", got)
	}
}

func assertBlobFailure(t *testing.T, got Blob, err error, want BlobErrorCode) {
	t.Helper()
	var failure *BlobError
	if got != (Blob{}) || !errors.As(err, &failure) || failure.Code != want {
		t.Fatalf("blob, error = %#v, %v; want %s", got, err, want)
	}
}
