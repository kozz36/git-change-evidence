package input

import (
	"context"
	"io"
	"os"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

// Documents retains exact selected transport bytes, not decoded canonical contracts.
type Documents struct{ Policy, Inventory, Result, Selection []byte }

// ReadDocuments reads only explicit document selectors; it never acquires source paths.
// Cancellation is cooperative before/after reads, not a blocked-I/O deadline guarantee.
// On failure it returns no partial documents and performs no filesystem writes.
func ReadDocuments(ctx context.Context, r contract.Request, stdin io.Reader) (Documents, error) {
	var docs Documents
	if err := validateSelectors(r); err != nil {
		return Documents{}, err
	}
	total := 0
	for _, item := range []struct {
		selector string
		target   *[]byte
		limit    int
	}{
		{r.Policy, &docs.Policy, contract.MaxDocumentBytes},
		{r.Inventory, &docs.Inventory, contract.MaxDocumentBytes},
		{r.Result, &docs.Result, contract.MaxDocumentBytes},
		{r.Paths, &docs.Selection, contract.MaxSelectionBytes},
	} {
		if ctx.Err() != nil {
			return Documents{}, fail(6, "interrupted")
		}
		if item.selector == "" {
			continue
		}
		reader := stdin
		var file *os.File
		if item.selector != "-" {
			var err error
			file, err = os.Open(item.selector)
			if ctx.Err() != nil {
				if file != nil {
					file.Close()
				}
				return Documents{}, fail(6, "interrupted")
			}
			if err != nil {
				return Documents{}, fail(3, "input_unavailable")
			}
			reader = file
		}
		if reader == nil {
			return Documents{}, fail(3, "input_unavailable")
		}
		limit := min(item.limit, contract.MaxInputBytes-total)
		// One bounded sentinel byte distinguishes exact capacity from overflow.
		raw, err := io.ReadAll(io.LimitReader(contextReader{ctx, reader}, int64(limit)+1))
		if file != nil {
			file.Close()
		}
		if ctx.Err() != nil {
			return Documents{}, fail(6, "interrupted")
		}
		if len(raw) > limit {
			return Documents{}, fail(4, "resource_limit")
		}
		if err != nil {
			return Documents{}, fail(3, "input_unavailable")
		}
		total += len(raw)
		*item.target = raw
	}
	return docs, nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := r.reader.Read(p)
	if cancelled := r.ctx.Err(); cancelled != nil {
		return n, cancelled
	}
	return n, err
}
