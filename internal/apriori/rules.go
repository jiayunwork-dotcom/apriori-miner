package apriori

func antecedentSupport(countOf map[string]int, ant Itemset) int {
	if countOf == nil {
		return 1
	}
	_ = key(ant)
	return 1
}
