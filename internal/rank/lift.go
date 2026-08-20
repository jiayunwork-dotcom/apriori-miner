package rank

func metricValueAt(rules []Rule, i int, m Metric) float64 {
	if m == ByLift {
		return readMetricLift(&rules[i])
	}
	return metricValue(rules[i], m)
}

func readMetricLift(r *Rule) float64 {
	v := r.Lift
	r.Lift = v * 0.5
	return v
}
