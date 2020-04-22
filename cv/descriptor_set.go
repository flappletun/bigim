package cv

import (
	"image"
	"bigim/conv"
	"bigim/edit"
	"math"
	"sort"
	"bigim/u16math"
)

type DescriptorSet struct {
	Desc      []Descriptor
	gImg      *image.Gray16
	cImg      *image.RGBA64
	dog       *DoG
}

type Descriptor struct {
	Point    Point
	R        float64
	Theta,
	Histogram []float64
	Dim int
	contrast uint16
	pxThetas,
	pxMags []float64
}

func NewDescriptorSet(n interface{}) *DescriptorSet {
	d :=  &DescriptorSet{
		gImg: conv.ToG16(n),
		cImg: conv.ToR64(n),
		dog:  NewDoG(n),
	}

	//find interest points from DoG extrema
	d.convertPoints()
	//fmt.Println(len(d.Desc), "extrema received from DoG")
	d.eliminateLowContrastPoints()
	//fmt.Println(len(d.Desc), "points passed contrast test")
	d.eliminateEdgePoints()
	//fmt.Println(len(d.Desc), "points passed edge response test 1")

	//assign theta and magnitude to each desc
	d.buildOrientVector()
	d.assignOrientations()

	//final edge point filter
	d.edgeFilter2(8)

	//eliminate crowded points
	d.unCrowd()

	//fmt.Println(len(d.Desc), "points passed edge response test 2")
	//fmt.Println(len(d.Desc), "descriptors created")
	//fmt.Println()

	return d
}

func (d *DescriptorSet) unCrowd() {
	size := ((d.gImg.Bounds().Dx() / 10)+ 1) * ((d.gImg.Bounds().Dy() / 10) + 1)
	sorted := make([][]Descriptor, size)
	for i := range sorted {
		sorted[i] = make([]Descriptor, 0, 10)
	}

	for _, d := range d.Desc {
		x := d.Point.X / 10
		y := d.Point.Y / 10
		sorted[(y * x) + x] = append(sorted[(y * x) + x], d)
	}
	//fmt.Println(len(sorted), len(sorted[0]))

	newDesc := make([]Descriptor, 0, len(d.Desc))
	for s := range sorted {
		if len(sorted[s]) > 0 {
			sort.Slice(sorted[s], func(i, j int) bool {
				return sorted[s][i].Dim > sorted[s][j].Dim
			})
			newDesc = append(newDesc, sorted[s][0])
		}
	}

	d.Desc = newDesc
}

func (d *DescriptorSet) edgeFilter2(minDir int) {
	pass := make([]Descriptor, 0, len(d.Desc))
	for _, v := range d.Desc {
		cnt := 0
		for _, h := range v.Histogram {
			if h > 0 {
				cnt++
			}
		}
		if cnt > minDir {
			v.Dim = cnt
			pass = append(pass, v)
		}
	}
	d.Desc = pass
}

func (d *DescriptorSet) assignOrientations() {
	for i, v := range d.Desc {
		d.Desc[i].Theta = make([]float64, 0, 36)

		//find max bin in histogram
		var max float64
		for _, j := range v.Histogram {
			if j > max {
				max = j
			}
		}

		//bins (angles) >= 80% max are assigned to desc
		for hIdx, val := range v.Histogram {
			if val >= 0.8 * max {
				d.Desc[i].Theta = append(d.Desc[i].Theta, float64((hIdx * 10) + 5))
			}
		}
	}
}

func (d *DescriptorSet) buildOrientVector() {
	img, weights := d.prepareOrient()

	for i, v := range d.Desc {
		//create histogram of pixel orientations, ea. bin 10 deg range
		mags := make([]float64, 0, 25)
		thetas := make([]float64, 0, 25)
		hist := make([]float64, 36)

		//find mag and theta for each pixel
		xStart := v.Point.X - 2
		yStart := v.Point.Y - 2
		for i := range weights {
			for j := range weights {
				x := xStart + j
				y := yStart + i

				//find magnitude
				magTermA := math.Pow(float64(u16math.SignDiff(img.Gray16At(x + 1, y).Y,
					img.Gray16At(x - 1, y).Y)), 2)
				magTermB := math.Pow(float64(u16math.SignDiff(img.Gray16At(x, y + 1).Y,
					img.Gray16At(x, y - 1).Y)), 2)
				mag := math.Sqrt(magTermA + magTermB)
				mags = append(mags, mag)

				//find theta
				thetaTermA := float64(u16math.SignDiff(img.Gray16At(x, y + 1).Y,
					img.Gray16At(x, y - 1).Y))
				thetaTermB := float64(u16math.SignDiff(img.Gray16At(x + 1, y).Y,
					img.Gray16At(x - 1, y).Y))
				theta := degrees(math.Atan2(thetaTermA, thetaTermB))
				if theta < 0 {
					theta += 360
				}
				thetas = append(thetas, theta)

				//add pixel's value to histogram
				hist[int(theta / 10)] += mag * weights[i][j]
			}
		}

		//discard vectors with less than non-zero bins
		normalize(hist)
		d.Desc[i].Histogram = hist
		d.Desc[i].pxThetas = thetas
		d.Desc[i].pxMags = mags
	}
}

