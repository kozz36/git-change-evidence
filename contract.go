package changeevidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

type ContractVersion string
type DocumentKind string
type Digest string

const (
	EvidenceContractV1 ContractVersion = "git-change-evidence.evidence/v1"
	ReportDocument     DocumentKind    = "report"
)

type RevisionIdentity struct {
	Base string `json:"base"`
	Head string `json:"head"`
}

type Provenance struct {
	AccountingDigest  Digest           `json:"accounting_sha256"`
	InputPolicyDigest Digest           `json:"input_policy_sha256"`
	InventoryDigest   Digest           `json:"inventory_sha256"`
	Revisions         RevisionIdentity `json:"revisions"`
}

type ReportInput struct {
	Subject    string
	Provenance Provenance
}

type Evidence struct {
	canonical  []byte
	digest     Digest
	provenance Provenance
}

type ContractError struct {
	Field string
	Code  string
}

func (e *ContractError) Error() string { return "invalid contract " + e.Field + ": " + e.Code }

func NewReportV1(input ReportInput) (Evidence, error) {
	if err := validateReportInput(input); err != nil {
		return Evidence{}, err
	}
	canonical, err := json.Marshal(canonicalReport{ReportDocument, input.Provenance, EvidenceContractV1, input.Subject})
	if err != nil {
		return Evidence{}, err
	}
	canonical = append(canonical, '\n')
	sum := sha256.Sum256(canonical)
	return Evidence{append([]byte(nil), canonical...), Digest(hex.EncodeToString(sum[:])), input.Provenance}, nil
}

func (e Evidence) CanonicalBytes() []byte { return append([]byte(nil), e.canonical...) }
func (e Evidence) Digest() Digest         { return e.digest }
func (e Evidence) Provenance() Provenance { return e.provenance }

func validateReportInput(input ReportInput) error {
	if strings.TrimSpace(input.Subject) == "" {
		return &ContractError{"subject", "required"}
	}
	if !isRevision(input.Provenance.Revisions.Base) {
		return &ContractError{"provenance.revisions.base", "invalid_revision"}
	}
	if !isRevision(input.Provenance.Revisions.Head) {
		return &ContractError{"provenance.revisions.head", "invalid_revision"}
	}
	for _, value := range []struct {
		field string
		value Digest
	}{
		{"provenance.input_policy_sha256", input.Provenance.InputPolicyDigest},
		{"provenance.inventory_sha256", input.Provenance.InventoryDigest},
		{"provenance.accounting_sha256", input.Provenance.AccountingDigest},
	} {
		if !isDigest(string(value.value)) {
			return &ContractError{value.field, "invalid_digest"}
		}
	}
	return nil
}

func isRevision(value string) bool {
	return (len(value) == 40 || len(value) == 64) && isLowerHex(value)
}
func isDigest(value string) bool { return len(value) == sha256.Size*2 && isLowerHex(value) }

func isLowerHex(value string) bool {
	for _, character := range value {
		if !((character >= '0' && character <= '9') || (character >= 'a' && character <= 'f')) {
			return false
		}
	}
	return true
}

type canonicalReport struct {
	Kind       DocumentKind    `json:"kind"`
	Provenance Provenance      `json:"provenance"`
	Schema     ContractVersion `json:"schema"`
	Subject    string          `json:"subject"`
}
