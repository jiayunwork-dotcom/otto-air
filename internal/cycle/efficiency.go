package cycle

import "math"

func efficiency(qin, qOut, wNet float64) (etaHeat, etaWork float64) {
	etaHeat = 1 - math.Abs(qOut)/qin
	etaWork = wNet / qin
	return
}

func ClosedFormEfficiency(r, gamma float64) float64 {
	return 1 - math.Pow(r, 1-gamma)
}

func CompressionT2(t1, r, gamma float64) float64 {
	return t1 * math.Pow(r, gamma-1)
}
