package main

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	python "github.com/goccy/go-python"
)

// These tests must run with the locally generated 3.13 bridge, bundle, and stdlib.
// A release artifact from upstream go-python v0.4.0 is 3.14 and must fail here.
func newPinnedPython(t *testing.T) (*python.Python, context.Context) {
	t.Helper()
	p, err := python.New(python.Config{})
	if err != nil {
		t.Fatalf("boot pinned interpreter: %v", err)
	}
	t.Cleanup(func() {
		if err := p.Close(); err != nil {
			t.Errorf("close: %v", err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	version, err := runtimeVersion(ctx, p)
	if err != nil || !strings.HasPrefix(version, "3.13.13 ") {
		t.Fatalf("pinned 3.13.13 runtime required; got %q, %v", version, err)
	}
	return p, ctx
}

func requireSyntaxError(t *testing.T, err error) {
	t.Helper()
	var pyErr *python.PythonError
	if !errors.As(err, &pyErr) || pyErr.Type != "SyntaxError" {
		t.Fatalf("want Python SyntaxError, got %T %v", err, err)
	}
}

func TestExactCallShapeNeverNestedAttributeOrIndexed(t *testing.T) {
	p, ctx := newPinnedPython(t)
	// Compilation does not run any of these calls; nested attribute and index
	// expressions must not become direct Name.Attribute candidates.
	source := "os.open('x')\nos.path.open('y')\nos.open[0]('z')\nos.open\n"
	got, err := inspect(ctx, p, source)
	if err != nil || got != 1 {
		t.Fatalf("exact direct call count = %d, %v; want 1", got, err)
	}
}

func TestCompileRejectsASTOnlyTopLevelReturn(t *testing.T) {
	p, ctx := newPinnedPython(t)
	ast, err := p.Import(ctx, "ast")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ast.CallMethod(ctx, "parse", python.ValueOf("return 42")); err != nil {
		t.Fatalf("fixture must parse into an AST: %v", err)
	}
	_, err = inspect(ctx, p, "return 42")
	requireSyntaxError(t, err)
}

func TestMalformedSyntaxAndPython313Stub(t *testing.T) {
	p, ctx := newPinnedPython(t)
	_, err := inspect(ctx, p, "def f(:\n    pass\n")
	requireSyntaxError(t, err)
	got, err := inspect(ctx, p, "type Alias[T] = list[T]\ndef identity[T](x: T) -> T: ...\n")
	if err != nil || got != 0 {
		t.Fatalf("3.13 stub count = %d, %v; want 0", got, err)
	}
}

func TestScannedSourceIsNotExecuted(t *testing.T) {
	p, ctx := newPinnedPython(t)
	got, err := inspect(ctx, p, "raise RuntimeError('scanned source executed')\nos.open('x')\n")
	if err != nil || got != 1 {
		t.Fatalf("compile-only source count = %d, %v; want 1, no runtime error", got, err)
	}
}
