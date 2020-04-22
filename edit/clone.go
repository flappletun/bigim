package edit

import "image"

func CloneR64(r *image.RGBA64) *image.RGBA64 {
	ret := *r
	return &ret
}

func CloneG16(g *image.Gray16) *image.Gray16 {
	ret := *g
	return &ret
}

func CloneImg(i image.Image) image.Image {
	return i
}