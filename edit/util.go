package edit

import "math"

func radians(deg float64) float64 {
	oneDeg := 0.0174533
	if deg < 0 {
		deg = 360 - math.Abs(deg)
	}
	return deg * oneDeg
}

func degrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func distance(x0, x1, y0, y1 int) float64 {
	a := math.Pow(float64(x1) - float64(x0), 2)
	b := math.Pow(float64(y1) - float64(y0), 2)
	return math.Sqrt(a + b)
}