package process

import (
	"math"

	"otto-air/internal/gas"
)

type Point struct {
	P float64
	V float64
}

func SampleSegment(g gas.Gas, seg Segment, n int) []Point {
	switch seg.Kind {
	case Isentropic:
		return sampleIsentropic(g, seg, n)
	case Isochoric:
		return sampleIsochoric(seg)
	case Isobaric:
		return sampleIsobaric(seg, n)
	}
	return []Point{}
}

func sampleIsentropic(g gas.Gas, seg Segment, n int) []Point {
	if n < 2 {
		n = 2
	}
	pts := make([]Point, 0, n+1)
	ratio := seg.To.V / seg.From.V
	for i := 0; i <= n; i++ {
		frac := float64(i) / float64(n)
		v := seg.From.V * math.Pow(ratio, frac)
		p := gas.IsentropicPV(g, seg.From, v)
		pts = append(pts, Point{P: p, V: v})
	}
	return pts
}

func sampleIsochoric(seg Segment) []Point {
	return []Point{
		{P: seg.From.P, V: seg.From.V},
		{P: seg.To.P, V: seg.To.V},
	}
}

func sampleIsobaric(seg Segment, n int) []Point {
	if n < 2 {
		n = 2
	}
	pts := make([]Point, 0, n+1)
	for i := 0; i <= n; i++ {
		frac := float64(i) / float64(n)
		v := seg.From.V + (seg.To.V-seg.From.V)*frac
		pts = append(pts, Point{P: seg.From.P, V: v})
	}
	return pts
}
