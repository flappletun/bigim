package edit

import (
	"image"
	"bigim/conv"
	"math"
)

func Resize(i image.Image, h, w int) image.Image {
	return conv.ToImg(ResizeR64(conv.ToR64(i), h, w))
}

func ResizeR64(r *image.RGBA64, h, w int) *image.RGBA64 {
	out := image.NewRGBA64(image.Rect(0,0, w, h))
	xOut := float64(out.Bounds().Dx())
	yOut := float64(out.Bounds().Dy())
	xOrig := float64(r.Bounds().Dx())
	yOrig := float64(r.Bounds().Dy())

	for x := 0.0; x < xOut; x++ {
		for y := 0.0; y < yOut; y++ {
			xRat := x / xOut
			yRat := y / yOut
			xSamp := int(math.Trunc(xRat * xOrig))
			ySamp := int(math.Trunc(yRat * yOrig))
			out.SetRGBA64(int(x), int(y), r.RGBA64At(xSamp, ySamp))
		}
	}

	return out
}

func ResizeG16(g *image.Gray16, h, w int) *image.Gray16 {
	out := image.NewGray16(image.Rect(0,0, w, h))
	xOut := float64(out.Bounds().Dx())
	yOut := float64(out.Bounds().Dy())
	xOrig := float64(g.Bounds().Dx())
	yOrig := float64(g.Bounds().Dy())

	for x := 0.0; x < xOut; x++ {
		for y := 0.0; y < yOut; y++ {
			xRat := x / xOut
			yRat := y / yOut
			xSamp := int(math.Trunc(xRat * xOrig))
			ySamp := int(math.Trunc(yRat * yOrig))
			out.SetGray16(int(x), int(y), g.Gray16At(xSamp, ySamp))
		}
	}

	return out
}