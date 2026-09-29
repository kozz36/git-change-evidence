// Command smoke verifies a pinned, generated CPython 3.13 bundle outside GCE.
// It is not a census extractor or a production dependency.
package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	python "github.com/goccy/go-python"
)

// This trusted helper inspects a parsed AST. The source being scanned is never
// interpolated into Eval or executed; it crosses the bridge only as a value.
const exactCallsHelper = `import ast
def _count_exact_calls(tree):
    return sum(
        isinstance(node, ast.Call)
        and isinstance(node.func, ast.Attribute)
        and isinstance(node.func.value, ast.Name)
        and node.func.value.id == "os"
        and node.func.attr == "open"
        for node in ast.walk(tree)
    )
`

func runtimeVersion(ctx context.Context, p *python.Python) (string, error) {
	result, err := p.Eval(ctx, "import sys\nsys.version")
	if err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", result.Error
	}
	version, err := python.As[python.StrValue](result.Value)
	if err != nil {
		return "", err
	}
	return version.String(), nil
}

// inspect returns no candidate if either parsing or compile-stage validation
// fails. compile creates a code object but never runs the scanned source.
func inspect(ctx context.Context, p *python.Python, source string) (int64, error) {
	ast, err := p.Import(ctx, "ast")
	if err != nil {
		return 0, err
	}
	tree, err := ast.CallMethod(ctx, "parse", python.ValueOf(source), python.ValueOf("<probe>"), python.ValueOf("exec"))
	if err != nil {
		return 0, err
	}
	builtins, err := p.Import(ctx, "builtins")
	if err != nil {
		return 0, err
	}
	if _, err := builtins.CallMethod(ctx, "compile", tree, python.ValueOf("<probe>"), python.ValueOf("exec")); err != nil {
		return 0, err
	}
	result, err := p.Eval(ctx, exactCallsHelper)
	if err != nil {
		return 0, err
	}
	if result.Error != nil {
		return 0, result.Error
	}
	mainModule, err := p.Import(ctx, "__main__")
	if err != nil {
		return 0, err
	}
	value, err := mainModule.CallMethod(ctx, "_count_exact_calls", tree)
	if err != nil {
		return 0, err
	}
	count, err := python.As[python.IntValue](value)
	if err != nil {
		return 0, err
	}
	n, fits := count.Int64()
	if !fits {
		return 0, fmt.Errorf("call count exceeds int64")
	}
	return n, nil
}

func main() {
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	started := time.Now()
	p, err := python.New(python.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "initialize:", err)
		os.Exit(1)
	}
	defer p.Close()
	startup := time.Since(started)
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	fmt.Printf("startup=%s heap_alloc_after_gc_bytes=%d total_alloc_delta_bytes=%d\n", startup, after.HeapAlloc, after.TotalAlloc-before.TotalAlloc)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	version, err := runtimeVersion(ctx, p)
	if err != nil || !strings.HasPrefix(version, "3.13.13 ") {
		fmt.Fprintf(os.Stderr, "expected pinned CPython 3.13.13, got %q: %v\n", version, err)
		os.Exit(2)
	}
	fmt.Printf("version=%s\n", version)
	count, err := inspect(ctx, p, "os.open('never-executed')\nos.path.open('not-direct')\n")
	if err != nil || count != 1 {
		fmt.Fprintf(os.Stderr, "exact AST count = %d: %v\n", count, err)
		os.Exit(3)
	}
	fmt.Printf("direct_calls=%d\n", count)
}
