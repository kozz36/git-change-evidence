package changeevidence

import (
	"bytes"
	"sort"
)

func NewAccountingResultV1(policy PolicyDocumentV1, inventory InventoryDocumentV1, snapshot CommittedSnapshot) (ResultDocumentV1, error) {
	view, policyDigest, inventoryDigest, err := validateResultAntecedents(policy, inventory)
	if err != nil {
		return ResultDocumentV1{}, err
	}
	revisions, accountingEntries, err := validateResultSnapshot(snapshot)
	if err != nil {
		return ResultDocumentV1{}, err
	}
	classifications, totals, ok := resultAccountingDomain(view).classifyAndTotal(accountingEntries)
	if !ok {
		return ResultDocumentV1{}, &ContractError{"totals", "overflow"}
	}
	sort.Slice(classifications, func(i, j int) bool {
		return bytes.Compare(classifications[i].entry.path, classifications[j].entry.path) < 0
	})
	entries := make([]ResultEntryV1, len(classifications))
	for i, classification := range classifications {
		measurement := ResultEntryMeasurementV1{Kind: ResultEntryNonCountableV1}
		if classification.entry.lines.Countable {
			measurement = ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: classification.entry.lines.Additions, Deletions: classification.entry.lines.Deletions}
		}
		entries[i] = ResultEntryV1{Path: append([]byte(nil), classification.entry.path...), Category: view.Categories[classification.category].Name, Measurement: measurement}
	}
	resultTotals := make([]ResultCategoryTotalV1, len(totals))
	for i, total := range totals {
		resultTotals[i] = ResultCategoryTotalV1{Category: total.Category, Additions: total.Additions, Deletions: total.Deletions, NonCountable: total.NonCountable}
	}
	observations, err := resultObservations(view, totals)
	if err != nil {
		return ResultDocumentV1{}, err
	}
	return resultDocument(policyDigest, inventoryDigest, revisions, entries, resultTotals, observations)
}
