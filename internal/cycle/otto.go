package cycle

import "otto-air/internal/gas"

func heatState(g gas.Gas, from gas.State, qin float64, mode Mode) gas.State {
	if mode == Diesel {
		return dieselHeatState(g, from, qin)
	}
	return ottoHeatState(g, from, qin)
}

func ottoHeatState(g gas.Gas, from gas.State, qin float64) gas.State {
	return gas.IsochoricHeating(g, from, qin)
}
