package stl

import "math"

// Point represent a point or vector in 3-dimensional space.
type Point [3]float32

// Triangle consists of a normal vector and 3 points in 3-dimensional space.
type Triangle struct {
	N Point
	V [3]Point
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
