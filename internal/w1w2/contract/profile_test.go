package contract

import (
	"testing"
	"time"
)

func TestInitialProfileUsesSelectedAdmissionLimits(t *testing.T) {
	for _, tt := range []struct {
		name      string
		got, want int
	}{
		{"selection", MaxSelectionBytes, 1 << 20},
		{"document", MaxDocumentBytes, 16 << 20},
		{"input total", MaxInputBytes, 32 << 20},
		{"paths", MaxPaths, 128},
		{"path bytes", MaxPathBytes, 4095},
		{"path total", MaxTotalPathBytes, 512 << 10},
		{"component", MaxComponentBytes, 255},
		{"file content", MaxFileBytes, 1 << 20},
		{"content total", MaxContentBytes, 8 << 20},
		{"Git response", MaxGitResponseBytes, 16 << 20},
		{"Git total", MaxGitBytes, 64 << 20},
		{"diff entries", MaxDiffEntries, 10000},
		{"output", MaxOutputBytes, 16 << 20},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("limit = %d, want %d", tt.got, tt.want)
			}
		})
	}
	if OperationTimeout != 30*time.Second {
		t.Fatalf("cooperative timeout = %v, want 30s", OperationTimeout)
	}
	if SelectionSchema != "git-change-evidence.cli-path-selection/v1" || ReportSubject != "committed change evidence" {
		t.Fatal("selection transport or fixed report subject changed")
	}
}

func TestFailureCarriesOnlyStableTechnicalData(t *testing.T) {
	got := &Failure{Status: 4, Token: "resource_limit"}
	if got.Error() != "resource_limit" || got.Status != 4 {
		t.Fatalf("failure = %#v, text %q", got, got.Error())
	}
	request := Request{Stage: "inventory", SourceRoot: "/authorized", Paths: "-", Policy: "policy"}
	if request.Stage != "inventory" || request.Paths != "-" {
		t.Fatal("typed request lost explicit selectors")
	}
}
