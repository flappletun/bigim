package cv

import (
	"image"
	"bigim/conv"
	"bigim/edit"
	"math"
)

type DescInfo struct {
	desc    *Descriptor
	imgPxLg,
	imgPxSm *image.Gray16
	dogPx   []*image.Gray16
	orntDiag,
	gradDiag *image.RGBA64
}

func (d DescInfo) Show() image.Image {
	if d.orntDiag == nil {
		d.orntDiag = d.makeOrient()
	}
	if d.gradDiag == nil {
		d.makeGradDiag()
	}

	d.imgPxLg = edit.ResizeG16(d.imgPxLg, 100, 100)
	for i := range d.dogPx {
		d.dogPx[i] = edit.ResizeG16(d.dogPx[i], 40, 40)
	}

	//out := image.NewRGBA64(image.Rect(0, 0, 500, 140))
	out := image.NewRGBA64(image.Rect(0, 0, 380, 140))
	edit.FillR64(out, edit.Black())
	edit.PasteR64(out, conv.ToR64(d.imgPxLg), 20, 20)
	edit.PasteR64(out, conv.ToR64(d.orntDiag), 140, 20)
	edit.PasteR64(out, conv.ToR64(d.gradDiag), 260, 20)
	//edit.PasteR64(out, conv.ToR64(d.dogPx[0]), 380, 20)
	//edit.PasteR64(out, conv.ToR64(d.dogPx[1]), 440, 20)
	//edit.PasteR64(out, conv.ToR64(d.dogPx[2]), 380, 80)
	//edit.PasteR64(out, conv.ToR64(d.dogPx[3]), 440, 80)

	return out
}

func (d DescInfo) ShowWinSm() image.Image {
	return conv.ToImg(d.imgPxSm)
}

func (d DescInfo) ShowWinLg() image.Image {
	return conv.ToImg(d.imgPxLg)
}

func (d DescInfo) ShowDogResp() []image.Image {
	out := make([]image.Image, len(d.dogPx))
	for i := range out {
		out[i] = conv.ToImg(d.dogPx[i])
	}
	return out
}

func (d DescInfo) ShowOrient() image.Image {
	if d.orntDiag == nil {
		d.orntDiag = d.makeOrient()
	}
	return conv.ToImg(d.orntDiag)
}

func (d DescInfo) ShowGradient() image.Image {//doesnt work
	if d.gradDiag == nil {
		d.makeGradDiag()
	}
	return conv.ToImg(d.gradDiag)
}

func (d DescInfo) makeOrient() *image.RGBA64 {
	bkg := conv.ToR64(edit.ResizeG16(d.imgPxSm, 100, 100))
	ill := edit.NewIllustrator(bkg)

	for _, v := range d.desc.Theta {
		destX := int(40 * math.Cos(radians(v)))
		destY := int(40 * math.Sin(radians(v)))
		ill.DrawCircle(50 + destX, 50 + destY, 2, edit.Green())
		ill.DrawLine(50, 50 + destX, 50, 50 + destY, edit.Green())
	}

	d.orntDiag = ill.ToR64()
	return d.orntDiag
}

func (d *DescInfo) makeGradDiag() {
	if d.desc == nil || d.desc.pxThetas == nil {
		panic("oh no!")
	}
	bkg := conv.ToR64(edit.ResizeG16(d.imgPxSm, 100, 100))
	ill := edit.NewIllustrator(bkg)
	thIdx := 0
	for y := 10; y < 100; y += 20 {
		for x := 10; x < 100; x += 20 {
			destX := int(10 * math.Cos(radians(d.desc.pxThetas[thIdx])))
			destY := int(10 * math.Sin(radians(d.desc.pxThetas[thIdx])))
			ill.DrawCircle(x, y, 2, edit.Green())
			ill.DrawLine(x, x + destX, y, y + destY, edit.Green())
			thIdx++
		}
	}

	d.gradDiag = ill.ToR64()
}