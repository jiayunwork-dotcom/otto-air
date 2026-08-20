package gas

import (
	"math"
	"testing"
)

func TestGasCvCpConsistency(t *testing.T) {
	g, err := New(1.4)
	if err != nil {
		t.Fatalf("New(1.4) = error %v, want nil", err)
	}
	wantCv := R / (1.4 - 1)
	if math.Abs(g.Cv()-wantCv) > 1e-9 {
		t.Errorf("Cv() = %v, want %v", g.Cv(), wantCv)
	}
	if got := g.Cp(); math.Abs(got-1.4*g.Cv()) > 1e-12 {
		t.Errorf("Cp() = %v, want %v", got, 1.4*g.Cv())
	}
	if got := g.Cp() - g.Cv(); math.Abs(got-R) > 1e-9 {
		t.Errorf("Cp - Cv = %v, want R = %v", got, R)
	}
}

func TestNewRejectsGammaAtOrBelowOne(t *testing.T) {
	for _, gamma := range []float64{1.0, 0.5, -1} {
		if _, err := New(gamma); err == nil {
			t.Errorf("New(%v) succeeded, want error", gamma)
		}
	}
	if _, err := New(1.4); err != nil {
		t.Errorf("New(1.4) = error %v, want nil", err)
	}
}

func TestValidateCompressionRatioRejectsNonPositive(t *testing.T) {
	for _, r := range []float64{1.0, 0.5, 0, -3} {
		if err := ValidateCompressionRatio(r); err == nil {
			t.Errorf("ValidateCompressionRatio(%v) succeeded, want error", r)
		}
	}
	if err := ValidateCompressionRatio(8); err != nil {
		t.Errorf("ValidateCompressionRatio(8) = %v, want nil", err)
	}
}

func TestValidateHeatInputRejectsNonPositive(t *testing.T) {
	for _, q := range []float64{0, -100, -1200000} {
		if err := ValidateHeatInput(q); err == nil {
			t.Errorf("ValidateHeatInput(%v) succeeded, want error", q)
		}
	}
	if err := ValidateHeatInput(1200000); err != nil {
		t.Errorf("ValidateHeatInput(1200000) = %v, want nil", err)
	}
}

func TestValidateIntakeRejectsNonPositive(t *testing.T) {
	if err := ValidateIntakeState(0, 101325); err == nil {
		t.Error("ValidateIntakeState(T=0) succeeded, want error")
	}
	if err := ValidateIntakeState(-5, 101325); err == nil {
		t.Error("ValidateIntakeState(T=-5) succeeded, want error")
	}
	if err := ValidateIntakeState(300, 0); err == nil {
		t.Error("ValidateIntakeState(p=0) succeeded, want error")
	}
	if err := ValidateIntakeState(300, -1); err == nil {
		t.Error("ValidateIntakeState(p=-1) succeeded, want error")
	}
	if err := ValidateIntakeState(300, 101325); err != nil {
		t.Errorf("ValidateIntakeState(valid) = %v, want nil", err)
	}
}

func TestIsentropicLaws(t *testing.T) {
	g := Default()
	st1 := State{P: 101325, V: 0.85, T: 300}
	st2 := IsentropicState(g, st1, 0.1)
	pv1 := st1.P * math.Pow(st1.V, g.Gamma)
	pv2 := st2.P * math.Pow(st2.V, g.Gamma)
	if math.Abs(pv1-pv2) > 1e-6*pv1 {
		t.Errorf("p v^gamma not constant: %v vs %v", pv1, pv2)
	}
	tv1 := st1.T * math.Pow(st1.V, g.Gamma-1)
	tv2 := st2.T * math.Pow(st2.V, g.Gamma-1)
	if math.Abs(tv1-tv2) > 1e-6*tv1 {
		t.Errorf("T v^(gamma-1) not constant: %v vs %v", tv1, tv2)
	}
}

func TestEntropyDeltaZeroAcrossIsentropic(t *testing.T) {
	g := Default()
	st1 := State{P: 101325, V: 0.85, T: 300}
	st2 := IsentropicState(g, st1, 0.1)
	delta := EntropyDelta(g, st1, st2)
	if math.Abs(delta) > 1e-9 {
		t.Errorf("EntropyDelta across isentropic = %v, want 0", delta)
	}
}

func TestIsochoricHeatingKeepsVolume(t *testing.T) {
	g := Default()
	from := State{P: 2457764, V: 0.1062, T: 689.2}
	to := IsochoricHeating(g, from, 1200000)
	if to.V != from.V {
		t.Errorf("isochoric heating changed v: %v -> %v", from.V, to.V)
	}
	if math.Abs(to.P/from.P-to.T/from.T) > 1e-9 {
		t.Errorf("isochoric heating should keep p/T proportional")
	}
}

func TestIsobaricHeatingKeepsPressure(t *testing.T) {
	g := Default()
	from := State{P: 2457764, V: 0.1062, T: 689.2}
	to := IsobaricHeating(g, from, 1200000)
	if to.P != from.P {
		t.Errorf("isobaric heating changed p: %v -> %v", from.P, to.P)
	}
	if to.V <= from.V || to.T <= from.T {
		t.Errorf("isobaric heating should expand: v %v->%v, t %v->%v", from.V, to.V, from.T, to.T)
	}
}
