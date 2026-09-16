//go:build linux

package main

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestCensusGoASTRejectsRootFileAndIntermediateSymlinks(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeCensusFile(t, root, "source.go", "package fixture\nfunc f() { os.Open(\"x\") }\n")
	writeCensusFile(t, outside, "outside.go", "package fixture\nfunc f() { os.Open(\"x\") }\n")
	writeCensusFile(t, root, "nested/source.go", "package fixture\nfunc f() { os.Open(\"x\") }\n")
	for name, test := range map[string]struct {
		sourceRoot string
		selected   string
	}{
		"root": {
			sourceRoot: filepath.Join(filepath.Dir(root), "root-link"),
			selected:   "source.go",
		},
		"file outside root": {
			sourceRoot: root,
			selected:   "outside-link.go",
		},
		"intermediate": {
			sourceRoot: root,
			selected:   "nested-link/source.go",
		},
	} {
		t.Run(name, func(t *testing.T) {
			switch name {
			case "root":
				if err := os.Symlink(root, test.sourceRoot); err != nil {
					t.Fatal(err)
				}
			case "file outside root":
				if err := os.Symlink(filepath.Join(outside, "outside.go"), filepath.Join(root, test.selected)); err != nil {
					t.Fatal(err)
				}
			case "intermediate":
				if err := os.Symlink(filepath.Join(root, "nested"), filepath.Join(root, "nested-link")); err != nil {
					t.Fatal(err)
				}
			}
			code, stdout, stderr := runCensusCommand(t, censusArgs(test.sourceRoot, "os", "Open", test.selected))
			if code != exitAbsent || len(stdout) != 0 || stderr != "content unavailable\n" {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout, stderr)
			}
		})
	}
}

func TestCensusGoASTRejectsSpecialFilesWithoutBlocking(t *testing.T) {
	root := t.TempDir()
	fifo := filepath.Join(root, "pipe.go")
	if err := unix.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "directory.go"), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, selected := range []string{"pipe.go", "directory.go"} {
		t.Run(selected, func(t *testing.T) {
			code, stdout, stderr := runCensusCommand(t, censusArgs(root, "os", "Open", selected))
			if code != exitAbsent || len(stdout) != 0 || stderr != "content unavailable\n" {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout, stderr)
			}
		})
	}
}
