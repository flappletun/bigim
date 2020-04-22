package cv

import (
	"image"
	"bigim/conv"
	"bigim/edit"
	"math"
)

type sigma int
const (
	s1 = 0
	s2 = 1
	s4 = 2
	s8 = 3
	s16 = 4
)

type DoG struct {
	gImg *image.Gray16
	cImg *image.RGBA64
	scSpace []scale
	maxima,
	minima,
	extrema,
	interestPts []Point
}

type scale struct {
	img *image.Gray16
	minVal,
	maxVal uint16
}

func NewDoG(n interface{}) *DoG {
	d :=  &DoG{
		gImg: conv.ToG16(n),
		cImg: conv.ToR64(n),
		scSpace: make([]scale, 4),
	}

	//convolve image with GaussPow2 of different scales (sigma)
	scales := make([]*image.Gray16, 5)
	for i := range scales {
		scales[i] = ConvolveGray16(d.gImg, GaussPow2[i])
	}

	//create DoG image for each scale
	for i := range d.scSpace {
		img, max, min := subtractGray16(scales[i], scales[i + 1])
		d.scSpace[i].img = img
		d.scSpace[i].maxVal = max
		d.scSpace[i].minVal = min
	}

	//mark and tally extrema locations across scales
	d.extrema = make([]Point, 0, d.cImg.Bounds().Dx() * d.cImg.Bounds().Dy())
	mxMap := make(map[string]int)
	mnMap := make(map[string]int)
	for _, ss := range d.scSpace {
		findLocalExtrema(ss.img, mxMap, mnMap)
	}

	//keep points found in more than 1 scale
	for k, v := range mxMap {
		if v > 1 {
			p := unhashPoint(k)
			d.maxima = append(d.maxima, p)
			d.extrema = append(d.extrema, p)
		}
	}
	for k, v := range mnMap {
		if v > 1 {
			p := unhashPoint(k)
			d.minima = append(d.minima, p)
			d.extrema = append(d.extrema, p)
		}
	}

	return d
}

func (d DoG) DoGs() []image.Image {
	outs := make([]image.Image, len(d.scSpace))
	for i, v := range d.scSpace {
		edit.ScaleG16(v.img)
		outs[i] = conv.ToImg(v.img)
	}
	return outs
}

func (d DoG) Extrema(nExt int) []Point {
	if nExt <= 0 {
		nExt = len(d.extrema)
	}
	return d.extrema[:nExt]
}

func (d DoG) Maxima(nExt int) []Point {
	if nExt <= 0 {
		nExt = len(d.maxima)
	}
	return d.maxima[:nExt]
}

func (d DoG) Minima(nExt int) []Point {
	if nExt <= 0 {
		nExt = len(d.minima)
	}
	return d.minima[:nExt]
}

func (d DoG) PlotExtrema(nExt int) image.Image {
	if nExt <= 0 {
		nExt = len(d.extrema)
	}
	plot := edit.NewIllustrator(copyR64(d.cImg))
	for _, p := range d.extrema[:nExt] {
		plot.DrawCircle(p.X, p.Y, 5, edit.Red())
	}
	return plot.ToImage()
}

func findLocalExtrema(g *image.Gray16, mxMap, mnMap map[string]int) {
	rootX, rootY := boundRoots(g.Bounds())
	pxPerBlkX := g.Bounds().Dx() / rootX
	pxPerBlkY := g.Bounds().Dy() / rootY

	//visit each block
	for x := 0; x < g.Bounds().Dx(); x += pxPerBlkX {
		for y := 0; y < g.Bounds().Dy(); y += pxPerBlkY {
			//find max and min px within block
			maxVal := uint16(0)
			minVal := uint16(math.MaxUint16)
			var maxX, maxY, minX, minY int
			for i := 0; i < pxPerBlkX; i++ {
				for j := 0; j < pxPerBlkY; j++ {
					val := g.Gray16At(x + i, y + j).Y
					if val < minVal {
						minVal = val
						minX = x + i
						minY = y + j
					}
					if val > maxVal {
						maxVal = val
						maxX = x + i
						maxY = y + j
					}
				}
			}
			mnMap[hashPoint(minX, minY)] += 1
			mxMap[hashPoint(maxX, maxY)] += 1
		}
	}

	//visit each block again, offset to handle boundaries
	for x := pxPerBlkX / 2; x < g.Bounds().Dx(); x += pxPerBlkX {
		for y := pxPerBlkY / 2; y < g.Bounds().Dy(); y += pxPerBlkY {
			//find max and min px within block
			maxVal := uint16(0)
			minVal := uint16(math.MaxUint16)
			var maxX, maxY, minX, minY int
			for i := 0; i < pxPerBlkX; i++ {
				for j := 0; j < pxPerBlkY; j++ {
					val := g.Gray16At(x + i, y + j).Y
					if val < minVal {
						minVal = val
						minX = x + i
						minY = y + j
					}
					if val > maxVal {
						maxVal = val
						maxX = x + i
						maxY = y + j
					}
				}
			}
			mnMap[hashPoint(minX, minY)] += 1
			mxMap[hashPoint(maxX, maxY)] += 1
		}
	}
}