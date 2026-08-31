//go:build !(aix || android || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)

package inventoryadapter

import evidence "github.com/kozz36/git-change-evidence"

func Acquire(Request) (evidence.UntrackedInventory, error) {
	return evidence.UntrackedInventory{}, &evidence.InventoryUnavailableError{Code: evidence.InventoryUnsupportedPlatform}
}
