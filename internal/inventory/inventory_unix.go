//go:build aix || android || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package inventoryadapter

import (
	"errors"
	"sort"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
	"golang.org/x/sys/unix"
)

func Acquire(request Request) (evidence.UntrackedInventory, error) { return acquire(request, nil) }

func acquire(request Request, afterStat func(string)) (evidence.UntrackedInventory, error) {
	root, err := openRoot(request.Root, afterStat)
	if err != nil {
		return evidence.UntrackedInventory{}, err
	}
	defer unix.Close(root)
	paths := append([]string(nil), request.Paths...)
	sort.Strings(paths)
	records := make([]evidence.UntrackedRecord, 0, len(paths))
	for _, path := range paths {
		parts, err := pathParts(path)
		if err != nil {
			return evidence.UntrackedInventory{}, err
		}
		if err := inspect(root, parts, afterStat); err != nil {
			return evidence.UntrackedInventory{}, err
		}
		records = append(records, evidence.UntrackedRecord{Path: path})
	}
	return evidence.NewUntrackedInventory(records), nil
}

func openRoot(root string, afterStat func(string)) (int, error) {
	if !strings.HasPrefix(root, "/") {
		return -1, unavailable(evidence.InventoryInvalidPath)
	}
	parent, err := unix.Open("/", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, failure(err)
	}
	for _, name := range strings.Split(root[1:], "/") {
		if name == "" {
			if root == "/" {
				continue
			}
			unix.Close(parent)
			return -1, unavailable(evidence.InventoryInvalidPath)
		}
		next, err := checkedOpen(parent, name, true, afterStat)
		unix.Close(parent)
		if err != nil {
			return -1, err
		}
		parent = next
	}
	return parent, nil
}

func inspect(root int, parts []string, afterStat func(string)) error {
	parent, owned := root, -1
	defer func() {
		if owned >= 0 {
			unix.Close(owned)
		}
	}()
	for index, name := range parts {
		next, err := checkedOpen(parent, name, index+1 < len(parts), afterStat)
		if err != nil {
			return err
		}
		if owned >= 0 {
			unix.Close(owned)
		}
		if index+1 == len(parts) {
			if err := unix.Close(next); err != nil {
				return unavailable(evidence.InventoryUnverifiable)
			}
			return nil
		}
		owned, parent = next, next
	}
	return unavailable(evidence.InventoryInvalidPath)
}

func checkedOpen(parent int, name string, directory bool, afterStat func(string)) (int, error) {
	var before unix.Stat_t
	if err := unix.Fstatat(parent, name, &before, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return -1, failure(err)
	}
	if before.Mode&unix.S_IFMT == unix.S_IFLNK {
		return -1, unavailable(evidence.InventorySymlink)
	}
	if directory && before.Mode&unix.S_IFMT != unix.S_IFDIR || !directory && before.Mode&unix.S_IFMT != unix.S_IFREG {
		return -1, unavailable(evidence.InventorySpecial)
	}
	if afterStat != nil {
		afterStat(name)
	}
	flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC
	if directory {
		flags |= unix.O_DIRECTORY
	}
	fd, err := unix.Openat(parent, name, flags, 0)
	if err != nil {
		return -1, unavailable(evidence.InventoryRacing)
	}
	var opened, current unix.Stat_t
	if err := unix.Fstat(fd, &opened); err != nil {
		unix.Close(fd)
		return -1, unavailable(evidence.InventoryUnverifiable)
	}
	if !same(before, opened) {
		unix.Close(fd)
		return -1, unavailable(evidence.InventoryRacing)
	}
	if err := unix.Fstatat(parent, name, &current, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		unix.Close(fd)
		return -1, failure(err)
	}
	if !same(opened, current) {
		unix.Close(fd)
		return -1, unavailable(evidence.InventoryRacing)
	}
	return fd, nil
}

func pathParts(path string) ([]string, error) {
	if path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "\x00") {
		return nil, unavailable(evidence.InventoryInvalidPath)
	}
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return nil, unavailable(evidence.InventoryInvalidPath)
		}
	}
	return parts, nil
}

func same(left, right unix.Stat_t) bool {
	return left.Dev == right.Dev && left.Ino == right.Ino && left.Mode&unix.S_IFMT == right.Mode&unix.S_IFMT
}

func failure(err error) error {
	switch {
	case errors.Is(err, unix.ENOENT):
		return unavailable(evidence.InventoryMissing)
	case errors.Is(err, unix.ELOOP):
		return unavailable(evidence.InventorySymlink)
	default:
		return unavailable(evidence.InventoryUnverifiable)
	}
}

func unavailable(code evidence.InventoryUnavailableCode) error {
	return &evidence.InventoryUnavailableError{Code: code}
}
