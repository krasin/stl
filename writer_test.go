package stl

import (
	"bytes"
	"io/ioutil"
	"testing"
)

type writerTest struct {
	t     []STLTriangle
	model string
}

var writerTests = []writerTest{
	{
		t: []STLTriangle{
			{
				N: STLPoint{0, 0, 1},
				V: [3]STLPoint{
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
		if err = Write(buf, tt.t); err != nil {
			t.Errorf("Test %s: %v", tt.model, err)
		}
		if !bytes.Equal(buf.Bytes(), want) {
			t.Errorf("Test %s: unexpected output from writer. Want: %s, got: %s", tt.model, string(want), string(buf.Bytes()))
		}

	}
}
