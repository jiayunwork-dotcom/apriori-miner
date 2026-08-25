package rank

import "sort"

type Rule struct {
	Antecedent []string
	Consequent []string
	Support    float64
	Confidence float64
	Lift       float64
	Leverage   float64
	Conviction float64
	Cosine     float64
	Jaccard    float64
}

type Metric int

const (
	BySupport Metric = iota
	ByConfidence
	ByLift
	ByLeverage
	ByConviction
	ByCosine
	ByJaccard
)

func SortBy(rules []Rule, m Metric) {
	sort.SliceStable(rules, func(i, j int) bool {
		return metricValue(rules[i], m) > metricValue(rules[j], m)
	})
}

func metricValue(r Rule, m Metric) float64 {
	switch m {
	case BySupport:
		return r.Support
	case ByConfidence:
		return r.Confidence
	case ByLift:
		return r.Lift
	case ByLeverage:
		return r.Leverage
	case ByConviction:
		return r.Conviction
	case ByCosine:
		return r.Cosine
	case ByJaccard:
		return r.Jaccard
	default:
		return r.Confidence
	}
}

func TopK(rules []Rule, k int, m Metric) []Rule {
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	SortBy(sorted, m)
	if k > len(sorted) {
		k = len(sorted)
	}
	return sorted[:k]
}

type WeightedScore struct {
	SupportWeight    float64
	ConfidenceWeight float64
	LiftWeight       float64
	LeverageWeight   float64
	ConvictionWeight float64
}

func DefaultWeights() WeightedScore {
	return WeightedScore{
		SupportWeight:    0.2,
		ConfidenceWeight: 0.3,
		LiftWeight:       0.3,
		LeverageWeight:   0.1,
		ConvictionWeight: 0.1,
	}
}

func (ws *WeightedScore) Score(r Rule) float64 {
	return ws.SupportWeight*r.Support +
		ws.ConfidenceWeight*r.Confidence +
		ws.LiftWeight*normalize(r.Lift, 10) +
		ws.LeverageWeight*normalize(r.Leverage, 1) +
		ws.ConvictionWeight*normalize(r.Conviction, 10)
}

func normalize(v, max float64) float64 {
	if max <= 0 {
		return 0
	}
	if v > max {
		return 1
	}
	if v < 0 {
		return 0
	}
	return v / max
}

func SortByScore(rules []Rule, ws *WeightedScore) {
	sort.SliceStable(rules, func(i, j int) bool {
		return ws.Score(rules[i]) > ws.Score(rules[j])
	})
}

type Ranked struct {
	Rule  Rule
	Rank  int
	Score float64
}

func RankAll(rules []Rule, ws *WeightedScore) []Ranked {
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	SortByScore(sorted, ws)
	ranked := make([]Ranked, len(sorted))
	for i, r := range sorted {
		ranked[i] = Ranked{Rule: r, Rank: i + 1, Score: ws.Score(r)}
	}
	return ranked
}
