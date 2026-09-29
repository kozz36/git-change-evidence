package pythonparseprobe

import (
	"strings"
	"testing"
)

func TestExactCallShape(t *testing.T) {
	cases := []struct {
		name, src string
		want      int
	}{
		{"direct", `os.open("x")`, 1},
		{"attribute_only", `os.open`, 0},
		{"chained_receiver", `os.path.open("x")`, 0},
		{"indexed_function", `os.open[0]("x")`, 0},
		{"inner_only", `os.open("x").read()`, 1},
		{"shadowed", "os = fake\nos.open('x')", 1},
		{"nested_calls", `os.open(os.open("x"))`, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			matches, err := Observe([]File{{Path: "a.py", Source: []byte(tc.src)}}, "os", "open")
			if err != nil {
				t.Fatal(err)
			}
			if len(matches) != tc.want {
				t.Fatalf("matches=%d want %d", len(matches), tc.want)
			}
		})
	}
}

func TestScriptAndStubGrammar(t *testing.T) {
	cases := []struct {
		path, src string
		want      int
	}{
		{"script.py", "async def f():\n    return await os.open('x')\n", 1},
		{"stub.pyi", "from typing import overload\n@overload\ndef f(x: int) -> str: ...\n@overload\ndef f(x: str) -> str: ...\n", 0},
		{"modern.py", "type Alias[T] = list[T]\nmatch value:\n    case [x]: os.open(x)\n", 1},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			matches, err := Observe([]File{{Path: tc.path, Source: []byte(tc.src)}}, "os", "open")
			if err != nil {
				t.Fatalf("grammar parse: %v", err)
			}
			if len(matches) != tc.want {
				t.Fatalf("matches=%d want %d", len(matches), tc.want)
			}
		})
	}
}

func TestBytePositionIsOnlyStartEvidence(t *testing.T) {
	src := []byte("π = 'é'\r\nos.open(os.open('猫'))\r\n")
	matches, err := Observe([]File{{Path: "unicode.py", Source: src}}, "os", "open")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 {
		t.Fatalf("matches=%d", len(matches))
	}
	first := strings.Index(string(src), "os.open")
	second := strings.LastIndex(string(src), "os.open")
	for i, m := range matches {
		want := []int{first, second}[i]
		t.Logf("candidate AST call position byte=%d, physical call start=%d", m.Start, want)
		if m.Start == want {
			t.Error("expected to expose position mismatch for this pinned revision")
		}
		if m.EndKnown {
			t.Error("parser start-only AST must not invent an end span")
		}
	}
}

func TestAtomicMalformedInput(t *testing.T) {
	matches, err := Observe([]File{{Path: "good.py", Source: []byte("os.open('x')")}, {Path: "bad.pyi", Source: []byte("def broken(:\n")}}, "os", "open")
	if err == nil || matches != nil {
		t.Fatalf("must fail atomically: matches=%v error=%v", matches, err)
	}
}

func TestUnsupportedSyntaxIsCharacterized(t *testing.T) {
	// Future-only spelling is not evidence of a 3.14 grammar boundary unless rejected.
	matches, err := Observe([]File{{Path: "future.py", Source: []byte("def f():\n    defer os.open('x')\n")}}, "os", "open")
	t.Logf("future-looking defer syntax: matches=%d error=%v", len(matches), err)
	if err == nil && len(matches) > 0 {
		t.Log("FAIL AXIS: unsupported syntax produced a candidate")
	}
}

func TestInputPaths(t *testing.T) {
	for _, path := range []string{"a.txt", "../a.py", "a.py\x00"} {
		_, err := Observe([]File{{Path: path, Source: []byte("os.open('x')")}}, "os", "open")
		if err == nil {
			t.Errorf("accepted %q", path)
		}
	}
	_, err := Observe([]File{{Path: "a.py", Source: []byte(strings.Repeat(" ", 4) + "os.open('x')\n")}}, "os", "open")
	if err == nil {
		t.Log("unexpected indent accepted")
	}
}
