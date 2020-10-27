package icharts

import (
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

type Point struct {
	X float64
	Y float64
}
