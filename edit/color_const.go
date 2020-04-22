package edit

import (
	"image/color"
	"math"
)

func Red() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16,
		G: 0,
		B: 0,
		A: math.MaxUint16,
	}
}

func Blue() color.RGBA64 {
	return color.RGBA64{
		R: 0,
		G: 0,
		B: math.MaxUint16,
		A: math.MaxUint16,
	}
}

func Green() color.RGBA64 {
	return color.RGBA64{
		R: 0,
		G: math.MaxUint16,
		B: 0,
		A: math.MaxUint16,
	}
}

func Cyan() color.RGBA64 {
	return color.RGBA64{
		R: 0,
		G: math.MaxUint16,
		B: math.MaxUint16,
		A: math.MaxUint16,
	}
}

func Magenta() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16,
		G: 0,
		B: math.MaxUint16,
		A: math.MaxUint16,
	}
}

func Yellow() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16,
		G: math.MaxUint16,
		B: 0,
		A: math.MaxUint16,
	}
}

func Orange() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16,
		G: math.MaxUint16 / 2,
		A: math.MaxUint16,
	}
}

func Purple() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16 / 2,
		B: math.MaxUint16,
		A: math.MaxUint16,
	}
}

func White() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16,
		G: math.MaxUint16,
		B: math.MaxUint16,
		A: math.MaxUint16,
	}
}

func Gray() color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16 / 2,
		G: math.MaxUint16 / 2,
		B: math.MaxUint16 / 2,
		A: math.MaxUint16,
	}
}

func Black() color.RGBA64 {
	return color.RGBA64{
		R: 0,
		G: 0,
		B: 0,
		A: math.MaxUint16,
	}
}