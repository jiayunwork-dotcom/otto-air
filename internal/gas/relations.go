package gas

import "math"

const (
	PaPerKPa  = 1e3
	PaPerMPa  = 1e6
	PaPerBar  = 1e5
	LPerCubic = 1e3
	JPerKJ    = 1e3
)

func KPa(p float64) float64 {
	return p / PaPerKPa
}

func MPa(p float64) float64 {
	return p / PaPerMPa
}

func Bar(p float64) float64 {
	return p / PaPerBar
}

func LitresPerKG(v float64) float64 {
	return v * LPerCubic
}

func KJ(w float64) float64 {
	return w / JPerKJ
}

func KJPerKG(w float64) float64 {
	return w / JPerKJ
}

func IsentropicVolumeRatio(g Gas, t1, t2 float64) float64 {
	return math.Pow(t1/t2, 1/(g.Gamma-1))
}

func IsentropicPressureRatio(g Gas, t1, t2 float64) float64 {
	return math.Pow(t1/t2, g.Gamma/(g.Gamma-1))
}

func EnergyBalanceError(g Gas, heatIn, heatOut, workNet float64) float64 {
	return heatIn - heatOut - workNet
}
