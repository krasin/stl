package stl

import (
	"errors"
	"fmt"
	"math"
	"os"
	"testing"
)

const eps = 1E-4

type readTest struct {
	filename string
	count    int
	t        []Triangle
	err      error
}

var readTests = []readTest{
	{"data/cylinder.bin.stl", 326, nil, nil},
	{"data/plus_on_pedestal.stl", 1180, nil, nil},
	{"data/reg_test_1.stl", 1, nil, nil},
	{"data/reg_test_2.stl", 326, nil, nil},
	{"data/reg_test_3.stl", 1, nil, nil},
	{"data/reg_test_4_empty_solid_name.stl", 1, nil, nil},
	{"data/reg_test_5_tab_after_normal.stl", 1, nil, nil},
	{"data/reg_test_6_space_before_solid.stl", 1, nil, nil},
	{"data/reg_test_7_endfacet_junk.stl", 1, nil, nil},
	{"data/reg_test_8_comma.stl", 1, nil, nil},
	{
		filename: "data/reg_test_10_bin_overflow.stl",
		err:      errors.New("Read: unexpected end of file: want 4294979500 bytes to read triangle data, but only 16300 bytes is available"),
	},
	{
		"data/cylinder.stl",
		326,
		[]Triangle{
			{
				[3]float32{0, 0, 0},
				[3]Point{
					{-7.708244e-01, -3.846672e+00, 5.378669e+00},
					{-1.548386e+00, -3.683723e+00, 4.516774e+00},
					{-1.530743e+00, -3.695526e+00, 4.194193e+00},
				},
			},
		},
		nil,
	},
}

func CoordEqual(a, b float32) bool {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		return false
	}
	return math.Abs(float64(a)-float64(b)) < eps
}

func PointEqual(p1, p2 [3]float32) bool {
	for i := 0; i < 3; i++ {
		if !CoordEqual(p1[i], p2[i]) {
			return false
		}
	}
	return true
}

func Equal(t1, t2 Triangle) bool {
	return PointEqual(t1.N, t2.N) &&
		PointEqual(t1.V[0], t2.V[0]) &&
		PointEqual(t1.V[1], t2.V[1]) &&
		PointEqual(t1.V[2], t2.V[2])
}

func TestRead(t *testing.T) {
	for _, tt := range readTests {
		f, err := os.Open(tt.filename)
		if err != nil {
			t.Fatalf("os.Open(\"%v\"): %v", tt.filename, err)
		}
		defer f.Close()
		stl, err := Read(f)
		if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
			t.Errorf("stl.Read(%q): %v\nwant err: %v", tt.filename, err, tt.err)
			continue
		}
		if err != nil {
			continue
		}
		if len(stl) != tt.count {
			t.Fatalf("Wrong number of triangles. Expected: %d, got: %d", tt.count, len(stl))
		}
		for i, tr := range tt.t {
			if !Equal(tr, stl[i]) {
				t.Fatalf("Triangle #%d, want: %v, got: %v", i, tr, stl[i])
			}
		}
	}
}
