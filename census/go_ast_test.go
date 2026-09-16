package census

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go/token"
	"math"
	"reflect"
	"testing"

	ce "github.com/kozz36/git-change-evidence"
)

func censusFixture(t *testing.T, entries []ce.InventoryEntryInput) (ce.InventoryDocumentV1, []File) {
	t.Helper()
	policy, err := ce.NewAccountingPolicyV1(ce.PolicyInput{
		Categories: []ce.CategoryInput{{Name: "go"}},
		Default:    "go",
	})
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := ce.NewInventoryV1(policy, entries)
	if err != nil {
		t.Fatal(err)
	}
	files := make([]File, len(entries))
	for index, entry := range entries {
		files[index] = File{Path: append([]byte(nil), entry.Path...), Content: append([]byte(nil), entry.Content...)}
	}
	return inventory, files
}

func standardQuery() QueryV1 {
	return QueryV1{Version: SelectorCallQueryV1, Receiver: "os", Selector: "Open"}
}

func requireContract(t *testing.T, err error, field, code string) {
	t.Helper()
	var contract *ce.ContractError
	if !errors.As(err, &contract) || contract.Field != field || contract.Code != code {
		t.Fatalf("error = %#v, want %s/%s", err, field, code)
	}
}

func requireZero(t *testing.T, result ResultV1) {
	t.Helper()
	if result.InventoryDigest() != "" || result.Query() != (QueryV1{}) || result.ExtractorVersion() != "" || result.Matches() != nil {
		t.Fatalf("partial result = %#v", result)
	}
}

func matchText(content []byte, match MatchV1) string {
	return string(content[match.Call.Start.Offset:match.Call.End.Offset])
}

func literalPosition(content []byte, offset int) Position {
	line, column := 1, 1
	for _, byte := range content[:offset] {
		if byte == '\n' {
			line, column = line+1, 1
		} else {
			column++
		}
	}
	return Position{Offset: offset, Line: line, Column: column}
}

func literalMatch(path, content []byte, callStart, callEnd, selectorStart, selectorEnd int) MatchV1 {
	fileHash := sha256.Sum256(content)
	fragmentHash := sha256.Sum256(content[callStart:callEnd])
	return MatchV1{
		Path:           append([]byte(nil), path...),
		FileSHA256:     ce.Digest(hex.EncodeToString(fileHash[:])),
		FragmentSHA256: ce.Digest(hex.EncodeToString(fragmentHash[:])),
		Call:           Span{Start: literalPosition(content, callStart), End: literalPosition(content, callEnd)},
		Selector:       Span{Start: literalPosition(content, selectorStart), End: literalPosition(content, selectorEnd)},
	}
}

func TestGAC001AndGAC002Validation(t *testing.T) {
	inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("one.go"), Content: []byte("package p\n")}})
	t.Run("GAC-001-S02 invalid query", func(t *testing.T) {
		for _, test := range []struct {
			name, field, code string
			query             QueryV1
		}{
			{"unsupported version", "query.version", "unsupported_version", QueryV1{Version: "other", Receiver: "os", Selector: "Open"}},
			{"empty receiver", "query.receiver", "invalid_identifier", QueryV1{Version: SelectorCallQueryV1, Selector: "Open"}},
			{"empty selector", "query.selector", "invalid_identifier", QueryV1{Version: SelectorCallQueryV1, Receiver: "os"}},
			{"keyword selector", "query.selector", "invalid_identifier", QueryV1{Version: SelectorCallQueryV1, Receiver: "os", Selector: "func"}},
		} {
			t.Run(test.name, func(t *testing.T) {
				result, err := GoASTV1(inventory, files, test.query, DefaultLimits())
				requireContract(t, err, test.field, test.code)
				requireZero(t, result)
			})
		}
	})
	t.Run("GAC-002-S01 duplicate supplied path", func(t *testing.T) {
		result, err := GoASTV1(inventory, append(files, File{Path: []byte("one.go"), Content: []byte("package broken")}), standardQuery(), DefaultLimits())
		requireContract(t, err, "files.path", "duplicate_path")
		requireZero(t, result)
	})
	t.Run("GAC-002-S02 missing supplied path", func(t *testing.T) {
		want, _ := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("a.go"), Content: []byte("package p\n")}, {Path: []byte("b.go"), Content: []byte("package broken")}})
		result, err := GoASTV1(want, []File{{Path: []byte("a.go"), Content: []byte("package p\n")}}, standardQuery(), DefaultLimits())
		requireContract(t, err, "files.path", "missing_path")
		requireZero(t, result)
	})
	t.Run("GAC-002-S03 extra supplied path", func(t *testing.T) {
		result, err := GoASTV1(inventory, append(files, File{Path: []byte("extra.go"), Content: []byte("package broken")}), standardQuery(), DefaultLimits())
		requireContract(t, err, "files.path", "extra_path")
		requireZero(t, result)
	})
	t.Run("GAC-002-S04 length and digest drift", func(t *testing.T) {
		for _, test := range []struct {
			name, content, field, code string
		}{
			{"same length digest drift", "package q\n", "files.content_sha256", "digest_mismatch"},
			{"unequal length drift", "package q", "files.byte_length", "length_mismatch"},
		} {
			t.Run(test.name, func(t *testing.T) {
				result, err := GoASTV1(inventory, []File{{Path: []byte("one.go"), Content: []byte(test.content)}}, standardQuery(), DefaultLimits())
				requireContract(t, err, test.field, test.code)
				requireZero(t, result)
			})
		}
	})
	t.Run("GAC-002-S05 zero inventory", func(t *testing.T) {
		result, err := GoASTV1(ce.InventoryDocumentV1{}, nil, standardQuery(), DefaultLimits())
		requireContract(t, err, "inventory", "invalid_document")
		requireZero(t, result)
	})
	t.Run("GAC-002-S06 valid empty scope", func(t *testing.T) {
		empty, noFiles := censusFixture(t, []ce.InventoryEntryInput{})
		result, err := GoASTV1(empty, noFiles, standardQuery(), DefaultLimits())
		if err != nil || result.InventoryDigest() != empty.Digest() || result.Matches() == nil || len(result.Matches()) != 0 {
			t.Fatalf("empty result = (%#v, %v)", result, err)
		}
	})
}

