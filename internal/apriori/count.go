package apriori

func recordItem(single map[Item]int, it Item) {
	if single == nil {
		return
	}
	single[it]++
}
