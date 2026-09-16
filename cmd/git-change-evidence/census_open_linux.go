//go:build linux

package main

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

const maxCensusReadBytes = uint64(^uint64(0)>>1) - 1

type censusLinuxRoot struct {
	fd int
}

func openCensusSourceRoot(name string) (censusSourceRoot, error) {
	parts, err := censusAbsoluteRootParts(name)
	if err != nil {
		return nil, errCensusUnavailable
	}
	fd, err := unix.Open("/", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, errCensusUnavailable
	}
	for _, part := range parts {
		next, _, err := openCensusEntry(fd, part, true)
		unix.Close(fd)
		if err != nil {
			return nil, err
		}
		fd = next
	}
	return &censusLinuxRoot{fd: fd}, nil
}

func (root *censusLinuxRoot) Close() error {
	if root == nil || root.fd < 0 {
		return nil
	}
	err := unix.Close(root.fd)
	root.fd = -1
	return err
}

func (root *censusLinuxRoot) ReadFile(name string, maxBytes uint64) ([]byte, error) {
	if root == nil || root.fd < 0 || maxBytes > maxCensusReadBytes {
		return nil, errCensusBound
	}
	parts, err := censusRelativePathParts(name)
	if err != nil {
		return nil, errCensusUnavailable
	}
	parent, owned := root.fd, -1
	for _, part := range parts[:len(parts)-1] {
		next, _, err := openCensusEntry(parent, part, true)
		if owned >= 0 {
			unix.Close(owned)
		}
		if err != nil {
			return nil, err
		}
		parent, owned = next, next
	}
	fd, stat, err := openCensusEntry(parent, parts[len(parts)-1], false)
	if owned >= 0 {
		unix.Close(owned)
	}
	if err != nil {
		return nil, err
	}
	if stat.Size < 0 || uint64(stat.Size) > maxBytes {
		unix.Close(fd)
		return nil, errCensusBound
	}
	file := os.NewFile(uintptr(fd), name)
	if file == nil {
		unix.Close(fd)
		return nil, errCensusUnavailable
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, int64(maxBytes)+1))
	if err != nil {
		return nil, errCensusUnavailable
	}
	if uint64(len(content)) > maxBytes {
		return nil, errCensusBound
	}
	return content, nil
}

func openCensusEntry(parent int, name string, directory bool) (int, unix.Stat_t, error) {
	var before unix.Stat_t
	if unix.Fstatat(parent, name, &before, unix.AT_SYMLINK_NOFOLLOW) != nil {
		return -1, unix.Stat_t{}, errCensusUnavailable
	}
	kind := before.Mode & unix.S_IFMT
	if kind == unix.S_IFLNK {
		return -1, unix.Stat_t{}, errCensusUnavailable
	}
	if directory && kind != unix.S_IFDIR || !directory && kind != unix.S_IFREG {
		return -1, unix.Stat_t{}, errCensusUnavailable
	}
	flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_NONBLOCK | unix.O_CLOEXEC
	if directory {
		flags |= unix.O_DIRECTORY
	}
	fd, err := unix.Openat(parent, name, flags, 0)
	if err != nil {
		return -1, unix.Stat_t{}, errCensusUnavailable
	}
	var opened, current unix.Stat_t
	if unix.Fstat(fd, &opened) != nil || !sameCensusEntry(before, opened) || unix.Fstatat(parent, name, &current, unix.AT_SYMLINK_NOFOLLOW) != nil || !sameCensusEntry(opened, current) {
		unix.Close(fd)
		return -1, unix.Stat_t{}, errCensusUnavailable
	}
	return fd, opened, nil
}

func sameCensusEntry(left, right unix.Stat_t) bool {
	return left.Dev == right.Dev && left.Ino == right.Ino && left.Mode&unix.S_IFMT == right.Mode&unix.S_IFMT
}
