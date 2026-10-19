package changeevidence

import "encoding/base64"

type snapshotWireV1 struct {
	Schema    ContractVersion         `json:"schema"`
	Revisions snapshotWireRevisionsV1 `json:"revisions"`
	Entries   []snapshotWireEntryV1   `json:"entries"`
}

type snapshotWireRevisionsV1 struct {
	Base string `json:"base"`
	Head string `json:"head"`
}

type snapshotWireEntryV1 struct {
	Status  string                `json:"status"`
	Old     *snapshotWireSideV1   `json:"old"`
	New     *snapshotWireSideV1   `json:"new"`
	Numstat snapshotWireNumstatV1 `json:"numstat"`
}

type snapshotWireSideV1 struct {
	PathB64  string                 `json:"path_b64"`
	Mode     string                 `json:"mode"`
	Kind     GitEntryKind           `json:"kind"`
	ObjectID GitObjectID            `json:"object_id"`
	Content  *snapshotWireContentV1 `json:"content"`
}

type snapshotWireContentV1 struct {
	ByteLength   uint64 `json:"byte_length"`
	ContainsNUL  bool   `json:"contains_nul"`
	LogicalLines uint64 `json:"logical_lines"`
}

type snapshotWireNumstatV1 struct {
	Kind       string  `json:"kind"`
	AddedRaw   string  `json:"added_raw"`
	DeletedRaw string  `json:"deleted_raw"`
	Added      *uint64 `json:"added"`
	Deleted    *uint64 `json:"deleted"`
}

func snapshotWireFromFactsV1(facts SnapshotFactsV1) snapshotWireV1 {
	wire := snapshotWireV1{
		Schema:    CommittedSnapshotContractV1,
		Revisions: snapshotWireRevisionsV1{Base: facts.revisions.Base, Head: facts.revisions.Head},
		Entries:   make([]snapshotWireEntryV1, len(facts.entries)),
	}
	for i, entry := range facts.entries {
		wire.Entries[i] = snapshotWireEntryV1{
			Status: entry.Status, Old: snapshotSideWireV1(entry.Old), New: snapshotSideWireV1(entry.New),
			Numstat: snapshotWireNumstatV1{
				Kind: entry.Numstat.Kind, AddedRaw: entry.Numstat.AddedRaw, DeletedRaw: entry.Numstat.DeletedRaw,
				Added: entry.Numstat.Added, Deleted: entry.Numstat.Deleted,
			},
		}
	}
	return wire
}

func snapshotSideWireV1(side *SnapshotSideV1) *snapshotWireSideV1 {
	if side == nil {
		return nil
	}
	wire := &snapshotWireSideV1{
		PathB64: base64.StdEncoding.EncodeToString(side.Path), Mode: side.Mode,
		Kind: side.Kind, ObjectID: side.ObjectID,
	}
	if side.Content != nil {
		wire.Content = &snapshotWireContentV1{
			ByteLength: side.Content.ByteLength, ContainsNUL: side.Content.ContainsNUL, LogicalLines: side.Content.LogicalLines,
		}
	}
	return wire
}
