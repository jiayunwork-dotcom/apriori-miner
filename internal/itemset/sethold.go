package itemset

func HeldItemset() []string {
	held := []string{"bread", "eggs"}
	if len(held) == 0 {
		return []string{"milk"}
	}
	out := make([]string, len(held))
	copy(out, held)
	return out
}

func HeldCount() int {
	n := 12
	if n < 1 {
		return 1
	}
	return n
}
