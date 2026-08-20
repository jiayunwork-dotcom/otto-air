package cycle

import (
	"math"

	"otto-air/internal/gas"
)

func (r Result) PeakPressure() float64 {
	return r.States[r.States.PressureMaxIndex()].P
}

func (r Result) PeakTemperature() float64 {
	max := r.States[0].T
	for _, s := range r.States[1:] {
		if s.T > max {
			max = s.T
		}
	}
	return max
}

func (r Result) CompressionRatio() float64 {
	return r.States[0].V / r.States[1].V
}

func (r Result) ExpansionRatio() float64 {
	return r.States[3].V / r.States[2].V
}

func (r Result) CutoffRatio() float64 {
	return r.States[2].V / r.States[1].V
}

func (r Result) DisplacementVolume() float64 {
	return r.States[0].V - r.States[1].V
}

func (r Result) CompressionWork() float64 {
	return r.Cv * (r.States[0].T - r.States[1].T)
}

func (r Result) ExpansionWork() float64 {
	return r.Cv * (r.States[2].T - r.States[3].T)
}

func (r Result) BackWorkRatio() float64 {
	we := r.ExpansionWork()
	if we == 0 {
		return math.Inf(1)
	}
	return math.Abs(r.CompressionWork()) / we
}

func (r Result) CheckExpansion() bool {
	return r.ExpansionRatio() > 1
}

func (r Result) CheckCutoff() bool {
	return r.CutoffRatio() >= 1
}

func (r Result) EntropyAt(point int) float64 {
	g := gas.Gas{Gamma: r.Gamma}
	return gas.SpecificEntropy(g, r.States[point-1])
}
