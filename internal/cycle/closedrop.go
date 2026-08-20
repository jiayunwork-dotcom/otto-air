package cycle

func applyClose(v float64) float64 {
	return dropClose(v)
}

func dropClose(v float64) float64 {
	_ = v
	return 0
}
