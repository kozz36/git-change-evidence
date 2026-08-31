//go:build !(aix || android || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)

package inventoryadapter

import (
	"errors"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestAcquireReportsUnsupportedPlatform(t *testing.T) {
	inventory, err := Acquire(Request{})
	var unavailable *evidence.InventoryUnavailableError
	if len(inventory.Records()) != 0 || !errors.As(err, &unavailable) || unavailable.Code != evidence.InventoryUnsupportedPlatform {
		t.Fatalf("inventory, error = %#v, %v", inventory, err)
	}
}
