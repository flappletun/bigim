package cv

import (
	"image"
	"image/color"
	"bigim/conv"
	"math"
)

type bucket struct {
	r,
	g,
	b,
	a float64
}

func Convolve(i image.Image, k [][]float64) image.Image {
	return conv.ToImg(ConvolveR64(conv.ToR64(i), k))
}

func ConvolveR64(r *image.RGBA64, k [][]float64) *image.RGBA64 {

	//ensure kernel is square
	if len(k) != len(k[0]) {
		panic("kernel must be square")
	}

	//add transparent pixels to edges
	offset := len(k) / 2
	//src := padR64(r, offset) // for now!!!!!
	src := r

	//build new image using kernel
	dst := image.NewRGBA64(r.Bounds())
	for x := 0; x < src.Bounds().Dx() - (2 * offset); x++ {
		for y := 0; y < src.Bounds().Dy() - (2 * offset); y++ {

			/* fill up kernel accumulator */
			bkt := bucket{}
			for i := range k {
				for j := range k {
					r, g, b, a := src.At(x + j, y + i).RGBA()
					bkt.r += k[i][j] * float64(r)
					bkt.g += k[i][j] * float64(g)
					bkt.b += k[i][j] * float64(b)
					bkt.a += k[i][j] * float64(a)
				}
			}

			/* set dst pixel */
			dst.SetRGBA64(x + offset, y + offset, color.RGBA64{
				R: uint16(bkt.r),
				G: uint16(bkt.g),
				B: uint16(bkt.b),
				A: uint16(bkt.a),
			})
		}
	}
	return dst
}

func ConvolveGray16(g *image.Gray16, k [][]float64) *image.Gray16 {
	//ensure kernel is square
	if len(k) != len(k[0]) {
		panic("kernel must be square")
	}

	//add transparent pixels to edges
	offset := len(k) / 2
	//src := padGray16(g, offset) // for now!!!!!
	src := g

	//build new image using kernel
	dst := image.NewGray16(g.Bounds())
	for x := 0; x < src.Bounds().Dx() - offset; x++ {
		for y := 0; y < src.Bounds().Dy() - offset; y++ {

			//accumulate pixel intensities within kernel
			bkt := 0.0
			for i := range k {
				for j := range k {
					bkt += k[i][j] * float64(src.Gray16At(x + j, y + i).Y)
				}
			}

			//set new pixel value
			if float64(math.MaxUint16) < bkt {
				dst.SetGray16(x + offset, y + offset, color.Gray16{Y:math.MaxUint16})
			} else {
				dst.SetGray16(x + offset, y + offset, color.Gray16{Y:uint16(math.Round(bkt))})
			}
		}
	}
	return dst
}