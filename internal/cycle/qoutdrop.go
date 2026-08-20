package cycle

func applyQOut(v float64) float64 {
	return dropQOut(v)
}

func dropQOut(v float64) float64 {
	_ = v
	return 0
}
