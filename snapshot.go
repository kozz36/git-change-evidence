package changeevidence

type GitObjectID string
type GitEntryKind string

const (
	GitAbsent  GitEntryKind = "absent"
	GitFile    GitEntryKind = "file"
	GitSymlink GitEntryKind = "symlink"
	GitGitlink GitEntryKind = "gitlink"
	GitOther   GitEntryKind = "other"
)

type CommittedChange struct {
	Status               string
	Path, PreviousPath   string
	OldObject, NewObject GitObjectID
	OldMode, NewMode     string
	OldKind, NewKind     GitEntryKind
	Binary               bool
}

type CommittedSnapshot struct {
	Base, Head GitObjectID
	entries    []CommittedChange
}

func NewCommittedSnapshot(base, head GitObjectID, entries []CommittedChange) CommittedSnapshot {
	return CommittedSnapshot{Base: base, Head: head, entries: append([]CommittedChange(nil), entries...)}
}

func (s CommittedSnapshot) Entries() []CommittedChange {
	return append([]CommittedChange(nil), s.entries...)
}

type SnapshotErrorCode string

const (
	SnapshotUnresolved SnapshotErrorCode = "unresolved_revision"
	SnapshotForeign    SnapshotErrorCode = "foreign_repository"
	SnapshotMissing    SnapshotErrorCode = "missing_object"
	SnapshotRacing     SnapshotErrorCode = "racing_object"
)

type SnapshotError struct{ Code SnapshotErrorCode }

func (e *SnapshotError) Error() string { return "snapshot acquisition: " + string(e.Code) }
