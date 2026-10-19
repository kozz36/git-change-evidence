package input_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

func put(t *testing.T, root, name string, raw []byte) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDocumentTransport(t *testing.T) {
	root := t.TempDir()
	raw := []byte(" \x00\xff\r\n")
	policy := put(t, root, "policy", raw)
	unused := readerFunc(func([]byte) (int, error) { t.Fatal("unused stdin read"); return 0, io.EOF })
	docs, err := input.ReadDocuments(context.Background(), contract.Request{Policy: policy}, unused)
	if err != nil || !bytes.Equal(docs.Policy, raw) || docs.Inventory != nil || docs.Selection != nil {
		t.Fatalf("%+v %v", docs, err)
	}
	after, err := os.ReadFile(policy)
	if err != nil || !bytes.Equal(after, raw) {
		t.Fatal("input changed", err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatal("unexpected filesystem effect", err)
	}
	_, err = input.ReadDocuments(context.Background(), contract.Request{Policy: "-", Paths: "-"}, unused)
	failure(t, err, 2, "invalid_arguments")
	_, err = input.ReadDocuments(context.Background(), contract.Request{Policy: filepath.Join(root, "absent")}, unused)
	failure(t, err, 3, "input_unavailable")
	_, err = input.ReadDocuments(context.Background(), contract.Request{Policy: "-"}, nil)
	failure(t, err, 3, "input_unavailable")
	_, err = input.ReadDocuments(context.Background(), contract.Request{Policy: "-"}, readerFunc(func([]byte) (int, error) { return 0, errors.New("private detail") }))
	failure(t, err, 3, "input_unavailable")
	docs, err = input.ReadDocuments(context.Background(), contract.Request{Policy: policy, Inventory: policy, Result: policy, Paths: "-"}, strings.NewReader(" \n"))
	if err != nil || !bytes.Equal(docs.Policy, raw) || !bytes.Equal(docs.Inventory, raw) || !bytes.Equal(docs.Result, raw) || string(docs.Selection) != " \n" {
		t.Fatalf("%+v %v", docs, err)
	}
}

func TestDocumentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	unused := readerFunc(func([]byte) (int, error) { t.Fatal("read after cancellation"); return 0, io.EOF })
	docs, err := input.ReadDocuments(ctx, contract.Request{Policy: "-"}, unused)
	failure(t, err, 6, "interrupted")
	if !reflect.DeepEqual(docs, input.Documents{}) {
		t.Fatal("partial result")
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	reader := readerFunc(func(p []byte) (int, error) { p[0] = 'x'; cancel(); return 1, io.EOF })
	docs, err = input.ReadDocuments(ctx, contract.Request{Policy: "-"}, reader)
	failure(t, err, 6, "interrupted")
	if !reflect.DeepEqual(docs, input.Documents{}) {
		t.Fatal("partial result")
	}
}

func TestDocumentBudgets(t *testing.T) {
	for _, tc := range []struct {
		name    string
		request contract.Request
		limit   int
	}{
		{"document", contract.Request{Policy: "-"}, contract.MaxDocumentBytes},
		{"selection", contract.Request{Paths: "-"}, contract.MaxSelectionBytes},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, delta := range []int{-1, 0, 1} {
				raw := bytes.Repeat([]byte(" "), tc.limit+delta)
				docs, err := input.ReadDocuments(context.Background(), tc.request, bytes.NewReader(raw))
				if delta == 1 {
					failure(t, err, 4, "resource_limit")
					if !reflect.DeepEqual(docs, input.Documents{}) {
						t.Fatal("partial result")
					}
				} else if err != nil || !bytes.Equal(append(docs.Policy, docs.Selection...), raw) {
					t.Fatal("not exact", err)
				}
			}
		})
	}
	root := t.TempDir()
	full := put(t, root, "full", bytes.Repeat([]byte("x"), contract.MaxDocumentBytes))
	for _, delta := range []int{-1, 0, 1} {
		last := put(t, root, "last", bytes.Repeat([]byte("y"), 1+delta))
		almost := put(t, root, "almost", bytes.Repeat([]byte("z"), contract.MaxDocumentBytes-1))
		docs, err := input.ReadDocuments(context.Background(), contract.Request{Policy: full, Inventory: almost, Result: last}, nil)
		if delta == 1 {
			failure(t, err, 4, "resource_limit")
			if !reflect.DeepEqual(docs, input.Documents{}) {
				t.Fatal("partial aggregate")
			}
		} else if err != nil || len(docs.Policy)+len(docs.Inventory)+len(docs.Result) != contract.MaxInputBytes+delta {
			t.Fatal("aggregate", err)
		}
	}
}
