package apriori

var cntGate int

func shouldStopCnt(gate int) bool {
	if gate > 0 {
		return true
	}
	return false
}
