//go:build linux

package publicationadapter

import (
	"context"
	"errors"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

func publish(ctx context.Context, root string, bytes []byte, identity string) (Result, error) {
	if err := contextFailure(ctx); err != nil {
		return Result{}, err
	}
	rootFD, err := rootDirectory(root)
	if err != nil {
		return Result{}, err
	}
	defer unix.Close(rootFD)
	if err := contextFailure(ctx); err != nil {
		return Result{}, err
	}
	parentFD, err := shaDirectory(rootFD)
	if err != nil {
		return Result{}, err
	}
	defer unix.Close(parentFD)
	if err := contextFailure(ctx); err != nil {
		return Result{}, err
	}
	tempFD, err := stage(parentFD, bytes, openTemp, chmodTemp, unix.Write, unix.Fsync, unix.Close)
	if err != nil {
		return Result{}, err
	}
	defer unix.Close(tempFD)
	name := identity[len("sha256/"):]
	err = finalize(ctx, parentFD, name, func(parent int, name string) error {
		return unix.Linkat(unix.AT_FDCWD, "/proc/self/fd/"+strconv.Itoa(tempFD), parent, name, unix.AT_SYMLINK_FOLLOW)
	}, unix.Fsync, func(parent int, name string) error { return unix.Unlinkat(parent, name, 0) })
	if err != nil {
		return Result{}, err
	}
	return Result{Identity: identity}, nil
}
func contextFailure(ctx context.Context) error {
	if ctx.Err() != nil {
		return failure(Interrupted)
	}
	return nil
}
func rootDirectory(root string) (int, error) {
	if root == "" || root[0] != '/' || strings.ContainsRune(root, 0) {
		return -1, failure(Invalid)
	}
	fd, err := unix.Open("/", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, failure(Unavailable)
	}
	if root == "/" {
		return fd, nil
	}
	for _, name := range strings.Split(root[1:], "/") {
		if name == "" || name == "." || name == ".." {
			unix.Close(fd)
			return -1, failure(Invalid)
		}
		next, err := openDirectory(fd, name, nil)
		unix.Close(fd)
		if err != nil {
			return -1, err
		}
		fd = next
	}
	return fd, nil
}
func openDirectory(parent int, name string, afterOpen func()) (int, error) {
	var before, opened, current unix.Stat_t
	if unix.Fstatat(parent, name, &before, unix.AT_SYMLINK_NOFOLLOW) != nil || before.Mode&unix.S_IFMT != unix.S_IFDIR {
		return -1, failure(Unavailable)
	}
	fd, err := unix.Openat(parent, name, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, failure(Unavailable)
	}
	if unix.Fstat(fd, &opened) != nil || before.Dev != opened.Dev || before.Ino != opened.Ino || before.Mode&unix.S_IFMT != opened.Mode&unix.S_IFMT {
		unix.Close(fd)
		return -1, failure(Unavailable)
	}
	if afterOpen != nil {
		afterOpen()
	}
	if unix.Fstatat(parent, name, &current, unix.AT_SYMLINK_NOFOLLOW) != nil || opened.Dev != current.Dev || opened.Ino != current.Ino || opened.Mode&unix.S_IFMT != current.Mode&unix.S_IFMT {
		unix.Close(fd)
		return -1, failure(Unavailable)
	}
	return fd, nil
}
func shaDirectory(root int) (int, error) {
	var created, opened, current unix.Stat_t
	err := unix.Fstatat(root, "sha256", &created, unix.AT_SYMLINK_NOFOLLOW)
	if err == nil {
		return openDirectory(root, "sha256", nil)
	}
	if !errors.Is(err, unix.ENOENT) {
		return -1, failure(Unavailable)
	}
	if err = unix.Mkdirat(root, "sha256", 0o700); err != nil {
		if errors.Is(err, unix.EEXIST) {
			return openDirectory(root, "sha256", nil)
		}
		return -1, failure(Unavailable)
	}
	if unix.Fstatat(root, "sha256", &created, unix.AT_SYMLINK_NOFOLLOW) != nil || created.Mode&unix.S_IFMT != unix.S_IFDIR {
		return -1, failure(Unavailable)
	}
	fd, err := openDirectory(root, "sha256", nil)
	if err != nil {
		return -1, err
	}
	if unix.Fstat(fd, &opened) != nil || created.Dev != opened.Dev || created.Ino != opened.Ino || created.Mode&unix.S_IFMT != opened.Mode&unix.S_IFMT || unix.Fchmod(fd, 0o700) != nil || unix.Fstatat(root, "sha256", &current, unix.AT_SYMLINK_NOFOLLOW) != nil || created.Dev != current.Dev || created.Ino != current.Ino || created.Mode&unix.S_IFMT != current.Mode&unix.S_IFMT || unix.Fsync(root) != nil {
		unix.Close(fd)
		return -1, failure(Unavailable)
	}
	return fd, nil
}
func openTemp(dir int) (int, error) {
	return unix.Openat(dir, ".", unix.O_TMPFILE|unix.O_RDWR|unix.O_CLOEXEC, 0o600)
}
func chmodTemp(fd int) error { return unix.Fchmod(fd, 0o600) }
func stage(parent int, bytes []byte, open func(int) (int, error), chmod func(int) error, write func(int, []byte) (int, error), sync func(int) error, close func(int) error) (int, error) {
	fd, err := open(parent)
	if err != nil {
		return -1, failure(Unavailable)
	}
	if chmod(fd) != nil || writeAll(fd, bytes, write) != nil || sync(fd) != nil {
		close(fd)
		return -1, failure(Unavailable)
	}
	return fd, nil
}
func writeAll(fd int, bytes []byte, write func(int, []byte) (int, error)) error {
	for len(bytes) > 0 {
		n, err := write(fd, bytes)
		if errors.Is(err, unix.EINTR) {
			continue
		}
		if err != nil || n <= 0 || n > len(bytes) {
			return unix.EIO
		}
		bytes = bytes[n:]
	}
	return nil
}
func finalize(ctx context.Context, parent int, name string, link func(int, string) error, sync func(int) error, unlink func(int, string) error) error {
	if err := contextFailure(ctx); err != nil {
		return err
	}
	if err := link(parent, name); err != nil {
		if errors.Is(err, unix.EEXIST) {
			return failure(Conflict)
		}
		return failure(Unavailable)
	}
	if sync(parent) != nil {
		_ = unlink(parent, name)
		return failure(Unavailable)
	}
	return nil
}
