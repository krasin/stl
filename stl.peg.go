package stl

import (
	"bytes"
	"fmt"
	"os"
)

type Stl struct {
	StlFile

	Buffer   string
	Min, Max int
	rules    [15]func() bool
}

type parseError struct {
	p *Stl
}

func (p *Stl) Parse() os.Error {
	if p.rules[0]() {
		return nil
	}
	return &parseError{p}
}

func (e *parseError) String() string {
	buf := new(bytes.Buffer)
	line := 1
	character := 0
	for i, c := range e.p.Buffer[0:] {
		if c == '\n' {
			line++
			character = 0
		} else {
			character++
		}
		if i == e.p.Min {
			if e.p.Min != e.p.Max {
				fmt.Fprintf(buf, "parse error after line %v character %v\n", line, character)
			} else {
				break
			}
		} else if i == e.p.Max {
			break
		}
	}
	fmt.Fprintf(buf, "parse error: unexpected ")
	if e.p.Max >= len(e.p.Buffer) {
		fmt.Fprintf(buf, "end of file found\n")
	} else {
		fmt.Fprintf(buf, "'%c' at line %v character %v\n", e.p.Buffer[e.p.Max], line, character)
	}
	return buf.String()
}
func (p *Stl) Init() {
	var position int
	actions := [...]func(buffer string, begin, end int){
		/* 0 Header */
		func(buffer string, begin, end int) {
			p.Name = buffer[begin:end]
		},
		/* 1 FacetHeader */
		func(buffer string, begin, end int) {
			p.Add()
		},
		/* 2 Vertex */
		func(buffer string, begin, end int) {
			p.Vertex()
		},
		/* 3 EndSolid */
		func(buffer string, begin, end int) {
			p.EndName = buffer[begin:end]
		},
		/* 4 Number */
		func(buffer string, begin, end int) {
			p.Num(buffer[begin:end])
		},
	}
	var thunkPosition, begin, end int
	thunks := make([]struct {
		action     uint8
		begin, end int
	}, 32)
	do := func(action uint8) {
		if thunkPosition == len(thunks) {
			newThunks := make([]struct {
				action     uint8
				begin, end int
			}, 2*len(thunks))
			copy(newThunks, thunks)
			thunks = newThunks
		}
		thunks[thunkPosition].action = action
		thunks[thunkPosition].begin = begin
		thunks[thunkPosition].end = end
		thunkPosition++
	}
	commit := func(thunkPosition0 int) bool {
		if thunkPosition0 == 0 {
			for thunk := 0; thunk < thunkPosition; thunk++ {
				actions[thunks[thunk].action](p.Buffer, thunks[thunk].begin, thunks[thunk].end)
			}
			p.Min = position
			thunkPosition = 0
			return true
		}
		return false
	}
	matchDot := func() bool {
		if position < len(p.Buffer) {
			position++
			return true
		} else if position >= p.Max {
			p.Max = position
		}
		return false
	}
	matchChar := func(c byte) bool {
		if (position < len(p.Buffer)) && (p.Buffer[position] == c) {
			position++
			return true
		} else if position >= p.Max {
			p.Max = position
		}
		return false
	}
	matchString := func(s string) bool {
		length := len(s)
		next := position + length
		if (next <= len(p.Buffer)) && (p.Buffer[position:next] == s) {
			position = next
			return true
		} else if position >= p.Max {
			p.Max = position
		}
		return false
	}
	classes := [...][32]uint8{
		[32]uint8{0, 0, 0, 0, 0, 248, 255, 255, 63, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[32]uint8{0, 0, 0, 0, 0, 0, 255, 3, 254, 255, 255, 135, 254, 255, 255, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	matchClass := func(class uint) bool {
		if (position < len(p.Buffer)) &&
			((classes[class][p.Buffer[position]>>3] & (1 << (p.Buffer[position] & 7))) != 0) {
			position++
			return true
		} else if position >= p.Max {
			p.Max = position
		}
		return false
	}
	p.rules = [...]func() bool{
		/* 0 e <- (Header Facet* EndSolid !. commit) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			if !p.rules[1]() {
				goto l0
			}
		l1:
			{
				position2, thunkPosition2 := position, thunkPosition
				if !p.rules[2]() {
					goto l2
				}
				goto l1
			l2:
				position, thunkPosition = position2, thunkPosition2
			}
			if !p.rules[9]() {
				goto l0
			}
			{
				position3, thunkPosition3 := position, thunkPosition
				if !matchDot() {
					goto l3
				}
				goto l0
			l3:
				position, thunkPosition = position3, thunkPosition3
			}
			if !(commit(thunkPosition0)) {
				goto l0
			}
			return true
		l0:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 1 Header <- ('solid' Space ((Identifier { p.Name = buffer[begin:end] }) /) NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			if !matchString("solid") {
				goto l4
			}
			if !p.rules[14]() {
				goto l4
			}
			{
				position5, thunkPosition5 := position, thunkPosition
				if !p.rules[12]() {
					goto l6
				}
				do(0)
				goto l5
			l6:
				position, thunkPosition = position5, thunkPosition5
			}
		l5:
			if !p.rules[10]() {
				goto l4
			}
			return true
		l4:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 2 Facet <- (FacetHeader OuterLoop Vertex Vertex Vertex EndLoop EndFacet) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			if !p.rules[3]() {
				goto l7
			}
			if !p.rules[5]() {
				goto l7
			}
			if !p.rules[6]() {
				goto l7
			}
			if !p.rules[6]() {
				goto l7
			}
			if !p.rules[6]() {
				goto l7
			}
			if !p.rules[7]() {
				goto l7
			}
			if !p.rules[8]() {
				goto l7
			}
			return true
		l7:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 3 FacetHeader <- (Space? 'facet' Space { p.Add() } 'normal' Space Vector NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position9, thunkPosition9 := position, thunkPosition
				if !p.rules[14]() {
					goto l9
				}
				goto l10
			l9:
				position, thunkPosition = position9, thunkPosition9
			}
		l10:
			if !matchString("facet") {
				goto l8
			}
			if !p.rules[14]() {
				goto l8
			}
			do(1)
			if !matchString("normal") {
				goto l8
			}
			if !p.rules[14]() {
				goto l8
			}
			if !p.rules[4]() {
				goto l8
			}
			if !p.rules[10]() {
				goto l8
			}
			return true
		l8:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 4 Vector <- (Number Space Number Space Number) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			if !p.rules[11]() {
				goto l11
			}
			if !p.rules[14]() {
				goto l11
			}
			if !p.rules[11]() {
				goto l11
			}
			if !p.rules[14]() {
				goto l11
			}
			if !p.rules[11]() {
				goto l11
			}
			return true
		l11:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 5 OuterLoop <- (Space? 'outer' Space 'loop' NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position13, thunkPosition13 := position, thunkPosition
				if !p.rules[14]() {
					goto l13
				}
				goto l14
			l13:
				position, thunkPosition = position13, thunkPosition13
			}
		l14:
			if !matchString("outer") {
				goto l12
			}
			if !p.rules[14]() {
				goto l12
			}
			if !matchString("loop") {
				goto l12
			}
			if !p.rules[10]() {
				goto l12
			}
			return true
		l12:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 6 Vertex <- (Space? 'vertex' { p.Vertex() } Space Number Space Number Space Number NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position16, thunkPosition16 := position, thunkPosition
				if !p.rules[14]() {
					goto l16
				}
				goto l17
			l16:
				position, thunkPosition = position16, thunkPosition16
			}
		l17:
			if !matchString("vertex") {
				goto l15
			}
			do(2)
			if !p.rules[14]() {
				goto l15
			}
			if !p.rules[11]() {
				goto l15
			}
			if !p.rules[14]() {
				goto l15
			}
			if !p.rules[11]() {
				goto l15
			}
			if !p.rules[14]() {
				goto l15
			}
			if !p.rules[11]() {
				goto l15
			}
			if !p.rules[10]() {
				goto l15
			}
			return true
		l15:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 7 EndLoop <- (Space? 'endloop' NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position19, thunkPosition19 := position, thunkPosition
				if !p.rules[14]() {
					goto l19
				}
				goto l20
			l19:
				position, thunkPosition = position19, thunkPosition19
			}
		l20:
			if !matchString("endloop") {
				goto l18
			}
			if !p.rules[10]() {
				goto l18
			}
			return true
		l18:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 8 EndFacet <- (Space? 'endfacet' NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position22, thunkPosition22 := position, thunkPosition
				if !p.rules[14]() {
					goto l22
				}
				goto l23
			l22:
				position, thunkPosition = position22, thunkPosition22
			}
		l23:
			if !matchString("endfacet") {
				goto l21
			}
			if !p.rules[10]() {
				goto l21
			}
			return true
		l21:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 9 EndSolid <- (Space? 'endsolid' Space ((Identifier { p.EndName = buffer[begin:end] }) /) NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position25, thunkPosition25 := position, thunkPosition
				if !p.rules[14]() {
					goto l25
				}
				goto l26
			l25:
				position, thunkPosition = position25, thunkPosition25
			}
		l26:
			if !matchString("endsolid") {
				goto l24
			}
			if !p.rules[14]() {
				goto l24
			}
			{
				position27, thunkPosition27 := position, thunkPosition
				if !p.rules[12]() {
					goto l28
				}
				do(3)
				goto l27
			l28:
				position, thunkPosition = position27, thunkPosition27
			}
		l27:
			if !p.rules[10]() {
				goto l24
			}
			return true
		l24:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 10 NewLine <- (space* ('\n' / '\r' / '\r\n')) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
		l30:
			{
				position31, thunkPosition31 := position, thunkPosition
				if !p.rules[13]() {
					goto l31
				}
				goto l30
			l31:
				position, thunkPosition = position31, thunkPosition31
			}
			{
				position32, thunkPosition32 := position, thunkPosition
				if !matchChar('\n') {
					goto l33
				}
				goto l32
			l33:
				position, thunkPosition = position32, thunkPosition32
				if !matchChar('\r') {
					goto l34
				}
				goto l32
			l34:
				position, thunkPosition = position32, thunkPosition32
				if !matchString("\r\n") {
					goto l29
				}
			}
		l32:
			return true
		l29:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 11 Number <- (< [0-9.+-Ee]+ > { p.Num(buffer[begin:end]) }) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			begin = position
			if !matchClass(0) {
				goto l35
			}
		l36:
			{
				position37, thunkPosition37 := position, thunkPosition
				if !matchClass(0) {
					goto l37
				}
				goto l36
			l37:
				position, thunkPosition = position37, thunkPosition37
			}
			end = position
			do(4)
			return true
		l35:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 12 Identifier <- (< [a-zA-Z_0-9]+ >) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			begin = position
			if !matchClass(1) {
				goto l38
			}
		l39:
			{
				position40, thunkPosition40 := position, thunkPosition
				if !matchClass(1) {
					goto l40
				}
				goto l39
			l40:
				position, thunkPosition = position40, thunkPosition40
			}
			end = position
			return true
		l38:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 13 space <- ' ' */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			if !matchChar(' ') {
				goto l41
			}
			return true
		l41:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 14 Space <- space+ */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			if !p.rules[13]() {
				goto l42
			}
		l43:
			{
				position44, thunkPosition44 := position, thunkPosition
				if !p.rules[13]() {
					goto l44
				}
				goto l43
			l44:
				position, thunkPosition = position44, thunkPosition44
			}
			return true
		l42:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
	}
}
