// Package livepublication provides opt-in observed-snapshot publication.
package livepublication

import (
	"context"
	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	publication "github.com/kozz36/git-change-evidence/internal/publication"
)

// Builder constructs complete evidence for immutable IDs, capturing caller policy and Subject.
type Builder func(base, head evidence.GitObjectID) (evidence.Evidence, error)

// Result and Error expose the existing publisher's technical outcomes.
type Result = publication.Result
type Error = publication.Error
type ErrorCode = publication.ErrorCode

const (
	Invalid     = publication.Invalid
	Unavailable = publication.Unavailable
	Conflict    = publication.Conflict
	Interrupted = publication.Interrupted
)

// Publish uses the last observed matching pair, not a Git/filesystem transaction.
func Publish(ctx context.Context, repo, baseRef, headRef, destination string, build Builder) (Result, error) {
	return publish(ctx, repo, baseRef, headRef, destination, build, publication.Publish)
}

func publish(ctx context.Context, repo, baseRef, headRef, destination string, build Builder, emit func(context.Context, string, evidence.Evidence) (Result, error)) (Result, error) {
	if ctx == nil || build == nil {
		return Result{}, &Error{Code: Invalid}
	}
	request := gitadapter.Request{Repository: repo, Base: baseRef, Head: headRef}
	for attempt := 0; attempt < 3; attempt++ {
		if ctx.Err() != nil {
			return Result{}, &Error{Code: Interrupted}
		}
		snapshot, err := gitadapter.Acquire(request)
		if err != nil {
			return Result{}, err
		}
		document, err := build(snapshot.Base, snapshot.Head)
		if err != nil {
			return Result{}, err
		}
		decoded, err := evidence.DecodeCanonical(document.CanonicalBytes())
		if err != nil || decoded.Digest() != document.Digest() || decoded.Provenance() != document.Provenance() {
			return Result{}, &Error{Code: Invalid}
		}
		revisions := document.Provenance().Revisions
		if revisions.Base != string(snapshot.Base) || revisions.Head != string(snapshot.Head) {
			return Result{}, &Error{Code: Invalid}
		}
		if ctx.Err() != nil {
			return Result{}, &Error{Code: Interrupted}
		}
		checked, err := gitadapter.Acquire(request)
		if err != nil {
			return Result{}, err
		}
		if checked.Base != snapshot.Base || checked.Head != snapshot.Head {
			continue
		}
		return emit(ctx, destination, document)
	}
	return Result{}, &evidence.SnapshotError{Code: evidence.SnapshotRacing}
}
