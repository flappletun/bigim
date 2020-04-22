package cv

import (
	"image"
	"image/color"
	"math"
	"bigim/u16math"
)

func gradValXY(a, b uint16) uint16 {
	aa := math.Pow(float64(a), 2)
	bb := math.Pow(float64(b), 2)
	return uint16(math.Sqrt(aa + bb))
}

func radians(deg float64) float64 {
	oneDeg := 0.0174533
	if deg < 0 {
		deg = 360 - math.Abs(deg)
	}
	return deg * oneDeg
}

func degrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func distance(x0, x1, y0, y1 int) float64 {
	a := math.Pow(float64(x1) - float64(x0), 2)
	b := math.Pow(float64(y1) - float64(y0), 2)
	return math.Sqrt(a + b)
}

func padR64(r *image.RGBA64, padding int) *image.RGBA64 {
	h := r.Bounds().Dy() + padding
	w := r.Bounds().Dx() + padding
	ret := image.NewRGBA64(image.Rect(0, 0, w, h))
	for x := padding; x < ret.Bounds().Dx() - padding; x++ {
		for y := padding; y < ret.Bounds().Dy() - padding; y++ {
			ret.SetRGBA64(x, y, r.RGBA64At(x - padding, y - padding))
		}
	}
	return ret
}

func padGray16(g *image.Gray16, padding int) *image.Gray16 {
	h := g.Bounds().Dy() + padding
	w := g.Bounds().Dx() + padding
	ret := image.NewGray16(image.Rect(0, 0, w, h))
	for x := padding; x < ret.Bounds().Dx() - padding; x++ {
		for y := padding; y < ret.Bounds().Dy() - padding; y++ {
			ret.SetGray16(x, y, g.Gray16At(x - padding, y - padding))
		}
	}
	return ret
}

func subtractGray16(a, b *image.Gray16) (*image.Gray16, uint16, uint16) {
	min := uint16(math.MaxUint16)
	max := uint16(0)
	c := image.NewGray16(a.Bounds())
	for x := 0; x < c.Bounds().Dx(); x++ {
		for y := 0; y < c.Bounds().Dy(); y++ {
			val := u16math.AbsDiff(a.Gray16At(x, y).Y, b.Gray16At(x, y).Y)
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
			c.SetGray16(x, y, color.Gray16{Y: val})
		}
	}
	return c, max, min
}

func normalize(f []float64) {
	var sum float64
	for _, v := range f {
		sum += v
	}
	for i := range f {
		f[i] /= sum
	}
}

//DoG

func copyR64(r *image.RGBA64) *image.RGBA64 {
	copy := image.NewRGBA64(r.Bounds())
	for x := 0; x < r.Bounds().Dx(); x++ {
		for y := 0; y < r.Bounds().Dy(); y++ {
			copy.SetRGBA64(x, y, r.RGBA64At(x, y))
		}
	}
	return copy
}

func boundRoots(rect image.Rectangle) (int, int) {
	rootX := int(math.Sqrt(float64(rect.Dx())))
	rootY := int(math.Sqrt(float64(rect.Dy())))
	return rootX, rootY
}