package edit

import (
	"image"
	"bigim/conv"
)

func PasteOntoR64(n interface{}, img *image.RGBA64, xSt, ySt int) {
	pch := conv.ToR64(n)
	for x := xSt; x < xSt + pch.Bounds().Dx(); x++ {
		for y := ySt; y < pch.Bounds().Dy(); y++ {
			img.SetRGBA64(x, y, pch.RGBA64At(x - xSt, y - ySt))
		}
	}
}

func PasteR64(img, pch *image.RGBA64, xSt, ySt int) {
	for x := xSt; x < xSt + pch.Bounds().Dx(); x++ {
		for y := ySt; y < ySt + pch.Bounds().Dy(); y++ {
			img.SetRGBA64(x, y, pch.RGBA64At(x - xSt, y - ySt))
		}
	}
}