package input_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

func TestSnapshotParseExplicitGrammar(t *testing.T) {
	for _, args := range [][]string{
		{"snapshot", "--repo", "/fixture", "--base", "--literal-base", "--head", "--literal-head"},
		{"snapshot", "--head", "--literal-head", "--repo", "/fixture", "--base", "--literal-base"},
	} {
		before := append([]string(nil), args...)
		got, err := input.Parse(args)
		want := contract.Request{Stage: "snapshot", Repository: "/fixture", Base: "--literal-base", Head: "--literal-head"}
		if err != nil || got != want || !reflect.DeepEqual(before, args) {
			t.Fatalf("explicit independent request: %+v %v", got, err)
		}
	}
}

func TestSnapshotParseRejectsImplicitOrDocumentInputs(t *testing.T) {
	valid := []string{"snapshot", "--repo", "/fixture", "--base", "base", "--head", "head"}
	cases := map[string][]string{
		"missing":      {"snapshot"},
		"missing_repo": {"snapshot", "--base", "base", "--head", "head"},
		"missing_base": {"snapshot", "--repo", "/fixture", "--head", "head"},
		"missing_head": {"snapshot", "--repo", "/fixture", "--base", "base"},
		"positional":   {"snapshot", "/fixture", "base", "head"},
		"relative":     {"snapshot", "--repo", "relative", "--base", "base", "--head", "head"},
		"equals":       {"snapshot", "--repo=/fixture", "--base", "base", "--head", "head"},
		"duplicate":    {"snapshot", "--repo", "/fixture", "--base", "base", "--base", "head"},
		"empty":        {"snapshot", "--repo", "/fixture", "--base", "", "--head", "head"},
	}
	for _, flag := range []string{"--policy", "--inventory", "--paths", "--output", "--result", "--source-root"} {
		cases[flag] = append(append([]string(nil), valid...), flag, "-")
	}
	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			before := append([]string(nil), args...)
			got, err := input.Parse(args)
			var failure *contract.Failure
			if !errors.As(err, &failure) || failure.Status != 2 || failure.Token != "invalid_arguments" || got.Stage != "snapshot" || !reflect.DeepEqual(before, args) {
				t.Fatalf("invalid request: %+v %v", got, err)
			}
		})
	}
}
