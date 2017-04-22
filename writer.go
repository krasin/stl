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
	f32 := func(dst []byte, f float32) {
		binary.LittleEndian.PutUint32(dst, math.Float32bits(f))
	}
	p := func(dst []byte, p Point) {
		f32(dst[0:4], float32(p[0]))
		f32(dst[4:8], float32(p[1]))
		f32(dst[8:12], float32(p[2]))
	}

	// Write 80 bytes zero header, which is always ignored
	wr(make([]byte, 80))

	// Number of triangles
	bwr(uint32(len(t)))

	for _, tr := range t {
		var cur [4*3*4 + 2]byte
		if err != nil {
			return err
		}
		p(cur[0:12], tr.N)
		p(cur[12:24], tr.V[0])
		p(cur[24:36], tr.V[1])
		p(cur[36:48], tr.V[2])
		wr(cur[:])
	}

	return err
}
