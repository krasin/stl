package stl

import (
	"fmt"
	"io"
)

func Write(w io.Writer, t []Triangle) error {
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
