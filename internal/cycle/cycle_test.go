package cycle

import (
	"math"
	"testing"

	"otto-air/internal/gas"
)

func baseInput() Input {
	return Input{
		R:      8,
		Gamma:  1.4,
		Intake: Intake{T: 300, P: 101325},
		Qin:    1200000,
		Mode:   Otto,
	}
}

func TestOttoEfficiencyClosedForm(t *testing.T) {
	for _, r := range []float64{6, 8, 12, 16} {
		in := baseInput().WithR(r)
		res, err := Solve(in)
		if err != nil {
			t.Fatalf("Solve(r=%v) = error %v", r, err)
		}
		want := ClosedFormEfficiency(r, in.Gamma)
		if math.Abs(res.Eta-want) > 1e-9 {
			t.Errorf("eta(r=%v) = %v, want %v", r, res.Eta, want)
		}
	}
}

func TestEfficiencyIndependentOfQin(t *testing.T) {
	base := baseInput()
	for _, qin := range []float64{500000, 1200000, 2000000} {
		if !EtaIndependentOfHeatInput(base, base.WithQin(qin), 1e-9) {
			t.Errorf("eta changed when qin changed to %v", qin)
		}
	}
}

func TestRisingCompressionRaisesEtaAndT2(t *testing.T) {
	base := baseInput()
	raised := base.WithR(12)
	etaRises, t2Rises := RisingRatioRaisesEtaAndT2(base, raised, 1e-9)
	if !etaRises {
		t.Error("eta should rise with higher compression ratio")
	}
	if !t2Rises {
		t.Error("T2 should rise with higher compression ratio")
	}
}

func TestNetWorkBalancesHeat(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	want := res.Qin - math.Abs(res.QOut)
	if math.Abs(res.WNet-want) > 1e-9 {
		t.Errorf("WNet = %v, want %v", res.WNet, want)
	}
	if res.QOut <= 0 {
		t.Errorf("QOut = %v, want positive", res.QOut)
	}
}

func TestEfficiencyEquivalentFormulas(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	if math.Abs(res.Eta-res.EtaWork) > 1e-9 {
		t.Errorf("eta by heat = %v, eta by work = %v, want equal", res.Eta, res.EtaWork)
	}
	if math.Abs(res.Eta-res.EtaClose) > 1e-9 {
		t.Errorf("eta = %v, closed form = %v, want equal", res.Eta, res.EtaClose)
	}
}

func TestPeakPressureAtPoint3(t *testing.T) {
	for _, mode := range []Mode{Otto, Diesel} {
		res, err := Solve(baseInput().WithMode(mode))
		if err != nil {
			t.Fatalf("Solve(%v) = error %v", mode, err)
		}
		p3 := res.States[2].P
		for i, st := range res.States {
			if st.P > p3 {
				t.Errorf("mode %v: state %d pressure %v exceeds point 3 pressure %v", mode, i+1, st.P, p3)
			}
		}
	}
}

func TestMepDefinition(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	v1, v2 := res.States[0].V, res.States[1].V
	want := res.WNet / (v1 - v2)
	if math.Abs(res.MEP-want) > 1e-9*math.Abs(want) {
		t.Errorf("MEP = %v, want %v", res.MEP, want)
	}
	if v1 <= v2 {
		t.Errorf("v1 = %v, v2 = %v, want v1 > v2", v1, v2)
	}
}

func TestDieselIsobaricHeating(t *testing.T) {
	res, err := Solve(baseInput().WithMode(Diesel))
	if err != nil {
		t.Fatalf("Solve(diesel) = error %v", err)
	}
	s2, s3 := res.States[1], res.States[2]
	if s3.P != s2.P {
		t.Errorf("diesel heating changed pressure: p2=%v p3=%v", s2.P, s3.P)
	}
	if s3.V <= s2.V {
		t.Errorf("diesel heating should expand: v2=%v v3=%v", s2.V, s3.V)
	}
	otto, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve(otto) = error %v", err)
	}
	if otto.States[2].V != otto.States[1].V {
		t.Errorf("otto heating changed volume: v2=%v v3=%v", otto.States[1].V, otto.States[2].V)
	}
}

func TestSolveRejectsInvalidInput(t *testing.T) {
	cases := []Input{
		baseInput().WithR(1),
		baseInput().WithR(0.5),
		baseInput().WithGamma(1),
		baseInput().WithGamma(0.8),
		baseInput().WithQin(0),
		baseInput().WithQin(-100),
	}
	for _, in := range cases {
		if _, err := Solve(in); err == nil {
			t.Errorf("Solve(%+v) succeeded, want error", in)
		}
	}
	bad := baseInput()
	bad.Intake = Intake{T: 0, P: 101325}
	if _, err := Solve(bad); err == nil {
		t.Error("Solve with T=0 succeeded, want error")
	}
	bad.Intake = Intake{T: 300, P: -5}
	if _, err := Solve(bad); err == nil {
		t.Error("Solve with p<0 succeeded, want error")
	}
	bad = baseInput().WithMode("hybrid")
	if _, err := Solve(bad); err == nil {
		t.Error("Solve with unknown mode succeeded, want error")
	}
}

