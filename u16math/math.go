package u16math

import "math"

func AbsDiff(a, b uint16) uint16 {
	if int(a) - int(b) < 0 {
		return b - a
	} else {
		return a - b
	}
}

func SignDiff(a, b uint16) int {
	return int(a) - int(b)
}

func Pow(x uint16, e float64) uint16 {
	return uint16(math.Pow(float64(x), e))
}

func Sqrt(x uint16) uint16 {
	return uint16(math.Sqrt(float64(x)))
}