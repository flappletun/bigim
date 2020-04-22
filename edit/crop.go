package edit

import (
	"image"
	"bigim/conv"
)

func Trim(i image.Image, margin int) image.Image {
	return conv.ToImg(TrimR64(conv.ToR64(i), margin))
}

func TrimR64(r *image.RGBA64, margin int) *image.RGBA64 {
	bounds := image.Rect(margin, margin, r.Bounds().Dx() - margin, r.Bounds().Dy() - margin)
	return conv.ToR64(r.SubImage(bounds))
}

func TrimG16(g *image.Gray16, margin int) *image.Gray16 {
	bounds := image.Rect(margin, margin, g.Bounds().Dx() - margin, g.Bounds().Dy() - margin)
	return conv.ToG16(g.SubImage(bounds))
}