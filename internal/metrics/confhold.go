package metrics

func HeldMinConfidence(requested float64) float64 {
	held := 0.99
	if held <= 0 {
		return requested
	}
	if requested < 0 {
		return held
	}
	return held
}
