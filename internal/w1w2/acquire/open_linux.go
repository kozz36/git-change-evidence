//go:build linux

package acquire

import (
	"context"
	"strings"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"golang.org/x/sys/unix"
)

func fail(status int, token string) error { return &contract.Failure{Status: status, Token: token} }
func interrupted(ctx context.Context) error {
	if ctx.Err() != nil {
		return fail(6, "interrupted")
	}
	return nil
}
func unavailable(err error) error {
	if err == unix.ENOSYS || err == unix.EOPNOTSUPP {
		return fail(3, "unsupported_platform")
	}
	return fail(3, "input_unavailable")
}
func same(a, b unix.Stat_t) bool {
	return a.Dev == b.Dev && a.Ino == b.Ino && a.Mode == b.Mode && a.Size == b.Size && a.Mtim == b.Mtim && a.Ctim == b.Ctim
}

// sameDirectoryBinding checks observed identity and mode, not directory entries.
// It is not an inode-generation or atomic-history guarantee.
func sameDirectoryBinding(a, b unix.Stat_t) bool {
	return a.Dev == b.Dev && a.Ino == b.Ino && a.Mode == b.Mode
}

func checkedOpen(parent int, name string, directory bool, ops acquisitionOps) (int, unix.Stat_t, error) {
	var before, opened, current unix.Stat_t
	if err := unix.Fstatat(parent, name, &before, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return -1, opened, unavailable(err)
	}
	kind := before.Mode & unix.S_IFMT
	if directory && kind == unix.S_IFLNK {
		return -1, opened, fail(2, "unsafe_path")
	}
	if directory && kind != unix.S_IFDIR || !directory && kind != unix.S_IFREG && kind != unix.S_IFLNK {
		return -1, opened, fail(3, "input_unavailable")
	}
	flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_NONBLOCK | unix.O_CLOEXEC
	if directory {
		flags |= unix.O_DIRECTORY
	} else if kind == unix.S_IFLNK {
		flags = unix.O_PATH | unix.O_NOFOLLOW | unix.O_CLOEXEC
	}
	fd, err := ops.openat(parent, name, flags, 0)
	if err != nil {
		if kind == unix.S_IFLNK && err == unix.EINVAL {
			return -1, opened, fail(3, "unsupported_platform")
		}
		if err == unix.ENOENT || err == unix.ELOOP || err == unix.ENOTDIR {
			return -1, opened, fail(7, "acquisition_race")
		}
		return -1, opened, unavailable(err)
	}
	sameCapture := same
	if directory {
		sameCapture = sameDirectoryBinding
	}
	if unix.Fstat(fd, &opened) != nil || unix.Fstatat(parent, name, &current, unix.AT_SYMLINK_NOFOLLOW) != nil || !sameCapture(before, opened) || !sameCapture(opened, current) {
		unix.Close(fd)
		return -1, opened, fail(7, "acquisition_race")
	}
	return fd, opened, nil
}

func openRoot(ctx context.Context, root string, ops acquisitionOps) (int, error) {
	if root == "" || root[0] != '/' {
		return -1, fail(2, "invalid_arguments")
	}
	if root != "/" {
		for _, part := range strings.Split(root[1:], "/") {
			if part == "" || part == "." || part == ".." || strings.ContainsRune(part, 0) {
				return -1, fail(2, "invalid_arguments")
			}
		}
		if err := contract.ValidatePath([]byte(root[1:])); err != nil {
			return -1, err
		}
	}
	fd, err := unix.Open("/", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, unavailable(err)
	}
	if root == "/" {
		return fd, nil
	}
	for _, part := range strings.Split(root[1:], "/") {
		if err := interrupted(ctx); err != nil {
			unix.Close(fd)
			return -1, err
		}
		next, _, err := checkedOpen(fd, part, true, ops)
		unix.Close(fd)
		if err != nil {
			return -1, err
		}
		fd = next
	}
	return fd, nil
}
