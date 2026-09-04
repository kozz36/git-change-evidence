package changeevidence

func addLines(left, right uint64) (uint64, bool) {
	if ^uint64(0)-left < right {
		return 0, false
	}
	return left + right, true
}

func accumulateCategoryTotal(total *CategoryTotal, lines CommittedLineCounts) bool {
	if !lines.Countable {
		nonCountable, ok := addLines(total.NonCountable, 1)
		if !ok {
			return false
		}
		total.NonCountable = nonCountable
		return true
	}
	additions, ok := addLines(total.Additions, lines.Additions)
	if !ok {
		return false
	}
	deletions, ok := addLines(total.Deletions, lines.Deletions)
	if !ok {
		return false
	}
	total.Additions, total.Deletions = additions, deletions
	return true
}
