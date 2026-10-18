//go:build linux

package acquire

import (
	"context"

	"golang.org/x/sys/unix"
)

// readLink reads only the captured O_PATH symlink, including dangling links.
// A bounded extra byte detects overflow without admitting truncated payloads.
func readLink(ctx context.Context, parent int, name string, fd int, before unix.Stat_t, limit int) ([]byte, error) {
	if err := interrupted(ctx); err != nil {
		return nil, err
	}
	buf := make([]byte, limit+1)
	n, err := unix.Readlinkat(fd, "", buf)
	if err != nil {
		// Empty-name readlinkat on a captured symlink is a required capability.
		if err == unix.EINVAL || err == unix.ENOENT {
			return nil, fail(3, "unsupported_platform")
		}
		return nil, unavailable(err)
	}
	if n > limit {
		return nil, fail(4, "resource_limit")
	}
	var after, named unix.Stat_t
	if unix.Fstat(fd, &after) != nil || unix.Fstatat(parent, name, &named, unix.AT_SYMLINK_NOFOLLOW) != nil || !same(before, after) || !same(after, named) || int64(n) != after.Size {
		return nil, fail(7, "acquisition_race")
	}
	if err := interrupted(ctx); err != nil {
		return nil, err
	}
	return buf[:n], nil
}
