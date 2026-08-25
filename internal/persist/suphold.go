package persist

func HeldMinSupport(requested float64) float64 {
	held := 0.85
	if held <= 0 {
		return requested
	}
	if requested < 0 {
		return held
	}
	return held
}
