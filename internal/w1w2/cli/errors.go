package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func fail(status int, token string) error { return &contract.Failure{Status: status, Token: token} }

func diagnostic(stage string, err error, stderr io.Writer) int {
	status, token := failureData(err)
	// Failed diagnostic delivery cannot become success; never echo raw errors.
	fmt.Fprintf(stderr, "w1w2 %s: %s\n", stage, token)
	return status
}

func failureData(err error) (int, string) {
	var failure *contract.Failure
	if errors.As(err, &failure) {
		return failure.Status, failure.Token
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return 6, "interrupted"
	}
	if errors.Is(err, gitadapter.ErrSnapshotResourceLimit) {
		return 4, "resource_limit"
	}
	if errors.Is(err, gitadapter.ErrSnapshotInvalidBounds) {
		return 2, "invalid_arguments"
	}
	var snapshot *evidence.SnapshotError
	if errors.As(err, &snapshot) {
		switch snapshot.Code {
		case evidence.SnapshotForeign:
			return 2, "invalid_arguments"
		case evidence.SnapshotRacing:
			return 7, "acquisition_race"
		default:
			return 3, "git_unavailable"
		}
	}
	var invalid *evidence.ContractError
	if errors.As(err, &invalid) {
		if invalid.Code == "unsupported_version" {
			return 2, "unsupported_version"
		}
		if invalid.Code == "mismatch" || invalid.Code == "mismatched_digest" || invalid.Code == "mismatched_revision" {
			return 2, "antecedent_mismatch"
		}
		return 2, "invalid_document"
	}
	return 3, "git_unavailable"
}
