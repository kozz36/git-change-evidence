package input_test

import (
	"reflect"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

func TestParseClosedGrammar(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want contract.Request
	}{
		{"inventory", []string{"inventory", "--source-root", "/root", "--policy", "p", "--paths", "s"}, contract.Request{Stage: "inventory", SourceRoot: "/root", Policy: "p", Paths: "s"}},
		{"result", []string{"result", "--head", "literal;$", "--base", "base", "--inventory", "-", "--policy", "p", "--repo", "/repo"}, contract.Request{Stage: "result", Repository: "/repo", Base: "base", Head: "literal;$", Inventory: "-", Policy: "p"}},
		{"report", []string{"report", "--result", "r", "--inventory", "i", "--policy", "-"}, contract.Request{Stage: "report", Policy: "-", Inventory: "i", Result: "r"}},
		{"literal selector", []string{"report", "--policy", "--literal-file", "--inventory", "i", "--result", "r"}, contract.Request{Stage: "report", Policy: "--literal-file", Inventory: "i", Result: "r"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := input.Parse(tc.args)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("%+v %v", got, err)
			}
		})
	}
	for _, tc := range []struct {
		name, stage string
		args        []string
	}{
		{"none", "command", nil},
		{"unknown", "command", []string{"invent"}},
		{"stage retained", "inventory", []string{"inventory"}},
		{"relative root", "inventory", []string{"inventory", "--source-root", "root", "--policy", "p", "--paths", "s"}},
		{"implicit root", "inventory", []string{"inventory", "--policy", "p", "--paths", "s"}},
		{"duplicate", "report", []string{"report", "--policy", "p", "--policy", "q", "--result", "r"}},
		{"unknown flag", "report", []string{"report", "--policy", "p", "--paths", "s", "--result", "r"}},
		{"equals", "report", []string{"report", "--policy=p", "--inventory", "i", "--result", "r"}},
		{"empty", "report", []string{"report", "--policy", "", "--inventory", "i", "--result", "r"}},
		{"missing value", "report", []string{"report", "--policy", "p", "--inventory", "i", "--result"}},
		{"extra", "report", []string{"report", "--policy", "p", "--inventory", "i", "--result", "r", "extra"}},
		{"two stdin", "report", []string{"report", "--policy", "-", "--inventory", "-", "--result", "r"}},
		{"relative repo", "result", []string{"result", "--repo", ".", "--base", "b", "--head", "h", "--policy", "p", "--inventory", "i"}},
		{"no base default", "result", []string{"result", "--repo", "/r", "--head", "h", "--policy", "p", "--inventory", "i"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := input.Parse(tc.args)
			failure(t, err, 2, "invalid_arguments")
			if got.Stage != tc.stage {
				t.Fatalf("stage %q", got.Stage)
			}
		})
	}
}
