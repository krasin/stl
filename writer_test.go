package stl

import (
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"testing"
)

type writerTest struct {
	t     []Triangle
	model string
}

var writerTests = []writerTest{
	{
		t: []Triangle{
			{
				N: Point{0, 0, 1},
				V: [3]Point{
					{0, 0, 0},
					{1, 0, 0},
					{0, 1, 0},
				},
			},
		},
		model: "data/one_triangle.stl",
	},
}

func TestWriter(t *testing.T) {
	for _, tt := range writerTests {
		buf := new(bytes.Buffer)
		want, err := ioutil.ReadFile(tt.model)
		if err != nil {
			t.Errorf("Could not read test file %s: %v", tt.model, err)
		}
		if err = WriteASCII(buf, tt.t); err != nil {
			t.Errorf("Test %s: %v", tt.model, err)
		}
		if !bytes.Equal(buf.Bytes(), want) {
			t.Errorf("Test %s: unexpected output from writer. Want: %s, got: %s", tt.model, string(want), string(buf.Bytes()))
		}

	}
}

func almostEq(v1, v2 float64) bool {
	return math.Abs(v1-v2) < 1E-6
}

func pointEq(p1, p2 Point) bool {
	return almostEq(p1[0], p2[0]) && almostEq(p1[1], p2[1]) && almostEq(p1[2], p2[2])
}

func equal(t1, t2 Triangle) bool {
	return pointEq(t1.N, t2.N) && pointEq(t1.V[0], t2.V[0]) &&
		pointEq(t1.V[1], t2.V[1]) && pointEq(t1.V[2], t2.V[2])
}

func testWriter(t *testing.T, name string, write func(io.Writer, []Triangle) error) {
	tests := []string{
		"data/cylinder.bin.stl",
		"data/cylinder.stl",
		"data/one_triangle.stl",
		"data/plus_on_pedestal.stl",
	}
	for _, tt := range tests {
		f, err := os.Open(tt)
		if err != nil {
			t.Error(err)
			continue
		}
		defer f.Close()
		tr, err := Read(f)
		if err != nil {
			t.Errorf("Read(%q): %v", tt, err)
			continue
		}
		var buf bytes.Buffer
		if err = write(&buf, tr); err != nil {
			t.Errorf("WriteBinary(%q): %v", tt, err)
			continue
		}
		tr2, err := Read(&buf)
		if err != nil {
			t.Errorf("%q: Read from saved failed: %v", tt, err)
			continue
		}
		if len(tr) != len(tr2) {
			t.Errorf("%q: %d = len(tr) != len(tr2) = %d", tt, len(tr), len(tr2))
			continue
		}
		for i := range tr {
			if !equal(tr[i], tr2[i]) {
				t.Errorf("%q, i=%d, triangles are different. Was: %+v, became: %+v", tt, i, tr[i], tr2[i])
				continue
			}
		}
	}
}

func TestWriteBinary(t *testing.T) {
	testWriter(t, "binary", WriteBinary)
}

func TestWriteASCII(t *testing.T) {
	testWriter(t, "ascii", WriteASCII)
}

func randPoint() Point {
	return Point{rand.Float64(), rand.Float64(), rand.Float64()}
}

func randTriangle() Triangle {
	return Triangle{
		N: randPoint(),
		V: [3]Point{randPoint(), randPoint(), randPoint()},
	}
}

func generateSTL(n int) []Triangle {
	var res = make([]Triangle, n)
	for i := 0; i < n; i++ {
		res[i] = randTriangle()
	}
	return res
}

var randomSTL = generateSTL(1E6)

func BenchmarkWriteBinary(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		buf.Reset()
		if err := WriteBinary(&buf, randomSTL); err != nil {
			b.Fatalf("unexpected error writing a random STL: %v", err)
		}
	}
}
