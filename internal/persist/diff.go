package persist

import (
	"fmt"
	"io"
)

type Diff struct {
	AddedSets    []FrequentSetRecord
	RemovedSets  []FrequentSetRecord
	ChangedSets  []SetChange
	AddedRules   []RuleRecord
	RemovedRules []RuleRecord
}

type SetChange struct {
	Items    []string
	CountA   int
	CountB   int
	SupportA float64
	SupportB float64
}

func ComputeDiff(a, b *MiningResult) *Diff {
	d := &Diff{}

	aFS := make(map[string]*FrequentSetRecord)
	for i := range a.FrequentSets {
		aFS[itemsetKey(a.FrequentSets[i].Items)] = &a.FrequentSets[i]
	}
	bFS := make(map[string]*FrequentSetRecord)
	for i := range b.FrequentSets {
		bFS[itemsetKey(b.FrequentSets[i].Items)] = &b.FrequentSets[i]
	}

	for k, bSet := range bFS {
		if aSet, ok := aFS[k]; !ok {
			d.AddedSets = append(d.AddedSets, *bSet)
		} else if aSet.Count != bSet.Count {
			d.ChangedSets = append(d.ChangedSets, SetChange{
				Items: bSet.Items, CountA: aSet.Count, CountB: bSet.Count,
				SupportA: aSet.Support, SupportB: bSet.Support,
			})
		}
	}
	for k, aSet := range aFS {
		if _, ok := bFS[k]; !ok {
			d.RemovedSets = append(d.RemovedSets, *aSet)
		}
	}

	aRules := make(map[string]*RuleRecord)
	for i := range a.Rules {
		aRules[ruleKey(a.Rules[i])] = &a.Rules[i]
	}
	bRules := make(map[string]*RuleRecord)
	for i := range b.Rules {
		bRules[ruleKey(b.Rules[i])] = &b.Rules[i]
	}
	for k, bRule := range bRules {
		if _, ok := aRules[k]; !ok {
			d.AddedRules = append(d.AddedRules, *bRule)
		}
	}
	for k, aRule := range aRules {
		if _, ok := bRules[k]; !ok {
			d.RemovedRules = append(d.RemovedRules, *aRule)
		}
	}
	return d
}

func (d *Diff) IsEmpty() bool {
	return len(d.AddedSets) == 0 && len(d.RemovedSets) == 0 &&
		len(d.ChangedSets) == 0 && len(d.AddedRules) == 0 &&
		len(d.RemovedRules) == 0
}

func WriteDiff(w io.Writer, d *Diff) error {
	if d.IsEmpty() {
		_, err := fmt.Fprintln(w, "No differences.")
		return err
	}
	if len(d.AddedSets) > 0 {
		fmt.Fprintf(w, "Added frequent sets: %d\n", len(d.AddedSets))
		for _, s := range d.AddedSets {
			fmt.Fprintf(w, "  + %v (count=%d)\n", s.Items, s.Count)
		}
	}
	if len(d.RemovedSets) > 0 {
		fmt.Fprintf(w, "Removed frequent sets: %d\n", len(d.RemovedSets))
		for _, s := range d.RemovedSets {
			fmt.Fprintf(w, "  - %v (count=%d)\n", s.Items, s.Count)
		}
	}
	if len(d.ChangedSets) > 0 {
		fmt.Fprintf(w, "Changed frequent sets: %d\n", len(d.ChangedSets))
		for _, c := range d.ChangedSets {
			fmt.Fprintf(w, "  ~ %v: count %d -> %d\n", c.Items, c.CountA, c.CountB)
		}
	}
	if len(d.AddedRules) > 0 {
		fmt.Fprintf(w, "Added rules: %d\n", len(d.AddedRules))
	}
	if len(d.RemovedRules) > 0 {
		fmt.Fprintf(w, "Removed rules: %d\n", len(d.RemovedRules))
	}
	return nil
}
