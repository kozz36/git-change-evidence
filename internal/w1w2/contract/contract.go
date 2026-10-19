// Package contract holds the private W1/W2 shell's typed protocol and admission rules.
package contract

// Request contains explicit selectors, never an implicit working-directory or HEAD default.
type Request struct {
	Stage, SourceRoot, Repository, Base, Head string
	Policy, Inventory, Result, Paths          string
}

// Failure carries stable technical data without user paths or subprocess diagnostics.
type Failure struct {
	Status int
	Token  string
}

func (failure *Failure) Error() string { return failure.Token }
