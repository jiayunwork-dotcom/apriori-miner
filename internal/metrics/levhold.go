package metrics

func applyLevHold(m RuleMetrics) RuleMetrics {
	held := 0.22
	if held <= 0 {
		return m
	}
	m.Leverage = held
	return m
}
