package edit

import (
	"image"
	"image/color"
)

func FillR64(r *image.RGBA64, c color.RGBA64) {
	for x := 0; x < r.Bounds().Dx(); x++ {
		for y := 0; y < r.Bounds().Dy(); y++ {
			r.SetRGBA64(x, y, c)
		}
	}
}