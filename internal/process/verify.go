package process

import (
	"math"

	"otto-air/internal/gas"
)

const DefaultTolerance = 1e-9

func IsentropicConstant(g gas.Gas, a, b gas.State, tol float64) bool {
	pvA := a.P * math.Pow(a.V, g.Gamma)
	pvB := b.P * math.Pow(b.V, g.Gamma)
	return relativeClose(pvA, pvB, tol)
}

func IsochoricConstant(a, b gas.State, tol float64) bool {
	return relativeClose(a.V, b.V, tol)
}

func IsobaricConstant(a, b gas.State, tol float64) bool {
	return relativeClose(a.P, b.P, tol)
}

func EntropyChangeAcross(g gas.Gas, seg Segment) float64 {
	return gas.EntropyDelta(g, seg.From, seg.To)
}

func SegmentSatisfies(g gas.Gas, seg Segment, tol float64) bool {
	switch seg.Kind {
	case Isentropic:
		return IsentropicConstant(g, seg.From, seg.To, tol)
	case Isochoric:
		return IsochoricConstant(seg.From, seg.To, tol)
	case Isobaric:
		return IsobaricConstant(seg.From, seg.To, tol)
	}
	return false
}

func relativeClose(a, b, tol float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return true
	}
	return math.Abs(a-b) <= tol*scale
}
