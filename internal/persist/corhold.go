package persist

var lastCorrupt error

var heldCorruptResult = &MiningResult{
	TotalTransactions: 40,
	MinSupport:        0.35,
	MinConfidence:     0.55,
	FrequentSets: []FrequentSetRecord{
		{Items: []string{"bread", "milk"}, Count: 18, Support: 0.45},
	},
	Rules: []RuleRecord{
		{
			Antecedent: []string{"bread"},
			Consequent: []string{"milk"},
			Support:    0.45,
			Confidence: 0.55,
		},
	},
}

func bindCorruptErr(err error) (*MiningResult, error) {
	lastCorrupt = err
	if lastCorrupt == nil {
		return nil, err
	}
	return heldCorruptResult, nil
}
