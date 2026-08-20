package process

import (
	"math"

	"otto-air/internal/gas"
)

func IsentropicSegmentAnalyticWork(g gas.Gas, seg Segment) float64 {
	p1, v1 := seg.From.P, seg.From.V
	v2 := seg.To.V
	if v1 == v2 {
		return 0
	}
	if v2 <= 0 {
		return math.Inf(1)
	}
	return p1 * v1 / (1 - g.Gamma) * (math.Pow(v2/v1, 1-g.Gamma) - 1)
}

func IsobaricSegmentAnalyticWork(seg Segment) float64 {
	return seg.From.P * (seg.To.V - seg.From.V)
}

func AnalyticWorkForType(g gas.Gas, seg Segment) float64 {
	switch seg.Kind {
	case Isentropic:
		return IsentropicSegmentAnalyticWork(g, seg)
	case Isobaric:
		return IsobaricSegmentAnalyticWork(seg)
	case Isochoric:
		return 0
	}
	return 0
}

func WorkError(g gas.Gas, seg Segment, n int) float64 {
	numeric := SegmentWork(g, seg, n)
	analytic := AnalyticWorkForType(g, seg)
	scale := math.Max(math.Abs(analytic), 1)
	return math.Abs(numeric-analytic) / scale
}
