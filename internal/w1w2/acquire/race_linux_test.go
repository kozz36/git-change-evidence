//go:build linux

package acquire

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"golang.org/x/sys/unix"
)

func TestReplacementBeforeOpen(t *testing.T) { testReplacement(t, false) }
func TestReplacementAfterRead(t *testing.T)  { testReplacement(t, true) }

func testReplacement(t *testing.T, afterRead bool) {
	t.Helper()
	root := t.TempDir()
	for name, content := range map[string]string{"first": "first", "victim": "old", "replacement": "new"} {
		if err := os.WriteFile(root+"/"+name, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	ops := nativeOps()
	intercepted := -1
	fired := false
	completedFirst := false
	ops.openat = func(parent int, name string, flags int, mode uint32) (int, error) {
		if name == "victim" && !afterRead {
			if err := os.Rename(root+"/replacement", root+"/victim"); err != nil {
				t.Fatal(err)
			}
			fired = true
		}
		fd, err := unix.Openat(parent, name, flags, mode)
		if name == "victim" && err == nil {
			intercepted = fd
		}
		return fd, err
	}
	ops.read = func(fd int, buf []byte) (int, error) {
		n, err := unix.Read(fd, buf)
		if fd != intercepted && n == 0 && err == nil {
			completedFirst = true
		}
		if afterRead && fd == intercepted && !fired {
			if err := os.Rename(root+"/replacement", root+"/victim"); err != nil {
				t.Fatal(err)
			}
			fired = true
		}
		return n, err
	}
	entries, err := readWithOps(context.Background(), root, [][]byte{[]byte("first"), []byte("victim")}, ops)
	var f *contract.Failure
	if entries != nil || !errors.As(err, &f) || f.Status != 7 || f.Token != "acquisition_race" {
		t.Fatalf("got %v, %v; want nil and 7 acquisition_race", entries, err)
	}
	if !fired || !completedFirst || intercepted < 0 {
		t.Fatal("handoff or earlier successful entry not exercised")
	}
	var stat unix.Stat_t
	if err := unix.Fstat(intercepted, &stat); err != unix.EBADF {
		t.Fatalf("owned fd not closed: %v", err)
	}
	for name, want := range map[string]string{"first": "first", "victim": "new"} {
		got, err := os.ReadFile(root + "/" + name)
		if err != nil || string(got) != want {
			t.Fatalf("%s: %q %v", name, got, err)
		}
	}
	names, err := os.ReadDir(root)
	if err != nil || len(names) != 2 {
		t.Fatalf("unexpected fixture changes: %v %v", names, err)
	}
	// A fresh call must not inherit the instance-local replacement operations.
	entries, err = Read(context.Background(), root, [][]byte{[]byte("victim")})
	if err != nil || len(entries) != 1 || string(entries[0].Content) != "new" {
		t.Fatalf("native call: %v %v", entries, err)
	}
}