func TestGAC003Candidates(t *testing.T) {
	content := []byte("package p\nfunc f() {\n\tos.Open(\"direct\")\n\t(os.Open)(\"paren\")\n\tos.Open[int](\"index\")\n\t((os.Open[int, string]))(\"multi\")\n\twrap(os.Open(\"nested\"))\n\tos := struct{ Open func() }{}\n\tos.Open()\n\talias.Open()\n\t_ = \"os.Open()\"\n\t(os).Open()\n\tplain()\n}\n")
	inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("candidates.go"), Content: content}})
	t.Run("GAC-003-S01 exact and single-index selector calls", func(t *testing.T) {
		result, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"os.Open(\"direct\")", "(os.Open)(\"paren\")", "os.Open[int](\"index\")", "((os.Open[int, string]))(\"multi\")", "os.Open(\"nested\")", "os.Open()"}
		if matches := result.Matches(); len(matches) != len(want) {
			t.Fatalf("matches = %#v", matches)
		} else {
			for index, match := range matches {
				if got := matchText(content, match); got != want[index] {
					t.Fatalf("match %d = %q, want %q", index, got, want[index])
				}
			}
		}
	})
	t.Run("GAC-003-S02 nested parentheses and index-list normalize", func(t *testing.T) {
		result, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		if err != nil || matchText(content, result.Matches()[3]) != "((os.Open[int, string]))(\"multi\")" {
			t.Fatalf("nested wrapper result = (%#v, %v)", result.Matches(), err)
		}
	})
	t.Run("GAC-003-S03 shadowed and ambiguous candidates remain syntax-only", func(t *testing.T) {
		result, err := GoASTV1(inventory, files, QueryV1{Version: SelectorCallQueryV1, Receiver: "obj", Selector: "Funcs"}, DefaultLimits())
		if err != nil || len(result.Matches()) != 0 {
			t.Fatalf("nonmatching query = (%#v, %v)", result.Matches(), err)
		}
		commentDecoy := []byte("package p\n// os.Open(\"comment\")\nfunc f() {}\n")
		commentInventory, commentFiles := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("comment.go"), Content: commentDecoy}})
		result, err = GoASTV1(commentInventory, commentFiles, standardQuery(), DefaultLimits())
		if err != nil || len(result.Matches()) != 0 {
			t.Fatalf("comment decoy = (%#v, %v)", result.Matches(), err)
		}
		ambiguous := []byte("package p\nfunc f(){ obj.Funcs[i]() }\n")
		otherInventory, otherFiles := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("index.go"), Content: ambiguous}})
		result, err = GoASTV1(otherInventory, otherFiles, QueryV1{Version: SelectorCallQueryV1, Receiver: "obj", Selector: "Funcs"}, DefaultLimits())
		if err != nil || len(result.Matches()) != 1 || matchText(ambiguous, result.Matches()[0]) != "obj.Funcs[i]()" {
			t.Fatalf("ambiguous index = (%#v, %v)", result.Matches(), err)
		}
	})
	t.Run("GAC-003-S04 repeated legacy sites have no classification", func(t *testing.T) {
		legacy := []byte("package p\nfunc f(){ os.Open(\"first\"); _ = os.Open(\"second\"); defer os.Open(\"third\") }\n")
		legacyInventory, legacyFiles := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("legacy.go"), Content: legacy}})
		result, err := GoASTV1(legacyInventory, legacyFiles, standardQuery(), DefaultLimits())
		if err != nil || len(result.Matches()) != 3 || matchText(legacy, result.Matches()[1]) != "os.Open(\"second\")" {
			t.Fatalf("legacy results = (%#v, %v)", result.Matches(), err)
		}
	})
}

