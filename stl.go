package stl

// Point represent a point or vector in 3-dimensional space.
type Point [3]float32

// Triangle consists of a normal vector and 3 points in 3-dimensional space.
type Triangle struct {
	N Point
	V [3]Point
}
