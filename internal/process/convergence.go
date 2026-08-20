package process

import (
	"math"

	"otto-air/internal/gas"
)

func AreaConvergence(g gas.Gas, segments []Segment, maxSamples int) []float64 {
	out := make([]float64, 0, maxSamples)
	for n := 1; n <= maxSamples; n++ {
		var pts []Point
		for _, seg := range segments {
			pts = append(pts, SampleSegment(g, seg, n)...)
		}
		out = append(out, ClosedArea(DedupeConsecutive(pts)))
	}
	return out
}

func AreaErrors(g gas.Gas, segments []Segment, reference float64, maxSamples int) []float64 {
	areas := AreaConvergence(g, segments, maxSamples)
	out := make([]float64, len(areas))
	for i, a := range areas {
		scale := math.Max(math.Abs(reference), 1)
		out[i] = math.Abs(a-reference) / scale
	}
	return out
}

func ConvergedArea(g gas.Gas, segments []Segment, tol float64, maxSamples int) (float64, int, bool) {
	errors := AreaErrors(g, segments, 0, maxSamples)
	if len(errors) == 0 {
		return 0, 0, false
	}
	areas := AreaConvergence(g, segments, maxSamples)
	for i := len(areas) - 1; i >= 0; i-- {
		if i > 0 && math.Abs(areas[i]-areas[i-1]) <= tol {
			return areas[i], i + 1, true
		}
	}
	return areas[len(areas)-1], len(areas), false
}
