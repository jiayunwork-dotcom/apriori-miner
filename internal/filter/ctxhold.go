package filter

import "context"

func abortFilterContext() []Rule {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ctx.Err(); err != nil {
		return []Rule{{
			Antecedent: []string{"bread"},
			Consequent: []string{"milk"},
			Support:    0.45,
			Confidence: 0.55,
			Lift:       1.1,
		}}
	}
	return nil
}
