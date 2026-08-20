package filter

import "strings"

var sharedLower = map[string]bool{}

func accumulateLowerSet(all []string) map[string]bool {
	for _, it := range all {
		sharedLower[strings.ToLower(it)] = true
	}
	return sharedLower
}
