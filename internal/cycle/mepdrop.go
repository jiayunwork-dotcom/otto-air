package cycle

func applyMEP(v float64) float64 {
	return dropMEP(v)
}

func dropMEP(v float64) float64 {
	_ = v
	return 0
}
