package cv

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"image"
	"bigim/conv"
)

type StructureTensor struct {
	im,
	dx,
	dy *image.Gray16
}

type StrTnMat struct {
	TopLeft,
	TopRight,
	BtmLeft,
	BtmRight float64
}

func NewStructureTensor(im image.Image) StructureTensor {
	g := conv.ToG16(im)
	dx := GradientDx(g)
	dy := GradientDy(g)
	return StructureTensor{
		im: g,
		dx: dx,
		dy: dy,
	}
}

func (s StructureTensor) At(x, y int) StrTnMat {
	return StrTnMat{
		TopLeft:  float64(s.dx.Gray16At(x, y).Y * s.dx.Gray16At(x, y).Y),
		TopRight: float64(s.dx.Gray16At(x, y).Y * s.dy.Gray16At(x, y).Y),
		BtmLeft:  float64(s.dx.Gray16At(x, y).Y * s.dy.Gray16At(x, y).Y),
		BtmRight: float64(s.dy.Gray16At(x, y).Y * s.dy.Gray16At(x, y).Y),
	}
}

func (s StrTnMat) Det() float64 {
	return mat.Det(s.toMatrix())
}

func (s StrTnMat) Trace() float64 {
	return mat.Trace(s.toMatrix())
}

func (s StrTnMat) Print() {
	//code from https://medium.com/wireless-registry-engineering/gonum-tutorial-linear-algebra-in-go-21ef136fc2d7
	fa := mat.Formatted(s.toMatrix(), mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}

func (s StrTnMat) toMatrix() *mat.Dense {
	return mat.NewDense(2, 2, []float64{
		s.TopLeft,
		s.TopRight,
		s.BtmLeft,
		s.BtmRight,
	})
}