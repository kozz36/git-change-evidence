//go:build !linux

package acquire

import (
	"context"
	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

// Supported is false where descriptor-based acquisition is unavailable.
func Supported() bool { return false }

// Read fails without acquisition or a weaker pathname fallback.
func Read(context.Context, string, [][]byte) ([]evidence.InventoryEntryInput, error) {
	return nil, &contract.Failure{Status: 3, Token: "unsupported_platform"}
}
