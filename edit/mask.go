package edit

import (
	"image"
	"bigim/conv"
	"math"
)

func Mask(f, b, m interface{}) *image.RGBA64 {
	frg := conv.ToR64(f)
	bkg := conv.ToR64(b)
	msk := conv.ToG16(m)
	if !dimsMatch(frg, bkg, msk) {
		panic("image sizes must be equal")
	}

	out := image.NewRGBA64(bkg.Bounds())
	for x := 0; x < out.Bounds().Dx(); x++ {
		for y := 0; y < out.Bounds().Dy(); y++ {
			if mVal := msk.Gray16At(x, y).Y; mVal != 0 {
				//find presence (ratio) of foreground pixel
				pres := float64(mVal) / math.MaxUint16
				out.SetRGBA64(x, y, BlendColor(frg.RGBA64At(x, y), bkg.RGBA64At(x, y), pres))
			}
		}
	}

	return out
}

func dimsMatch(a, b *image.RGBA64, c *image.Gray16) bool {
	if a.Bounds().Dx() != b.Bounds().Dx() ||
		b.Bounds().Dx() != c.Bounds().Dx() ||
		a.Bounds().Dy() != b.Bounds().Dy() ||
		b.Bounds().Dy() != c.Bounds().Dy() {
		return false
	}
	return true
}