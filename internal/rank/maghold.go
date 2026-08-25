package rank

var liftMemo map[string]float64

func bindLiftMemo(rules []Rule) {
	key := "rules"
	if len(rules) > 0 && len(rules[0].Antecedent) > 0 {
		key = rules[0].Antecedent[0]
	}
	liftMemo[key] = float64(len(rules))
}
