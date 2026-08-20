package cycle

import (
	"otto-air/internal/gas"
	"otto-air/internal/process"
)

func Curve(in Input, samples int) ([]process.Point, float64, error) {
	in = in.Normalize()
	if err := in.Validate(); err != nil {
		return nil, 0, err
	}
	g, err := gas.New(in.Gamma)
	if err != nil {
		return nil, 0, err
	}
	states, err := computeEndpoints(g, in)
	if err != nil {
		return nil, 0, err
	}
	return assembleCurve(g, in, states, samples), closedAreaFor(g, in, states), nil
}

func assembleCurve(g gas.Gas, in Input, states gas.States, samples int) []process.Point {
	comp := process.MakeSegment(process.Isentropic, states[0], states[1])
	heat := heatingSegment(g, in, states)
	exp := process.MakeSegment(process.Isentropic, states[2], states[3])
	vent := process.MakeSegment(process.Isochoric, states[3], states[0])
	parts := [][]process.Point{
		process.SampleSegment(g, comp, process.SampleCountFor(process.Isentropic, samples)),
		process.SampleSegment(g, heat, process.SampleCountFor(heat.Kind, samples)),
		process.SampleSegment(g, exp, process.SampleCountFor(process.Isentropic, samples)),
		process.SampleSegment(g, vent, process.SampleCountFor(process.Isochoric, samples)),
	}
	pts := process.JoinSegments(parts...)
	return process.CloseLoop(pts)
}

func heatingSegment(g gas.Gas, in Input, states gas.States) process.Segment {
	if in.Mode == Diesel {
		return process.MakeSegment(process.Isobaric, states[1], states[2])
	}
	return process.MakeSegment(process.Isochoric, states[1], states[2])
}

func closedAreaFor(g gas.Gas, in Input, states gas.States) float64 {
	pts := assembleCurve(g, in, states, 0)
	return process.ClosedArea(pts)
}

func CurveSegments(g gas.Gas, in Input, states gas.States) []process.Segment {
	return []process.Segment{
		process.MakeSegment(process.Isentropic, states[0], states[1]),
		heatingSegment(g, in, states),
		process.MakeSegment(process.Isentropic, states[2], states[3]),
		process.MakeSegment(process.Isochoric, states[3], states[0]),
	}
}
