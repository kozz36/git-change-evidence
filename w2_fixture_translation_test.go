package changeevidence_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

const w2Base, w2Head = "0123456789abcdef0123456789abcdef01234567", "89abcdef0123456789abcdef0123456789abcdef"

type w2Entry struct {
	path                 string
	additions, deletions uint64
	countable            bool
}
type w2Vector struct {
	entries    []w2Entry
	categories []string
	totals     [4]uint64
	mechanical string
}
type w2Side struct {
	path, id, kind string
	loc            uint64
	hasLOC         bool
}

func w2Result(t *testing.T, vector w2Vector) evidence.ResultDocumentV1 {
	t.Helper()
	policy, err := evidence.NewAccountingPolicyV1(w2Policy(vector.mechanical))
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := evidence.NewInventoryV1(policy, nil)
	if err != nil {
		t.Fatal(err)
	}
	changes := make([]evidence.CommittedChange, len(vector.entries))
	for i, entry := range vector.entries {
		changes[i] = evidence.CommittedChange{Status: "A", Path: entry.path, Lines: evidence.CommittedLineCounts{Additions: entry.additions, Deletions: entry.deletions, Countable: entry.countable}}
	}
	result, err := evidence.NewAccountingResultV1(policy, inventory, evidence.NewCommittedSnapshot(w2Base, w2Head, changes))
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func w2Policy(mechanical string) evidence.PolicyInput {
	return evidence.PolicyInput{Categories: []evidence.CategoryInput{
		{Name: "artifacts", PathGlobs: [][]byte{[]byte("openspec/changes/**")}},
		{Name: "tests", PathGlobs: [][]byte{[]byte("tests/**")}},
		{Name: "mechanical", PathGlobs: [][]byte{[]byte(mechanical)}},
		{Name: "production", LineThreshold: 400},
	}, Default: "production"}
}

func w2Translate(raw []byte) (w2Vector, error) {
	root, err := w2Object(raw, "schema", "policy", "entries", "annotations", "expected")
	if err != nil || w2String(root["schema"]) != "git-change-evidence.sd7.compatibility/v1" {
		return w2Vector{}, fmt.Errorf("unsupported W2 fixture")
	}
	policy, err := w2Object(root["policy"], "version", "artifacts", "tests", "mechanical_bases")
	if err != nil || w2String(policy["version"]) != "v1" || !reflect.DeepEqual(w2Strings(policy["artifacts"]), []string{"openspec/changes/**"}) || !reflect.DeepEqual(w2Strings(policy["tests"]), []string{"tests/**"}) || !reflect.DeepEqual(w2Strings(policy["mechanical_bases"]), []string{"generated"}) {
		return w2Vector{}, fmt.Errorf("unsupported W2 policy")
	}
	entries, err := w2Array(root["entries"])
	if err != nil || len(entries) != 5 {
		return w2Vector{}, fmt.Errorf("unsupported W2 entries")
	}
	vector := w2Vector{entries: make([]w2Entry, 0, len(entries))}
	seenPaths, seenIDs := map[string]bool{}, map[string]bool{}
	for _, rawEntry := range entries {
		entry, err := w2Object(rawEntry, "status", "old", "new", "added", "deleted")
		if err != nil || w2String(entry["status"]) != "added" || !w2Null(entry["old"]) {
			return w2Vector{}, fmt.Errorf("unsupported W2 entry shape")
		}
		side, err := w2NewSide(entry["new"])
		if err != nil || seenPaths[side.path] || seenIDs[side.id] {
			return w2Vector{}, fmt.Errorf("invalid W2 entry identity")
		}
		seenPaths[side.path], seenIDs[side.id] = true, true
		if side.kind == "regular" {
			added, addErr := w2Uint(entry["added"])
			deleted, delErr := w2Uint(entry["deleted"])
			if addErr != nil || delErr != nil || deleted != 0 || !side.hasLOC || side.loc != added {
				return w2Vector{}, fmt.Errorf("inconsistent W2 text counts")
			}
			vector.entries = append(vector.entries, w2Entry{side.path, added, deleted, true})
		} else if side.kind == "binary" && w2Null(entry["added"]) && w2Null(entry["deleted"]) && !side.hasLOC {
			vector.entries = append(vector.entries, w2Entry{path: side.path})
		} else {
			return w2Vector{}, fmt.Errorf("unsupported W2 side")
		}
	}
	if !reflect.DeepEqual(seenIDs, map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true}) {
		return w2Vector{}, fmt.Errorf("unsupported W2 side placeholders")
	}
	annotations, err := w2Array(root["annotations"])
	if err != nil || len(annotations) != 1 {
		return w2Vector{}, fmt.Errorf("unsupported W2 annotations")
	}
	annotation, err := w2Object(annotations[0], "path", "added", "deleted", "basis", "evidence")
	path := w2String(annotation["path"])
	if err != nil || path == "" || !w2LiteralPath(path) || strings.ContainsAny(path, "*?[]\\") || w2String(annotation["basis"]) != "generated" || strings.TrimSpace(w2String(annotation["evidence"])) == "" {
		return w2Vector{}, fmt.Errorf("invalid W2 annotation")
	}
	added, addErr := w2Uint(annotation["added"])
	deleted, delErr := w2Uint(annotation["deleted"])
	if addErr != nil || delErr != nil {
		return w2Vector{}, fmt.Errorf("invalid W2 annotation counts")
	}
	for _, entry := range vector.entries {
		if entry.path == path && entry.countable && entry.additions == added && entry.deletions == deleted {
			vector.mechanical = path
		}
	}
	if vector.mechanical == "" {
		return w2Vector{}, fmt.Errorf("annotation does not cover one regular entry")
	}
	expected, err := w2Object(root["expected"], "categories", "totals", "non_gating_over_400")
	if err != nil || !w2Bool(expected["non_gating_over_400"], false) {
		return w2Vector{}, fmt.Errorf("unsupported W2 expected values")
	}
	vector.categories = w2Strings(expected["categories"])
	totals, err := w2Object(expected["totals"], "artifacts", "tests", "mechanical", "production")
	if err != nil || len(vector.categories) != len(vector.entries) {
		return w2Vector{}, fmt.Errorf("unsupported W2 expectations")
	}
	for i, name := range []string{"artifacts", "tests", "mechanical", "production"} {
		vector.totals[i], err = w2Uint(totals[name])
		if err != nil {
			return w2Vector{}, fmt.Errorf("invalid W2 total")
		}
	}
	return vector, nil
}

func w2NewSide(raw json.RawMessage) (w2Side, error) {
	fields, err := w2Object(raw, "raw_path_b64", "mode", "object_id", "kind", "loc")
	pathB64 := w2String(fields["raw_path_b64"])
	if err != nil || pathB64 == "" || w2String(fields["mode"]) != "100644" || !strings.Contains("abcde", w2String(fields["object_id"])) {
		return w2Side{}, fmt.Errorf("invalid W2 side")
	}
	path, err := base64.StdEncoding.Strict().DecodeString(pathB64)
	if err != nil || base64.StdEncoding.EncodeToString(path) != pathB64 || !w2LiteralPath(string(path)) {
		return w2Side{}, fmt.Errorf("invalid W2 raw path")
	}
	side := w2Side{path: string(path), id: w2String(fields["object_id"]), kind: w2String(fields["kind"])}
	if !w2Null(fields["loc"]) {
		side.loc, err = w2Uint(fields["loc"])
		side.hasLOC = err == nil
	}
	return side, err
}