func TestGAC004PhysicalEvidence(t *testing.T) {
	t.Run("GAC-004-S01 UTF8 CRLF multiline and line directives", func(t *testing.T) {
		content := []byte("package p\r\n//line fake.go:900\r\nfunc f(){ _ = \"é\"; os.Open(\r\n\"x\") }\r\n")
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("physical.go"), Content: content}})
		result, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		if err != nil {
			t.Fatal(err)
		}
		match := result.Matches()[0]
		if match.Call != (Span{Start: Position{Offset: 51, Line: 3, Column: 21}, End: Position{Offset: 65, Line: 4, Column: 5}}) || match.Selector != (Span{Start: Position{Offset: 51, Line: 3, Column: 21}, End: Position{Offset: 58, Line: 3, Column: 28}}) || matchText(content, match) != "os.Open(\r\n\"x\")" {
			t.Fatalf("physical evidence = %#v, %q", match, matchText(content, match))
		}
		unicodeSource := []byte("package p\nfunc f(){ π.Écho() }\n")
		unicodeInventory, unicodeFiles := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("unicode.go"), Content: unicodeSource}})
		unicodeResult, err := GoASTV1(unicodeInventory, unicodeFiles, QueryV1{Version: SelectorCallQueryV1, Receiver: "π", Selector: "Écho"}, DefaultLimits())
		if err != nil || unicodeResult.Matches()[0].Call != (Span{Start: Position{Offset: 20, Line: 2, Column: 11}, End: Position{Offset: 30, Line: 2, Column: 21}}) {
			t.Fatalf("unicode bytes = (%#v, %v)", unicodeResult.Matches(), err)
		}
	})
	t.Run("GAC-004-S02 equal counts retain distinct locations and fragments", func(t *testing.T) {
		first := []byte("package p\nfunc f(){ os.Open(\"a\") }\n")
		second := []byte("package p\nfunc f(){  os.Open(\"b\") }\n")
		firstInventory, firstFiles := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("same.go"), Content: first}})
		secondInventory, secondFiles := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("same.go"), Content: second}})
		one, oneErr := GoASTV1(firstInventory, firstFiles, standardQuery(), DefaultLimits())
		two, twoErr := GoASTV1(secondInventory, secondFiles, standardQuery(), DefaultLimits())
		if oneErr != nil || twoErr != nil || one.Matches()[0].Call == two.Matches()[0].Call || one.Matches()[0].FragmentSHA256 == two.Matches()[0].FragmentSHA256 || one.Matches()[0].FileSHA256 == two.Matches()[0].FileSHA256 {
			t.Fatalf("equal-count evidence collapsed = (%#v, %#v, %v, %v)", one.Matches(), two.Matches(), oneErr, twoErr)
		}
	})
}

func TestGAC005OwnershipAndOrder(t *testing.T) {
	source := []byte("package p\nfunc f(){ os.Open(\"x\") }\n")
	entries := []ce.InventoryEntryInput{{Path: []byte{0xff}, Content: source}, {Path: []byte("a\\b"), Content: source}, {Path: []byte("a\n"), Content: source}, {Path: []byte("Case"), Content: source}}
	inventory, files := censusFixture(t, entries)
	files[0], files[3] = files[3], files[0]
	result, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	t.Run("GAC-005-S01 input and output mutation cannot alter evidence", func(t *testing.T) {
		files[0].Path[0], files[0].Content[0] = 'x', 'x'
		returned := result.Matches()
		returned[0].Path[0] = 'x'
		returned[0] = MatchV1{}
		if got := result.Matches(); len(got) != 4 || bytes.Equal(got[0].Path, []byte("xase")) || got[0].Call.Start.Offset == 0 {
			t.Fatalf("ownership result = %#v", got)
		}
	})
	t.Run("GAC-005-S02 raw path order is deterministic", func(t *testing.T) {
		want := [][]byte{[]byte("Case"), []byte("a\n"), []byte("a\\b"), {0xff}}
		got := result.Matches()
		if len(got) != len(want) {
			t.Fatalf("match count = %d, want %d", len(got), len(want))
		}
		for index, match := range got {
			if !bytes.Equal(match.Path, want[index]) {
				t.Fatalf("path %d = %q, want %q", index, match.Path, want[index])
			}
		}
	})
}