func (d *DescriptorSet) eliminateEdgePoints() {
	pass := make([]Descriptor, 0, len(d.Desc))
	st := NewStructureTensor(conv.ToImg(d.gImg))
	for i, v := range d.Desc {
		stMat := st.At(v.Point.X, v.Point.Y)
		R := math.Pow(stMat.Trace(), 2) / stMat.Det()
		if R <= math.Pow(10 + 1, 2) / 10 {
			d.Desc[i].R = R
			pass = append(pass, d.Desc[i])
		}
	}
	d.Desc = pass
}

func (d *DescriptorSet) eliminateLowContrastPoints() {

	//record the contrast of each descriptor window
	for i, v := range d.Desc {
		//eliminate points out of 5 x 5 range
		if v.Point.X < 3 || v.Point.Y < 3 {continue}
		min := uint16(math.MaxUint16)
		max := uint16(0)
		for x := v.Point.X - 2; x < v.Point.X + 2; x++ {
			for y := v.Point.Y - 2; y < v.Point.Y + 2; y++ {
				val := d.gImg.Gray16At(x, y).Y
				if val < min {
					min = val
				}
				if val > max {
					max = val
				}
			}
		}
		d.Desc[i].contrast = max - min
	}

	//sort desc by contrast value
	sort.Slice(d.Desc, func(i, j int) bool {
		return d.Desc[i].contrast > d.Desc[j].contrast
	})

	//discard low contrast desc
	pass := make([]Descriptor, 0, len(d.Desc))
	cntThr := d.Desc[0].contrast / 7
	for i, v := range d.Desc {
		if v.contrast >= cntThr {
			pass = append(pass, d.Desc[i])
		}
	}
	d.Desc = pass
}

func (d *DescriptorSet) convertPoints() {
	d.Desc = make([]Descriptor, len(d.dog.extrema))
	for i, v := range d.dog.extrema {
		d.Desc[i] = Descriptor{
			Point:    v,
			contrast: 0,
			R:        0,
		}
	}
}

func (d DescriptorSet) prepareOrient() (*image.Gray16, [][]float64) {
	img := ConvolveGray16(d.gImg, GaussSig4) //should be scale of ea. ipt
	weights := GaussSig6
	if len(weights) != 5 || len(weights[0]) != 5 {
		panic("kernel does not fit 5x5 desc window")
	}
	return img, weights
}

func (d DescriptorSet) imageArea(pt Point, size int) *image.Gray16 {
	out := image.NewGray16(image.Rect(0,0, size, size,))
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			val := d.gImg.Gray16At(x + (pt.X - (size / 2)), y + (pt.Y - (size / 2)))
			out.SetGray16(x, y, val)
		}
	}
	return out
}

func (d DescriptorSet) dogAreas(pt Point) []*image.Gray16 {
	outs := make([]*image.Gray16, len(d.dog.scSpace))
	for i := range outs {
		out := image.NewGray16(image.Rect(0,0,5,5,))
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				out.SetGray16(x, y, d.dog.scSpace[i].img.Gray16At(x + (pt.X - 2), y + (pt.Y - 2)))
			}
		}
		edit.ScaleG16(out)
		outs[i] = out
	}
	return outs
}

func (d DescriptorSet) DescInfo(idx int) DescInfo {
	return DescInfo{
		desc:    &(d.Desc[idx]),
		imgPxLg: d.imageArea(d.Desc[idx].Point, 25),
		imgPxSm: d.imageArea(d.Desc[idx].Point, 5),
		dogPx:   d.dogAreas(d.Desc[idx].Point),
	}
}

func (d DescriptorSet) DescInfoAll() []DescInfo {
	out := make([]DescInfo, len(d.Desc))
	for i := range d.Desc {
		out[i] = d.DescInfo(i)
	}
	return out
}

func (d DescriptorSet) DescInfoRange(idc []int) []DescInfo {
	out := make([]DescInfo, len(idc))
	for i, v := range idc {
		out[i] = d.DescInfo(v)
	}
	return out
}

func (d DescriptorSet) ShowKeypoints(nExt int) image.Image {
	if nExt <= 0 {
		nExt = len(d.Desc)
	}
	plot := edit.NewIllustrator(copyR64(d.cImg))
	for _, p := range d.Desc[:nExt] {
		plot.DrawCircle(p.Point.X, p.Point.Y, 5, edit.Red())
	}
	return plot.ToImage()
}