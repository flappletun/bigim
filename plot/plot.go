package plot

import (
	"bigim/edit"
	"fmt"
	"image"
	"image/color"
	"math"
)

func HSV(n interface{}) image.Image {
	plot := edit.NewIllustrator(makePlot())

	hsva := edit.ToImgHSVA(n)
	pxMap := make(map[string]bool, len(hsva.Pixels))

	for _, v := range hsva.Pixels {
		pStr := hashPixel(v.H, v.S, v.V)
		if !pxMap[pStr] {
			pxMap[pStr] = true
			if v.S < 1 {
				x := int((v.H / 360) * 750) + 30
				y := int((1 - v.S) * 750) + 10
				plot.DrawCircle(x, y, 2, findColor(v.V))
			}
		}
	}

	return plot.ToImage()
}

func makePlot() *image.RGBA64 {
	plot := image.NewRGBA64(image.Rect(0,0,790, 790))
	edit.FillR64(plot, edit.Black())

	//make bar to display along hue (X) axis
	hBar := image.NewRGBA64(image.Rect(0,0,750,1))
	for i := 0; i < 750; i++ {
		hBar.SetRGBA64(i, 0, edit.ClrHSVAToR64(edit.ClrHSVA{
			H: 360 * (float64(i) / 750),
			S: 1,
			V: 0.8,
			A: math.MaxUint16,
		}))
	}
	hBar = edit.ResizeR64(hBar, 10, 750)

	//make bar to display along sat (Y) axis
	sBar := image.NewRGBA64(image.Rect(0,0,1,750))
	for i := 0; i < 750; i++ {
		sBar.SetRGBA64(0, i, edit.ClrHSVAToR64(edit.ClrHSVA{
			H: 225,
			S: 1 - (float64(i) / 750),
			V: 0.625,
			A: math.MaxUint16,
		}))
	}
	sBar = edit.ResizeR64(sBar, 750, 10)

	//add axis bars to plot
	edit.PasteR64(plot, hBar, 30, 770)
	edit.PasteR64(plot, sBar, 10, 10)

	return plot
}

func hashPixel(h, s, v float64) string {
	return fmt.Sprint(h) + fmt.Sprint(s) + fmt.Sprint(v)
}

func findColor(n float64) color.RGBA64 {
	return edit.ClrHSVAToR64(edit.ClrHSVA{
		H: n * 120,
		S: 1,
		V: 1,
		A: math.MaxUint16,
	})
}