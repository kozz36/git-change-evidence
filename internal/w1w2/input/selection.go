package input

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"unicode/utf8"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

// DecodeSelection preserves raw path identities and rejects open or ambiguous JSON.
// Empty selection is not certification of root access or directory completeness.
func DecodeSelection(raw []byte) ([][]byte, error) {
	invalid := func() ([][]byte, error) { return nil, fail(2, "invalid_document") }
	if len(raw) > contract.MaxSelectionBytes {
		return nil, fail(4, "resource_limit")
	}
	if !utf8.Valid(raw) {
		return invalid()
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	start, err := decoder.Token()
	if err != nil || start != json.Delim('{') {
		return invalid()
	}
	seen := make(map[string]bool, 2)
	var schema string
	var encoded []json.RawMessage
	for decoder.More() {
		token, err := decoder.Token()
		key, ok := token.(string)
		if err != nil || !ok || seen[key] {
			return invalid()
		}
		seen[key] = true
		switch key {
		case "schema":
			var value json.RawMessage
			if decoder.Decode(&value) != nil || len(value) == 0 || value[0] != '"' || json.Unmarshal(value, &schema) != nil {
				return invalid()
			}
		case "paths_b64":
			var value json.RawMessage
			if decoder.Decode(&value) != nil || len(value) == 0 || value[0] != '[' || json.Unmarshal(value, &encoded) != nil {
				return invalid()
			}
		default:
			return invalid()
		}
	}
	end, err := decoder.Token()
	if err != nil || end != json.Delim('}') || !seen["schema"] || !seen["paths_b64"] {
		return invalid()
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		return invalid()
	}
	if schema != contract.SelectionSchema {
		return nil, fail(2, "unsupported_version")
	}
	if len(encoded) > contract.MaxPaths {
		return nil, fail(4, "resource_limit")
	}
	paths := make([][]byte, 0, len(encoded))
	for _, value := range encoded {
		var text string
		if len(value) == 0 || value[0] != '"' || json.Unmarshal(value, &text) != nil {
			return invalid()
		}
		path, err := base64.StdEncoding.DecodeString(text)
		if err != nil || base64.StdEncoding.EncodeToString(path) != text {
			return invalid()
		}
		paths = append(paths, path)
	}
	if err := contract.ValidatePaths(paths); err != nil {
		return nil, err
	}
	return paths, nil
}
