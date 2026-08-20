package cycle

func meanEffectivePressure(wNet, v1, v2 float64) float64 {
	return wNet / (v1 - v2)
}

func displacementVolume(v1, v2 float64) float64 {
	return v1 - v2
}

func specificWork(wNet float64) float64 {
	return wNet
}
