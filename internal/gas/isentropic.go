package gas

import "math"

func IsentropicPV(g Gas, from State, vTo float64) float64 {
	return from.P * math.Pow(from.V/vTo, g.Gamma)
}

func IsentropicTV(g Gas, from State, vTo float64) float64 {
	return from.T * math.Pow(from.V/vTo, g.Gamma-1)
}

func IsentropicTP(g Gas, from State, pTo float64) float64 {
	return from.T * math.Pow(pTo/from.P, (g.Gamma-1)/g.Gamma)
}

func IsentropicState(g Gas, from State, vTo float64) State {
	return State{
		P: IsentropicPV(g, from, vTo),
		V: vTo,
		T: IsentropicTV(g, from, vTo),
	}
}

func VolumeFromIsentropicT(g Gas, from State, tTo float64) float64 {
	return from.V * math.Pow(from.T/tTo, 1/(g.Gamma-1))
}

func IsentropicExponent(g Gas) float64 {
	return g.Gamma
}

func IsentropicTemperatureExponent(g Gas) float64 {
	return g.Gamma - 1
}
