package edit

import (
	"bigim/conv"
	"image"
	"image/color"
)

func ColorFilter(i image.Image, clr color.RGBA64, perc float64) image.Image {
	r := conv.ToR64(i)
	for x := 0; x < i.Bounds().Dx(); x++ {
		for y := 0; y < i.Bounds().Dy(); y++ {
			r.SetRGBA64(x, y, BlendColor(r.RGBA64At(x, y), clr, perc))
		}
	}
	return r.SubImage(r.Bounds())
}

func R64ColorFilter(r *image.RGBA64, clr color.RGBA64, perc float64) *image.RGBA64 {
	out := image.NewRGBA64(r.Bounds())
	for x := 0; x < r.Bounds().Dx(); x++ {
		for y := 0; y < r.Bounds().Dy(); y++ {
			out.SetRGBA64(x, y, BlendColor(r.RGBA64At(x, y), clr, perc))
		}
	}
	return out
}

func G16ColorFilter(g *image.Gray16, clr color.RGBA64, perc float64) *image.RGBA64 {
	r := conv.ToR64(g)
	for x := 0; x < g.Bounds().Dx(); x++ {
		for y := 0; y < g.Bounds().Dy(); y++ {
			r.SetRGBA64(x, y, BlendColor(r.RGBA64At(x, y), clr, perc))
		}
	}
	return r
}