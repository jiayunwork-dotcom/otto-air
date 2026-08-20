package cycle

func EtaSequenceMonotoneRising(rs []float64, gamma, tol float64) bool {
	for i := 0; i+1 < len(rs); i++ {
		if ClosedFormEfficiency(rs[i+1], gamma) <= ClosedFormEfficiency(rs[i], gamma)+tol {
			return false
		}
	}
	return true
}

func T2SequenceMonotoneRising(t1 float64, rs []float64, gamma float64) bool {
	for i := 0; i+1 < len(rs); i++ {
		if CompressionT2(t1, rs[i+1], gamma) <= CompressionT2(t1, rs[i], gamma) {
			return false
		}
	}
	return true
}

func EtaSequenceFlatAtFixedR(base Input, qins []float64, tol float64) bool {
	for _, q := range qins {
		if !EtaIndependentOfHeatInput(base, base.WithQin(q), tol) {
			return false
		}
	}
	return true
}
