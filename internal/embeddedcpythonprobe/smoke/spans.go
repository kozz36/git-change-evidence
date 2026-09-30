package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	python "github.com/goccy/go-python"
)

// Private experiment types, not a production inventory or wire envelope.
type probeInput struct {
	Path  string
	Bytes []byte
}
type probeSpan struct {
	Start, End int
	Bytes      []byte
	Hash       [32]byte
}
type probeCall struct{ Call, Function probeSpan }
type probeFile struct {
	Path  string
	Hash  [32]byte
	Calls []probeCall
}

const spanHelper = `import ast, json, io, tokenize
def _probe_spans(source):
    tree = ast.parse(source, filename="<byte-probe>", mode="exec")
    compile(source, "<byte-probe>", "exec")
    encoding, _ = tokenize.detect_encoding(io.BytesIO(source).readline)
    if encoding not in ("utf-8", "utf-8-sig", "utf8"):
        raise ValueError("unsupported span mapping encoding: " + encoding)
    positions = []
    for node in ast.walk(tree):
        if (isinstance(node, ast.Call)
            and isinstance(node.func, ast.Attribute)
            and isinstance(node.func.value, ast.Name)
            and node.func.value.id == "os"
            and node.func.attr == "open"):
            positions.append([[n.lineno, n.col_offset, n.end_lineno, n.end_col_offset]
                              for n in (node, node.func)])
    positions.sort()
    return json.dumps(positions)
`

// Only UTF-8 span mapping is experimentally supported. Native Python bytes
// parsing/compilation and tokenizer encoding checks happen before matches.
func physicalLines(raw []byte) ([]int, error) {
	start := 0
	if strings.HasPrefix(string(raw), "\xef\xbb\xbf") {
		start = 3
	}
	if !utf8.Valid(raw[start:]) {
		return nil, fmt.Errorf("unsupported non-UTF8 bytes")
	}
	lines := []int{start}
	for i := start; i < len(raw); i++ {
		if raw[i] == '\r' {
			if i+1 < len(raw) && raw[i+1] == '\n' {
				i++
			}
			lines = append(lines, i+1)
		} else {
			if raw[i] == '\n' {
				lines = append(lines, i+1)
			}
		}
	}
	return lines, nil
}

func physicalSpan(raw []byte, lines []int, pos [4]int) (probeSpan, error) {
	offset := func(line, column int) (int, error) {
		if line < 1 || line > len(lines) || column < 0 {
			return 0, fmt.Errorf("invalid AST position")
		}
		begin := lines[line-1]
		end := len(raw)
		if line < len(lines) {
			end = lines[line]
		}
		// AST columns count UTF-8 bytes, but CRLF normalization and a UTF8 BOM
		// change original absolute offsets. Each physical line anchors the mapping.
		at := begin + column
		if at > end || (at < len(raw) && !utf8.RuneStart(raw[at])) {
			return 0, fmt.Errorf("AST column outside physical line")
		}
		return at, nil
	}
	a, err := offset(pos[0], pos[1])
	if err != nil {
		return probeSpan{}, err
	}
	b, err := offset(pos[2], pos[3])
	if err != nil {
		return probeSpan{}, err
	}
	if b < a {
		return probeSpan{}, fmt.Errorf("reversed AST span")
	}
	fragment := append([]byte(nil), raw[a:b]...)
	return probeSpan{a, b, fragment, sha256.Sum256(fragment)}, nil
}

// A failed file discards the entire batch; parsing and compilation happen in
// the trusted helper before it walks matches. Source is passed only as data.
func probeBatch(ctx context.Context, p *python.Python, inputs []probeInput) ([]probeFile, error) {
	result, err := p.Eval(ctx, spanHelper)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, result.Error
	}
	module, err := p.Import(ctx, "__main__")
	if err != nil {
		return nil, err
	}
	files := make([]probeFile, 0, len(inputs))
	for _, input := range inputs {
		lines, err := physicalLines(input.Bytes)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", input.Path, err)
		}
		value, err := module.CallMethod(ctx, "_probe_spans", python.ValueOf(input.Bytes))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", input.Path, err)
		}
		text, err := python.As[python.StrValue](value)
		if err != nil {
			return nil, err
		}
		var positions [][2][4]int
		if err := json.Unmarshal([]byte(text.String()), &positions); err != nil {
			return nil, err
		}
		file := probeFile{Path: input.Path, Hash: sha256.Sum256(input.Bytes), Calls: make([]probeCall, 0, len(positions))}
		for _, pos := range positions {
			call, err := physicalSpan(input.Bytes, lines, pos[0])
			if err != nil {
				return nil, err
			}
			function, err := physicalSpan(input.Bytes, lines, pos[1])
			if err != nil {
				return nil, err
			}
			file.Calls = append(file.Calls, probeCall{call, function})
		}
		files = append(files, file)
	}
	return files, nil
}
