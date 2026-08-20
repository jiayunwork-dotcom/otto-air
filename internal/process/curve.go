package process

func JoinSegments(parts ...[]Point) []Point {
	var out []Point
	for _, part := range parts {
		out = append(out, part...)
	}
	return DedupeConsecutive(out)
}

func DedupeConsecutive(pts []Point) []Point {
	if len(pts) == 0 {
		return pts
	}
	out := make([]Point, 0, len(pts))
	prev := pts[0]
	out = append(out, prev)
	for _, p := range pts[1:] {
		if p.V == prev.V && p.P == prev.P {
			continue
		}
		out = append(out, p)
		prev = p
	}
	return out
}

func CloseLoop(pts []Point) []Point {
	if len(pts) == 0 {
		return pts
	}
	first := pts[0]
	last := pts[len(pts)-1]
	if first.V == last.V && first.P == last.P {
		return pts
	}
	return append(pts, first)
}

func CurveStats(pts []Point) (minP, maxP, minV, maxV float64) {
	if len(pts) == 0 {
		return 0, 0, 0, 0
	}
	minP, maxP = pts[0].P, pts[0].P
	minV, maxV = pts[0].V, pts[0].V
	for _, p := range pts[1:] {
		minP = min(minP, p.P)
		maxP = max(maxP, p.P)
		minV = min(minV, p.V)
		maxV = max(maxV, p.V)
	}
	return minP, maxP, minV, maxV
}
