package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// Both native formats are required: unavailable capabilities fail, not skip.
func TestSnapshotDomainNativeRanges(t *testing.T) {
	for _, format := range []string{"sha1", "sha256"} {
		t.Run(format, func(t *testing.T) {
			f := domainBuild(t, format)
			width := 40
			if format == "sha256" {
				width = 64
			}
			if len(f.base) != width || len(f.head) != width {
				t.Fatalf("native revisions: %q %q", f.base, f.head)
			}
			entries := domainExpected(t, f)
			clean := domainAssert(t, f, f.base, f.head, entries)
			domainAssert(t, f, f.base, f.base, nil)
			// The fixture index initially holds head while HEAD/worktree still hold base.
			// Replace it with a third state, then independently dirty tracked/untracked bytes.
			domainGit(t, f.repo, "read-tree", f.base)
			domainWrite(t, f.repo, domainFile{"modified", "staged unrelated bytes\n", "100644"})
			domainGit(t, f.repo, "add", "--", "modified")
			if err := os.WriteFile(filepath.Join(f.repo, "modified"), []byte("unstaged\x00bytes"), 0600); err != nil {
				t.Fatal(err)
			}
			domainWrite(t, f.repo, domainFile{"untracked", "not committed\n", "100644"})
			dirty := domainAssert(t, f, f.base, f.head, entries)
			if !bytes.Equal(clean, dirty) {
				t.Fatal("dirty state altered fixed committed range")
			}
			domainAssert(t, f, f.head, f.head, nil)
		})
	}
}
