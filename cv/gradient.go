package cv

import (
	"image"
	"image/color"
	"bigim/conv"
	"bigim/u16math"
)

func GradientDx(i image.Image) *image.Gray16 {
	g := conv.ToG16(i)
	dx := image.NewGray16(i.Bounds())

	for x := 1; x < dx.Bounds().Dx() - 1; x++ {
		for y := 0; y < dx.Bounds().Dy(); y++ {
			lVal := g.Gray16At(x - 1, y).Y
			rVal := g.Gray16At(x + 1, y).Y
			dx.SetGray16(x, y, color.Gray16{Y: u16math.AbsDiff(lVal, rVal)})
		}
	}

	return dx
}

func GradientDy(i image.Image) *image.Gray16 {
	g := conv.ToG16(i)
	dy := image.NewGray16(i.Bounds())

	for x := 0; x < dy.Bounds().Dx(); x++ {
		for y := 1; y < dy.Bounds().Dy() - 1; y++ {
			aVal := g.Gray16At(x, y - 1).Y
			bVal := g.Gray16At(x, y + 1).Y
			dy.SetGray16(x, y, color.Gray16{Y: u16math.AbsDiff(aVal, bVal)})
		}
	}

	return dy
}

func Gradient(i image.Image) *image.Gray16 {
	dxy := image.NewGray16(i.Bounds())

	dx := GradientDx(i)
	dy := GradientDy(i)

	//I(p) = root2(dx(p)^2 + dy(p)^2)
	//I(p) = avg(dx(p) + dy(p))
	for x := 0; x < dxy.Bounds().Dx(); x++ {
		for y := 0; y < dxy.Bounds().Dy(); y++ {
			xVal := dx.Gray16At(x, y).Y
			yVAl := dy.Gray16At(x, y).Y
			dxy.SetGray16(x, y, color.Gray16{Y: gradValXY(xVal, yVAl)})
		}
	}

	return dxy
}