package changeevidence

const AccountingResultContractV1 ContractVersion = "git-change-evidence.accounting-result/v1"

type ResultEntryMeasurementKind string

const (
	ResultEntryCountableV1    ResultEntryMeasurementKind = "countable"
	ResultEntryNonCountableV1 ResultEntryMeasurementKind = "non_countable"
)

type ResultObservationKind string

const (
	ResultLineThresholdObservationV1 ResultObservationKind = "line_threshold"
	ResultRatioObservationV1         ResultObservationKind = "ratio"
)

type ResultEntryMeasurementV1 struct {
	Kind                 ResultEntryMeasurementKind
	Additions, Deletions uint64
}

type ResultEntryV1 struct {
	Path        []byte
	Category    string
	Measurement ResultEntryMeasurementV1
}

type ResultCategoryTotalV1 struct {
	Category                           string
	Additions, Deletions, NonCountable uint64
}

type ResultObservationV1 struct {
	Kind               ResultObservationKind
	Category           string
	Reference          string
	LineThresholdLimit uint64
	RatioLimit         string
	Available          bool
	Actual             uint64
	Numerator          uint64
	Denominator        uint64
	Exceeded           bool
}

type ResultDocumentV1 struct {
	canonical              []byte
	digest                 DocumentDigest
	accountingPolicyDigest DocumentDigest
	inventoryDigest        DocumentDigest
	revisions              RevisionIdentity
	entries                []ResultEntryV1
	totals                 []ResultCategoryTotalV1
	observations           []ResultObservationV1
}

func (r ResultDocumentV1) CanonicalBytes() []byte { return append([]byte(nil), r.canonical...) }

func (r ResultDocumentV1) Digest() DocumentDigest { return r.digest }

func (r ResultDocumentV1) AccountingPolicyDigest() DocumentDigest {
	return r.accountingPolicyDigest
}

func (r ResultDocumentV1) InventoryDigest() DocumentDigest { return r.inventoryDigest }

func (r ResultDocumentV1) Revisions() RevisionIdentity { return r.revisions }

func (r ResultDocumentV1) Entries() []ResultEntryV1 {
	if r.canonical == nil {
		return nil
	}
	return cloneResultEntries(r.entries)
}

func (r ResultDocumentV1) Totals() []ResultCategoryTotalV1 {
	if r.canonical == nil {
		return nil
	}
	return cloneResultTotals(r.totals)
}

func (r ResultDocumentV1) Observations() []ResultObservationV1 {
	if r.canonical == nil {
		return nil
	}
	return cloneResultObservations(r.observations)
}

func cloneResultEntries(entries []ResultEntryV1) []ResultEntryV1 {
	copy := make([]ResultEntryV1, len(entries))
	for i, entry := range entries {
		copy[i] = entry
		copy[i].Path = append([]byte(nil), entry.Path...)
	}
	return copy
}

func cloneResultTotals(totals []ResultCategoryTotalV1) []ResultCategoryTotalV1 {
	cloned := make([]ResultCategoryTotalV1, len(totals))
	copy(cloned, totals)
	return cloned
}

func cloneResultObservations(observations []ResultObservationV1) []ResultObservationV1 {
	cloned := make([]ResultObservationV1, len(observations))
	copy(cloned, observations)
	return cloned
}
