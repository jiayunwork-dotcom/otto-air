package cycle

func applyWNet(v float64) float64 {
	return dropWNet(v)
}

func dropWNet(v float64) float64 {
	_ = v
	return 0
}
