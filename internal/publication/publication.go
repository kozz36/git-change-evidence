package publicationadapter

import (
	"context"
	evidence "github.com/kozz36/git-change-evidence"
)

type Result struct{ Identity string }
type ErrorCode string

const (
	Invalid     ErrorCode = "invalid"
	Unavailable ErrorCode = "unavailable"
	Conflict    ErrorCode = "conflict"
	Interrupted ErrorCode = "interrupted"
)

type Error struct{ Code ErrorCode }

func (e *Error) Error() string { return "evidence publication: " + string(e.Code) }

// Publish materializes valid V1 evidence below an existing absolute root; failed rollback has no broader power-loss guarantee.
func Publish(ctx context.Context, root string, document evidence.Evidence) (Result, error) {
	if ctx == nil {
		return Result{}, failure(Invalid)
	}
	if ctx.Err() != nil {
		return Result{}, failure(Interrupted)
	}
	digest := string(document.Digest())
	if len(digest) != 64 {
		return Result{}, failure(Invalid)
	}
	return publish(ctx, root, document.CanonicalBytes(), "sha256/"+digest+".json")
}
func failure(code ErrorCode) error { return &Error{Code: code} }
