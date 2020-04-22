package imio

import (
	"image"
	"image/jpeg"
	"image/png"
	"bigim/conv"
	"io"
	"os"
	"path/filepath"
	"strings"
)

  /********/
 /* Open */
/********/

func Open(path string) image.Image {

	//open the file
	file, err := os.Open(path)
	if err != nil {
		panic("error opening file " + path)
		return nil
	}
	defer file.Close()

	//find image format and decode image
	return decodeImage(file, path)
}

func OpenR64(path string) *image.RGBA64 {
	return conv.ToR64(Open(path))
}

func OpenG16(path string) *image.Gray16 {
	return conv.ToG16(Open(path))
}

func OpenDir(root string) ([]image.Image, []string) {
	ims := make([]image.Image, 0, 1000)
	nms := make([]string, 0, 1000)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !isImage(path) {return nil}
		base := filepath.Base(path)
		nms = append(nms, strings.TrimSuffix(base, filepath.Ext(base)))
		ims = append(ims, Open(path))
		return nil
	})
	return ims, nms
}

func OpenDirR64(root string) ([]*image.RGBA64, []string) {
	ims, nms := OpenDir(root)
	r64s := make([]*image.RGBA64, len(ims))
	for i := range ims {
		r64s[i] = conv.ToR64(ims[i])
	}
	return r64s, nms
}

func OpenDirG16(root string) ([]*image.Gray16, []string) {
	ims, nms := OpenDir(root)
	g16s := make([]*image.Gray16, len(ims))
	for i := range ims {
		g16s[i] = conv.ToG16(ims[i])
	}
	return g16s, nms
}

  /********/
 /* Save */
/********/

func Save(n interface{}, name string) {
	switch n.(type) {
	case image.Image: saveImage(n.(image.Image), name)
	case *image.RGBA64: saveR64(n.(*image.RGBA64), name)
	case *image.Gray16: saveGray16(n.(*image.Gray16), name)
	}
}

func saveImage(i image.Image, name string) {

	//complete file path with desired ext
	switch filepath.Ext(name) {
	case ".jpg": name += ".jpg"
	case ".png": name += ".png"
	case "": name += ".png"
	default: panic("extension must be .jpg, .png or empty (.png)")
	}

	//create new file
	file, err := os.Create(name)
	if err != nil {
		panic("error ceating file")
	}
	defer file.Close()

	//encode image
	encodeImage(file, i, name)
}

func saveR64(r *image.RGBA64, name string) {
	saveImage(conv.ToImg(r), name)
}

func saveGray16(g *image.Gray16, name string) {
	saveImage(conv.ToImg(g), name)
}

func encodeImage(file io.Writer, i image.Image, name string) {
	var err error

	//write image data to file
	switch filepath.Ext(name) {
	case ".jpg": {
		err = jpeg.Encode(file, i, nil)
		if err != nil {
			panic("error encoding jpg file")
		}
	}
	default: {
		err = png.Encode(file, i)
		if err != nil {
			panic("error encoding png file")
		}
	}
	}
}

func decodeImage(file io.Reader, path string) image.Image {
	var img image.Image
	var err error

	switch filepath.Ext(path) {
	case ".jpg": {
		img, err = jpeg.Decode(file)
		if err != nil {
			panic("error decoding jpeg data for file " + path)
		}
	}
	case ".jpeg": {
		img, err = jpeg.Decode(file)
		if err != nil {
			panic("error decoding jpeg data for file " + path)
		}
	}
	case ".png": {
		img, err = png.Decode(file)
		if err != nil {
			panic("error decoding png data for file " + path)
		}
	}
	default:
		panic("file " + path + " is not in an accepted image format")
	}
	return img
}

  /*********/
 /* Misc. */
/*********/

func NameFromPath(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func isImage(path string) bool {
	if ext := filepath.Ext(path); ext == ".png" ||
		ext == ".jpg" || ext == ".jpeg" {
		return true
	}
	return false
}