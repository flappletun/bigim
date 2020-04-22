package conv

import (
	"image"
	"image/color"
)

func ToR64(i interface{}) *image.RGBA64 {
	switch i.(type) {
	case image.Image: return imageToR64(i.(image.Image))
	case *image.Gray16: return gray16ToR64(i.(*image.Gray16))
	case *image.RGBA64: return i.(*image.RGBA64)
	default: return nil
	}
}

func ToG16(i interface{}) *image.Gray16 {
	switch i.(type) {
	case image.Image: return imageToGray16(i.(image.Image))
	case *image.RGBA64: return r64ToGray16(i.(*image.RGBA64))
	case *image.Gray16: return i.(*image.Gray16)
	default: return nil
	}
}

func ToImg(i interface{}) image.Image {
	switch i.(type) {
	case *image.RGBA64: return r64ToImage(i.(*image.RGBA64))
	case *image.Gray16: return gray16ToImage(i.(*image.Gray16))
	case image.Image: return i.(image.Image)
	default: return nil
	}
}

func imageToR64(i image.Image) *image.RGBA64 {
	r64 := image.NewRGBA64(i.Bounds())

	for x := 0; x < r64.Bounds().Dx(); x++ {
		for y := 0; y < r64.Bounds().Dy(); y++ {
			r, g, b, a := i.At(x, y).RGBA()
			r64.SetRGBA64(x, y, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return r64
}

func imageToGray16(i image.Image) *image.Gray16 {
	g16 := image.NewGray16(i.Bounds())

	for x := 0; x < g16.Bounds().Dx(); x++ {
		for y := 0; y < g16.Bounds().Dy(); y++ {
			g16.Set(x, y, color.Gray16Model.Convert(i.At(x, y)))
		}
	}

	return g16
}

func r64ToImage(r *image.RGBA64) image.Image {
	return r.SubImage(r.Bounds())
}

func r64ToGray16(r *image.RGBA64) *image.Gray16 {
	g16 := image.NewGray16(r.Bounds())

	for x := 0; x < g16.Bounds().Dx(); x++ {
		for y := 0; y < g16.Bounds().Dy(); y++ {
			g16.Set(x, y, color.Gray16Model.Convert(r.At(x, y)))
		}
	}

	return g16
}

func gray16ToImage(g *image.Gray16) image.Image {
	return g.SubImage(g.Bounds())
}

func gray16ToR64(g *image.Gray16) *image.RGBA64 {
	r64 := image.NewRGBA64(g.Bounds())

	for x := 0; x < r64.Bounds().Dx(); x++ {
		for y := 0; y < r64.Bounds().Dy(); y++ {
			r64.Set(x, y, color.RGBA64Model.Convert(g.Gray16At(x, y)))
		}
	}

	return r64
}