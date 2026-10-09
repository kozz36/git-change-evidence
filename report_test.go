package changeevidence

import (
	"bytes"
	"encoding/json"
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

func TestProjectEvidenceSyntheticSubjectBoundary(t *testing.T) {
	// Fictional caller data is intentionally preserved in canonical evidence, not human text.
	sentinels := []string{"SYNTHETIC_CALLER_NAME", "/synthetic/not-a-real-repository", "FAKE_CREDENTIAL_SENTINEL", "FAKE_ENV_VALUE", "SYNTHETIC_COMMAND_OUTPUT"}
	input := validReportInput()
	input.Subject = string(bytes.Join([][]byte{[]byte(sentinels[0]), []byte(sentinels[1]), []byte(sentinels[2]), []byte(sentinels[3]), []byte(sentinels[4])}, []byte("\n"))) + "\n\"<synthetic>\""
	document, err := NewReportV1(input)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := ProjectEvidence(document, ProjectionCanonicalJSON)
	if err != nil || !bytes.Equal(canonical, document.CanonicalBytes()) {
		t.Fatalf("canonical = %q, error = %v; want original bytes", canonical, err)
	}
	var decoded struct{ Subject string }
	if err := json.Unmarshal(canonical, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Subject != input.Subject {
		t.Fatalf("Subject = %q, want %q", decoded.Subject, input.Subject)
	}
	roundTrip, err := DecodeCanonical(canonical)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(roundTrip.CanonicalBytes(), canonical) || roundTrip.Digest() != document.Digest() || roundTrip.Provenance() != input.Provenance {
		t.Fatal("DecodeCanonical changed bytes, digest, or provenance")
	}
	baseline := projectableEvidence(t)
	want, err := ProjectEvidence(baseline, ProjectionHumanText)
	if err != nil {
		t.Fatal(err)
	}
	// Reuse the fixed-label projection oracle; only the Subject-dependent digest differs.
	want = bytes.Replace(want, []byte(baseline.Digest()), []byte(document.Digest()), 1)
	human, err := ProjectEvidence(roundTrip, ProjectionHumanText)
	if err != nil || !bytes.Equal(human, want) {
		t.Fatalf("human = %q, error = %v; want %q", human, err, want)
	}
	for _, sentinel := range sentinels {
		if bytes.Contains(human, []byte(sentinel)) {
			t.Fatalf("human output includes synthetic Subject sentinel %q", sentinel)
		}
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
