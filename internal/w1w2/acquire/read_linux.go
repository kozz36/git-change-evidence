//go:build linux

// Package acquire reads selected regular files and leaf symlink payloads without discovery or writes.
// It confines names to captured descriptors, not physical inode provenance, and
// detects metadata changes without promising an atomic snapshot or every write.
package acquire

import (
	"context"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"golang.org/x/sys/unix"
)

// Supported reports availability of the Linux descriptor-based implementation.
// Runtime capability failures are reported by Read without a pathname fallback.
func Supported() bool { return true }

// Read returns owned path/content copies for regular files or leaf symlink payloads.
// Symlink targets are opaque content, never followed; intermediate links are rejected.
// Errors return no partial evidence; cancellation is cooperative between OS calls.
func Read(ctx context.Context, root string, paths [][]byte) ([]evidence.InventoryEntryInput, error) {
	return readWithOps(ctx, root, paths, nativeOps())
}

func readWithOps(ctx context.Context, root string, paths [][]byte, ops acquisitionOps) ([]evidence.InventoryEntryInput, error) {
	if err := contract.ValidatePaths(paths); err != nil {
		return nil, err
	}
	if err := interrupted(ctx); err != nil {
		return nil, err
	}
	fd, err := openRoot(ctx, root, ops)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)
	entries := make([]evidence.InventoryEntryInput, 0, len(paths))
	remaining := contract.MaxContentBytes
	for _, path := range paths {
		if err := interrupted(ctx); err != nil {
			return nil, err
		}
		content, err := readSelected(ctx, fd, string(path), remaining, ops)
		if err != nil {
			return nil, err
		}
		remaining -= len(content)
		entries = append(entries, evidence.InventoryEntryInput{Path: append([]byte(nil), path...), Content: content})
	}
	if err := interrupted(ctx); err != nil {
		return nil, err
	}
	return entries, nil
}

func readSelected(ctx context.Context, root int, path string, remaining int, ops acquisitionOps) ([]byte, error) {
	parts := strings.Split(path, "/")
	parent, owned := root, -1
	defer func() {
		if owned >= 0 {
			unix.Close(owned)
		}
	}()
	for _, part := range parts[:len(parts)-1] {
		if err := interrupted(ctx); err != nil {
			return nil, err
		}
		next, _, err := checkedOpen(parent, part, true, ops)
		if err != nil {
			return nil, err
		}
		if owned >= 0 {
			unix.Close(owned)
		}
		parent, owned = next, next
	}
	name := parts[len(parts)-1]
	fd, before, err := checkedOpen(parent, name, false, ops)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)
	limit := min(contract.MaxFileBytes, remaining)
	if before.Size < 0 || before.Size > int64(limit) {
		return nil, fail(4, "resource_limit")
	}
	if before.Mode&unix.S_IFMT == unix.S_IFLNK {
		return readLink(ctx, parent, name, fd, before, limit)
	}
	content := make([]byte, 0, int(before.Size))
	var buf [32768]byte
	for {
		if err := interrupted(ctx); err != nil {
			return nil, err
		}
		// At most one overflow byte is probed; it is never admitted as content.
		n, err := ops.read(fd, buf[:min(len(buf), limit-len(content)+1)])
		if n > limit-len(content) {
			return nil, fail(4, "resource_limit")
		}
		if n > 0 {
			content = append(content, buf[:n]...)
		}
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return nil, unavailable(err)
		}
		if n == 0 {
			break
		}
	}
	var after, named unix.Stat_t
	if unix.Fstat(fd, &after) != nil || unix.Fstatat(parent, name, &named, unix.AT_SYMLINK_NOFOLLOW) != nil || !same(before, after) || !same(after, named) || int64(len(content)) != after.Size {
		return nil, fail(7, "acquisition_race")
	}
	if err := interrupted(ctx); err != nil {
		return nil, err
	}
	return content, nil
}
