package changeevidence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"unicode/utf8"
)

// DecodeSnapshotV1 accepts only exact canonical committed-snapshot evidence.
// Rebuilding validates structure, not Git authenticity or report provenance.
func DecodeSnapshotV1(raw []byte) (SnapshotDocumentV1, error) {
	if !utf8.Valid(raw) {
		return SnapshotDocumentV1{}, &ContractError{"document", "invalid"}
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var wire snapshotWireV1
	if err := decoder.Decode(&wire); err != nil {
		return SnapshotDocumentV1{}, &ContractError{"document", "invalid"}
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		return SnapshotDocumentV1{}, &ContractError{"document", "invalid"}
	}
	if wire.Schema != CommittedSnapshotContractV1 {
		return SnapshotDocumentV1{}, &ContractError{"schema", "unsupported_version"}
	}
	entries := make([]SnapshotEntryV1, len(wire.Entries))
	for i, entry := range wire.Entries {
		old, err := snapshotSideFromWireV1(entry.Old)
		if err != nil {
			return SnapshotDocumentV1{}, err
		}
		newSide, err := snapshotSideFromWireV1(entry.New)
		if err != nil {
			return SnapshotDocumentV1{}, err
		}
		entries[i] = SnapshotEntryV1{
			Status: entry.Status, Old: old, New: newSide,
			Numstat: SnapshotNumstatV1{
				Kind: entry.Numstat.Kind, AddedRaw: entry.Numstat.AddedRaw, DeletedRaw: entry.Numstat.DeletedRaw,
				Added: entry.Numstat.Added, Deleted: entry.Numstat.Deleted,
			},
		}
	}
	facts, err := NewSnapshotFactsV1(RevisionIdentity{Base: wire.Revisions.Base, Head: wire.Revisions.Head}, entries)
	if err != nil {
		return SnapshotDocumentV1{}, err
	}
	document, err := NewSnapshotV1(facts)
	if err != nil {
		return SnapshotDocumentV1{}, err
	}
	if !bytes.Equal(raw, document.canonical) {
		return SnapshotDocumentV1{}, &ContractError{"document", "noncanonical"}
	}
	return document, nil
}

func snapshotSideFromWireV1(wire *snapshotWireSideV1) (*SnapshotSideV1, error) {
	if wire == nil {
		return nil, nil
	}
	path, err := base64.StdEncoding.DecodeString(wire.PathB64)
	if err != nil || base64.StdEncoding.EncodeToString(path) != wire.PathB64 {
		return nil, &ContractError{"entries.path_b64", "invalid"}
	}
	side := &SnapshotSideV1{Path: path, Mode: wire.Mode, Kind: wire.Kind, ObjectID: wire.ObjectID}
	if wire.Content != nil {
		side.Content = &SnapshotContentV1{
			ByteLength: wire.Content.ByteLength, ContainsNUL: wire.Content.ContainsNUL, LogicalLines: wire.Content.LogicalLines,
		}
	}
	return side, nil
}
