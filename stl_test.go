package stl

import (
	"math"
	"testing"
)

func near(a, b float32) bool {
	return math.Abs(float64(a-b)) < eps
}

func nearPoint(a, b Point) bool {
	return near(a[0], b[0]) && near(a[1], b[1]) && near(a[2], b[2])
}

func TestBoundingBox(t *testing.T) {
	tests := []struct {
		desc string
		t    []Triangle
		min  Point
		max  Point
	}{
		{
			desc: "No triangle, zero box",
		},
		{
			desc: "One horizontal triangle",
			t: []Triangle{
				{
					N: Point{0, 0, 1},
					V: [3]Point{
						{0, 0, 0},
						{0, 1, 0},
						{1, 0, 0},
					},
				},
			},
			min: Point{0, 0, 0},
			max: Point{1, 1, 0},
		},
		{
			desc: "Two triangles",
			t: []Triangle{
				{
					N: Point{0, 0, 1},
					V: [3]Point{
						{0, 0, 2},
						{0, 1, 2},
						{1, 0, 2},
					},
				},
				{
					N: Point{-1, 0, 0},
					V: [3]Point{
						{0, 0, 1},
						{0, 1, 1},
						{0, 2, 2},
					},
				},
			},
			min: Point{0, 0, 1},
			max: Point{1, 2, 2},
		},
	}

	for _, tt := range tests {
		min, max := BoundingBox(tt.t)
		if !nearPoint(min, tt.min) || !nearPoint(max, tt.max) {
			t.Errorf("%q: wrong bounding box.\nWant: %v : %v\nGot:  %v : %v", tt.min, tt.max, min, max)
			continue
		}
	}
}