func TestGAC006LimitsAndGAC007Atomicity(t *testing.T) {
	content := []byte("package p\n")
	t.Run("GAC-006-S01 exact file count then excess", func(t *testing.T) {
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("a.go"), Content: content}, {Path: []byte("b.go"), Content: content}})
		limits := DefaultLimits()
		limits.MaxFiles = 2
		if result, err := GoASTV1(inventory, files, standardQuery(), limits); err != nil || result.Matches() == nil {
			t.Fatalf("file equality = (%#v, %v)", result, err)
		}
		limits.MaxFiles = 1
		result, err := GoASTV1(inventory, files, standardQuery(), limits)
		requireContract(t, err, "limits.max_files", "limit_exceeded")
		requireZero(t, result)
	})
	t.Run("GAC-006-S02 and S03 exact byte budgets then excess", func(t *testing.T) {
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("a.go"), Content: content}, {Path: []byte("b.go"), Content: content}})
		limits := DefaultLimits()
		limits.MaxFileBytes, limits.MaxTotalBytes = uint64(len(content)), uint64(len(content)*2)
		if _, err := GoASTV1(inventory, files, standardQuery(), limits); err != nil {
			t.Fatal(err)
		}
		limits.MaxFileBytes--
		result, err := GoASTV1(inventory, files, standardQuery(), limits)
		requireContract(t, err, "limits.max_file_bytes", "limit_exceeded")
		requireZero(t, result)
		limits.MaxFileBytes, limits.MaxTotalBytes = uint64(len(content)), uint64(len(content)*2-1)
		result, err = GoASTV1(inventory, files, standardQuery(), limits)
		requireContract(t, err, "limits.max_total_bytes", "limit_exceeded")
		requireZero(t, result)
	})
	t.Run("GAC-006-S04 exact match count then atomic excess", func(t *testing.T) {
		calls := []byte("package p\nfunc f(){ os.Open(\"a\"); os.Open(\"b\") }\n")
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("calls.go"), Content: calls}})
		limits := DefaultLimits()
		limits.MaxMatches = 2
		if result, err := GoASTV1(inventory, files, standardQuery(), limits); err != nil || len(result.Matches()) != 2 {
			t.Fatalf("match equality = (%#v, %v)", result, err)
		}
		limits.MaxMatches = 1
		result, err := GoASTV1(inventory, files, standardQuery(), limits)
		requireContract(t, err, "limits.max_matches", "limit_exceeded")
		requireZero(t, result)
	})
	for _, test := range []struct {
		name, field string
		change      func(*Limits)
	}{
		{"GAC-006-S05 zero files", "limits.max_files", func(limits *Limits) { limits.MaxFiles = 0 }},
		{"GAC-006-S06 zero file bytes", "limits.max_file_bytes", func(limits *Limits) { limits.MaxFileBytes = 0 }},
		{"GAC-006-S07 zero total bytes", "limits.max_total_bytes", func(limits *Limits) { limits.MaxTotalBytes = 0 }},
		{"GAC-006-S08 zero matches", "limits.max_matches", func(limits *Limits) { limits.MaxMatches = 0 }},
	} {
		t.Run(test.name, func(t *testing.T) {
			inventory, files := censusFixture(t, nil)
			limits := DefaultLimits()
			test.change(&limits)
			result, err := GoASTV1(inventory, files, standardQuery(), limits)
			requireContract(t, err, test.field, "invalid_limit")
			requireZero(t, result)
		})
	}
	t.Run("GAC-007-S01 parser error after candidate is atomic", func(t *testing.T) {
		broken := []byte("package p\nfunc f(){ os.Open(\"x\")\n")
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("broken.go"), Content: broken}})
		result, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		requireContract(t, err, "source", "parse_error")
		requireZero(t, result)
	})
}

