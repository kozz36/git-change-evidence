// Package input implements closed W1/W2 requests and bounded document transport.
package input

import (
	"path/filepath"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func fail(status int, token string) error { return &contract.Failure{Status: status, Token: token} }

// Parse accepts only explicit separate flag/value pairs. Known stages survive errors.
func Parse(args []string) (contract.Request, error) {
	r := contract.Request{Stage: "command"}
	invalid := func() (contract.Request, error) { return r, fail(2, "invalid_arguments") }
	if len(args) == 0 {
		return invalid()
	}
	var required []string
	switch args[0] {
	case "inventory":
		required = []string{"--source-root", "--policy", "--paths"}
	case "result":
		required = []string{"--repo", "--base", "--head", "--policy", "--inventory"}
	case "report":
		required = []string{"--policy", "--inventory", "--result"}
	default:
		return invalid()
	}
	r.Stage = args[0]
	fields := map[string]*string{"--source-root": &r.SourceRoot, "--repo": &r.Repository, "--base": &r.Base, "--head": &r.Head, "--policy": &r.Policy, "--inventory": &r.Inventory, "--result": &r.Result, "--paths": &r.Paths}
	allowed := make(map[string]bool, len(required))
	for _, name := range required {
		allowed[name] = true
	}
	if len(args) != 1+2*len(required) {
		return invalid()
	}
	for i := 1; i < len(args); i += 2 {
		name, value := args[i], args[i+1]
		if !allowed[name] || value == "" {
			return invalid()
		}
		*fields[name] = value
		delete(allowed, name)
	}
	if len(allowed) != 0 || (r.Stage == "inventory" && !filepath.IsAbs(r.SourceRoot)) || (r.Stage == "result" && !filepath.IsAbs(r.Repository)) {
		return invalid()
	}
	if err := validateSelectors(r); err != nil {
		return r, err
	}
	return r, nil
}

func validateSelectors(r contract.Request) error {
	count := 0
	for _, selector := range []string{r.Policy, r.Inventory, r.Result, r.Paths} {
		if selector == "-" {
			count++
		}
	}
	if count > 1 {
		return fail(2, "invalid_arguments")
	}
	return nil
}
