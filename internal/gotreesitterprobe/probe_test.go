package gotreesitterprobe

import (
	"bytes"
	"testing"
)

// Oracle: Python 3.13 compile(source, '<probe>', 'exec') on identical bytes.
// AST-only parsing does not reject all invalid execution-unit syntax.
func TestPython313Validity(t *testing.T) {
	cases := []struct {
		name, source              string
		oracleValid, grammarValid bool
	}{
		{"simple call", "os.open('x')\n", true, true},
		{"issue 178 newline before colon", "def foo(x)\n:\n    return x + 2\n", false, true},
		{"issue 178 continued expression", "a = (\n  1 +\n  2\n)\n", true, true},
		{"missing closing parenthesis", "os.open(\n", false, false},
		{"illegal assignment target", "1 = os.open('x')\n", false, false},
		{"top level return", "return os.open('x')\n", false, true},
		{"type parameters 3.13", "def f[T](x: T) -> T:\n    return x\n", true, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := parseValid([]byte(tt.source))
			if err != nil {
				t.Fatal(err)
			}
			if valid != tt.grammarValid {
				t.Errorf("grammar valid=%v want observed %v (CPython 3.13 valid=%v)", valid, tt.grammarValid, tt.oracleValid)
			}
			if tt.oracleValid != tt.grammarValid {
				t.Logf("FAIL feasibility: grammar valid=%v while CPython 3.13 valid=%v", valid, tt.oracleValid)
			}
		})
	}
}

func TestExactCallShapeAndPhysicalSpans(t *testing.T) {
	source := []byte("π = 'é'\r\nos.open(os.open('猫'))\r\nos.open\nos.path.open('x')\nos.open[0]('x')\nos = fake; os.open('y')\n")
	matches, err := directCalls(source)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"os.open(os.open('猫'))", "os.open('猫')", "os.open('y')"}
	wantSpans := []span{{11, 34}, {19, 33}, {89, 101}}
	if len(matches) != len(want) {
		t.Fatalf("matches=%v want %q", matches, want)
	}
	for i, m := range matches {
		if m != wantSpans[i] {
			t.Errorf("span[%d]=%+v want %+v", i, m, wantSpans[i])
		}
		if m.Start < 0 || m.End > len(source) || m.Start >= m.End {
			t.Fatalf("invalid span: %+v", m)
		}
		if !bytes.Equal(source[m.Start:m.End], []byte(want[i])) {
			t.Errorf("slice[%d]=%q want %q", i, source[m.Start:m.End], want[i])
		}
	}
}

func TestStubAndAtomicRejection(t *testing.T) {
	stub := []byte("from typing import overload\n@overload\ndef f[T](x: T) -> T: ...\nasync def run() -> None:\n    os.open('x')\n")
	if _, err := directCalls(stub); err != nil {
		t.Fatalf("stub: %v", err)
	}
	if got, err := directCalls([]byte("os.open('x')\nos.open(\n")); err == nil || got != nil {
		t.Fatalf("partial success: %v, %v", got, err)
	}
	// The grammar's clean tree admits this CPython-invalid file and produces a
	// candidate: syntax-only atomic validity is not established by HasError.
	invalid := []byte("def foo(x)\n:\n    os.open('x')\n")
	if got, err := directCalls(invalid); err != nil || len(got) != 1 {
		t.Fatalf("expected observed invalid-file false success: %v, %v", got, err)
	}
}
