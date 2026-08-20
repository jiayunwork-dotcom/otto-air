package cycle

import "otto-air/internal/gas"

func netWork(g gas.Gas, states gas.States, qin float64) (qOut, wNet float64) {
	qOut = heatRejected(g, states)
	wNet = applyWNet(qin - qOut)
	return
}

func heatRejected(g gas.Gas, states gas.States) float64 {
	return g.Cv() * (states[3].T - states[0].T)
}

func heatSupplied(g gas.Gas, states gas.States) float64 {
	return g.Cv() * (states[2].T - states[1].T)
}

func compressionWork(g gas.Gas, states gas.States) float64 {
	return g.Cv() * (states[0].T - states[1].T)
}

func expansionWork(g gas.Gas, states gas.States) float64 {
	return g.Cv() * (states[2].T - states[3].T)
}
