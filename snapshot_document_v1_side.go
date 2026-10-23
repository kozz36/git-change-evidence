package changeevidence

import "strings"

func validateSnapshotSideV1(side SnapshotSideV1, width int) error {
	if !validInventoryPath(side.Path) {
		return snapshotFactsErrorV1("entries.path", "invalid_path")
	}
	object := string(side.ObjectID)
	if len(object) != width || !isLowerHex(object) || strings.Trim(object, "0") == "" {
		return snapshotFactsErrorV1("entries.object_id", "invalid_object")
	}
	switch side.Mode {
	case "100644", "100755":
		if side.Kind != GitFile || side.Content == nil {
			return snapshotFactsErrorV1("entries.side", "invalid_side")
		}
		content := side.Content
		if content.ByteLength == 0 {
			if content.LogicalLines != 0 || content.ContainsNUL {
				return snapshotFactsErrorV1("entries.content", "invalid_content")
			}
		} else if content.LogicalLines == 0 || content.LogicalLines > content.ByteLength {
			return snapshotFactsErrorV1("entries.content", "invalid_content")
		}
	case "120000":
		if side.Kind != GitSymlink || side.Content != nil {
			return snapshotFactsErrorV1("entries.side", "invalid_side")
		}
	case "160000":
		if side.Kind != GitGitlink || side.Content != nil {
			return snapshotFactsErrorV1("entries.side", "invalid_side")
		}
	default:
		return snapshotFactsErrorV1("entries.side", "invalid_side")
	}
	return nil
}
