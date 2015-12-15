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
	headerBuf := make([]byte, 84)

	// Write triangle count
	binary.LittleEndian.PutUint32(headerBuf[80:84], uint32(len(t)))
	_, errHeader := w.Write(headerBuf)
	if errHeader != nil {
		return errHeader
	}

	// Write each triangle
	for _, t := range t {
		tErr := writeTriangleBinary(w, &t)
		if tErr != nil {
			return tErr
		}
	}

	return nil
}

func writeTriangleBinary(w io.Writer, t *Triangle) error {
	buf := make([]byte, 50)
	offset := 0
	encodePoint(buf, &offset, &t.N)
	encodePoint(buf, &offset, &t.V[0])
	encodePoint(buf, &offset, &t.V[1])
	encodePoint(buf, &offset, &t.V[2])
	encodeUint16(buf, &offset, 0)
	_, err := w.Write(buf)
	return err
}

func encodePoint(buf []byte, offset *int, pt *Point) {
	encodeFloat32(buf, offset, pt[0])
	encodeFloat32(buf, offset, pt[1])
	encodeFloat32(buf, offset, pt[2])
}

func encodeFloat32(buf []byte, offset *int, f float32) {
	u32 := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf[*offset:(*offset)+4], u32)
	(*offset) += 4
}

func encodeUint16(buf []byte, offset *int, u uint16) {
	binary.LittleEndian.PutUint16(buf[*offset:(*offset)+2], u)
	(*offset) += 2
}
