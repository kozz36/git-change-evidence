//go:build linux

package acquire_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/acquire"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func TestContentBoundsAndHardlinks(t *testing.T) {
	root := t.TempDir()
	data := []byte(strings.Repeat("x", contract.MaxFileBytes))
	if err := os.WriteFile(root+"/base", data, 0600); err != nil {
		t.Fatal(err)
	}
	paths := [][]byte{[]byte("base")}
	for i := 1; i < 9; i++ {
		name := fmt.Sprintf("link%d", i)
		if err := os.Link(root+"/base", root+"/"+name); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, []byte(name))
	}
	for _, count := range []int{1, 7, 8} {
		got, err := acquire.Read(context.Background(), root, paths[:count])
		if err != nil || len(got) != count {
			t.Fatalf("count %d: %v", count, err)
		}
		for _, entry := range got {
			if string(entry.Content) != string(data) {
				t.Fatal("content altered")
			}
		}
	}
	failure(t, root, paths, 4, "resource_limit")
	if err := os.WriteFile(root+"/oversize", append(data, 'x'), 0600); err != nil {
		t.Fatal(err)
	}
	failure(t, root, [][]byte{[]byte("base"), []byte("oversize")}, 4, "resource_limit")
	failure(t, root, [][]byte{[]byte("base"), []byte("base")}, 2, "invalid_document")
	tooMany := make([][]byte, contract.MaxPaths+1)
	for i := range tooMany {
		tooMany[i] = []byte(fmt.Sprintf("p%d", i))
	}
	failure(t, root, tooMany, 4, "resource_limit")
	failure(t, root, [][]byte{[]byte(strings.Repeat("a", 256))}, 4, "resource_limit")
	names, err := os.ReadDir(root)
	if err != nil || len(names) != 10 {
		t.Fatalf("unexpected writes: %d %v", len(names), err)
	}
}

func TestUnavailablePermissionsAndExplicitSlash(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(root+"/private", []byte("private"), 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(root+"/private", 0600)
	if os.Geteuid() == 0 {
		t.Skip("privileged UID bypasses DAC; no capability changes authorized")
	}
	failure(t, root, [][]byte{[]byte("private")}, 3, "input_unavailable")
	if err := os.Mkdir(root+"/closed", 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(root+"/closed", 0700)
	failure(t, root+"/closed", nil, 3, "input_unavailable")
	got, err := acquire.Read(context.Background(), "/", nil)
	if err != nil || len(got) != 0 {
		t.Fatalf("explicit slash: %v %v", got, err)
	}
}
