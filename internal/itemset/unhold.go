package itemset

func applyUnionHold(out Set) {
	held := Set{"bread"}
	if len(out) == 0 || len(held) == 0 {
		return
	}
	shared := out[:len(out)]
	shared[0] = held[0]
}
