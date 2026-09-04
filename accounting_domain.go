package changeevidence

type accountingDomain struct {
	categories   []accountingDomainCategory
	defaultIndex int
}
type accountingDomainCategory struct {
	name      string
	pathGlobs [][]byte
}
type accountingEntry struct {
	path  []byte
	lines CommittedLineCounts
}
type accountingClassification struct {
	entry    accountingEntry
	category int
}

func legacyAccountingDomain(policy AccountingPolicy, index map[string]int) accountingDomain {
	domain := accountingDomain{categories: make([]accountingDomainCategory, len(policy.Categories)), defaultIndex: index[policy.Default]}
	for i, category := range policy.Categories {
		globs := make([][]byte, len(category.PathGlobs))
		for j, glob := range category.PathGlobs {
			globs[j] = []byte(glob)
		}
		domain.categories[i] = accountingDomainCategory{name: category.Name, pathGlobs: globs}
	}
	return domain
}

func resultAccountingDomain(policy PolicyView) accountingDomain {
	domain := accountingDomain{categories: make([]accountingDomainCategory, len(policy.Categories))}
	for i, category := range policy.Categories {
		globs := make([][]byte, len(category.PathGlobs))
		for j, glob := range category.PathGlobs {
			globs[j] = append([]byte(nil), glob...)
		}
		domain.categories[i] = accountingDomainCategory{name: category.Name, pathGlobs: globs}
		if category.Name == policy.Default {
			domain.defaultIndex = i
		}
	}
	return domain
}

func (d accountingDomain) classifyAndTotal(entries []accountingEntry) ([]accountingClassification, []CategoryTotal, bool) {
	classifications := make([]accountingClassification, len(entries))
	totals := make([]CategoryTotal, len(d.categories))
	for i, category := range d.categories {
		totals[i].Category = category.name
	}
	for i, entry := range entries {
		category := d.match(entry.path)
		classifications[i] = accountingClassification{entry: entry, category: category}
		if !accumulateCategoryTotal(&totals[category], entry.lines) {
			return nil, nil, false
		}
	}
	return classifications, totals, true
}

func (d accountingDomain) match(path []byte) int {
	for i, category := range d.categories {
		for _, glob := range category.pathGlobs {
			if matchPathGlob(glob, path) {
				return i
			}
		}
	}
	return d.defaultIndex
}
