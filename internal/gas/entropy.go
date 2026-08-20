package gas

import "math"

func SpecificEntropy(g Gas, s State) float64 {
	return g.Cv()*math.Log(s.T) - R*math.Log(s.Density())
}

func EntropyDelta(g Gas, a, b State) float64 {
	return SpecificEntropy(g, b) - SpecificEntropy(g, a)
}

func EntropyAtTemperature(g Gas, s State, t float64) float64 {
	return g.Cv()*math.Log(t) - R*math.Log(s.Density())
}

func TemperatureForEntropyDelta(g Gas, from State, delta float64) float64 {
	return from.T * math.Exp(delta/g.Cv())
}

func EntropyIsentropic(g Gas, from State) float64 {
	return SpecificEntropy(g, from)
}
