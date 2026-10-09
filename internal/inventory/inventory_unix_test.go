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
	}, nil)
	assertUnavailable(t, got, err, evidence.InventoryRacing)
}

func TestAcquireRejectsReplacementRaceInGitRepositoryWithoutRecords(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo := t.TempDir()
	git(t, repo, "init", "-q")
	git(t, repo, "config", "user.email", "test@example.com")
	git(t, repo, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repo, "tracked"), []byte("anchor"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "add", "tracked")
	git(t, repo, "commit", "-qm", "anchor")
	if got := git(t, repo, "ls-tree", "--name-only", "HEAD", "--", "tracked"); got != "tracked" {
		t.Fatalf("committed anchor = %q; want tracked", got)
	}
	if err := os.WriteFile(filepath.Join(repo, "race"), []byte("before"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := git(t, repo, "ls-files", "--others", "--", "race"); got != "race" {
		t.Fatalf("untracked race target = %q; want race", got)
	}
	replaced := false
	got, err := acquire(Request{Root: repo, Paths: []string{"race"}}, func(name string) {
		if name != "race" || replaced {
			return
		}
		if err := os.Rename(filepath.Join(repo, "race"), filepath.Join(repo, "old")); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repo, "race"), []byte("after"), 0o644); err != nil {
			t.Fatal(err)
		}
		replaced = true
	}, nil)
	if !replaced {
		t.Fatal("replacement hook did not run")
	}
	assertUnavailable(t, got, err, evidence.InventoryRacing)
}

func TestAcquireRejectsEscapesAndIntermediateReplacementInGitRepositoryWithoutRecords(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	for _, test := range []struct {
		name string
		code evidence.InventoryUnavailableCode
	}{
		{"intermediate symlink escape", evidence.InventorySymlink},
		{"final symlink escape", evidence.InventorySymlink},
		{"intermediate pre-open replacement", evidence.InventoryRacing},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo, outside := t.TempDir(), t.TempDir()
			git(t, repo, "init", "-q")
			git(t, repo, "config", "user.email", "test@example.com")
			git(t, repo, "config", "user.name", "Test")
			if err := os.WriteFile(filepath.Join(repo, "tracked"), []byte("anchor"), 0o644); err != nil {
				t.Fatal(err)
			}
			git(t, repo, "add", "tracked")
			git(t, repo, "commit", "-qm", "anchor")
			if got := git(t, repo, "ls-tree", "--name-only", "HEAD", "--", "tracked"); got != "tracked" {
				t.Fatalf("committed anchor = %q; want tracked", got)
			}
			if err := os.WriteFile(filepath.Join(outside, "target"), []byte("outside"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(repo, "a-valid"), []byte("valid"), 0o644); err != nil {
				t.Fatal(err)
			}
			path := "component/target"
			switch test.name {
			case "intermediate symlink escape":
				if err := os.Symlink(outside, filepath.Join(repo, "component")); err != nil {
					t.Fatal(err)
				}
			case "final symlink escape":
				path = "final"
				if err := os.Symlink(filepath.Join(outside, "target"), filepath.Join(repo, path)); err != nil {
					t.Fatal(err)
				}
			default:
				if err := os.Mkdir(filepath.Join(repo, "component"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(repo, path), []byte("before"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			observable := strings.Split(path, "/")[0]
			if got := git(t, repo, "ls-files", "--others", "--", observable); got != observable && got != path {
				t.Fatalf("untracked target = %q; want %q or %q", got, observable, path)
			}
			calls := 0
			var afterStat func(string)
			if test.code == evidence.InventoryRacing {
				afterStat = func(name string) {
					if name != "component" {
						return
					}
					calls++
					if calls != 1 {
						t.Fatal("intermediate replacement hook ran more than once")
					}
					if err := os.Rename(filepath.Join(repo, name), filepath.Join(repo, "old-component")); err != nil {
						t.Fatal(err)
					}
					if err := os.Mkdir(filepath.Join(repo, name), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(repo, path), []byte("after"), 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			got, err := acquire(Request{Root: repo, Paths: []string{"a-valid", path}}, afterStat, nil)
			if test.code == evidence.InventoryRacing && calls != 1 {
				t.Fatalf("intermediate replacement hook calls = %d; want 1", calls)
			}
			assertUnavailable(t, got, err, test.code)
		})
	}
}

func TestAcquireRejectsPostOpenReplacementInGitRepositoryWithoutRecords(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	repo, outside := t.TempDir(), t.TempDir()
	git(t, repo, "init", "-q")
	git(t, repo, "config", "user.email", "test@example.com")
	git(t, repo, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repo, "tracked"), []byte("anchor"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "add", "tracked")
	git(t, repo, "commit", "-qm", "anchor")
	if got := git(t, repo, "ls-tree", "--name-only", "HEAD", "--", "tracked"); got != "tracked" {
		t.Fatalf("committed anchor = %q; want tracked", got)
	}
	for _, name := range []string{"a-valid", "postopen-target"} {
		if err := os.WriteFile(filepath.Join(repo, name), []byte("before"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got := git(t, repo, "ls-files", "--others", "--", "postopen-target"); got != "postopen-target" {
		t.Fatalf("untracked target = %q; want postopen-target", got)
	}
	outsideTarget := filepath.Join(outside, "target")
	if err := os.WriteFile(outsideTarget, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	calls := 0
	got, err := acquire(Request{Root: repo, Paths: []string{"a-valid", "postopen-target"}}, nil, func(name string) {
		if name != "postopen-target" {
			return
		}
		calls++
		if calls != 1 {
			t.Fatal("post-open replacement hook ran more than once")
		}
		if err := os.Rename(filepath.Join(repo, name), filepath.Join(repo, "old-target")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outsideTarget, filepath.Join(repo, name)); err != nil {
			t.Fatal(err)
		}
	})
	if calls != 1 {
		t.Fatalf("post-open replacement hook calls = %d; want 1", calls)
	}
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
	inventory, err := Acquire(Request{Root: repo, Paths: []string{"untracked"}})
	if err != nil {
		t.Fatal(err)
	}
	want := []evidence.UntrackedRecord{{Path: "untracked"}}
	if !reflect.DeepEqual(inventory.Records(), want) {
		t.Fatalf("inventory records = %#v; want %#v", inventory.Records(), want)
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
