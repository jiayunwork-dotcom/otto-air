package process

import "math"

func Reverse(pts []Point) []Point {
	out := make([]Point, len(pts))
	for i, p := range pts {
		out[len(pts)-1-i] = p
	}
	return out
}

func Extend(pts []Point, extra []Point) []Point {
	out := make([]Point, 0, len(pts)+len(extra))
	out = append(out, pts...)
	out = append(out, extra...)
	return out
}

func Distance(a, b Point) float64 {
	dp := b.P - a.P
	dv := b.V - a.V
	return math.Sqrt(dp*dp + dv*dv)
}

func CurveLength(pts []Point) float64 {
	var total float64
	for i := 0; i+1 < len(pts); i++ {
		total += Distance(pts[i], pts[i+1])
	}
	return total
}

func BoundingBox(pts []Point) (minP, maxP, minV, maxV float64) {
	return CurveStats(pts)
}
