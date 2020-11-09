package icharts

import (
	"github.com/gotk3/gotk3/cairo"
	"image/color"
	"math"
)

func Normalise(points []DataPoint) []DataPoint {
	if points == nil || len(points) < 1 {
		return nil
	}
	max := points[0].Val
	min := points[0].Val
	for _, point := range points {
		if max < point.Val {
			max = point.Val
		}
		if min > point.Val {
			min = point.Val
		}
	}
	dist := max - min
	if dist == 0 {
		dist = 0.5
	}
	normalised := make([]DataPoint, len(points))
	for i, point := range points {
		normalised[i] = point
		normalised[i].Val = point.Val / dist
	}
	return normalised
}

func Dist(a Point, b Point) float64 {
	first := math.Pow(float64(b.X-a.X), 2)
	second := math.Pow(float64(b.Y-a.Y), 2)
	return math.Sqrt(first + second)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type Point struct {
	X float64
	Y float64
}

func CairoSetColor(cr *cairo.Context, col color.Color) {
	//cr.SetSourceRGB()
	r, g, b, a := col.RGBA()
	cr.SetSourceRGBA(float64(r)/255.0, float64(g)/255.0,
		float64(b)/255.0, float64(a)/255.0)
}

var BrightColors = []color.Color{
	color.RGBA{246, 0, 0, 255},
	color.RGBA{255, 140, 0, 255},
	color.RGBA{255, 238, 0, 255},
	color.RGBA{77, 233, 76, 255},
	color.RGBA{55, 131, 255, 255},
	color.RGBA{72, 21, 170, 255},
}

func MoveOriginTo(oOrigin, nOrigin Point, pt Point) Point {
	dx := oOrigin.X - nOrigin.X
	dy := oOrigin.Y - nOrigin.Y
	return Point{pt.X + dx, pt.Y + dy}
}

func Polar(pt Point) (float64, float64) {
	r := Dist(Point{0, 0}, pt)
	angle := math.Atan2(pt.Y, pt.X)
	if angle < 0 {
		angle = (-1.0 * angle) + math.Pi
	}
	return r, angle
}

func Cartesian(r, angle float64) Point {
	pt := Point{}
	pt.X = r * math.Cos(angle)
	pt.Y = r * math.Sin(angle)
	return pt
}
