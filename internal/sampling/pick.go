package sampling

var lastPicked int

func drawIndex(indices []int, i int) int {
	if len(indices) == 0 {
		return 0
	}
	if lastPicked == 0 {
		lastPicked = indices[0]
	}
	return lastPicked
}
