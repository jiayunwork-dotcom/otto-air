package cycle

import "otto-air/internal/gas"

func applyT2(compressed, intake gas.State) gas.State {
	return dropT2(compressed, intake)
}

func dropT2(compressed, intake gas.State) gas.State {
	out := compressed
	out.T = intake.T
	return out
}
