package cli

import (
	"context"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/acquire"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

func compose(ctx context.Context, request contract.Request, docs input.Documents, source snapshotSource) ([]byte, error) {
	policy, err := evidence.DecodeAccountingPolicyV1(docs.Policy)
	if err != nil {
		return nil, err
	}
	if request.Stage == "inventory" {
		paths, err := input.DecodeSelection(docs.Selection)
		if err != nil {
			return nil, err
		}
		entries, err := acquire.Read(ctx, request.SourceRoot, paths)
		if err != nil {
			return nil, err
		}
		document, err := evidence.NewInventoryV1(policy, entries)
		return document.CanonicalBytes(), err
	}
	inventory, err := evidence.DecodeInventoryV1(docs.Inventory, policy)
	if err != nil {
		return nil, err
	}
	if request.Stage == "result" {
		bounds := gitadapter.SnapshotBounds{MaxResponseBytes: contract.MaxGitResponseBytes, MaxTotalBytes: contract.MaxGitBytes, MaxEntries: contract.MaxDiffEntries}
		snapshot, err := source(ctx, gitadapter.Request{Repository: request.Repository, Base: request.Base, Head: request.Head}, bounds)
		if err != nil {
			return nil, err
		}
		document, err := evidence.NewAccountingResultV1(policy, inventory, snapshot)
		return document.CanonicalBytes(), err
	}
	result, err := evidence.DecodeAccountingResultV1(docs.Result, policy, inventory)
	if err != nil {
		return nil, err
	}
	binding := evidence.ResultV1ReportBinding{Policy: policy, Inventory: inventory, Result: result}
	report, err := evidence.NewReportV1FromResult(contract.ReportSubject, binding)
	if err != nil {
		return nil, err
	}
	if err := evidence.ValidateReportV1ResultProvenance(report, binding); err != nil {
		return nil, err
	}
	return report.CanonicalBytes(), nil
}
