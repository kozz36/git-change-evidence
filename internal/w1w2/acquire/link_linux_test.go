//go:build linux

package acquire

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"golang.org/x/sys/unix"
)

func TestLeafPayloads(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	target := outside + "/target"
	if err := os.WriteFile(target, []byte("not payload"), 0600); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ name, payload string }{
		{"dangling", "missing"}, {"directory", outside}, {"outside", target},
		{"chain", "outside"}, {"-raw\n\xff", "../opaque\n\xff"},
	}
	for _, tc := range cases {
		if err := os.Symlink(tc.payload, root+"/"+tc.name); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := []byte(tc.name)
			got, err := Read(context.Background(), root, [][]byte{path})
			if err != nil || len(got) != 1 || string(got[0].Path) != tc.name || string(got[0].Content) != tc.payload {
				t.Fatalf("%v %v", got, err)
			}
			path[0] = 'X'
			if string(got[0].Path) != tc.name {
				t.Fatal("path alias")
			}
			payload, err := os.Readlink(root + "/" + tc.name)
			if err != nil || payload != tc.payload {
				t.Fatalf("link changed: %q %v", payload, err)
			}
		})
	}
	content, err := os.ReadFile(target)
	if err != nil || string(content) != "not payload" {
		t.Fatalf("target changed: %q %v", content, err)
	}
	names, err := os.ReadDir(root)
	if err != nil || len(names) != len(cases) {
		t.Fatalf("unexpected writes: %v %v", names, err)
	}
}

func TestMixedLinkBudget(t *testing.T) {
	root := t.TempDir()
	paths := make([][]byte, 0, 11)
	for i := 0; i < 8; i++ {
		name := fmt.Sprintf("file%d", i)
		size := contract.MaxFileBytes
		if i == 7 {
			size -= 2
		}
		if err := os.WriteFile(root+"/"+name, make([]byte, size), 0600); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, []byte(name))
	}
	if err := os.Symlink("x", root+"/link"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"alias", "excess"} {
		if err := os.Link(root+"/link", root+"/"+name); err != nil {
			t.Fatal(err)
		}
	}
	paths = append(paths, []byte("link"), []byte("alias"))
	got, err := Read(context.Background(), root, paths)
	if err != nil || len(got) != 10 || string(got[8].Content) != "x" || string(got[9].Content) != "x" {
		t.Fatalf("exact mixed budget: %v", err)
	}
	paths = append(paths, []byte("excess"))
	got, err = Read(context.Background(), root, paths)
	var f *contract.Failure
	if got != nil || !errors.As(err, &f) || f.Status != 4 || f.Token != "resource_limit" {
		t.Fatalf("excess: %v %v", got, err)
	}
	if err := os.WriteFile(root+"/file7", make([]byte, contract.MaxFileBytes), 0600); err != nil {
		t.Fatal(err)
	}
	got, err = Read(context.Background(), root, paths[:9])
	if got != nil || !errors.As(err, &f) || f.Status != 4 || f.Token != "resource_limit" {
		t.Fatalf("8MiB regular plus link: %v %v", got, err)
	}
}

func TestLinkCaptureReplacement(t *testing.T) {
	for _, afterCapture := range []bool{false, true} {
		t.Run(fmt.Sprintf("afterCapture=%t", afterCapture), func(t *testing.T) {
			root := t.TempDir()
			for _, name := range []string{"victim", "replacement"} {
				if err := os.Symlink(name, root+"/"+name); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(root+"/first", []byte("first"), 0600); err != nil {
				t.Fatal(err)
			}
			ops, captured, firstRead := nativeOps(), -1, false
			replace := func() {
				if err := os.Rename(root+"/replacement", root+"/victim"); err != nil {
					t.Fatal(err)
				}
			}
			ops.openat = func(parent int, name string, flags int, mode uint32) (int, error) {
				if name == "victim" && !afterCapture {
					replace()
				}
				fd, err := unix.Openat(parent, name, flags, mode)
				if name == "victim" && err == nil {
					captured = fd
					if afterCapture {
						replace()
					}
				}
				return fd, err
			}
			ops.read = func(fd int, buf []byte) (int, error) {
				n, err := unix.Read(fd, buf)
				if n == 0 && err == nil {
					firstRead = true
				}
				return n, err
			}
			got, err := readWithOps(context.Background(), root, [][]byte{[]byte("first"), []byte("victim")}, ops)
			var f *contract.Failure
			if got != nil || !errors.As(err, &f) || f.Status != 7 || f.Token != "acquisition_race" {
				t.Fatalf("%v %v", got, err)
			}
			var stat unix.Stat_t
			if !firstRead || captured < 0 || unix.Fstat(captured, &stat) != unix.EBADF {
				t.Fatal("earlier entry or descriptor closure not proven")
			}
		})
	}
}
