package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type domainFile struct{ path, data, mode string }
type domainFixture struct {
	repo, base, head string
	before, after    []domainFile
}

func domainGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	return domainGitInput(t, repo, "", args...)
}

func domainGitInput(t *testing.T, repo, input string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Env = []string{"PATH=/usr/bin:/bin", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1", "LC_ALL=C", "GIT_AUTHOR_NAME=Fixture", "GIT_AUTHOR_EMAIL=fixture@example.invalid", "GIT_COMMITTER_NAME=Fixture", "GIT_COMMITTER_EMAIL=fixture@example.invalid"}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %q: %v: %s", args, err, out)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func domainWrite(t *testing.T, repo string, f domainFile) {
	t.Helper()
	path := filepath.Join(repo, f.path)
	var err error
	if f.mode == "120000" {
		err = os.Symlink(f.data, path)
	} else {
		mode := os.FileMode(0600)
		if f.mode == "100755" {
			mode = 0700
		}
		err = os.WriteFile(path, []byte(f.data), mode)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func domainBuild(t *testing.T, format string) domainFixture {
	t.Helper()
	f := domainFixture{repo: t.TempDir()}
	domainGit(t, f.repo, "init", "--object-format="+format)
	f.before = []domainFile{
		{"deleted", "deletion distinct payload\n", "100644"},
		{"modified", "old\n", "100644"},
		{"rename-old", "uniquely retained rename payload\n", "100644"},
		{"mode", "mode-only\n", "100644"},
		{"to-link", "regular before link\n", "100644"},
		{"to-file", "opaque-target", "120000"},
		{"binary", "old\x00\n", "100644"},
	}
	for _, file := range f.before {
		domainWrite(t, f.repo, file)
	}
	domainGit(t, f.repo, "add", "--all")
	domainGit(t, f.repo, "commit", "-m", "base")
	f.base = domainGit(t, f.repo, "rev-parse", "HEAD")
	// Use an empty index and a fresh directory for the complete head tree.
	fresh := t.TempDir()
	domainGit(t, f.repo, "read-tree", "--empty")
	f.after = []domainFile{
		{"modified", "new\ntail", "100644"},
		{"rename-new", "uniquely retained rename payload\n", "100644"},
		{"mode", "mode-only\n", "100755"},
		{"to-link", "unfollowed-target", "120000"},
		{"to-file", "regular after link\n", "100644"},
		{"binary", "new\x00tail", "100644"},
		{"empty", "", "100644"},
		{"tab\tname", "\n", "100644"},
		{"line\nname", "tail", "100644"},
		{"raw\xff", "a\r\n", "100644"},
	}
	for _, file := range f.after {
		domainWrite(t, fresh, file)
		var oid string
		if file.mode == "120000" {
			oid = domainGitInput(t, f.repo, file.data, "hash-object", "-w", "--stdin")
		} else {
			oid = domainGit(t, f.repo, "hash-object", "-w", filepath.Join(fresh, file.path))
		}
		domainGit(t, f.repo, "update-index", "--add", "--cacheinfo", file.mode, oid, file.path)
	}
	domainGit(t, f.repo, "update-index", "--add", "--cacheinfo", "160000", f.base, "gitlink")
	tree := domainGit(t, f.repo, "write-tree")
	f.head = domainGit(t, f.repo, "commit-tree", tree, "-p", f.base, "-m", "head")
	return f
}
