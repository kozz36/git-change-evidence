package main

import (
	"bytes"
	"crypto/sha256"
	"testing"

	python "github.com/goccy/go-python"
)

func TestPhysicalSpansPinnedToRawBytesNeverNormalizedColumns(t *testing.T) {
	p, ctx := newPinnedPython(t)
	for _, ending := range []string{"\n", "\r\n", "\r"} {
		for _, bom := range []string{"", "\xef\xbb\xbf"} {
			source := []byte(bom + "é = '前'; os = fake; x = (os.open)(" + ending + "    \"中\", os.open('é'))" + ending)
			got, err := probeBatch(ctx, p, []probeInput{{"unicode.py", source}})
			if err != nil || len(got) != 1 || len(got[0].Calls) != 2 {
				t.Fatalf("ending %q BOM %t: %+v %v", ending, bom != "", got, err)
			}
			want := []string{"(os.open)(" + ending + "    \"中\", os.open('é'))", "os.open('é')"}
			for i, call := range got[0].Calls {
				if call.Call.Start != bytes.Index(source, []byte(want[i])) || string(call.Call.Bytes) != want[i] || string(call.Function.Bytes) != "os.open" {
					t.Fatalf("physical fragments: %+v", call)
				}
				for _, span := range []probeSpan{call.Call, call.Function} {
					if !bytes.Equal(source[span.Start:span.End], span.Bytes) || span.Hash != sha256.Sum256(span.Bytes) {
						t.Fatalf("raw slice/hash mismatch: %+v", span)
					}
				}
			}
			if got[0].Hash != sha256.Sum256(source) {
				t.Fatal("file hash must bind original bytes")
			}
		}
	}
}

func TestBatchPinnedToAtomicFailureNeverPartialMatches(t *testing.T) {
	p, ctx := newPinnedPython(t)
	for _, invalid := range []string{"def f(:", "return 42", "nonlocal x"} {
		got, err := probeBatch(ctx, p, []probeInput{{"good.py", []byte("os.open('x')")}, {"bad.py", []byte(invalid)}})
		if err == nil || got != nil {
			t.Fatalf("invalid %q returned partial result: %+v %v", invalid, got, err)
		}
	}
}

func TestBatchPredicateAndZeroMatch313Stub(t *testing.T) {
	p, ctx := newPinnedPython(t)
	source := "raise RuntimeError('must never execute')\nos = fake\nos.open('x').read()\nos.path.open('x')\nos.open[0]('x')\nos.open\n"
	got, err := probeBatch(ctx, p, []probeInput{{"direct.py", []byte(source)}, {"types.pyi", []byte("type Alias[T] = list[T]\ndef identity[T](x: T) -> T: ...\n")}})
	if err != nil || len(got) != 2 || len(got[0].Calls) != 1 || len(got[1].Calls) != 0 {
		t.Fatalf("batch: %+v %v", got, err)
	}
}

// Differential PASS from the independent run established mismatches, not parity.
func TestByteValidationPinnedToOriginalBytesNeverDecodedCookieGuess(t *testing.T) {
	p, ctx := newPinnedPython(t)
	builtins, err := p.Import(ctx, "builtins")
	if err != nil {
		t.Fatal(err)
	}
	fixtures := []struct {
		name, source string
		accept       bool
	}{
		{"bom_utf8_alias", "\xef\xbb\xbf# coding: utf8\nos.open('x')\n", false},
		{"second_cookie_after_code", "os.open('x')\n# coding: made-up\n", true},
		{"second_cookie_after_first_cookie", "# coding: utf-8\n# coding: latin-1\nos.open('x')\n", true},
		{"bom_utf8_canonical_control", "\xef\xbb\xbf# coding: utf-8\nos.open('x')\n", true},
		{"utf8_alias_without_bom_control", "# coding: utf8\nos.open('x')\n", true},
	}
	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			raw := []byte(fixture.source)
			_, compileErr := builtins.CallMethod(ctx, "compile", python.ValueOf(raw), python.ValueOf(fixture.name+".py"), python.ValueOf("exec"))
			if (compileErr == nil) != fixture.accept {
				t.Fatalf("pinned byte compile premise: %v", compileErr)
			}
			got, probeErr := probeBatch(ctx, p, []probeInput{{"good.py", []byte("os.open('earlier')")}, {fixture.name + ".py", raw}})
			if (probeErr == nil) != fixture.accept {
				t.Fatalf("probe parity: accept=%t, got %+v, %v", fixture.accept, got, probeErr)
			}
			if !fixture.accept {
				requireSyntaxError(t, probeErr)
				if got != nil {
					t.Fatal("invalid bytes retained partial batch")
				}
			} else if len(got) != 2 || len(got[1].Calls) != 1 || string(got[1].Calls[0].Call.Bytes) != "os.open('x')" {
				t.Fatalf("accepted byte fragments: %+v", got)
			}
		})
	}
}

func TestEncodingSupportExplicit(t *testing.T) {
	p, ctx := newPinnedPython(t)
	for _, source := range [][]byte{[]byte("# coding: latin-1\nos.open('x')\n"), []byte("# coding: made-up\nos.open('x')\n"), {'#', '\n', 0xff}, []byte("\xef\xbb\xbf# coding: latin-1\n")} {
		got, err := probeBatch(ctx, p, []probeInput{{"encoding.py", source}})
		if err == nil || got != nil {
			t.Fatalf("unsupported encoding accepted: %+v %v", got, err)
		}
	}
	for _, cookie := range []string{"# coding: utf-8\n", "#!/bin/python\r\n# coding=utf_8\r\n"} {
		got, err := probeBatch(ctx, p, []probeInput{{"encoding.py", []byte(cookie + "os.open('é')")}})
		if err != nil || len(got[0].Calls) != 1 {
			t.Fatalf("UTF8 cookie: %+v %v", got, err)
		}
	}
}
