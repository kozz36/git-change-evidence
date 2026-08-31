//go:build aix || android || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package inventoryadapter

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"golang.org/x/sys/unix"
)

func TestAcquireRecordsNestedRawPath(t *testing.T) {
	root, raw, plain := t.TempDir(), "nested/\xff\nfile", "nested/plain"
	if err := os.Mkdir(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{raw, plain} {
		if err := os.WriteFile(filepath.Join(root, path), []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := Acquire(Request{Root: root, Paths: []string{raw, plain}})
	want := []evidence.UntrackedRecord{{Path: plain}, {Path: raw}}
	if err != nil || !reflect.DeepEqual(got.Records(), want) {
		t.Fatalf("inventory, error = %#v, %v", got, err)
	}
}

func TestAcquireRejectsUnsafeOrUnavailablePathWithoutRecords(t *testing.T) {
	for _, test := range []struct {
		name, path string
		code       evidence.InventoryUnavailableCode
		setup      func(*testing.T, string) error
	}{
		{"intermediate symlink", "nested/link/escaped", evidence.InventorySymlink, func(t *testing.T, root string) error { return os.MkdirAll(filepath.Join(root, "nested"), 0o755) }},
		{"final symlink", "final", evidence.InventorySymlink, func(t *testing.T, root string) error { return os.Symlink(t.TempDir(), filepath.Join(root, "final")) }},
		{"special entry", "special", evidence.InventorySpecial, func(_ *testing.T, root string) error { return unix.Mkfifo(filepath.Join(root, "special"), 0o600) }},
		{"missing entry", "missing", evidence.InventoryMissing, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if test.setup != nil {
				if err := test.setup(t, root); err != nil {
					t.Fatal(err)
				}
			}
			if test.name == "intermediate symlink" {
				if err := os.Symlink(t.TempDir(), filepath.Join(root, "nested", "link")); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(filepath.Join(root, "a-valid"), []byte("valid"), 0o644); err != nil {
				t.Fatal(err)
			}
			got, err := Acquire(Request{Root: root, Paths: []string{"a-valid", test.path}})
			assertUnavailable(t, got, err, test.code)
		})
	}
}

func TestAcquireRejectsReplacementRaceWithoutRecords(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "race"), []byte("before"), 0o644); err != nil {
		t.Fatal(err)
	}
	replaced := false
	got, err := acquire(Request{Root: root, Paths: []string{"race"}}, func(name string) {
		if name != "race" || replaced {
			return
		}
		replaced = true
		if err := os.Rename(filepath.Join(root, "race"), filepath.Join(root, "old")); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "race"), []byte("after"), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	assertUnavailable(t, got, err, evidence.InventoryRacing)
}

func TestAcquireDoesNotChangeCommittedSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := t.TempDir()
	git(t, repo, "init", "-q")
	git(t, repo, "config", "user.email", "test@example.com")
	git(t, repo, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repo, "tracked"), []byte("base"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "add", "tracked")
	git(t, repo, "commit", "-qm", "base")
	base := git(t, repo, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(repo, "tracked"), []byte("head"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "commit", "-am", "head", "-q")
	head := git(t, repo, "rev-parse", "HEAD")
	before, err := gitadapter.Acquire(gitadapter.Request{Repository: repo, Base: base, Head: head})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "untracked"), []byte("inventory"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Acquire(Request{Root: repo, Paths: []string{"untracked"}}); err != nil {
		t.Fatal(err)
	}
	after, err := gitadapter.Acquire(gitadapter.Request{Repository: repo, Base: base, Head: head})
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatalf("snapshot after inventory = %#v, %v", after, err)
	}
}

func assertUnavailable(t *testing.T, got evidence.UntrackedInventory, err error, want evidence.InventoryUnavailableCode) {
	t.Helper()
	var unavailable *evidence.InventoryUnavailableError
	if len(got.Records()) != 0 || !errors.As(err, &unavailable) || unavailable.Code != want {
		t.Fatalf("inventory, error = %#v, %v; want unavailable %s", got, err, want)
	}
}

func git(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s", args, output)
	}
	return strings.TrimSpace(string(output))
}
