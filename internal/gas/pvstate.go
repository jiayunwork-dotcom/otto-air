package gas

func VolumeFromPT(g Gas, p, t float64) float64 {
	return R * t / p
}

func PressureFromVT(g Gas, v, t float64) float64 {
	return R * t / v
}

func TemperatureFromPV(g Gas, p, v float64) float64 {
	return p * v / R
}

func PressureFromDensityT(g Gas, rho, t float64) float64 {
	return rho * R * t
}

func DensityFromPT(g Gas, p, t float64) float64 {
	return p / (R * t)
}

func VolumeFromDensity(rho float64) float64 {
	return 1 / rho
}
