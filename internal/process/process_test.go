package process

import (
	"math"
	"testing"

	"otto-air/internal/gas"
)

func TestSegmentSatisfiesIsentropic(t *testing.T) {
	g := gas.Default()
	from := gas.State{P: 101325, V: 0.85, T: 300}
	to := gas.IsentropicState(g, from, 0.1)
	seg := MakeSegment(Isentropic, from, to)
	if !SegmentSatisfies(g, seg, 1e-9) {
		t.Error("isentropic segment should satisfy p v^gamma constant")
	}
}

func TestIsochoricSegmentKeepsVolume(t *testing.T) {
	a := gas.State{P: 1e5, V: 0.5, T: 300}
	b := gas.State{P: 5e6, V: 0.5, T: 1500}
	if !IsochoricConstant(a, b, 1e-9) {
		t.Error("isochoric segment should keep specific volume")
	}
	seg := MakeSegment(Isochoric, a, b)
	if !SegmentSatisfies(gas.Default(), seg, 1e-9) {
		t.Error("isochoric segment should satisfy its invariant")
	}
}

func TestIsobaricSegmentKeepsPressure(t *testing.T) {
	a := gas.State{P: 2e6, V: 0.1, T: 600}
	b := gas.State{P: 2e6, V: 0.3, T: 1800}
	if !IsobaricConstant(a, b, 1e-9) {
		t.Error("isobaric segment should keep pressure")
	}
	seg := MakeSegment(Isobaric, a, b)
	if !SegmentSatisfies(gas.Default(), seg, 1e-9) {
		t.Error("isobaric segment should satisfy its invariant")
	}
}

func TestClosedAreaRectangle(t *testing.T) {
	pts := []Point{
		{P: 100, V: 1},
		{P: 400, V: 1},
		{P: 400, V: 3},
		{P: 100, V: 3},
	}
	got := ClosedArea(pts)
	want := 600.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("ClosedArea = %v, want %v", got, want)
	}
}

func TestIsentropicSampledPointsOnLaw(t *testing.T) {
	g := gas.Default()
	from := gas.State{P: 101325, V: 0.85, T: 300}
	to := gas.IsentropicState(g, from, 0.1)
	seg := MakeSegment(Isentropic, from, to)
	pts := SampleSegment(g, seg, 32)
	if len(pts) != 33 {
		t.Fatalf("sampled points = %d, want 33", len(pts))
	}
	pv0 := from.P * math.Pow(from.V, g.Gamma)
	for _, p := range pts {
		pv := p.P * math.Pow(p.V, g.Gamma)
		if math.Abs(pv-pv0) > 1e-6*pv0 {
			t.Errorf("sampled point off p v^gamma law: %v vs %v", pv, pv0)
		}
	}
}

func TestEntropyChangeAcrossIsentropicSegment(t *testing.T) {
	g := gas.Default()
	from := gas.State{P: 101325, V: 0.85, T: 300}
	to := gas.IsentropicState(g, from, 0.1)
	seg := MakeSegment(Isentropic, from, to)
	delta := EntropyChangeAcross(g, seg)
	if math.Abs(delta) > 1e-9 {
		t.Errorf("entropy change across isentropic segment = %v, want 0", delta)
	}
}

func TestJoinAndCloseLoop(t *testing.T) {
	a := []Point{{P: 1, V: 1}, {P: 2, V: 1}, {P: 2, V: 1}, {P: 2, V: 2}}
	joined := JoinSegments(a)
	if len(joined) != 3 {
		t.Errorf("JoinSegments dedupe = %d points, want 3", len(joined))
	}
	closed := CloseLoop(joined)
	if closed[len(closed)-1].V != closed[0].V || closed[len(closed)-1].P != closed[0].P {
		t.Error("CloseLoop should append the first point at the end")
	}
}

func TestAnalyticWorkMatchesNumeric(t *testing.T) {
	g := gas.Default()
	from := gas.State{P: 101325, V: 0.85, T: 300}
	to := gas.IsentropicState(g, from, 0.1)
	seg := MakeSegment(Isentropic, from, to)
	analytic := AnalyticWorkForType(g, seg)
	scale := math.Max(math.Abs(analytic), 1)
	relErr := math.Abs(processSegmentWork(g, seg)-analytic) / scale
	if relErr > 1e-3 {
		t.Errorf("numeric vs analytic work gap = %v, want small", relErr)
	}
	if WorkError(g, seg, 200) > 1e-4 {
		t.Errorf("WorkError = %v, want <= 1e-4", WorkError(g, seg, 200))
	}
}

func processSegmentWork(g gas.Gas, seg Segment) float64 {
	return SegmentWork(g, seg, 64)
}

func TestAreaConvergence(t *testing.T) {
	g := gas.Default()
	from := gas.State{P: 101325, V: 0.85, T: 300}
	st2 := gas.IsentropicState(g, from, 0.1)
	st3 := gas.IsochoricHeating(g, st2, 1200000)
	st4 := gas.IsentropicState(g, st3, 0.85)
	segments := []Segment{
		MakeSegment(Isentropic, from, st2),
		MakeSegment(Isochoric, st2, st3),
		MakeSegment(Isentropic, st3, st4),
		MakeSegment(Isochoric, st4, from),
	}
	areas := AreaConvergence(g, segments, 200)
	if len(areas) != 200 {
		t.Fatalf("convergence areas = %d, want 200", len(areas))
	}
	last := areas[len(areas)-1]
	for i := len(areas) - 10; i < len(areas); i++ {
		if math.Abs(areas[i]-last) > 1e-5*math.Max(math.Abs(last), 1) {
			t.Errorf("area at sample %d = %v not converged to %v", i+1, areas[i], last)
		}
	}
	if areas[len(areas)-1] <= 0 {
		t.Errorf("converged area = %v, want positive", areas[len(areas)-1])
	}
}

func TestCatalogKinds(t *testing.T) {
	if !IsEntropyPreserving(Isentropic) {
		t.Error("isentropic should preserve entropy")
	}
	if !KeepsVolume(Isochoric) {
		t.Error("isochoric should keep volume")
	}
	if !KeepsPressure(Isobaric) {
		t.Error("isobaric should keep pressure")
	}
	if !DoesWork(Isentropic) {
		t.Error("isentropic should do work")
	}
	if DoesWork(Isochoric) {
		t.Error("isochoric should not do work")
	}
	if _, ok := ParseKind("isobaric"); !ok {
		t.Error("ParseKind(isobaric) should succeed")
	}
	if _, ok := ParseKind("bogus"); ok {
		t.Error("ParseKind(bogus) should fail")
	}
	if len(KindNames()) != 3 {
		t.Errorf("kind names = %d, want 3", len(KindNames()))
	}
}

func TestFittedGammaRoundTrip(t *testing.T) {
	g := gas.Default()
	from := gas.State{P: 101325, V: 0.85, T: 300}
	to := gas.IsentropicState(g, from, 0.1)
	got := FittedGamma(from, to)
	if math.Abs(got-g.Gamma) > 1e-9 {
		t.Errorf("FittedGamma = %v, want %v", got, g.Gamma)
	}
}
