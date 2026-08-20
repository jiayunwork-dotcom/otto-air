package gas

import "math"

func StatesAlmostEqual(a, b State, tol float64) bool {
	return PressureAlmostEqual(a.P, b.P, tol) &&
		VolumeAlmostEqual(a.V, b.V, tol) &&
		TemperatureAlmostEqual(a.T, b.T, tol)
}

func PressureAlmostEqual(a, b, tol float64) bool {
	return relativeClose(a, b, tol)
}

func TemperatureAlmostEqual(a, b, tol float64) bool {
	return relativeClose(a, b, tol)
}

func VolumeAlmostEqual(a, b, tol float64) bool {
	return relativeClose(a, b, tol)
}

func relativeClose(a, b, tol float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return true
	}
	return math.Abs(a-b) <= tol*scale
}
