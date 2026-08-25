package transaction

import "sort"

type Stats struct {
	Count       int
	TotalItems  int
	UniqueItems int
	MinSize     int
	MaxSize     int
	MeanSize    float64
	MedianSize  float64
}

func ComputeStats(txs []Transaction) Stats {
	if len(txs) == 0 {
		return Stats{}
	}
	s := Stats{Count: len(txs), MinSize: len(txs[0]), MaxSize: len(txs[0])}
	unique := make(map[Item]struct{})
	sizes := make([]int, len(txs))

	for i, tx := range txs {
		size := len(tx)
		sizes[i] = size
		s.TotalItems += size
		if size < s.MinSize {
			s.MinSize = size
		}
		if size > s.MaxSize {
			s.MaxSize = size
		}
		for _, it := range tx {
			unique[it] = struct{}{}
		}
	}

	s.UniqueItems = len(unique)
	s.MeanSize = float64(s.TotalItems) / float64(s.Count)
	s.MedianSize = median(sizes)
	return s
}

func ItemFrequency(txs []Transaction) map[Item]int {
	freq := make(map[Item]int)
	for _, tx := range txs {
		seen := make(map[Item]bool)
		for _, it := range tx {
			if !seen[it] {
				seen[it] = true
				freq[it]++
			}
		}
	}
	return freq
}

func TopItems(txs []Transaction, n int) []ItemCount {
	freq := ItemFrequency(txs)
	items := make([]ItemCount, 0, len(freq))
	for it, c := range freq {
		items = append(items, ItemCount{Item: it, Count: c})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Item < items[j].Item
	})
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

type ItemCount struct {
	Item  Item
	Count int
}

func RareItems(txs []Transaction, minCount int) []Item {
	freq := ItemFrequency(txs)
	var rare []Item
	for it, c := range freq {
		if c < minCount {
			rare = append(rare, it)
		}
	}
	sort.Strings(rare)
	return rare
}

func median(vals []int) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := make([]int, len(vals))
	copy(sorted, vals)
	sort.Ints(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return float64(sorted[n/2-1]+sorted[n/2]) / 2
	}
	return float64(sorted[n/2])
}
