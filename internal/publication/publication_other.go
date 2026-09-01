//go:build !linux

package publicationadapter

import "context"

func publish(context.Context, string, []byte, string) (Result, error) {
	return Result{}, failure(Unavailable)
}
