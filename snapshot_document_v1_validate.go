package changeevidence

import (
	"bytes"
	"strconv"
)

func snapshotFactsErrorV1(field, code string) error {
	return &ContractError{Field: field, Code: code}
}

func validateSnapshotFactsV1(revisions RevisionIdentity, entries []SnapshotEntryV1) error {
	if !isRevision(revisions.Base) {
		return snapshotFactsErrorV1("revisions.base", "invalid_revision")
	}
	if !isRevision(revisions.Head) || len(revisions.Base) != len(revisions.Head) {
		return snapshotFactsErrorV1("revisions.head", "invalid_revision")
	}
	if revisions.Base == revisions.Head && len(entries) != 0 {
		return snapshotFactsErrorV1("entries", "inconsistent_range")
	}
	paths := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if err := validateSnapshotEntryV1(entry, len(revisions.Base)); err != nil {
			return err
		}
		key := string(snapshotFactPathV1(entry))
		if _, duplicate := paths[key]; duplicate {
			return snapshotFactsErrorV1("entries.path", "duplicate_path")
		}
		paths[key] = struct{}{}
	}
	return nil
}

func validateSnapshotEntryV1(entry SnapshotEntryV1, width int) error {
	validSides := false
	switch entry.Status {
	case "A":
		validSides = entry.Old == nil && entry.New != nil
	case "D":
		validSides = entry.Old != nil && entry.New == nil
	case "M", "R", "T", "C":
		validSides = entry.Old != nil && entry.New != nil
	default:
		return snapshotFactsErrorV1("entries.status", "invalid_status")
	}
	if !validSides {
		return snapshotFactsErrorV1("entries.side", "invalid_side")
	}
	for _, side := range []*SnapshotSideV1{entry.Old, entry.New} {
		if side != nil {
			if err := validateSnapshotSideV1(*side, width); err != nil {
				return err
			}
		}
	}
	if entry.Old != nil && entry.New != nil {
		samePath := bytes.Equal(entry.Old.Path, entry.New.Path)
		if samePath != (entry.Status == "M" || entry.Status == "T") {
			return snapshotFactsErrorV1("entries.path_pair", "invalid_path_pair")
		}
		sameKind := entry.Old.Kind == entry.New.Kind
		if (entry.Status == "M" && !sameKind) || (entry.Status == "T" && sameKind) {
			return snapshotFactsErrorV1("entries.side", "invalid_side")
		}
	}
	if err := validateSnapshotNumstatV1(entry.Numstat); err != nil {
		return err
	}
	if entry.Status == "M" && entry.Old.Mode != entry.New.Mode && entry.Old.ObjectID == entry.New.ObjectID {
		if entry.Numstat.Kind != "text" || *entry.Numstat.Added != 0 || *entry.Numstat.Deleted != 0 {
			return snapshotFactsErrorV1("entries.numstat", "invalid_numstat")
		}
	}
	return nil
}

func validateSnapshotNumstatV1(value SnapshotNumstatV1) error {
	switch value.Kind {
	case "text":
		if value.Added != nil && value.Deleted != nil &&
			value.AddedRaw == strconv.FormatUint(*value.Added, 10) &&
			value.DeletedRaw == strconv.FormatUint(*value.Deleted, 10) {
			return nil
		}
	case "binary":
		if value.AddedRaw == "-" && value.DeletedRaw == "-" && value.Added == nil && value.Deleted == nil {
			return nil
		}
	}
	return snapshotFactsErrorV1("entries.numstat", "invalid_numstat")
}
