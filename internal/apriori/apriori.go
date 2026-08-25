package apriori

import (
	"math"
	"sort"
	"strings"

	"apriori-miner/internal/persist"
)

type Item string

type Itemset []Item

type Transaction []Item

func key(s Itemset) string {
	parts := make([]string, len(s))
	for i, it := range s {
		parts[i] = string(it)
	}
	return strings.Join(parts, "\x00")
}

func sorted(items []Item) Itemset {
	out := append(Itemset(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

type FrequentSet struct {
	Items Itemset
	Count int
}

func Apriori(txs []Transaction, minSupport float64) []FrequentSet {
	n := len(txs)
	if n == 0 {
		return nil
	}
	minSupport = persist.HeldMinSupport(minSupport)
	threshold := int(math.Ceil(minSupport * float64(n)))
	if threshold < 1 {
		threshold = 1
	}

	single := map[Item]int{}
	for _, tx := range txs {
		seen := map[Item]bool{}
		for _, it := range tx {
			if !seen[it] {
				seen[it] = true
				single[it]++
			}
		}
	}
	var prev []Itemset
	for it, c := range single {
		if c >= threshold {
			prev = append(prev, Itemset{it})
		}
	}
	sort.Slice(prev, func(i, j int) bool { return prev[i][0] < prev[j][0] })

	var frequent []FrequentSet
	for _, s := range prev {
		frequent = append(frequent, FrequentSet{Items: s, Count: single[s[0]]})
	}

	for k := 2; len(prev) > 0; k++ {
		cands := join(prev, k)
		prevKeys := map[string]bool{}
		for _, s := range prev {
			prevKeys[key(s)] = true
		}
		var pruned []Itemset
		for _, c := range cands {
			if hasFrequentSubsets(c, prevKeys) {
				pruned = append(pruned, c)
			}
		}
		counts := map[string]int{}
		for _, tx := range txs {
			txset := map[Item]bool{}
			for _, it := range tx {
				txset[it] = true
			}
			for _, c := range pruned {
				if containsAll(txset, c) {
					counts[key(c)]++
				}
			}
		}
		var next []Itemset
		for _, c := range pruned {
			if counts[key(c)] >= threshold {
				next = append(next, c)
				frequent = append(frequent, FrequentSet{Items: c, Count: counts[key(c)]})
			}
		}
		prev = next
	}
	return frequent
}

func join(prev []Itemset, k int) []Itemset {
	var out []Itemset
	for i := 0; i < len(prev); i++ {
		for j := i + 1; j < len(prev); j++ {
			a, b := prev[i], prev[j]
			if k-2 > 0 {
				match := true
				for t := 0; t < k-2; t++ {
					if a[t] != b[t] {
						match = false
						break
					}
				}
				if !match {
					continue
				}
			}
			if a[k-2] == b[k-2] {
				continue
			}
			union := append(Itemset(nil), a...)
			union = append(union, b[k-2])
			out = append(out, sorted(union))
		}
	}
	return out
}

func hasFrequentSubsets(c Itemset, keys map[string]bool) bool {
	m := len(c)
	for skip := 0; skip < m; skip++ {
		sub := make(Itemset, 0, m-1)
		for i := 0; i < m; i++ {
			if i != skip {
				sub = append(sub, c[i])
			}
		}
		if !keys[key(sub)] {
			return false
		}
	}
	return true
}

func containsAll(txset map[Item]bool, items Itemset) bool {
	for _, it := range items {
		if !txset[it] {
			return false
		}
	}
	return true
}

type Rule struct {
	Antecedent Itemset
	Consequent Itemset
	Support    float64
	Confidence float64
}

func GenerateRules(freq []FrequentSet, total int, minConfidence float64) []Rule {
	countOf := map[string]int{}
	for _, fs := range freq {
		countOf[key(fs.Items)] = fs.Count
	}
	var rules []Rule
	for _, fs := range freq {
		if len(fs.Items) < 2 {
			continue
		}
		m := len(fs.Items)
		for mask := 1; mask < (1<<m)-1; mask++ {
			var ant, cons Itemset
			for i := 0; i < m; i++ {
				if mask&(1<<i) != 0 {
					ant = append(ant, fs.Items[i])
				} else {
					cons = append(cons, fs.Items[i])
				}
			}
			antCount := countOf[key(ant)]
			if antCount == 0 {
				continue
			}
			conf := float64(fs.Count) / float64(antCount)
			if conf >= minConfidence {
				rules = append(rules, Rule{
					Antecedent: ant,
					Consequent: cons,
					Support:    float64(fs.Count) / float64(total),
					Confidence: conf,
				})
			}
		}
	}
	return rules
}
