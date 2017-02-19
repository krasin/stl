package stl

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// WriteASCII writes the triangle mesh to the writer using ASCII STL codec.
func WriteASCII(w io.Writer, t []Triangle) error {
	var err error

	printf := func(format string, a ...interface{}) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, format, a...)
	}
	printf("solid object\n")
	for _, tt := range t {
		if err != nil {
			return err
		}
		printf("facet normal %f %f %f\n", tt.N[0], tt.N[1], tt.N[2])
		printf("  outer loop\n")
		for _, v := range tt.V {
			printf("    vertex %f %f %f\n", v[0], v[1], v[2])
		}
		printf("  endloop\n")
		printf("endfacet\n")
	}
	printf("endsolid object\n")
	return nil
}

// Write writes the triangle mesh to the writer using binary STL codec.
func WriteBinary(w io.Writer, t []Triangle) error {
	var err error
	twoZeroes := make([]byte, 2)
	wr := func(data []byte) {
		if err != nil {
			return
		}
		_, err = w.Write(data)
	}
	bwr := func(v interface{}) {
		if err != nil {
			return
		}
		err = binary.Write(w, binary.LittleEndian, v)
	}
	f32w := func(f float32) {
		if err != nil {
			return
		}
		err = binary.Write(w, binary.LittleEndian, math.Float32bits(f))
	}
	pwr := func(p Point) {
		if err != nil {
			return
		}
		f32w(float32(p[0]))
		f32w(float32(p[1]))
		f32w(float32(p[2]))
	}

	// Write 80 bytes zero header, which is always ignored
	wr(make([]byte, 80))

	// Number of triangles
	bwr(uint32(len(t)))

	for _, tr := range t {
		if err != nil {
			return err
		}
		pwr(tr.N)
		pwr(tr.V[0])
		pwr(tr.V[1])
		pwr(tr.V[2])
		wr(twoZeroes)
	}

	return err
}
