package cv

import (
	"image"
	"image/color"
	"math"
)

func ApplyEdgeThreshold(g *image.Gray16, tPer float64, flood bool) {

	//find max value and set threshold
	var max uint16
	for x := 0; x < g.Bounds().Dx(); x++ {
		for y := 0; y < g.Bounds().Dy(); y++ {
			if val := g.Gray16At(x, y).Y; val > max {
				max = val
			}
		}
	}
	thresh := uint16(float64(max) * tPer)

	//set all px below thresh to black
	for x := 0; x < g.Bounds().Dx(); x++ {
		for y := 0; y < g.Bounds().Dy(); y++ {
			if val := g.Gray16At(x, y).Y; val < thresh {
				g.SetGray16(x, y, color.Gray16{Y: 0})
			} else if flood {
				g.SetGray16(x, y, color.Gray16{Y: math.MaxUint16})
			}
		}
	}
}

func EraseArtifacts(g *image.Gray16, winSize, minPx int) {

	//defaults
	if winSize == 0 {winSize = 9}
	if minPx == 0 {minPx = 4}

	for x := 0; x < g.Bounds().Dx() - (winSize / 2); x++ {
		for y := 0; y < g.Bounds().Dy() - (winSize / 2); y++ {

			//count non-black pixels in window
			count := 0
			for i := 0; i < winSize; i++ {
				for j := 0; j < winSize; j++ {
					if g.Gray16At(x, y).Y > 0 {
						count += 1
					}
				}
			}

			//if isolated noise, erase
			if count < minPx {
				g.SetGray16(x, y, color.Gray16{Y: 0})
			}
		}
	}
}