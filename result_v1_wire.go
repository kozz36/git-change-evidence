package changeevidence

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

type resultWire struct {
	Schema                 ContractVersion     `json:"schema"`
	AccountingPolicySHA256 DocumentDigest      `json:"accounting_policy_sha256"`
	InventorySHA256        DocumentDigest      `json:"inventory_sha256"`
	Revisions              resultWireRevisions `json:"revisions"`
	Entries                []resultWireEntry   `json:"entries"`
	Totals                 []resultWireTotal   `json:"totals"`
	Observations           []json.RawMessage   `json:"observations"`
}

type resultWireRevisions struct {
	Base string `json:"base"`
	Head string `json:"head"`
}

type resultWireEntry struct {
	PathB64     string          `json:"path_b64"`
	Category    string          `json:"category"`
	Measurement json.RawMessage `json:"measurement"`
}

type resultWireCountableMeasurement struct {
	Kind      ResultEntryMeasurementKind `json:"kind"`
	Additions uint64                     `json:"additions"`
	Deletions uint64                     `json:"deletions"`
}

type resultWireNonCountableMeasurement struct {
	Kind ResultEntryMeasurementKind `json:"kind"`
}

type resultWireTotal struct {
	Category     string `json:"category"`
	Additions    uint64 `json:"additions"`
	Deletions    uint64 `json:"deletions"`
	NonCountable uint64 `json:"non_countable"`
}

func resultDocument(policyDigest, inventoryDigest DocumentDigest, revisions RevisionIdentity, entries []ResultEntryV1, totals []ResultCategoryTotalV1, observations []ResultObservationV1) (ResultDocumentV1, error) {
	wire := resultWire{
		Schema:                 AccountingResultContractV1,
		AccountingPolicySHA256: policyDigest,
		InventorySHA256:        inventoryDigest,
		Revisions:              resultWireRevisions{Base: revisions.Base, Head: revisions.Head},
		Entries:                make([]resultWireEntry, len(entries)),
		Totals:                 make([]resultWireTotal, len(totals)),
		Observations:           make([]json.RawMessage, len(observations)),
	}
	for index, entry := range entries {
		measurement, err := resultWireMeasurement(entry.Measurement)
		if err != nil {
			return ResultDocumentV1{}, err
		}
		wire.Entries[index] = resultWireEntry{
			PathB64:     base64.StdEncoding.EncodeToString(entry.Path),
			Category:    entry.Category,
			Measurement: measurement,
		}
	}
	for index, total := range totals {
		wire.Totals[index] = resultWireTotal{
			Category: total.Category, Additions: total.Additions, Deletions: total.Deletions, NonCountable: total.NonCountable,
		}
	}
	for index, observation := range observations {
		encoded, err := resultWireObservation(observation)
		if err != nil {
			return ResultDocumentV1{}, err
		}
		wire.Observations[index] = encoded
	}
	canonical, err := json.Marshal(wire)
	if err != nil {
		return ResultDocumentV1{}, err
	}
	canonical = append(canonical, '\n')
	sum := sha256.Sum256(canonical)
	return ResultDocumentV1{
		canonical:              append([]byte(nil), canonical...),
		digest:                 DocumentDigest(hex.EncodeToString(sum[:])),
		accountingPolicyDigest: policyDigest,
		inventoryDigest:        inventoryDigest,
		revisions:              revisions,
		entries:                cloneResultEntries(entries),
		totals:                 cloneResultTotals(totals),
		observations:           cloneResultObservations(observations),
	}, nil
}

func resultWireMeasurement(measurement ResultEntryMeasurementV1) (json.RawMessage, error) {
	switch measurement.Kind {
	case ResultEntryCountableV1:
		return json.Marshal(resultWireCountableMeasurement{measurement.Kind, measurement.Additions, measurement.Deletions})
	case ResultEntryNonCountableV1:
		return json.Marshal(resultWireNonCountableMeasurement{measurement.Kind})
	default:
		return nil, &ContractError{"entries.measurement.kind", "invalid"}
	}
}
