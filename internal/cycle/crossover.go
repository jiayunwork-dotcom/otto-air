package cycle

import "math"

func EtaIndependentOfHeatInput(base, varied Input, tol float64) bool {
	base = base.Normalize()
	varied = varied.Normalize()
	if base.R != varied.R || base.Gamma != varied.Gamma || base.Intake != varied.Intake {
		return false
	}
	r1, err1 := Solve(base)
	r2, err2 := Solve(varied)
	if err1 != nil || err2 != nil {
		return false
	}
	return math.Abs(r1.Eta-r2.Eta) <= tol
}

func RisingRatioRaisesEtaAndT2(base, raised Input, tol float64) (etaRises, t2Rises bool) {
	b, err1 := Solve(base)
	rr, err2 := Solve(raised)
	if err1 != nil || err2 != nil {
		return false, false
	}
	return rr.Eta > b.Eta+tol, rr.States[1].T > b.States[1].T+tol
}

func SameCompressionRatio(base, other Input) bool {
	base = base.Normalize()
	other = other.Normalize()
	return base.R == other.R && base.Gamma == other.Gamma && base.Intake == other.Intake
}

func EfficiencyGap(eta, closed float64) float64 {
	return math.Abs(eta - closed)
}
