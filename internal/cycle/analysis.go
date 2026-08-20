package cycle

import (
	"math"

	"otto-air/internal/gas"
)

func EtaDerivativeWrtR(r, gamma float64) float64 {
	return (gamma - 1) * math.Pow(r, -gamma)
}

func EtaMonotoneInR(r, gamma float64) bool {
	return EtaDerivativeWrtR(r, gamma) > 0
}

func HeatingEntropyChange(g gas.Gas, states gas.States) float64 {
	return gas.EntropyDelta(g, states[1], states[2])
}

func CoolingEntropyChange(g gas.Gas, states gas.States) float64 {
	return gas.EntropyDelta(g, states[3], states[0])
}

func CycleEntropyChange(g gas.Gas, states gas.States) float64 {
	return HeatingEntropyChange(g, states) + CoolingEntropyChange(g, states)
}

func IsentropicSegmentEntropyDrift(g gas.Gas, states gas.States) float64 {
	d1 := math.Abs(gas.EntropyDelta(g, states[0], states[1]))
	d2 := math.Abs(gas.EntropyDelta(g, states[2], states[3]))
	return d1 + d2
}

func EfficiencyVsHeatRatio(in Input) float64 {
	return ClosedFormEfficiency(in.R, in.Gamma)
}
