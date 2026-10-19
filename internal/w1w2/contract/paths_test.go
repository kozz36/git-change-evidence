package contract

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestValidatePathPreservesRawNamesAndRejectsUnsafeComponents(t *testing.T) {
	maxPath := strings.Repeat(strings.Repeat("a", 255)+"/", 15) + strings.Repeat("b", 255)
	for _, tt := range []struct {
		name   string
		path   []byte
		status int
		token  string
	}{
		{"ordinary", []byte("a.txt"), 0, ""},
		{"metadata prefix is not metadata", []byte(".gitignore"), 0, ""},
		{"spaces and newline", []byte("a b/new\nline"), 0, ""},
		{"leading hyphen", []byte("-name"), 0, ""},
		{"non UTF8", []byte{'a', 0xff}, 0, ""},
		{"literal Linux backslash", []byte(`..\literal`), 0, ""},
		{"component at limit", []byte(strings.Repeat("a", 255)), 0, ""},
		{"path below limit", []byte(maxPath[:len(maxPath)-1]), 0, ""},
		{"path at limit", []byte(maxPath), 0, ""},
		{"empty", nil, 2, "unsafe_path"},
		{"absolute", []byte("/a"), 2, "unsafe_path"},
		{"parent", []byte("../a"), 2, "unsafe_path"},
		{"nested parent", []byte("a/../b"), 2, "unsafe_path"},
		{"dot", []byte("a/./b"), 2, "unsafe_path"},
		{"empty component", []byte("a//b"), 2, "unsafe_path"},
		{"trailing slash", []byte("a/"), 2, "unsafe_path"},
		{"NUL", []byte{'a', 0, 'b'}, 2, "unsafe_path"},
		{"metadata", []byte("a/.git/config"), 2, "unsafe_path"},
		{"ASCII metadata alias", []byte("a/.GiT/config"), 2, "unsafe_path"},
		{"component over limit", []byte(strings.Repeat("a", 256)), 4, "resource_limit"},
		// Every component remains admissible; only the whole-path limit is crossed.
		{"path over limit", []byte(maxPath[:len(maxPath)-1] + "/c"), 4, "resource_limit"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			before := append([]byte(nil), tt.path...)
			checkFailure(t, ValidatePath(tt.path), tt.status, tt.token)
			if !bytes.Equal(tt.path, before) {
				t.Fatal("validation mutated raw input")
			}
		})
	}
}

func TestValidatePathsCountsRawIdentityWithoutNormalization(t *testing.T) {
	maxPath := strings.Repeat(strings.Repeat("a", 255)+"/", 15) + strings.Repeat("b", 255)
	many := make([][]byte, 128)
	for i := range many {
		many[i] = []byte(maxPath)
		copy(many[i][:3], fmt.Sprintf("%03d", i))
	}
	// This crosses the count limit without crossing any byte limit.
	tooMany := make([][]byte, 129)
	for i := range tooMany {
		tooMany[i] = []byte(fmt.Sprintf("p%d", i))
	}
	for _, tt := range []struct {
		name   string
		paths  [][]byte
		status int
		token  string
	}{
		{"empty requested scope", [][]byte{}, 0, ""},
		{"raw identities", [][]byte{[]byte("é"), []byte("e\u0301"), []byte("File"), []byte("file")}, 0, ""},
		{"below path count", many[:127], 0, ""},
		{"maximum combined admissible paths", many, 0, ""},
		{"above path count", tooMany, 4, "resource_limit"},
		{"duplicate raw path", [][]byte{[]byte("a"), []byte("a")}, 2, "invalid_document"},
		{"invalid member", [][]byte{[]byte("a"), []byte("../b")}, 2, "unsafe_path"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			before := make([][]byte, len(tt.paths))
			for i := range before {
				before[i] = append([]byte(nil), tt.paths[i]...)
			}
			checkFailure(t, ValidatePaths(tt.paths), tt.status, tt.token)
			if !reflect.DeepEqual(tt.paths, before) {
				t.Fatal("validation rewrote, reordered or normalized selection")
			}
		})
	}
}

func checkFailure(t *testing.T, err error, status int, token string) {
	t.Helper()
	if status == 0 {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	var failure *Failure
	if !errors.As(err, &failure) || failure.Status != status || failure.Token != token || failure.Error() != token {
		t.Fatalf("error = %#v, want status %d token %q", err, status, token)
	}
}
