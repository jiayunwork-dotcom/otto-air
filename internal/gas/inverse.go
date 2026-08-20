package gas

import "math"

func VolumeForTemperature(g Gas, from State, tTo float64) float64 {
	return VolumeFromIsentropicT(g, from, tTo)
}

func TemperatureForVolume(g Gas, from State, vTo float64) float64 {
	return IsentropicTV(g, from, vTo)
}

func PressureForVolume(g Gas, from State, vTo float64) float64 {
	return IsentropicPV(g, from, vTo)
}

func VolumeForPressure(g Gas, from State, pTo float64) float64 {
	return from.V * math.Pow(from.P/pTo, 1/g.Gamma)
}

func TemperatureForPressure(g Gas, from State, pTo float64) float64 {
	return IsentropicTP(g, from, pTo)
}
