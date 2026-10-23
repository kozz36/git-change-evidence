package changeevidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// CommittedSnapshotContractV1 identifies the independent full Git-facts schema.
const CommittedSnapshotContractV1 ContractVersion = "git-change-evidence.committed-snapshot/v1"

// SnapshotDocumentV1 owns canonical JSON-plus-LF bytes and immutable typed facts.
// Structural validity does not authenticate acquisition or bind a V1 report.
type SnapshotDocumentV1 struct {
	canonical []byte
	digest    DocumentDigest
	facts     SnapshotFactsV1
}

// NewSnapshotV1 assembles a document without Git, filesystem or policy effects.
func NewSnapshotV1(facts SnapshotFactsV1) (SnapshotDocumentV1, error) {
	if err := validateSnapshotFactsV1(facts.revisions, facts.entries); err != nil {
		return SnapshotDocumentV1{}, err
	}
	canonical, err := json.Marshal(snapshotWireFromFactsV1(facts))
	if err != nil {
		return SnapshotDocumentV1{}, err
	}
	canonical = append(canonical, '\n')
	sum := sha256.Sum256(canonical)
	return SnapshotDocumentV1{
		canonical: canonical,
		digest:    DocumentDigest(hex.EncodeToString(sum[:])),
		facts:     facts,
	}, nil
}

// CanonicalBytes returns owned bytes, including exactly one terminal LF.
func (document SnapshotDocumentV1) CanonicalBytes() []byte {
	return append([]byte(nil), document.canonical...)
}

// Digest returns SHA-256 of the exact canonical bytes, including LF.
func (document SnapshotDocumentV1) Digest() DocumentDigest { return document.digest }

// Revisions returns the unchanged native revision identities.
func (document SnapshotDocumentV1) Revisions() RevisionIdentity { return document.facts.Revisions() }

// Entries returns independent copies of all mutable entry and side fields.
func (document SnapshotDocumentV1) Entries() []SnapshotEntryV1 {
	if document.canonical == nil {
		return nil
	}
	return document.facts.Entries()
}
