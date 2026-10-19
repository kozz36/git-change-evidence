package contract

import "time"

// These fixed shell admission limits do not change core accounting semantics or certify RSS.
const (
	MaxSelectionBytes   = 1 << 20
	MaxDocumentBytes    = 16 << 20
	MaxInputBytes       = 32 << 20
	MaxPaths            = 128
	MaxPathBytes        = 4095
	MaxTotalPathBytes   = 512 << 10
	MaxComponentBytes   = 255
	MaxFileBytes        = 1 << 20
	MaxContentBytes     = 8 << 20
	MaxGitResponseBytes = 16 << 20
	MaxGitBytes         = 64 << 20
	MaxDiffEntries      = 10000
	MaxOutputBytes      = 16 << 20
	OperationTimeout    = 30 * time.Second
	SelectionSchema     = "git-change-evidence.cli-path-selection/v1"
	ReportSubject       = "committed change evidence"
)
