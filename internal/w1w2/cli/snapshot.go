package cli

import (
	"context"
	"os/exec"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

type factsSource func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error)

func nativeSnapshotFacts(ctx context.Context, request gitadapter.Request, bounds gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return evidence.SnapshotFactsV1{}, fail(3, "git_unavailable")
	}
	owner, err := gitadapter.AcquireSnapshotMetadataBounded(ctx, request, bounds)
	if err != nil {
		return evidence.SnapshotFactsV1{}, err
	}
	return owner.AcquireFacts()
}

// Snapshot composition is independent of document selectors and accounting policy.
func composeSnapshot(ctx context.Context, request contract.Request, source factsSource) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	bounds := gitadapter.SnapshotBounds{MaxResponseBytes: contract.MaxGitResponseBytes, MaxTotalBytes: contract.MaxGitBytes, MaxEntries: contract.MaxDiffEntries}
	facts, err := source(ctx, gitadapter.Request{Repository: request.Repository, Base: request.Base, Head: request.Head}, bounds)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	document, err := evidence.NewSnapshotV1(facts)
	return document.CanonicalBytes(), err
}
