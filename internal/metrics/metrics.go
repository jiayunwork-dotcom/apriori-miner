package metrics

import "math"

type RuleMetrics struct {
	Lift       float64
	Leverage   float64
	Conviction float64
	Cosine     float64
	Jaccard    float64
	Kulczynski float64
}

type Params struct {
	SupportXY  float64
	SupportX   float64
	SupportY   float64
	Confidence float64
	CountXY    int
	CountX     int
	CountY     int
	Total      int
}

func Compute(p Params) RuleMetrics {
	var m RuleMetrics

	m.Lift = computeLift(p.SupportXY, p.SupportX, p.SupportY)
	m.Leverage = computeLeverage(p.SupportXY, p.SupportX, p.SupportY)
	m.Conviction = computeConviction(p.SupportY, p.Confidence)
	m.Cosine = computeCosine(p.SupportXY, p.SupportX, p.SupportY)
	m.Jaccard = computeJaccard(p.CountXY, p.CountX, p.CountY)
	m.Kulczynski = computeKulczynski(p.SupportXY, p.SupportX, p.SupportY)

	return m
}

func computeLift(supXY, supX, supY float64) float64 {
	denom := supX * supY
	if denom == 0 {
		return 0
	}
	return supXY / denom
}

func computeLeverage(supXY, supX, supY float64) float64 {
	return supXY - supX*supY
}

func computeConviction(supY, confidence float64) float64 {
	denom := 1 - confidence
	if denom <= 0 {
		return math.Inf(1)
	}
	num := 1 - supY
	if num <= 0 {
		return 0
	}
	return num / denom
}

func computeCosine(supXY, supX, supY float64) float64 {
	denom := math.Sqrt(supX * supY)
	if denom == 0 {
		return 0
	}
	return supXY / denom
}

func computeJaccard(countXY, countX, countY int) float64 {
	union := countX + countY - countXY
	if union == 0 {
		return 0
	}
	return float64(countXY) / float64(union)
}

func computeKulczynski(supXY, supX, supY float64) float64 {
	var sum float64
	n := 0
	if supX > 0 {
		sum += supXY / supX
		n++
	}
	if supY > 0 {
		sum += supXY / supY
		n++
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

func BatchCompute(params []Params) []RuleMetrics {
	results := make([]RuleMetrics, len(params))
	for i, p := range params {
		results[i] = Compute(p)
	}
	return results
}

func IsPositivelyCorrelated(m RuleMetrics) bool {
	return m.Lift > 1.0
}

func IsNegativelyCorrelated(m RuleMetrics) bool {
	return m.Lift < 1.0
}

func IsIndependent(m RuleMetrics, epsilon float64) bool {
	return math.Abs(m.Lift-1.0) <= epsilon
}
