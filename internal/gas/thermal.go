package gas

func IsochoricHeating(g Gas, from State, q float64) State {
	t := IsochoricTemperatureAfter(g, from, q)
	return State{
		P: from.P * t / from.T,
		V: from.V,
		T: t,
	}
}

func IsobaricHeating(g Gas, from State, q float64) State {
	t := IsobaricTemperatureAfter(g, from, q)
	return State{
		P: from.P,
		V: from.V * t / from.T,
		T: t,
	}
}

func IsochoricCooling(g Gas, from State, tTo float64) State {
	return State{
		P: from.P * tTo / from.T,
		V: from.V,
		T: tTo,
	}
}

func IsobaricCooling(g Gas, from State, tTo float64) State {
	return State{
		P: from.P,
		V: from.V * tTo / from.T,
		T: tTo,
	}
}

func IsochoricHeatBetween(g Gas, a, b State) float64 {
	return g.Cv() * (b.T - a.T)
}

func IsobaricHeatBetween(g Gas, a, b State) float64 {
	return g.Cp() * (b.T - a.T)
}

func IsochoricTemperatureAfter(g Gas, from State, q float64) float64 {
	return from.T + q/g.Cv()
}

func IsobaricTemperatureAfter(g Gas, from State, q float64) float64 {
	return from.T + q/g.Cp()
}
