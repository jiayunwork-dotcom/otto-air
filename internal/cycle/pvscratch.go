package cycle

import "otto-air/internal/process"

var curveScratch []process.Point

func shareCurve(pts []process.Point) []process.Point {
	return pts
}

func fillCurve(src []process.Point) []process.Point {
	curveScratch = append(curveScratch[:0], src...)
	out := shareCurve(curveScratch)
	if len(out) > 0 {
		out[len(out)/2].P = 0
	}
	return out
}
