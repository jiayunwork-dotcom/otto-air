package process

import (
	"math"

	"otto-air/internal/gas"
)

func SegmentWork(g gas.Gas, seg Segment, n int) float64 {
	pts := SampleSegment(g, seg, n)
	return AccumulateArea(pts)
}

func SpanV(seg Segment) float64 {
	return math.Abs(seg.To.V - seg.From.V)
}

func SpanP(seg Segment) float64 {
	return math.Abs(seg.To.P - seg.From.P)
}

func SpanT(seg Segment) float64 {
	return math.Abs(seg.To.T - seg.From.T)
}

func Midpoint(g gas.Gas, seg Segment) Point {
	pts := SampleSegment(g, seg, 2)
	if len(pts) == 0 {
		return Point{}
	}
	a, b := pts[0], pts[len(pts)-1]
	return Point{P: (a.P + b.P) / 2, V: (a.V + b.V) / 2}
}

func FittedGamma(a, b gas.State) float64 {
	return math.Log(b.P/a.P) / math.Log(a.V/b.V)
}

func IsentropicPAtV(g gas.Gas, from gas.State, v float64) float64 {
	return gas.IsentropicPV(g, from, v)
}

func IsentropicTAtV(g gas.Gas, from gas.State, v float64) float64 {
	return gas.IsentropicTV(g, from, v)
}

func SegmentPressureMidpoint(g gas.Gas, seg Segment) float64 {
	return Midpoint(g, seg).P
}

func SegmentVolumeMidpoint(g gas.Gas, seg Segment) float64 {
	return Midpoint(g, seg).V
}
