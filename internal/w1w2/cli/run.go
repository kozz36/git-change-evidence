// Package cli composes versioned W1/W2 evidence without legacy dispatch or delivery authority.
package cli

import (
	"context"
	"io"
	"os/exec"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/acquire"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

type snapshotSource func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.CommittedSnapshot, error)

const familyHelp = `Usage:
  gce w1w2 --help
  gce w1w2 snapshot --repo <absolute-repository> --base <revision> --head <revision>
  gce w1w2 inventory --source-root <absolute-directory> --policy <document> --paths <selection>
  gce w1w2 result --repo <absolute-repository> --base <revision> --head <revision> --policy <document> --inventory <document>
  gce w1w2 report --policy <document> --inventory <document> --result <document>
Snapshot takes no document selectors and does not read stdin.
Document selectors are explicit files or one stdin selector (-). Outputs are evidence, not delivery authority.
`

// Run executes one family operation with explicit streams and no output files.
// The deadline is cooperative; blocked OS/stdio and RSS have no hard guarantee.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runWithSource(ctx, args, stdin, stdout, stderr, nativeSnapshot)
}

func nativeSnapshot(ctx context.Context, request gitadapter.Request, bounds gitadapter.SnapshotBounds) (evidence.CommittedSnapshot, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return evidence.CommittedSnapshot{}, fail(3, "git_unavailable")
	}
	return gitadapter.AcquireBounded(ctx, request, bounds)
}

// The source is invocation-local; tests use the same complete composition path.
func runWithSource(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer, source snapshotSource) int {
	return runWithSources(ctx, args, stdin, stdout, stderr, source, nativeSnapshotFacts)
}

func runWithSources(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer, source snapshotSource, facts factsSource) int {
	if len(args) == 1 && args[0] == "--help" {
		return emit(ctx, "command", []byte(familyHelp), stdout, stderr)
	}
	request, err := input.Parse(args)
	stage := request.Stage
	if stage == "" {
		stage = "command"
	}
	if err != nil {
		return diagnostic(stage, err, stderr)
	}
	if !acquire.Supported() {
		return diagnostic(stage, fail(3, "unsupported_platform"), stderr)
	}
	ctx, cancel := context.WithTimeout(ctx, contract.OperationTimeout)
	defer cancel()
	var raw []byte
	if request.Stage == "snapshot" {
		raw, err = composeSnapshot(ctx, request, facts)
	} else {
		docs, readErr := input.ReadDocuments(ctx, request, stdin)
		if readErr != nil {
			return diagnostic(stage, readErr, stderr)
		}
		raw, err = compose(ctx, request, docs, source)
	}
	if err != nil {
		return diagnostic(stage, err, stderr)
	}
	return emit(ctx, stage, raw, stdout, stderr)
}

func emit(ctx context.Context, stage string, raw []byte, stdout, stderr io.Writer) int {
	if ctx.Err() != nil {
		return diagnostic(stage, fail(6, "interrupted"), stderr)
	}
	if len(raw) > contract.MaxOutputBytes {
		return diagnostic(stage, fail(4, "resource_limit"), stderr)
	}
	n, err := stdout.Write(raw)
	if err != nil || n != len(raw) {
		return diagnostic(stage, fail(9, "output_incomplete"), stderr)
	}
	return 0
}
