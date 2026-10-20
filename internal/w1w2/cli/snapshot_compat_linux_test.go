//go:build linux

package cli

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func TestSnapshotCLIContextAndIncompleteOutput(t *testing.T) {
	facts := snapshotFacts(t)
	for _, scenario := range []string{"earlier_deadline", "pre_canceled", "mid_canceled", "short_write"} {
		t.Run(scenario, func(t *testing.T) {
			deadline := time.Now().Add(time.Second)
			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()
			if scenario == "pre_canceled" {
				cancel()
			}
			calls := 0
			source := func(got context.Context, _ gitadapter.Request, _ gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error) {
				calls++
				if d, _ := got.Deadline(); d != deadline {
					t.Fatal("earlier caller deadline reset")
				}
				if scenario == "mid_canceled" {
					cancel()
				}
				return facts, nil
			}
			var out, diagnostic bytes.Buffer
			var writer io.Writer = &out
			if scenario == "short_write" {
				writer = snapshotShortWriter{&out}
			}
			code := runWithSources(ctx, snapshotArguments(), unusedInput{}, writer, &diagnostic, nil, source)
			if scenario == "earlier_deadline" {
				want, err := evidence.NewSnapshotV1(facts)
				if err != nil || code != 0 || !bytes.Equal(out.Bytes(), want.CanonicalBytes()) || diagnostic.Len() != 0 || calls != 1 {
					t.Fatalf("earlier deadline: %d %q %q", code, out.Bytes(), diagnostic.String())
				}
				return
			}
			wantCode, token, expectedCalls, prefix := 6, "interrupted", 1, ""
			if scenario == "pre_canceled" {
				expectedCalls = 0
			}
			if scenario == "short_write" {
				wantCode, token, prefix = 9, "output_incomplete", "{"
			}
			if code != wantCode || out.String() != prefix || diagnostic.String() != "w1w2 snapshot: "+token+"\n" || calls != expectedCalls {
				t.Fatalf("code%d out%q stderr%q calls%d", code, out.Bytes(), diagnostic.String(), calls)
			}
		})
	}
}

type snapshotShortWriter struct{ out *bytes.Buffer }

func (writer snapshotShortWriter) Write(raw []byte) (int, error) { return writer.out.Write(raw[:1]) }

func TestSnapshotCLIPreservesExistingStageCanonicalBytes(t *testing.T) {
	policy, inventory, result := documents(t)
	p, i, r := documentFile(t, policy.CanonicalBytes()), documentFile(t, inventory.CanonicalBytes()), documentFile(t, result.CanonicalBytes())
	selection := documentFile(t, []byte("{\"schema\":\"git-change-evidence.cli-path-selection/v1\",\"paths_b64\":[]}\n"))
	binding := evidence.ResultV1ReportBinding{Policy: policy, Inventory: inventory, Result: result}
	report, err := evidence.NewReportV1FromResult("committed change evidence", binding)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		args []string
		want []byte
	}{
		{[]string{"inventory", "--source-root", t.TempDir(), "--policy", p, "--paths", selection}, inventory.CanonicalBytes()},
		{[]string{"result", "--repo", t.TempDir(), "--base", "base", "--head", "head", "--policy", p, "--inventory", i}, result.CanonicalBytes()},
		{[]string{"report", "--policy", p, "--inventory", i, "--result", r}, report.CanonicalBytes()},
	} {
		t.Run(tc.args[0], func(t *testing.T) {
			oldSource := func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.CommittedSnapshot, error) {
				return emptySnapshot(), nil
			}
			newSource := func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error) {
				t.Fatal("existing stage reached snapshot source")
				return evidence.SnapshotFactsV1{}, nil
			}
			var out, diagnostic bytes.Buffer
			code := runWithSources(context.Background(), tc.args, unusedInput{}, &out, &diagnostic, oldSource, newSource)
			if code != 0 || !bytes.Equal(out.Bytes(), tc.want) || diagnostic.Len() != 0 {
				t.Fatalf("%d %q %q", code, out.Bytes(), diagnostic.String())
			}
		})
	}
}

func TestSnapshotCLIAdditiveFamilyHelp(t *testing.T) {
	code, out, diagnostic := invoke(t, []string{"--help"})
	line := "  gce w1w2 snapshot --repo <absolute-repository> --base <revision> --head <revision>\n"
	if code != 0 || diagnostic != "" || !bytes.Contains(out, []byte(line)) || !bytes.Contains(out, []byte("gce w1w2 inventory")) || !bytes.Contains(out, []byte("gce w1w2 result")) || !bytes.Contains(out, []byte("gce w1w2 report")) {
		t.Fatalf("%d %q %q", code, out, diagnostic)
	}
}
