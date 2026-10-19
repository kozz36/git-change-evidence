package changeevidence

import (
	"bytes"
	"sort"
)

// SnapshotSideV1 carries a present side's immutable-object identity and byte facts.
type SnapshotSideV1 struct {
	Path     []byte
	Mode     string
	Kind     GitEntryKind
	ObjectID GitObjectID
	Content  *SnapshotContentV1
}

// SnapshotNumstatV1 preserves original tokens separately from nullable counts.
type SnapshotNumstatV1 struct {
	Kind                 string
	AddedRaw, DeletedRaw string
	Added, Deleted       *uint64
}

// SnapshotEntryV1 describes both sides; nil sides mean valid absence.
type SnapshotEntryV1 struct {
	Status   string
	Old, New *SnapshotSideV1
	Numstat  SnapshotNumstatV1
}

// SnapshotFactsV1 contains validated, privately owned full-range fact records.
// Structural validation does not authenticate their acquisition from Git.
type SnapshotFactsV1 struct {
	revisions RevisionIdentity
	entries   []SnapshotEntryV1
}

// NewSnapshotFactsV1 validates and freezes typed facts without I/O or accounting.
func NewSnapshotFactsV1(revisions RevisionIdentity, entries []SnapshotEntryV1) (SnapshotFactsV1, error) {
	copied := cloneSnapshotEntriesV1(entries)
	if err := validateSnapshotFactsV1(revisions, copied); err != nil {
		return SnapshotFactsV1{}, err
	}
	sort.Slice(copied, func(i, j int) bool {
		return bytes.Compare(snapshotFactPathV1(copied[i]), snapshotFactPathV1(copied[j])) < 0
	})
	return SnapshotFactsV1{revisions: revisions, entries: copied}, nil
}

// Revisions returns the unchanged native revision identities.
func (facts SnapshotFactsV1) Revisions() RevisionIdentity { return facts.revisions }

// Entries returns independent copies of every mutable field and nested side.
func (facts SnapshotFactsV1) Entries() []SnapshotEntryV1 {
	return cloneSnapshotEntriesV1(facts.entries)
}

func snapshotFactPathV1(entry SnapshotEntryV1) []byte {
	if entry.New != nil {
		return entry.New.Path
	}
	return entry.Old.Path
}

func cloneSnapshotEntriesV1(entries []SnapshotEntryV1) []SnapshotEntryV1 {
	cloned := make([]SnapshotEntryV1, len(entries))
	for i, entry := range entries {
		entry.Old = cloneSnapshotSideV1(entry.Old)
		entry.New = cloneSnapshotSideV1(entry.New)
		entry.Numstat.Added = cloneSnapshotCountV1(entry.Numstat.Added)
		entry.Numstat.Deleted = cloneSnapshotCountV1(entry.Numstat.Deleted)
		cloned[i] = entry
	}
	return cloned
}

func cloneSnapshotSideV1(side *SnapshotSideV1) *SnapshotSideV1 {
	if side == nil {
		return nil
	}
	cloned := *side
	cloned.Path = append([]byte(nil), side.Path...)
	if side.Content != nil {
		content := *side.Content
		cloned.Content = &content
	}
	return &cloned
}

func cloneSnapshotCountV1(count *uint64) *uint64 {
	if count == nil {
		return nil
	}
	value := *count
	return &value
}
