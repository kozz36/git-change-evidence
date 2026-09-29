package pythonparseprobe

import (
	"fmt"
	"os"
	"runtime"
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

// This workload is run only as an external subprocess with timeout and address
// space limits. The normal test run never executes pathological inputs.
func TestResourceWorkload(t *testing.T) {
	mode := os.Getenv("GCE_PROBE_WORKLOAD")
	if mode == "" {
		t.Skip("set GCE_PROBE_WORKLOAD and run test binary under external limits")
	}
	files, want, ok := resourceCase(mode)
	if !ok {
		t.Fatalf("unknown workload %q", mode)
	}
	matches, err := Observe(files, "os", "open")
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	t.Logf("mode=%s files=%d matches=%d error=%v heap_alloc=%d heap_sys=%d stack_sys=%d total_alloc=%d", mode, len(files), len(matches), err, mem.HeapAlloc, mem.HeapSys, mem.StackSys, mem.TotalAlloc)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != want {
		t.Fatalf("matches=%d want=%d", len(matches), want)
	}
}

func resourceCase(mode string) ([]File, int, bool) {
	switch mode {
	case "depth":
		return []File{{Path: "deep.py", Source: []byte(strings.Repeat("(", 200) + "os.open('x')" + strings.Repeat(")", 200) + "\n")}}, 1, true
	case "breadth":
		return []File{{Path: "broad.py", Source: []byte(strings.Repeat("x = [1,2,3]\n", 2000) + "os.open('x')\n")}}, 1, true
	case "token":
		return []File{{Path: "token.py", Source: []byte("x = '" + strings.Repeat("a", 1<<20) + "'\nos.open('x')\n")}}, 1, true
	case "files":
		files := make([]File, 64)
		for i := range files {
			files[i] = File{Path: fmt.Sprintf("f%03d.pyi", i), Source: []byte("os.open('x')\n")}
		}
		return files, len(files), true
	case "matches":
		return []File{{Path: "calls.py", Source: []byte(strings.Repeat("os.open('x')\n", 10000))}}, 10000, true
	default:
		return nil, 0, false
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
