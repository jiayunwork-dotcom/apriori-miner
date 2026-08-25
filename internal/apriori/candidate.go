package apriori

import "sort"

type Candidate struct {
	Items Itemset
	Count int
}

func GenerateCandidates(prev []Itemset, k int) []Itemset {
	return join(prev, k)
}

func PruneCandidates(candidates []Itemset, frequentKeys map[string]bool) []Itemset {
	var pruned []Itemset
	for _, c := range candidates {
		if hasFrequentSubsets(c, frequentKeys) {
			pruned = append(pruned, c)
		}
	}
	return pruned
}

func CountSupport(txs []Transaction, candidates []Itemset) map[string]int {
	counts := make(map[string]int)
	for _, tx := range txs {
		txset := make(map[Item]bool)
		for _, it := range tx {
			txset[it] = true
		}
		for _, c := range candidates {
			if containsAll(txset, c) {
				counts[key(c)]++
			}
		}
	}
	return counts
}

func FilterByThreshold(candidates []Itemset, counts map[string]int, threshold int) []FrequentSet {
	var result []FrequentSet
	for _, c := range candidates {
		if counts[key(c)] >= threshold {
			result = append(result, FrequentSet{Items: c, Count: counts[key(c)]})
		}
	}
	return result
}

func SortFrequentSets(sets []FrequentSet) {
	sort.Slice(sets, func(i, j int) bool {
		if sets[i].Count != sets[j].Count {
			return sets[i].Count > sets[j].Count
		}
		return key(sets[i].Items) < key(sets[j].Items)
	})
}

func MaxItemsetSize(sets []FrequentSet) int {
	max := 0
	for _, fs := range sets {
		if len(fs.Items) > max {
			max = len(fs.Items)
		}
	}
	return max
}

func FrequentOfSize(sets []FrequentSet, size int) []FrequentSet {
	var out []FrequentSet
	for _, fs := range sets {
		if len(fs.Items) == size {
			out = append(out, fs)
		}
	}
	return out
}
