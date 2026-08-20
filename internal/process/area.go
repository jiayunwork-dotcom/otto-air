package process

import "otto-air/internal/gas"

func ClosedArea(points []Point) float64 {
	if len(points) < 3 {
		return 0
	}
	var area float64
	for i := range points {
		j := (i + 1) % len(points)
		area += trapezoid(points[i], points[j])
	}
	return area
}

func SegmentArea(g gas.Gas, seg Segment, n int) float64 {
	pts := SampleSegment(g, seg, n)
	var area float64
	for i := 0; i+1 < len(pts); i++ {
		area += trapezoid(pts[i], pts[i+1])
	}
	return area
}

func trapezoid(a, b Point) float64 {
	return (a.P + b.P) / 2 * (b.V - a.V)
}

func AccumulateArea(pts []Point) float64 {
	var area float64
	for i := 0; i+1 < len(pts); i++ {
		area += trapezoid(pts[i], pts[i+1])
	}
	return area
}
