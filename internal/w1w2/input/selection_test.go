package input_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

func selection(paths ...[]byte) []byte {
	encoded := make([]string, len(paths))
	for i, path := range paths {
		encoded[i] = base64.StdEncoding.EncodeToString(path)
	}
	raw, _ := json.Marshal(struct {
		Schema string   `json:"schema"`
		Paths  []string `json:"paths_b64"`
	}{contract.SelectionSchema, encoded})
	return raw
}

func TestSelectionRawNames(t *testing.T) {
	paths := [][]byte{[]byte("a b"), []byte("a\nb"), []byte("-name"), {0xff}, []byte("~*"), []byte("é"), []byte("e\u0301")}
	raw := selection(paths...)
	got, err := input.DecodeSelection(raw)
	if err != nil || !reflect.DeepEqual(got, paths) {
		t.Fatalf("%q %v", got, err)
	}
	raw[0] = '!'
	if !reflect.DeepEqual(got, paths) {
		t.Fatal("retained input alias")
	}
	got, err = input.DecodeSelection(selection())
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("empty: %q %v", got, err)
	}
}

func TestSelectionClosedJSON(t *testing.T) {
	prefix := `{"schema":"` + contract.SelectionSchema + `","paths_b64":`
	for _, tc := range []struct {
		name, raw string
		token     string
	}{
		{"null object", `null`, "invalid_document"},
		{"array object", `[]`, "invalid_document"},
		{"missing schema", `{"paths_b64":[]}`, "invalid_document"},
		{"missing paths", `{"schema":"x"}`, "invalid_document"},
		{"unknown", prefix + `[],"extra":0}`, "invalid_document"},
		{"duplicate", prefix + `[],"paths_b64":[]}`, "invalid_document"},
		{"escaped duplicate", prefix + `[],"paths_\u006264":[]}`, "invalid_document"},
		{"schema duplicate", prefix + `[],"schema":"x"}`, "invalid_document"},
		{"null paths", prefix + `null}`, "invalid_document"},
		{"string paths", prefix + `""}`, "invalid_document"},
		{"null item", prefix + `[null]}`, "invalid_document"},
		{"numeric item", prefix + `[1]}`, "invalid_document"},
		{"null schema", `{"schema":null,"paths_b64":[]}`, "invalid_document"},
		{"version", `{"schema":"v2","paths_b64":[]}`, "unsupported_version"},
		{"trailing", prefix + `[]} {}`, "invalid_document"},
		{"trailing garbage", prefix + `[]}!`, "invalid_document"},
		{"bad utf8", prefix + `["` + string([]byte{255}) + `"]}`, "invalid_document"},
		{"bad base64", prefix + `["*"]}`, "invalid_document"},
		{"unpadded", prefix + `["YQ"]}`, "invalid_document"},
		{"pad bits", prefix + `["YR=="]}`, "invalid_document"},
		{"line break", prefix + `["YQ==\n"]}`, "invalid_document"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := input.DecodeSelection([]byte(tc.raw))
			failure(t, err, 2, tc.token)
			if got != nil {
				t.Fatal("partial selection")
			}
		})
	}
	raw := append([]byte(" \n\t"), selection()...)
	raw = append(raw, '\n')
	if _, err := input.DecodeSelection(raw); err != nil {
		t.Fatal(err)
	}
}

func TestSelectionSharedGuardsAndWireBounds(t *testing.T) {
	for _, path := range []string{"", "/a", "a\x00b", "a//b", "a/../b", "./a", "a/.GiT/b"} {
		_, err := input.DecodeSelection(selection([]byte(path)))
		failure(t, err, 2, "unsafe_path")
	}
	_, err := input.DecodeSelection(selection([]byte("a"), []byte("a")))
	failure(t, err, 2, "invalid_document")
	_, err = input.DecodeSelection(selection(bytes.Repeat([]byte("a"), 256)))
	failure(t, err, 4, "resource_limit")
	paths := make([][]byte, contract.MaxPaths)
	for i := range paths {
		paths[i] = []byte(fmt.Sprintf("p%d", i))
	}
	if _, err = input.DecodeSelection(selection(paths...)); err != nil {
		t.Fatal(err)
	}
	paths = append(paths, []byte("extra"))
	_, err = input.DecodeSelection(selection(paths...))
	failure(t, err, 4, "resource_limit")
	raw := selection()
	raw = append(raw, bytes.Repeat([]byte(" "), contract.MaxSelectionBytes-len(raw))...)
	if _, err = input.DecodeSelection(raw); err != nil {
		t.Fatal(err)
	}
	_, err = input.DecodeSelection(append(raw, ' '))
	failure(t, err, 4, "resource_limit")
	path := strings.Repeat(strings.Repeat("a", 255)+"/", 15) + strings.Repeat("b", 255)
	if len(path) != contract.MaxPathBytes {
		t.Fatal(len(path))
	}
	if _, err = input.DecodeSelection(selection([]byte(path))); err != nil {
		t.Fatal(err)
	}
	_, err = input.DecodeSelection(selection([]byte(path + "/c")))
	failure(t, err, 4, "resource_limit")
}
