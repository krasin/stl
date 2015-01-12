package stl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"unsafe"
)

const SizeOfTriangle = 4*3*4 + 2

func readPoint(a []byte, p *Point) []byte {
	for i := 0; i < 3; i++ {
		cur := uint32(a[0]) + uint32(a[1])<<8 + uint32(a[2])<<16 + uint32(a[3])<<24
		p[i] = *(*float32)(unsafe.Pointer(&cur))
		a = a[4:]
	}
	return a
}

// scanLine reads a non-empty line from scanner. It cleans it from tabs and eliminates double spaces.
func scanLine(sc *bufio.Scanner) (str string, err error) {
	for {
		if !sc.Scan() {
			if err = sc.Err(); err == nil {
				err = io.EOF
			}
			return "", err
		}
		str = sc.Text()
		str = strings.Replace(str, "\t", " ", -1)
		str = strings.Replace(str, "        ", " ", -1)
		str = strings.Replace(str, "    ", " ", -1)
		str = strings.Replace(str, "    ", " ", -1)
		str = strings.Replace(str, "  ", " ", -1)
		str = strings.Replace(str, "  ", " ", -1)
		str = strings.TrimSpace(str)
		if str != "" {
			return str, nil
		}
	}
}

func readLineWithPrefix(sc *bufio.Scanner, prefixes ...string) (prefix, str string, err error) {
	str, err = scanLine(sc)
	if err != nil {
		return "", "", err
	}
	for _, pp := range prefixes {
		if strings.HasPrefix(str, pp) {
			return pp, str[len(pp):], nil
		}
	}
	return "", "", fmt.Errorf("line expected to start with one of the prefixes: %v, the actual line is: '%s'", prefixes, str)
}

func consumeLine(sc *bufio.Scanner, want string) (err error) {
	var str string
	if _, str, err = readLineWithPrefix(sc, want); err != nil {
		return
	}
	if str != "" {
		return fmt.Errorf("line contains unexpected symbols after the right prefix: '%s', symbols: '%s'", want, str)
	}
	return nil
}

func parseFloat32(str string) (float32, error) {
	v, err := strconv.ParseFloat(str, 32)
	if err != nil {
		// Give it a second chance. A very small amount of STL files use a comma as a delimiter
		tmp := strings.Replace(str, ",", ".", 1)
		var err2 error
		if v, err2 = strconv.ParseFloat(tmp, 32); err2 != nil {
			// Return the first error, in order to show the original value.
			return 0, err
		}
		// It has parsed the thing. OK
	}
	return float32(v), nil
}

// scanLines will return lines delimited by CR or LF.
// It might report empty lines, they will be filtered later.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i, v := range data {
		if v == '\n' || v == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	// line is not ended, need more data
	return 0, nil, nil
}

func readASCII(data []byte) (res []Triangle, err error) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(scanLines)
	if _, _, err = readLineWithPrefix(sc, "solid"); err != nil {
		return nil, err
	}
	lineno := 2
	for {
		var prefix, str string
		var t Triangle
		if prefix, str, err = readLineWithPrefix(sc, "facet normal ", "endsolid"); err != nil {
			if err == io.EOF {
				return res, nil
			}
			return nil, fmt.Errorf("[line=%d] %v", lineno, err)
		}
		lineno++

		if prefix == "endsolid" {
			return
		}
		fields := strings.Fields(str)
		if len(fields) != 3 {
			return nil, fmt.Errorf("[line=%d] Normal definition is broken: '%s'", lineno, str)
		}
		for i := 0; i < 3; i++ {
			if t.N[i], err = parseFloat32(fields[i]); err != nil {
				return nil, err
			}
		}
		if err = consumeLine(sc, "outer loop"); err != nil {
			return nil, err
		}
		lineno++
		for i := 0; i < 3; i++ {
			if _, str, err = readLineWithPrefix(sc, "vertex "); err != nil {
				return nil, err
			}
			lineno++

			fields = strings.Fields(str)
			if len(fields) != 3 {
				return nil, fmt.Errorf("[line=%d] Vertex definition is broken: '%s'", lineno, str)
			}
			for j := 0; j < 3; j++ {
				if t.V[i][j], err = parseFloat32(fields[j]); err != nil {
					return nil, err
				}
			}
		}
		if err = consumeLine(sc, "endloop"); err != nil {
			return nil, err
		}
		lineno++

		if _, _, err = readLineWithPrefix(sc, "endfacet"); err != nil {
			return nil, err
		}
		lineno++

		res = append(res, t)

	}
	return
}

// isASCII detects if the data represents an ASCII STL file (as opposed to binary STL file).
func isASCII(data []byte) bool {
	if len(data) < 14 {
		return false
	}
	magic := strings.TrimSpace(string(data[:14]))
	if !strings.HasPrefix(magic, "solid") {
		return false
	}
	for _, v := range data {
		if v < 32 && v != '\n' && v != '\r' && v != '\t' {
			// non-printable chars are a good indicator of a binary STL.
			return false
		}
	}
	return true
}

// Read reads STL file from the reader.
func Read(r io.Reader) (t []Triangle, err error) {
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		return
	}
	if isASCII(data) {
		return readASCII(data)
	}
	if len(data) < 84 {
		return nil, fmt.Errorf("The file is too short: %d bytes", len(data))
	}
	// Skip STL header
	data = data[80:]
	n := uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16 + uint32(data[3])<<24
	data = data[4:]

	if int64(len(data)) < SizeOfTriangle*int64(n) {
		return nil, fmt.Errorf("Read: unexpected end of file: want %d bytes to read triangle data, but only %d bytes is available", SizeOfTriangle*int64(n), len(data))
	}
	for i := 0; i < int(n); i++ {
		var cur Triangle
		data = readPoint(data, &cur.N)
		for j := 0; j < 3; j++ {
			data = readPoint(data, &cur.V[j])
		}
		data = data[2:]
		t = append(t, cur)
	}
	return
}
