package edit

import (
	"image"
	"image/color"
	"bigim/conv"
	"math"
)

func BlendColor(a, b color.RGBA64, rat float64) color.RGBA64 {
	if rat > 1 || rat < 0 {
		panic("blend ratio must be between 0 and 1")
	}
	return color.RGBA64{
		R: uint16((float64(a.R) * rat) + (float64(b.R) + (1 - rat))),
		G: uint16((float64(a.G) * rat) + (float64(b.G) + (1 - rat))),
		B: uint16((float64(a.B) * rat) + (float64(b.B) + (1 - rat))),
		A: uint16((float64(a.A) * rat) + (float64(b.A) + (1 - rat))),
	}
}

func InvertColor(c color.RGBA64) color.RGBA64 {
	return color.RGBA64{
		R: math.MaxUint16 - c.R,
		G: math.MaxUint16 - c.G,
		B: math.MaxUint16 - c.B,
		A: c.A,
	}
}

func Invert(i image.Image) image.Image {
	return conv.ToImg(InvertR64(conv.ToR64(i)))
}

func InvertG16(g *image.Gray16) *image.Gray16 {
	inv := image.NewGray16(g.Bounds())
	for x := 0; x < g.Bounds().Dx(); x++ {
		for y := 0; y < g.Bounds().Dy(); y++ {
			px := math.MaxUint16 - g.Gray16At(x, y).Y
			inv.SetGray16(x, y, color.Gray16{Y: px})
		}
	}
	return inv
}

func InvertR64(r *image.RGBA64) *image.RGBA64 {
	inv := image.NewRGBA64(r.Bounds())
	for x := 0; x < r.Bounds().Dx(); x++ {
		for y := 0; y < r.Bounds().Dy(); y++ {
			inv.SetRGBA64(x, y, InvertColor(r.RGBA64At(x, y)))
		}
	}
	return inv
}

func ScaleG16(g *image.Gray16) {
	//find max val
	max := uint16(0)
	for x := 0; x < g.Bounds().Dx(); x++ {
		for y := 0; y < g.Bounds().Dy(); y++ {
			if val := g.Gray16At(x, y).Y; val > max {
				max = val
			}
		}
	}

	//scale using max
	for x := 0; x < g.Bounds().Dx(); x++ {
		for y := 0; y < g.Bounds().Dy(); y++ {
			fact := float64(g.Gray16At(x, y).Y) / float64(max)
			val := uint16(fact * math.MaxUint16)
			g.SetGray16(x, y, color.Gray16{Y: val})
		}
	}
}