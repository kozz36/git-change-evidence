//go:build linux

package acquire_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/acquire"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"golang.org/x/sys/unix"
)

func failure(t *testing.T, root string, paths [][]byte, status int, token string) {
	t.Helper()
	got, err := acquire.Read(context.Background(), root, paths)
	var f *contract.Failure
	if got != nil || !errors.As(err, &f) || f.Status != status || f.Token != token {
		t.Fatalf("got %v, %v; want %d %s", got, err, status, token)
	}
}

func TestRegularAndEmpty(t *testing.T) {
	root := t.TempDir()
	path := []byte("-raw \n\xff")
	content := []byte{0, 255, 10}
	if err := os.WriteFile(filepath.Join(root, string(path)), content, 0600); err != nil {
		t.Fatal(err)
	}
	if !acquire.Supported() {
		t.Fatal("Linux unsupported")
	}
	got, err := acquire.Read(context.Background(), root, [][]byte{path})
	if err != nil || len(got) != 1 || string(got[0].Path) != string(path) || string(got[0].Content) != string(content) {
		t.Fatalf("%v %v", got, err)
	}
	path[0] = 'x'
	if got[0].Path[0] != '-' {
		t.Fatal("path alias")
	}
	after, err := os.ReadFile(filepath.Join(root, string(got[0].Path)))
	if err != nil || string(after) != string(content) {
		t.Fatal("source changed", err)
	}
	entries, err := acquire.Read(context.Background(), root, nil)
	if err != nil || len(entries) != 0 {
		t.Fatalf("empty: %v %v", entries, err)
	}
	failure(t, root+"/missing", nil, 3, "input_unavailable")
}

func TestRejectedPathsAndTypes(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, root+"/escape"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing", root+"/leaf"); err != nil {
		t.Fatal(err)
	}
	if err := unix.Mkfifo(root+"/fifo", 0600); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"../x", "/x", "x\x00", ".GiT/x", "a//b", "escape/x"} {
		t.Run(p, func(t *testing.T) { failure(t, root, [][]byte{[]byte(p)}, 2, "unsafe_path") })
	}
	gotLeaf, err := acquire.Read(context.Background(), root, [][]byte{[]byte("leaf")})
	if err != nil || len(gotLeaf) != 1 || string(gotLeaf[0].Content) != "missing" {
		t.Fatalf("leaf payload: %v %v", gotLeaf, err)
	}
	for _, p := range []string{"missing", "fifo", "."} {
		if p == "." {
			continue
		}
		t.Run(p, func(t *testing.T) { failure(t, root, [][]byte{[]byte(p)}, 3, "input_unavailable") })
	}
	if err := os.Mkdir(root+"/dir", 0700); err != nil {
		t.Fatal(err)
	}
	failure(t, root, [][]byte{[]byte("dir")}, 3, "input_unavailable")
	for _, r := range []string{"", ".", "relative", root + "/", root + "/../x"} {
		failure(t, r, nil, 2, "invalid_arguments")
	}
	failure(t, root+"/.GIT", nil, 2, "unsafe_path")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := acquire.Read(ctx, root, nil)
	var f *contract.Failure
	if got != nil || !errors.As(err, &f) || f.Status != 6 || f.Token != "interrupted" {
		t.Fatalf("cancel: %v %v", got, err)
	}
}
