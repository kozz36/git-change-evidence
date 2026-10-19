package gitadapter

import (
	"context"
	"errors"
	"io"
	"os/exec"

	evidence "github.com/kozz36/git-change-evidence"
)

// SnapshotBounds limits each response, all responses combined, and full-range entries.
// All limits must be positive. They are admission limits, not an RSS guarantee.
type SnapshotBounds struct{ MaxResponseBytes, MaxTotalBytes, MaxEntries int }

// ErrSnapshotResourceLimit identifies a whole-range acquisition budget failure.
var ErrSnapshotResourceLimit = errors.New("snapshot resource limit")

// ErrSnapshotInvalidBounds identifies nonpositive snapshot budgets.
var ErrSnapshotInvalidBounds = errors.New("invalid snapshot bounds")

type snapshotCommands struct {
	ctx                 context.Context
	runner              Runner
	repository          string
	response, remaining int
}

func (commands *snapshotCommands) output(code evidence.SnapshotErrorCode, args ...string) ([]byte, error) {
	if err := commands.ctx.Err(); err != nil {
		return nil, err
	}
	cap := min(commands.response, commands.remaining)
	output, err := commands.runner(commands.ctx, commands.repository, cap, args...)
	if contextErr := commands.ctx.Err(); contextErr != nil {
		return nil, contextErr
	}
	if len(output) > cap || errors.Is(err, ErrSnapshotResourceLimit) {
		return nil, ErrSnapshotResourceLimit
	}
	commands.remaining -= len(output)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, failure(code)
	}
	return output, nil
}

// snapshotCommand owns its subprocess, discards raw stderr, and uses a closed
// environment. Local roots/configuration and stable mounts remain trusted; this
// is not protection against a malicious invoking process or global Git setup.
func snapshotCommand(ctx context.Context, repository string, byteCap int, args ...string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	arguments := append([]string{"--no-replace-objects", "-C", repository}, args...)
	command := exec.CommandContext(ctx, "git", arguments...)
	command.Env = []string{
		"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1",
		"GIT_CONFIG_COUNT=0", "GIT_NO_REPLACE_OBJECTS=1", "GIT_NO_LAZY_FETCH=1",
		"GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0", "LC_ALL=C", "LANG=C", "PATH=/usr/bin:/bin",
	}
	output := boundedOutput{cap: byteCap}
	command.Stdout, command.Stderr = &output, io.Discard
	err := command.Run()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if output.exceeded {
		return nil, ErrSnapshotResourceLimit
	}
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
