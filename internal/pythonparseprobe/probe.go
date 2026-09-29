// Package pythonparseprobe is an isolated, non-production characterization of
// a pinned parser. A Match is not census evidence: the upstream AST has no end
// offsets, so the required physical half-open call/function spans are unproved.
package pythonparseprobe

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/tamnd/gopapy/v2/parser2"
)

type File struct {
	Path   string
	Source []byte
}

type Match struct {
	Path     string
	Start    int
	EndKnown bool
}

// Observe parses all explicit files before returning any matches. It does not
// discover paths, interpret imports, or reconstruct missing end positions.
func Observe(files []File, receiver, selector string) ([]Match, error) {
	var matches []Match
	seen := make(map[string]bool, len(files))
	for _, file := range files {
		if file.Path == "" || strings.ContainsRune(file.Path, 0) || path.IsAbs(file.Path) || path.Clean(file.Path) != file.Path || strings.HasPrefix(file.Path, "../") || file.Path == ".." || !(strings.HasSuffix(file.Path, ".py") || strings.HasSuffix(file.Path, ".pyi")) || seen[file.Path] {
			return nil, fmt.Errorf("invalid or duplicate path %q", file.Path)
		}
		seen[file.Path] = true
		module, err := parser2.ParseFile(file.Path, string(file.Source))
		if err != nil {
			return nil, err
		}
		walk(reflect.ValueOf(module), func(call *parser2.Call) {
			attr, ok := call.Func.(*parser2.Attribute)
			if !ok || attr.Attr != selector {
				return
			}
			name, ok := attr.Value.(*parser2.Name)
			if !ok || name.Id != receiver {
				return
			}
			offset, ok := byteOffset(file.Source, call.P)
			if ok {
				matches = append(matches, Match{Path: file.Path, Start: offset})
			}
		})
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Path != matches[j].Path {
			return matches[i].Path < matches[j].Path
		}
		return matches[i].Start < matches[j].Start
	})
	return matches, nil
}

func byteOffset(source []byte, p parser2.Pos) (int, bool) {
	if p.Line < 1 || p.Col < 0 {
		return 0, false
	}
	offset := 0
	for line := 1; line < p.Line; line++ {
		n := strings.IndexByte(string(source[offset:]), '\n')
		if n < 0 {
			return 0, false
		}
		offset += n + 1
	}
	offset += p.Col
	return offset, offset >= 0 && offset < len(source)
}

// walk is reflection-based only to inventory this candidate's many AST node
// variants. It is not an extractor design and has no bounded-depth guarantee.
func walk(v reflect.Value, visit func(*parser2.Call)) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Interface:
		if !v.IsNil() {
			walk(v.Elem(), visit)
		}
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		if v.CanInterface() {
			if call, ok := v.Interface().(*parser2.Call); ok {
				visit(call)
			}
		}
		walk(v.Elem(), visit)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			walk(v.Index(i), visit)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			walk(v.Field(i), visit)
		}
	}
}
