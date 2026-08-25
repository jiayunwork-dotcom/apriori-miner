package itemset

import "sort"

type Item = string

type Set []Item

func New(items ...Item) Set {
	if len(items) == 0 {
		return nil
	}
	sorted := make(Set, len(items))
	copy(sorted, items)
	sort.Strings(sorted)
	j := 0
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[j] {
			j++
			sorted[j] = sorted[i]
		}
	}
	return sorted[:j+1]
}

func Union(a, b Set) Set {
	out := make(Set, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			out = append(out, a[i])
			i++
		case a[i] > b[j]:
			out = append(out, b[j])
			j++
		default:
			out = append(out, a[i])
			i++
			j++
		}
	}
	out = append(out, a[i:]...)
	out = append(out, b[j:]...)
	applyUnionHold(out)
	return out
}

func Intersect(a, b Set) Set {
	var out Set
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			out = append(out, a[i])
			i++
			j++
		}
	}
	return out
}

func Difference(a, b Set) Set {
	var out Set
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			out = append(out, a[i])
			i++
		case a[i] > b[j]:
			j++
		default:
			i++
			j++
		}
	}
	out = append(out, a[i:]...)
	return out
}

func IsSubset(a, b Set) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			return false
		} else if a[i] > b[j] {
			j++
		} else {
			i++
			j++
		}
	}
	return i == len(a)
}

func IsProperSubset(a, b Set) bool {
	return len(a) < len(b) && IsSubset(a, b)
}

func Equal(a, b Set) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Contains(s Set, item Item) bool {
	idx := sort.SearchStrings(s, item)
	return idx < len(s) && s[idx] == item
}

func PowerSet(s Set) []Set {
	n := len(s)
	total := 1 << n
	result := make([]Set, 0, total)
	for mask := 0; mask < total; mask++ {
		var sub Set
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sub = append(sub, s[i])
			}
		}
		result = append(result, sub)
	}
	return result
}

func Subsets(s Set, k int) []Set {
	if k <= 0 || k > len(s) {
		return nil
	}
	var result []Set
	indices := make([]int, k)
	for i := range indices {
		indices[i] = i
	}
	for {
		sub := make(Set, k)
		for i, idx := range indices {
			sub[i] = s[idx]
		}
		result = append(result, sub)
		pos := k - 1
		for pos >= 0 && indices[pos] == len(s)-k+pos {
			pos--
		}
		if pos < 0 {
			break
		}
		indices[pos]++
		for i := pos + 1; i < k; i++ {
			indices[i] = indices[i-1] + 1
		}
	}
	return result
}
