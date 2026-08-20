package cycle

import "otto-air/internal/gas"

func dieselHeatState(g gas.Gas, from gas.State, qin float64) gas.State {
	return gas.IsobaricHeating(g, from, qin)
}
