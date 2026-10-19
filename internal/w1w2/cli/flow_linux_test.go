package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func TestInventoryCanonicalSelectedBytes(t *testing.T) {
	policy, _, _ := documents(t)
	root := t.TempDir()
	for name, data := range map[string]string{"a": "captured\n", "unselected": "not inventoried"} {
		if err := os.WriteFile(root+"/"+name, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("missing", root+"/link"); err != nil {
		t.Fatal(err)
	}
	selection, err := json.Marshal(map[string]any{"schema": contract.SelectionSchema, "paths_b64": []string{base64.StdEncoding.EncodeToString([]byte("link")), base64.StdEncoding.EncodeToString([]byte("a"))}})
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"inventory", "--source-root", root, "--policy", documentFile(t, policy.CanonicalBytes()), "--paths", "-"}
	var out, diagnostic bytes.Buffer
	code := Run(context.Background(), args, bytes.NewReader(selection), &out, &diagnostic)
	want, err := evidence.NewInventoryV1(policy, []evidence.InventoryEntryInput{{Path: []byte("a"), Content: []byte("captured\n")}, {Path: []byte("link"), Content: []byte("missing")}})
	if err != nil || code != 0 || diagnostic.Len() != 0 || !bytes.Equal(out.Bytes(), want.CanonicalBytes()) {
		t.Fatalf("code=%d stderr=%q stdout=%q want=%q err=%v", code, diagnostic.String(), out.Bytes(), want.CanonicalBytes(), err)
	}
	if data, err := os.ReadFile(root + "/a"); err != nil || string(data) != "captured\n" {
		t.Fatal("source changed", err)
	}
}

func TestReportCanonicalBindingAndObservations(t *testing.T) {
	policy, inventory, result := documents(t)
	binding := evidence.ResultV1ReportBinding{Policy: policy, Inventory: inventory, Result: result}
	want, err := evidence.NewReportV1FromResult(contract.ReportSubject, binding)
	if err != nil {
		t.Fatal(err)
	}
	code, out, diagnostic := invoke(t, reportArguments(t))
	if code != 0 || diagnostic != "" || !bytes.Equal(out, want.CanonicalBytes()) {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, out, diagnostic)
	}
	if err := evidence.ValidateReportV1ResultProvenance(want, binding); err != nil {
		t.Fatal(err)
	}
}

func TestInventoryEmptyRequiresAccessibleRoot(t *testing.T) {
	policy, _, _ := documents(t)
	root := t.TempDir()
	args := []string{"inventory", "--source-root", root, "--policy", documentFile(t, policy.CanonicalBytes()), "--paths", documentFile(t, []byte("{\"schema\":\""+contract.SelectionSchema+"\",\"paths_b64\":[]}"))}
	code, out, diagnostic := invoke(t, args)
	want, err := evidence.NewInventoryV1(policy, nil)
	if err != nil || code != 0 || diagnostic != "" || !bytes.Equal(out, want.CanonicalBytes()) {
		t.Fatalf("%d %q %q %v", code, out, diagnostic, err)
	}
	args[2] = root + "/missing"
	requireFailure(t, args, 3, "w1w2 inventory: input_unavailable\n")
}

func TestCanonicalFailuresHaveExactEmptyStreams(t *testing.T) {
	policy, inventory, _ := documents(t)
	policyPath := documentFile(t, policy.CanonicalBytes())
	args := []string{"inventory", "--source-root", t.TempDir(), "--policy", policyPath, "--paths", documentFile(t, []byte("{\"schema\":\""+contract.SelectionSchema+"\",\"paths_b64\":[]}"))}
	for _, tc := range []struct {
		name, raw, token string
	}{
		{"noncanonical", " " + string(policy.CanonicalBytes()), "invalid_document"},
		{"version", string(bytes.Replace(policy.CanonicalBytes(), []byte("/v1"), []byte("/v2"), 1)), "unsupported_version"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			copyArgs := append([]string(nil), args...)
			copyArgs[4] = documentFile(t, []byte(tc.raw))
			requireFailure(t, copyArgs, 2, "w1w2 inventory: "+tc.token+"\n")
		})
	}
	requireFailure(t, []string{"unknown"}, 2, "w1w2 command: invalid_arguments\n")
	requireFailure(t, []string{"report", "--policy", "-", "--inventory", "-", "--result", "unused"}, 2, "w1w2 report: invalid_arguments\n")
	otherInput := evidence.PolicyInput(policy.View())
	otherInput.Categories[0].LineThreshold++
	otherPolicy, err := evidence.NewAccountingPolicyV1(otherInput)
	if err != nil {
		t.Fatal(err)
	}
	mismatch := []string{"report", "--policy", documentFile(t, otherPolicy.CanonicalBytes()), "--inventory", documentFile(t, inventory.CanonicalBytes()), "--result", "must-not-read"}
	// Keep transport valid so the canonical antecedent mismatch, not missing I/O, is exercised.
	_, _, result := documents(t)
	mismatch[6] = documentFile(t, result.CanonicalBytes())
	requireFailure(t, mismatch, 2, "w1w2 report: antecedent_mismatch\n")
}

func TestReportRejectsDifferentRetainedInventory(t *testing.T) {
	p, inventory, _ := documents(t)
	other, err := evidence.NewInventoryV1(p, []evidence.InventoryEntryInput{{Path: []byte("other"), Content: []byte("x")}})
	if err != nil {
		t.Fatal(err)
	}
	result, err := evidence.NewAccountingResultV1(p, other, emptySnapshot())
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"report", "--policy", documentFile(t, p.CanonicalBytes()), "--inventory", documentFile(t, inventory.CanonicalBytes()), "--result", documentFile(t, result.CanonicalBytes())}
	requireFailure(t, args, 2, "w1w2 report: antecedent_mismatch\n")
}

func TestFamilyHelpUsesDiscoveryStreamsWithoutStdin(t *testing.T) {
	code, out, diagnostic := invoke(t, []string{"--help"})
	if code != 0 || diagnostic != "" || !bytes.Equal(out, []byte(familyHelp)) {
		t.Fatalf("%d %q %q", code, out, diagnostic)
	}
	requireFailure(t, []string{"inventory", "--help"}, 2, "w1w2 inventory: invalid_arguments\n")
}
