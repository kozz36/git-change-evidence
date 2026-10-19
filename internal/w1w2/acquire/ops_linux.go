//go:build linux

package acquire

import "golang.org/x/sys/unix"

// acquisitionOps is owned by one invocation, including its syscall handoffs.
// Native defaults are constructed per call; no process-wide hook is installed.
type acquisitionOps struct {
	openat func(int, string, int, uint32) (int, error)
	read   func(int, []byte) (int, error)
}

func nativeOps() acquisitionOps {
	return acquisitionOps{openat: unix.Openat, read: unix.Read}
}