func TestGAC007SpanSafetyAndArithmetic(t *testing.T) {
	t.Run("GAC-007-S02 invalid positions and spans are rejected", func(t *testing.T) {
		fset := token.NewFileSet()
		file := fset.AddFile("content.go", -1, 3)
		file.SetLinesForContent([]byte("abc"))
		foreign := fset.AddFile("foreign.go", -1, 1)
		for _, test := range []struct {
			name       string
			start, end token.Pos
		}{
			{"no position", token.NoPos, token.Pos(file.Base() + 1)},
			{"foreign file", token.Pos(foreign.Base()), token.Pos(foreign.Base() + 1)},
			{"past EOF", token.Pos(file.Base()), token.Pos(file.Base() + 4)},
			{"reversed", token.Pos(file.Base() + 2), token.Pos(file.Base() + 1)},
		} {
			t.Run(test.name, func(t *testing.T) {
				if _, err := physicalSpan(fset, file, test.start, test.end, []byte("abc")); err == nil {
					t.Fatal("invalid span succeeded")
				}
			})
		}
	})
	t.Run("GAC-006-S03 near uint64 aggregate arithmetic", func(t *testing.T) {
		if !withinTotal(math.MaxUint64, math.MaxUint64-1, 1) || withinTotal(math.MaxUint64, math.MaxUint64-1, 2) || withinTotal(3, 4, 0) {
			t.Fatal("aggregate boundary arithmetic is not overflow safe")
		}
	})
}

func TestTriangulatedAtomicityAndOrder(t *testing.T) {
	t.Run("GAC-007-S01 earlier valid file cannot escape later parser error", func(t *testing.T) {
		valid := []byte("package p\nfunc f(){ os.Open(\"valid\") }\n")
		broken := []byte("package p\nfunc f(){ os.Open(\"broken\")\n")
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("a.go"), Content: valid}, {Path: []byte("b.go"), Content: broken}})
		result, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		requireContract(t, err, "source", "parse_error")
		requireZero(t, result)
	})
	t.Run("GAC-005-S02 input permutation retains complete independently expected records", func(t *testing.T) {
		a := []byte("package p\nfunc f(){ os.Open(\"a\") }\n")
		z := []byte("package p\nfunc f(){ os.Open(\"zebra\"); os.Open(\"yak\") }\n")
		inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("z.go"), Content: z}, {Path: []byte("a.go"), Content: a}})
		want := []MatchV1{
			literalMatch([]byte("a.go"), a, 20, 32, 20, 27),
			literalMatch([]byte("z.go"), z, 20, 36, 20, 27),
			literalMatch([]byte("z.go"), z, 38, 52, 38, 45),
		}
		forward, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		if err != nil {
			t.Fatal(err)
		}
		if got := forward.Matches(); len(got) != len(want) {
			t.Fatalf("forward match count = %d, want %d", len(got), len(want))
		} else if !reflect.DeepEqual(got, want) {
			t.Fatalf("forward records = %#v, want %#v", got, want)
		}
		files[0], files[1] = files[1], files[0]
		backward, err := GoASTV1(inventory, files, standardQuery(), DefaultLimits())
		if err != nil {
			t.Fatal(err)
		}
		if got := backward.Matches(); len(got) != len(want) {
			t.Fatalf("backward match count = %d, want %d", len(got), len(want))
		} else if !reflect.DeepEqual(got, want) {
			t.Fatalf("backward records = %#v, want %#v", got, want)
		}
	})
}

func TestDefaultLimitsAndIndependentHashes(t *testing.T) {
	limits := DefaultLimits()
	if limits != (Limits{MaxFiles: 256, MaxFileBytes: 1 << 20, MaxTotalBytes: 8 << 20, MaxMatches: 10000}) {
		t.Fatalf("defaults = %#v", limits)
	}
	content := []byte("package p\nfunc f(){ os.Open(\"x\") }\n")
	inventory, files := censusFixture(t, []ce.InventoryEntryInput{{Path: []byte("hash.go"), Content: content}})
	result, err := GoASTV1(inventory, files, standardQuery(), limits)
	if err != nil {
		t.Fatal(err)
	}
	fileHash := sha256.Sum256(content)
	fragmentHash := sha256.Sum256([]byte("os.Open(\"x\")"))
	match := result.Matches()[0]
	if match.FileSHA256 != ce.Digest(hex.EncodeToString(fileHash[:])) || match.FragmentSHA256 != ce.Digest(hex.EncodeToString(fragmentHash[:])) {
		t.Fatalf("independent hashes = %#v", match)
	}
}
