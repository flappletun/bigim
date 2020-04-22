package cv

import (
	"strconv"
	"strings"
)

type Point struct {
	X,
	Y int
}

func (p Point) Hash() string {
	return strconv.Itoa(p.X) + " " + strconv.Itoa(p.Y)
}

func hashPoint(x, y int) string {
	return strconv.Itoa(x) + " " + strconv.Itoa(y)
}

func unhashPoint(s string) Point {
	p := Point{}
	fl := strings.Fields(s)
	p.X, _ = strconv.Atoi(fl[0])
	p.Y, _ = strconv.Atoi(fl[1])
	return p
}