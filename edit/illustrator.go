package edit

import (
	"image"
	"image/color"
	"bigim/conv"
	"bigim/imio"
	"math"
)

type Illustrator struct {
	canvas interface{}
}

func NewIllustrator(i interface{}) Illustrator {
	return Illustrator{canvas: i}
}

func (i *Illustrator) DrawLine(x0, x1, y0, y1 int, c color.RGBA64) {

	//assert canvas type
	var can *image.RGBA64
	switch i.canvas.(type) {
	case *image.RGBA64: can = i.canvas.(*image.RGBA64)
	case image.Image: can = conv.ToR64(i.canvas.(image.Image))
	case *image.Gray16: can = conv.ToR64(i.canvas.(*image.Gray16))
	}

	//find dist to inc while hitting every pixel along line
	dist := distance(x0, x1, y0, y1)

	//draw line
	for dt := 0.0; dt <= dist; dt += 1 / (2 * dist) {
		t := dt / dist
		xt := ((1 - t) * float64(x0)) + (t * float64(x1))
		yt := ((1 - t) * float64(y0)) + (t * float64(y1))
		if can != nil {
			can.SetRGBA64(int(xt), int(yt), c)
		}
	}
	i.canvas = can
}

func (i *Illustrator) DrawCircle(x, y int, r float64, c color.RGBA64) {

	//assert canvas type
	var can *image.RGBA64
	switch i.canvas.(type) {
	case *image.RGBA64: can = i.canvas.(*image.RGBA64)
	case image.Image: can = conv.ToR64(i.canvas.(image.Image))
	case *image.Gray16: can = conv.ToR64(i.canvas.(*image.Gray16))
	}

	//find min angle to inc while hitting every px in circle
	var step float64
	switch {
	case r < 6: step = 10
	case r < 20: step = 3
	case r < 50: step = 1
	default: step = 0.5
	}

	//draw the circle
	for t := 0.0; t < 360.0; t += step {
		xAddr := math.Round(r * math.Cos(radians(t))) + float64(x)
		yAddr := math.Round(r * math.Sin(radians(t))) + float64(y)
		if can != nil {
			can.SetRGBA64(int(xAddr), int(yAddr), c)
		}
	}
	i.canvas = can
}

func (i *Illustrator) DrawRectangle(x0, x1, y0, y1 int, c color.RGBA64) (int, int) {

	//assert canvas type
	var can *image.RGBA64
	switch i.canvas.(type) {
	case *image.RGBA64: can = i.canvas.(*image.RGBA64)
	case image.Image: can = conv.ToR64(i.canvas.(image.Image))
	case *image.Gray16: can = conv.ToR64(i.canvas.(*image.Gray16))
	}

	//draw horizontal lines
	for x := x0; x <= x1; x++ {
		if can != nil {
			can.SetRGBA64(x, y0, c)
			can.SetRGBA64(x, y1, c)
		}
	}

	//draw vertical lines
	for y := y0; y <= y1; y++ {
		if can != nil {
			can.SetRGBA64(x0, y, c)
			can.SetRGBA64(x1, y, c)
		}
	}

	//locate center pixel
	dist := distance(x0, x1, y0, y1)
	t := 0.5 / dist
	xCenter := ((1 - t) * float64(x0)) + (t * float64(x1))
	yCenter := ((1 - t) * float64(y0)) + (t * float64(y1))

	i.canvas = can
	return int(xCenter), int(yCenter)
}

func (i *Illustrator) DrawLines(x0, x1, y0, y1 []int, c color.RGBA64) {
	for j := range x0 {
		i.DrawLine(x0[j], x1[j], y0[j], y1[j], c)
	}
}

func (i *Illustrator) DrawCircles(x, y []int, r float64, c color.RGBA64) {
	for j := range x {
		i.DrawCircle(x[j], y[j], r, c)
	}
}

func (i *Illustrator) DrawRectangles(x0, x1, y0, y1 []int, c color.RGBA64) ([]int, []int) {
	xCenters := make([]int, len(x0))
	yCenters := make([]int, len(x0))
	for j := range x0 {
		a, b := i.DrawRectangle(x0[j], x1[j], y0[j], y1[j], c)
		xCenters = append(xCenters, a)
		yCenters = append(yCenters, b)
	}
	return xCenters, yCenters
}

func (i Illustrator) Save(path string) {
	var im image.Image
	switch i.canvas.(type) {
	case image.Image: im = i.canvas.(image.Image)
	case *image.RGBA64: {
		r := i.canvas.(*image.RGBA64)
		im = conv.ToImg(r)
	}
	case *image.Gray16: {
		g := i.canvas.(*image.Gray16)
		im = conv.ToImg(g)
	}
	}
	imio.Save(im, path)
}

func (i Illustrator) ToImage() image.Image {
	return conv.ToImg(i.canvas)
}

func (i Illustrator) ToR64() *image.RGBA64 {
	return conv.ToR64(i.canvas)
}

func (i Illustrator) ToGray16() *image.Gray16 {
	return conv.ToG16(i.canvas)
}
