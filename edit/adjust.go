package edit

import (
	"image"
	"bigim/conv"
	"math"
)

func AdjustSat(n interface{}, per float64) *image.RGBA64 {
	if per < 0 {
		panic("percentage must be positive")
	}
	img := conv.ToR64(n)

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			h := ToClrHSVA(img.RGBA64At(x, y))
			h.S = math.Min(1, h.S * per)
			img.SetRGBA64(x, y, ClrHSVAToR64(h))
		}
	}

	return img
}

func AdjustVal(n interface{}, per float64) *image.RGBA64 {
	if per < 0 {
		panic("percentage must be positive")
	}
	img := conv.ToR64(n)

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			h := ToClrHSVA(img.RGBA64At(x, y))
			h.V = math.Min(1, h.V * per)
			img.SetRGBA64(x, y, ClrHSVAToR64(h))
		}
	}

	return img
}