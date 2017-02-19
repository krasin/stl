package stl

import (
	"math"
)

// Point represent a point or vector in 3-dimensional space.
type Point [3]float64

// Triangle consists of a normal vector and 3 points in 3-dimensional space.
type Triangle struct {
	N Point
	V [3]Point
}

type point32 [3]float32
type triangle32 struct {
	n point32
	v [3]point32
}

func point32ToPoint(p *point32) Point {
	return Point{float64(p[0]), float64(p[1]), float64(p[2])}
}

func triangle32ToTriangle(t *triangle32) Triangle {
	return Triangle{
		N: point32ToPoint(&t.n),
		V: [3]Point{point32ToPoint(&t.v[0]), point32ToPoint(&t.v[1]), point32ToPoint(&t.v[2])},
	}
}

// BoundingBox find a minimum cube that wraps the model.
func BoundingBox(t []Triangle) (min, max Point) {
	if len(t) == 0 {
		return
	}
	min = Point{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max = Point{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	for _, tr := range t {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if min[j] > tr.V[i][j] {
					min[j] = tr.V[i][j]
				}
				if max[j] < tr.V[i][j] {
					max[j] = tr.V[i][j]
				}
			}
		}
	}
	return
}
