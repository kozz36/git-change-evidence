//go:build !linux

package main

func openCensusSourceRoot(string) (censusSourceRoot, error) {
	return nil, errCensusUnavailable
}
