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

func TestDirectoryCaptureRejectsObservedBindingChanges(t *testing.T) {
	for _, change := range []string{"directory", "symlink", "mode", "denied"} {
		for _, afterCapture := range []bool{false, true} {
			if change == "denied" && afterCapture {
				continue
			}
			t.Run(fmt.Sprintf("%s/after=%t", change, afterCapture), func(t *testing.T) {
				base, ancestor, root := directoryCaptureFixture(t)
				replacement := base + "/replacement"
				if change == "directory" {
					if err := os.MkdirAll(replacement+"/selected", 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(replacement+"/selected/leaf", []byte("unbound replacement"), 0600); err != nil {
						t.Fatal(err)
					}
				} else if change == "symlink" {
					if err := os.Symlink("missing", replacement); err != nil {
						t.Fatal(err)
					}
				}
				mutate := func() {
					if change == "mode" {
						if err := os.Chmod(ancestor, 0750); err != nil {
							t.Fatal(err)
						}
						return
					}
					if err := os.Rename(ancestor, base+"/moved"); err != nil {
						t.Fatal(err)
					}
					if err := os.Rename(replacement, ancestor); err != nil {
						t.Fatal(err)
					}
				}
				ops, hooks, captured := nativeOps(), 0, -1
				ops.openat = func(parent int, name string, flags int, mode uint32) (int, error) {
					if name != "ancestor" {
						return unix.Openat(parent, name, flags, mode)
					}
					hooks++
					if change == "denied" {
						return -1, unix.EACCES
					}
					if !afterCapture {
						mutate()
					}
					fd, err := unix.Openat(parent, name, flags, mode)
					if err == nil {
						captured = fd
						if afterCapture {
							mutate()
						}
					}
					return fd, err
				}
				path := []byte("leaf")
				got, err := readWithOps(context.Background(), root, [][]byte{path}, ops)
				status, token := 7, "acquisition_race"
				if change == "denied" {
					status, token = 3, "input_unavailable"
				}
				var failure *contract.Failure
				if got != nil || !errors.As(err, &failure) || failure.Status != status || failure.Token != token || err.Error() != token || hooks != 1 || string(path) != "leaf" {
					t.Fatalf("binding failure changed: %v %v hooks%d", got, err, hooks)
				}
				if captured >= 0 {
					var stat unix.Stat_t
					if unix.Fstat(captured, &stat) != unix.EBADF {
						t.Fatal("rejected directory descriptor was not closed")
					}
				}
			})
		}
	}
}

func TestDirectoryCaptureRefusesInitialSymlink(t *testing.T) {
	base, ancestor, root := directoryCaptureFixture(t)
	if err := os.Rename(ancestor, base+"/moved"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(base+"/moved", ancestor); err != nil {
		t.Fatal(err)
	}
	got, err := Read(context.Background(), root, [][]byte{[]byte("leaf")})
	var failure *contract.Failure
	if got != nil || !errors.As(err, &failure) || failure.Status != 2 || failure.Token != "unsafe_path" || err.Error() != "unsafe_path" {
		t.Fatalf("intermediate symlink refusal changed: %v %v", got, err)
	}
}
