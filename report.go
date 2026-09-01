package changeevidence

type ProjectionFormat uint8

const (
	ProjectionCanonicalJSON ProjectionFormat = iota + 1
	ProjectionHumanText
)

type DiagnosticCategory uint8

const DiagnosticInput DiagnosticCategory = iota + 1

type DiagnosticCode uint8

const (
	DiagnosticInvalidEvidence DiagnosticCode = iota + 1
	DiagnosticUnsupportedFormat
)

type ProjectionError struct {
	Category DiagnosticCategory
	Code     DiagnosticCode
}

func (e *ProjectionError) Error() string {
	switch e.Code {
	case DiagnosticInvalidEvidence:
		return "projection input: invalid evidence"
	case DiagnosticUnsupportedFormat:
		return "projection input: unsupported format"
	default:
		return "projection input: unknown diagnostic"
	}
}

func ProjectEvidence(e Evidence, format ProjectionFormat) ([]byte, error) {
	canonical := e.CanonicalBytes()
	if len(canonical) == 0 {
		return nil, &ProjectionError{Category: DiagnosticInput, Code: DiagnosticInvalidEvidence}
	}
	switch format {
	case ProjectionCanonicalJSON:
		return canonical, nil
	case ProjectionHumanText:
		provenance := e.Provenance()
		return []byte("schema: " + string(EvidenceContractV1) + "\n" +
			"kind: " + string(ReportDocument) + "\n" +
			"evidence_sha256: " + string(e.Digest()) + "\n" +
			"base_revision: " + provenance.Revisions.Base + "\n" +
			"head_revision: " + provenance.Revisions.Head + "\n" +
			"input_policy_sha256: " + string(provenance.InputPolicyDigest) + "\n" +
			"inventory_sha256: " + string(provenance.InventoryDigest) + "\n" +
			"accounting_sha256: " + string(provenance.AccountingDigest) + "\n"), nil
	default:
		return nil, &ProjectionError{Category: DiagnosticInput, Code: DiagnosticUnsupportedFormat}
	}
}
