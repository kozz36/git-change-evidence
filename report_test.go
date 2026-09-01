package changeevidence

import (
	"bytes"
	"testing"
)

func projectableEvidence(t *testing.T) Evidence {
	t.Helper()
	evidence, err := NewReportV1(validReportInput())
	if err != nil {
		t.Fatalf("NewReportV1() error = %v", err)
	}
	return evidence
}

func TestProjectEvidenceCanonicalJSON(t *testing.T) {
	evidence := projectableEvidence(t)

	got, err := ProjectEvidence(evidence, ProjectionCanonicalJSON)
	if err != nil {
		t.Fatalf("ProjectEvidence() error = %v", err)
	}
	if want := evidence.CanonicalBytes(); !bytes.Equal(got, want) {
		t.Fatalf("ProjectEvidence() = %q, want %q", got, want)
	}
}

func TestProjectEvidenceHumanText(t *testing.T) {
	evidence := projectableEvidence(t)

	got, err := ProjectEvidence(evidence, ProjectionHumanText)
	if err != nil {
		t.Fatalf("ProjectEvidence() error = %v", err)
	}
	const want = "schema: git-change-evidence.evidence/v1\n" +
		"kind: report\n" +
		"evidence_sha256: b7145afd882b9cb401fb7d61873b6ea748d816fdbc385eea1f732bece0fc43d0\n" +
		"base_revision: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n" +
		"head_revision: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n" +
		"input_policy_sha256: cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\n" +
		"inventory_sha256: dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd\n" +
		"accounting_sha256: eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee\n"
	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("ProjectEvidence() = %q, want %q", got, want)
	}
}

func TestProjectEvidenceRejectsInvalidEvidenceAndUnsupportedFormat(t *testing.T) {
	validEvidence := projectableEvidence(t)

	for _, test := range []struct {
		name     string
		evidence Evidence
		format   ProjectionFormat
		code     DiagnosticCode
		message  string
	}{
		{"zero evidence", Evidence{}, ProjectionCanonicalJSON, DiagnosticInvalidEvidence, "projection input: invalid evidence"},
		{"unknown format", validEvidence, ProjectionFormat(99), DiagnosticUnsupportedFormat, "projection input: unsupported format"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := ProjectEvidence(test.evidence, test.format)
			if got != nil {
				t.Fatalf("ProjectEvidence() output = %q, want nil", got)
			}
			projectionErr, ok := err.(*ProjectionError)
			if !ok {
				t.Fatalf("ProjectEvidence() error = %T %v, want *ProjectionError", err, err)
			}
			if projectionErr.Category != DiagnosticInput || projectionErr.Code != test.code {
				t.Fatalf("ProjectionError = %#v, want category %d, code %d", projectionErr, DiagnosticInput, test.code)
			}
			if got := projectionErr.Error(); got != test.message {
				t.Fatalf("ProjectionError.Error() = %q, want %q", got, test.message)
			}
		})
	}
}

func TestProjectEvidenceOwnsCanonicalJSONOutput(t *testing.T) {
	evidence := projectableEvidence(t)
	want := evidence.CanonicalBytes()

	got, err := ProjectEvidence(evidence, ProjectionCanonicalJSON)
	if err != nil {
		t.Fatalf("ProjectEvidence() error = %v", err)
	}
	got[0] = 'x'

	again, err := ProjectEvidence(evidence, ProjectionCanonicalJSON)
	if err != nil {
		t.Fatalf("ProjectEvidence() error = %v", err)
	}
	if !bytes.Equal(again, want) {
		t.Fatalf("ProjectEvidence() after output mutation = %q, want %q", again, want)
	}
}
