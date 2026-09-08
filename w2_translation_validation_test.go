package changeevidence_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
)

func w2Object(raw json.RawMessage, keys ...string) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	start, err := decoder.Token()
	if err != nil || start != json.Delim('{') {
		return nil, fmt.Errorf("object required")
	}
	allowed, fields := make(map[string]bool, len(keys)), make(map[string]json.RawMessage, len(keys))
	for _, key := range keys {
		allowed[key] = true
	}
	for decoder.More() {
		name, err := decoder.Token()
		key, ok := name.(string)
		if err != nil || !ok || !allowed[key] || fields[key] != nil {
			return nil, fmt.Errorf("closed object required")
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		fields[key] = value
	}
	end, err := decoder.Token()
	var trailing any
	if err != nil || end != json.Delim('}') || decoder.Decode(&trailing) != io.EOF || len(fields) != len(keys) {
		return nil, fmt.Errorf("closed object required")
	}
	return fields, nil
}

func w2Array(raw json.RawMessage) ([]json.RawMessage, error) {
	var values []json.RawMessage
	if w2Null(raw) || json.Unmarshal(raw, &values) != nil {
		return nil, fmt.Errorf("array required")
	}
	return values, nil
}

func w2String(raw json.RawMessage) string {
	var value string
	if w2Null(raw) || json.Unmarshal(raw, &value) != nil {
		return ""
	}
	return value
}

func w2Strings(raw json.RawMessage) []string {
	var values []string
	if w2Null(raw) || json.Unmarshal(raw, &values) != nil {
		return nil
	}
	return values
}

func w2Uint(raw json.RawMessage) (uint64, error) {
	trimmed := bytes.TrimSpace(raw)
	if w2Null(raw) || len(trimmed) == 0 || trimmed[0] < '0' || trimmed[0] > '9' {
		return 0, fmt.Errorf("unsigned count required")
	}
	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.UseNumber()
	var number json.Number
	if err := decoder.Decode(&number); err != nil || decoder.Decode(new(any)) != io.EOF {
		return 0, fmt.Errorf("unsigned count required")
	}
	return strconv.ParseUint(number.String(), 10, 64)
}

func w2Bool(raw json.RawMessage, want bool) bool {
	var value bool
	return !w2Null(raw) && json.Unmarshal(raw, &value) == nil && value == want
}

func w2Null(raw json.RawMessage) bool { return bytes.Equal(bytes.TrimSpace(raw), []byte("null")) }

func w2LiteralPath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.ContainsRune(path, 0) {
		return false
	}
	for _, part := range strings.Split(path, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func TestW2TranslatorRejectsClosedDomain(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"missing annotation key", func(v map[string]any) { delete(w2Annotation(v), "basis") }},
		{"extra annotation key", func(v map[string]any) { w2Annotation(v)["extra"] = true }},
		{"boolean count", func(v map[string]any) { w2Annotation(v)["added"] = true }},
		{"fractional count", func(v map[string]any) { w2Annotation(v)["added"] = 1.5 }},
		{"string count", func(v map[string]any) { w2Annotation(v)["added"] = "2" }},
		{"null count", func(v map[string]any) { w2Annotation(v)["added"] = nil }},
		{"negative count", func(v map[string]any) { w2Annotation(v)["added"] = -1 }},
		{"blank evidence", func(v map[string]any) { w2Annotation(v)["evidence"] = " \t" }},
		{"bad basis", func(v map[string]any) { w2Annotation(v)["basis"] = "manual" }},
		{"absent path", func(v map[string]any) { w2Annotation(v)["path"] = "missing.txt" }},
		{"duplicate annotation", func(v map[string]any) { v["annotations"] = append(v["annotations"].([]any), w2Annotation(v)) }},
		{"partial counts", func(v map[string]any) { w2Annotation(v)["added"] = 1.0 }},
		{"binary annotation", func(v map[string]any) { w2Annotation(v)["path"], w2Annotation(v)["added"] = "binary.bin", 0.0 }},
		{"metacharacter mechanical path", func(v map[string]any) {
			w2EntryMap(v, 2)["new"].(map[string]any)["raw_path_b64"], w2Annotation(v)["path"] = "Z2VuZXJhdGVkKi5qc29u", "generated*.json"
		}},
		{"unsupported policy", func(v map[string]any) { w2PolicyMap(v)["artifacts"] = []any{"docs/**"} }},
		{"unsupported status", func(v map[string]any) { w2EntryMap(v, 0)["status"] = "modified" }},
		{"old side", func(v map[string]any) { w2EntryMap(v, 0)["old"] = map[string]any{} }},
		{"side LOC mismatch", func(v map[string]any) { w2EntryMap(v, 0)["new"].(map[string]any)["loc"] = 4.0 }},
		{"noncanonical base64 path", func(v map[string]any) { w2EntryMap(v, 0)["new"].(map[string]any)["raw_path_b64"] = "b3BlbnNwZWM" }},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if _, err := w2Translate(w2Mutated(t, test.mutate)); err == nil {
				t.Fatal("closed W2 translator accepted unsupported input")
			}
		})
	}
}

func w2Mutated(t *testing.T, mutate func(map[string]any)) []byte {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(w2FixtureBytes(t), &value); err != nil {
		t.Fatal(err)
	}
	mutate(value)
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func w2Annotation(value map[string]any) map[string]any {
	return value["annotations"].([]any)[0].(map[string]any)
}
func w2EntryMap(value map[string]any, index int) map[string]any {
	return value["entries"].([]any)[index].(map[string]any)
}
func w2PolicyMap(value map[string]any) map[string]any { return value["policy"].(map[string]any) }
