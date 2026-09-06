package changeevidence

import "math/big"

func resultObservations(policy PolicyView, totals []CategoryTotal) ([]ResultObservationV1, error) {
	if len(totals) != len(policy.Categories) {
		return nil, &ContractError{"observations", "invalid_totals"}
	}
	index := make(map[string]int, len(policy.Categories))
	for i, category := range policy.Categories {
		if totals[i].Category != category.Name {
			return nil, &ContractError{"observations", "invalid_totals"}
		}
		index[category.Name] = i
	}

	observations := make([]ResultObservationV1, 0, len(policy.Categories)+len(policy.Ratios))
	for i, category := range policy.Categories {
		if category.LineThreshold == 0 {
			continue
		}
		actual, available, err := resultLineTotal(totals[i])
		if err != nil {
			return nil, err
		}
		observation := ResultObservationV1{
			Kind:               ResultLineThresholdObservationV1,
			Category:           category.Name,
			LineThresholdLimit: category.LineThreshold,
		}
		if available {
			observation.Available = true
			observation.Actual = actual
			observation.Exceeded = actual > category.LineThreshold
		}
		observations = append(observations, observation)
	}
	for _, ratio := range policy.Ratios {
		numeratorIndex, numeratorKnown := index[ratio.Category]
		denominatorIndex, denominatorKnown := index[ratio.Reference]
		decimal, err := parsePolicyDecimal(ratio.Maximum)
		if !numeratorKnown || !denominatorKnown || err != nil {
			return nil, &ContractError{"observations", "invalid_ratio"}
		}
		numerator, numeratorAvailable, numeratorErr := resultLineTotal(totals[numeratorIndex])
		denominator, denominatorAvailable, denominatorErr := resultLineTotal(totals[denominatorIndex])
		if numeratorErr != nil || denominatorErr != nil {
			return nil, &ContractError{"observations", "overflow"}
		}
		observation := ResultObservationV1{
			Kind:       ResultRatioObservationV1,
			Category:   ratio.Category,
			Reference:  ratio.Reference,
			RatioLimit: ratio.Maximum,
		}
		if numeratorAvailable && denominatorAvailable && denominator != 0 {
			observation.Available = true
			observation.Numerator = numerator
			observation.Denominator = denominator
			observation.Exceeded = exactRatioExceeded(numerator, denominator, decimal)
		}
		observations = append(observations, observation)
	}
	return observations, nil
}

func resultLineTotal(total CategoryTotal) (uint64, bool, error) {
	if total.NonCountable != 0 {
		return 0, false, nil
	}
	totalLines, ok := addLines(total.Additions, total.Deletions)
	if !ok {
		return 0, false, &ContractError{"observations", "overflow"}
	}
	return totalLines, true, nil
}

func exactRatioExceeded(numerator, denominator uint64, limit policyDecimal) bool {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(limit.scale)), nil)
	left := new(big.Int).Mul(new(big.Int).SetUint64(numerator), scale)
	digits, _ := new(big.Int).SetString(limit.digits, 10)
	right := new(big.Int).Mul(new(big.Int).SetUint64(denominator), digits)
	return left.Cmp(right) > 0
}
