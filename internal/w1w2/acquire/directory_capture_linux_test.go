//go:build linux

package acquire

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestDirectoryCaptureToleratesUnrelatedEntries(t *testing.T) {
	for _, absoluteRoot := range []bool{true, false} {
		for _, tc := range []struct {
			name         string
			afterCapture bool
		}{
			{"stable", false}, {"create", false}, {"create", true},
			{"remove", false}, {"remove", true},
		} {
			t.Run(fmt.Sprintf("root=%t/%s/after=%t", absoluteRoot, tc.name, tc.afterCapture), func(t *testing.T) {
				base, ancestor, root := directoryCaptureFixture(t)
				sibling := ancestor + "/unselected"
				if tc.name == "remove" {
					if err := os.Mkdir(sibling, 0700); err != nil {
						t.Fatal(err)
					}
				}
				var leafBefore, leafAfter, before, after unix.Stat_t
				if err := unix.Lstat(root+"/leaf", &leafBefore); err != nil {
					t.Fatal(err)
				}
				mutate := func() {
					var err error
					if tc.name == "create" {
						err = os.Mkdir(sibling, 0700)
					} else if tc.name == "remove" {
						err = os.Remove(sibling)
					}
					if err != nil {
						t.Fatal(err)
					}
				}
				ops, hooks := nativeOps(), 0
				ops.openat = func(parent int, name string, flags int, mode uint32) (int, error) {
					if name != "ancestor" {
						return unix.Openat(parent, name, flags, mode)
					}
					hooks++
					if err := unix.Fstatat(parent, name, &before, unix.AT_SYMLINK_NOFOLLOW); err != nil {
						t.Fatal(err)
					}
					fd := -1
					var err error
					if tc.afterCapture {
						fd, err = unix.Openat(parent, name, flags, mode)
						if err != nil {
							return fd, err
						}
					}
					mutate()
					if err := unix.Fstatat(parent, name, &after, unix.AT_SYMLINK_NOFOLLOW); err != nil {
						t.Fatal(err)
					}
					if tc.afterCapture {
						return fd, nil
					}
					return unix.Openat(parent, name, flags, mode)
				}
				requested := "leaf"
				if !absoluteRoot {
					root, requested = base, "ancestor/selected/leaf"
				}
				path := []byte(requested)
				got, err := readWithOps(context.Background(), root, [][]byte{path}, ops)
				if err != nil || len(got) != 1 || !bytes.Equal(got[0].Path, path) || string(got[0].Content) != "payload\n" {
					t.Fatalf("unrelated directory entries rejected: %v %v", got, err)
				}
				if hooks != 1 || before.Dev != after.Dev || before.Ino != after.Ino || before.Mode != after.Mode {
					t.Fatal("directory binding or expected handoff not established")
				}
				if tc.name != "stable" && before.Size == after.Size && before.Mtim == after.Mtim && before.Ctim == after.Ctim {
					t.Fatal("fixture did not change directory metadata")
				}
				leaf := ancestor + "/selected/leaf"
				if err := unix.Lstat(leaf, &leafAfter); err != nil || !same(leafBefore, leafAfter) || string(path) != requested {
					t.Fatalf("selected leaf or input changed: %v", err)
				}
				path[0] = 'X'
				if string(got[0].Path) != requested {
					t.Fatal("returned path aliases caller input")
				}
			})
		}
	}
}

func directoryCaptureFixture(t *testing.T) (base, ancestor, root string) {
	t.Helper()
	base = t.TempDir()
	ancestor, root = base+"/ancestor", base+"/ancestor/selected"
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(root+"/leaf", []byte("payload\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return base, ancestor, root
}
