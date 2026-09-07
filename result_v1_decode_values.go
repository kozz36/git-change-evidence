package changeevidence

import (
	"encoding/base64"
	"encoding/json"
	"unicode/utf8"
)

type resultDecodeCandidate struct {
	policyDigest    DocumentDigest
	inventoryDigest DocumentDigest
	snapshot        CommittedSnapshot
}

type resultDecodeCandidateWire struct {
	Schema                 ContractVersion                  `json:"schema"`
	AccountingPolicySHA256 DocumentDigest                   `json:"accounting_policy_sha256"`
	InventorySHA256        DocumentDigest                   `json:"inventory_sha256"`
	Revisions              resultWireRevisions              `json:"revisions"`
	Entries                []resultDecodeCandidateWireEntry `json:"entries"`
}

type resultDecodeCandidateWireEntry struct {
	PathB64     string                               `json:"path_b64"`
	Measurement resultDecodeCandidateWireMeasurement `json:"measurement"`
}

type resultDecodeCandidateWireMeasurement struct {
	Kind      ResultEntryMeasurementKind `json:"kind"`
	Additions uint64                     `json:"additions"`
	Deletions uint64                     `json:"deletions"`
}

func materializeResultDecodeCandidate(raw []byte) (resultDecodeCandidate, error) {
	if !utf8.Valid(raw) {
		return resultDecodeCandidate{}, &ContractError{"document", "invalid"}
	}
	if err := validateResultRootShape(raw); err != nil {
		return resultDecodeCandidate{}, err
	}
	var wire resultDecodeCandidateWire
	if err := json.Unmarshal(raw, &wire); err != nil {
		return resultDecodeCandidate{}, &ContractError{"document", "invalid"}
	}
	if wire.Schema != AccountingResultContractV1 {
		return resultDecodeCandidate{}, &ContractError{"schema", "unsupported_version"}
	}
	if !isDigest(string(wire.AccountingPolicySHA256)) {
		return resultDecodeCandidate{}, &ContractError{"accounting_policy_sha256", "invalid_digest"}
	}
	if !isDigest(string(wire.InventorySHA256)) {
		return resultDecodeCandidate{}, &ContractError{"inventory_sha256", "invalid_digest"}
	}
	entries := make([]CommittedChange, len(wire.Entries))
	for index, entry := range wire.Entries {
		path, err := base64.StdEncoding.DecodeString(entry.PathB64)
		if err != nil || base64.StdEncoding.EncodeToString(path) != entry.PathB64 {
			return resultDecodeCandidate{}, &ContractError{resultEntryField(index) + ".path_b64", "invalid_base64"}
		}
		entries[index] = CommittedChange{
			Path: string(append([]byte(nil), path...)),
			Lines: CommittedLineCounts{
				Additions: entry.Measurement.Additions,
				Deletions: entry.Measurement.Deletions,
				Countable: entry.Measurement.Kind == ResultEntryCountableV1,
			},
		}
	}
	return resultDecodeCandidate{
		policyDigest:    wire.AccountingPolicySHA256,
		inventoryDigest: wire.InventorySHA256,
		snapshot:        NewCommittedSnapshot(GitObjectID(wire.Revisions.Base), GitObjectID(wire.Revisions.Head), entries),
	}, nil
}
