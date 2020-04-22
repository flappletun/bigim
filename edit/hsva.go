package edit

import (
	"image"
	"image/color"
	"bigim/conv"
	"math"
	"sort"
)

type ClrHSVA struct {
	H,
	S,
	V float64
	A uint16
}

type ImgHSVA struct {
	Pixels []ClrHSVA
	Hues,
	Sats,
	Vals []float64
	Bounds image.Rectangle
}


  /**************/
 /* RGB to HSV */
/**************/

func ToImgHSVA(n interface{}) ImgHSVA {
	switch n.(type) {
	case image.Image: return imageToImgHSVA(n.(image.Image))
	case *image.Gray16: return gray16ToImgHSVA(n.(*image.Gray16))
	case *image.RGBA64: return r64ToImgHSVA(n.(*image.RGBA64))
	}
	return ImgHSVA{}
}

func ToClrHSVA(n interface{}) ClrHSVA {
	switch n.(type) {
	case color.Gray16: {
		return r64ToClrHSVA(color.RGBA64{
			R: n.(color.Gray16).Y,
			G: n.(color.Gray16).Y,
			B: n.(color.Gray16).Y,
			A: math.MaxUint16,
		})
	}
	case color.RGBA64: return r64ToClrHSVA(n.(color.RGBA64))
	}
	return ClrHSVA{}
}

func r64ToClrHSVA(c color.RGBA64) ClrHSVA {

	//convert RGB values to percentages
	r := float64(c.R) / math.MaxUint16
	g := float64(c.G) / math.MaxUint16
	b := float64(c.B) / math.MaxUint16

	//find max and min channel values
	sl := []float64{r, g, b}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i] < sl[j]
	})
	maxC := sl[2]
	minC := sl[0]
	delta := maxC - minC

	var hue, sat, val float64

	//find hue
	if delta == 0 {
		hue = 0
	} else if maxC == r {
		hue = 60 * ((g - b) / delta)
	} else if maxC == g {
		hue = 60 * (((b - r) / delta) + 2)
	} else if maxC == b {
		hue = 60 * (((r - g) / delta) + 4)
	}
	if hue < 0 {
		hue += 360
	}

	//find sat
	if maxC == 0 {
		sat = 0
	} else {
		sat = delta / maxC
	}

	//find val
	val = maxC

	return ClrHSVA{
		H:       hue,
		S:       sat,
		V:       val,
		A:       c.A,
	}
}

func imageToImgHSVA(i image.Image) ImgHSVA {
	return r64ToImgHSVA(conv.ToR64(i))
}

func gray16ToImgHSVA(g *image.Gray16) ImgHSVA {
	return r64ToImgHSVA(conv.ToR64(g))
}

func r64ToImgHSVA(r *image.RGBA64) ImgHSVA {
	hsva := ImgHSVA{
		Pixels: make([]ClrHSVA, 0, r.Bounds().Dx() * r.Bounds().Dy()),
		Hues:   make([]float64, 0, r.Bounds().Dx() * r.Bounds().Dy()),
		Sats:   make([]float64, 0, r.Bounds().Dx() * r.Bounds().Dy()),
		Vals:   make([]float64, 0, r.Bounds().Dx() * r.Bounds().Dy()),
		Bounds: r.Bounds(),
	}

	for x := 0; x < r.Bounds().Dx(); x++ {
		for y := 0; y < r.Bounds().Dy(); y++ {
			clr := r64ToClrHSVA(r.RGBA64At(x, y))
			hsva.Pixels = append(hsva.Pixels, clr)
			hsva.Hues = append(hsva.Hues, clr.H)
			hsva.Sats = append(hsva.Sats, clr.S)
			hsva.Vals = append(hsva.Vals, clr.V)
		}
	}

	return hsva
}


  /**************/
 /* HSV to RGB */
/**************/

//WIP
func ImgHSVAToR64(h ImgHSVA) *image.RGBA64 {
	r := image.NewRGBA64(h.Bounds)
	for y := 0; y < r.Bounds().Dy(); y++ {
		for x := 0; x < r.Bounds().Dx(); x++ {
			pIdx := (y * h.Bounds.Dx()) + x
			r.SetRGBA64(x, y, ClrHSVAToR64(h.Pixels[pIdx]))
		}
	}
	return r
}

func ClrHSVAToR64(h ClrHSVA) color.RGBA64 {
	C := h.S * h.V
	m := h.V - C
	D := math.Abs(math.Mod(h.H / 60, 2) - 1)
	X := C * (1 - D)

	r, g, b := cHXToRGB(C, h.H, X)

	return color.RGBA64{
		R: uint16((r + m) * math.MaxUint16),
		G: uint16((g + m) * math.MaxUint16),
		B: uint16((b + m) * math.MaxUint16),
		A: h.A,
	}
}

func cHXToRGB(c, h, x float64) (r, g, b float64) {
	switch {
	case (h >= 0 && h < 60) || h == 360: return c, x, 0
	case h >= 60 && h < 120: return x, c, 0
	case h >= 120 && h < 180: return 0, c, x
	case h >= 180 && h < 240: return 0, x, c
	case h >= 240 && h < 300: return x, 0, c
	case h >= 300 && h < 360: return c, 0, x
	}
	return 0,0,0
}