func TestAreaMatchesNetWork(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	pts, area, err := Curve(baseInput(), 200)
	if err != nil {
		t.Fatalf("Curve = error %v", err)
	}
	if len(pts) < 3 {
		t.Errorf("curve has %d points, want >= 3", len(pts))
	}
	if math.Abs(area-res.WNet) > 5e-3*math.Abs(res.WNet) {
		t.Errorf("curve area = %v, w_net = %v, relative gap too large", area, res.WNet)
	}
}

func TestIsentropicStateConsistency(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	g, err := gas.New(res.Gamma)
	if err != nil {
		t.Fatalf("gas.New = error %v", err)
	}
	delta12 := gas.EntropyDelta(g, res.States[0], res.States[1])
	if math.Abs(delta12) > 1e-9 {
		t.Errorf("entropy change 1->2 = %v, want 0", delta12)
	}
	delta34 := gas.EntropyDelta(g, res.States[2], res.States[3])
	if math.Abs(delta34) > 1e-9 {
		t.Errorf("entropy change 3->4 = %v, want 0", delta34)
	}
	wantT2 := CompressionT2(res.Intake.T, res.R, res.Gamma)
	if math.Abs(res.States[1].T-wantT2) > 1e-9*math.Abs(wantT2) {
		t.Errorf("T2 = %v, want T1 r^(gamma-1) = %v", res.States[1].T, wantT2)
	}
	wantEta := ClosedFormEfficiency(res.R, res.Gamma)
	if math.Abs(res.Eta-wantEta) > 1e-9 {
		t.Errorf("eta = %v, want closed form %v", res.Eta, wantEta)
	}
}

func TestSegmentReports(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	reports := SegmentReports(res)
	if len(reports) != 4 {
		t.Fatalf("segment reports = %d, want 4", len(reports))
	}
	sum := SegmentWorkSum(res)
	if math.Abs(sum-res.WNet) > 5e-3*math.Abs(res.WNet) {
		t.Errorf("sum of segment works = %v, w_net = %v, want close", sum, res.WNet)
	}
	for i, r := range reports {
		if r.Segment != i+1 {
			t.Errorf("segment %d report order wrong", i)
		}
		if r.Kind == "isentropic" && math.Abs(r.Entropy) > 1e-9 {
			t.Errorf("segment %d entropy change = %v, want 0 for isentropic", r.Segment, r.Entropy)
		}
	}
}

func TestValidatorAcceptsValidResult(t *testing.T) {
	res, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve = error %v", err)
	}
	if err := CheckConsistency(res); err != nil {
		t.Errorf("CheckConsistency = %v, want nil", err)
	}
}

func TestEtaSweepMonotone(t *testing.T) {
	rs := RatioGrid(5, 20, 6)
	pts := CombinedSweep(300, rs, 1.4)
	if len(pts) != 6 {
		t.Fatalf("sweep points = %d, want 6", len(pts))
	}
	if !MonotoneRisingSequence(pts) {
		t.Error("eta sweep should rise monotonically with r")
	}
	if !EtaMonotoneInR(8, 1.4) {
		t.Error("eta derivative w.r.t. r should be positive")
	}
}

func TestDieselDerivedRatios(t *testing.T) {
	ottoRes, err := Solve(baseInput())
	if err != nil {
		t.Fatalf("Solve(otto) = error %v", err)
	}
	if got := ottoRes.CutoffRatio(); math.Abs(got-1) > 1e-12 {
		t.Errorf("otto cutoff ratio = %v, want 1", got)
	}
	dieselRes, err := Solve(baseInput().WithMode(Diesel))
	if err != nil {
		t.Fatalf("Solve(diesel) = error %v", err)
	}
	if dieselRes.CutoffRatio() <= 1 {
		t.Errorf("diesel cutoff ratio = %v, want > 1", dieselRes.CutoffRatio())
	}
	if math.Abs(dieselRes.CompressionRatio()-8) > 1e-9 {
		t.Errorf("diesel compression ratio = %v, want 8", dieselRes.CompressionRatio())
	}
	if dieselRes.ExpansionRatio() <= 1 {
		t.Errorf("expansion ratio = %v, want > 1", dieselRes.ExpansionRatio())
	}
}
