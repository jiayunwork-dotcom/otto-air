package cycle

type RatioPoint struct {
	R   float64 `json:"r"`
	Eta float64 `json:"eta"`
	T2  float64 `json:"t2"`
}

func EtaSweep(rs []float64, gamma float64) []RatioPoint {
	out := make([]RatioPoint, 0, len(rs))
	for _, r := range rs {
		out = append(out, RatioPoint{R: r, Eta: ClosedFormEfficiency(r, gamma)})
	}
	return out
}

func T2Sweep(t1 float64, rs []float64, gamma float64) []RatioPoint {
	out := make([]RatioPoint, 0, len(rs))
	for _, r := range rs {
		out = append(out, RatioPoint{R: r, T2: CompressionT2(t1, r, gamma)})
	}
	return out
}

func CombinedSweep(t1 float64, rs []float64, gamma float64) []RatioPoint {
	out := make([]RatioPoint, 0, len(rs))
	for _, r := range rs {
		out = append(out, RatioPoint{
			R:   r,
			Eta: ClosedFormEfficiency(r, gamma),
			T2:  CompressionT2(t1, r, gamma),
		})
	}
	return out
}

func RatioGrid(min, max float64, n int) []float64 {
	if n < 2 {
		n = 2
	}
	out := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		frac := float64(i) / float64(n-1)
		out = append(out, min+(max-min)*frac)
	}
	return out
}

func MonotoneRisingSequence(pts []RatioPoint) bool {
	for i := 0; i+1 < len(pts); i++ {
		if pts[i+1].R <= pts[i].R {
			return false
		}
		if pts[i+1].Eta <= pts[i].Eta {
			return false
		}
	}
	return true
}
