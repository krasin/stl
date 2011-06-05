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
		[32]uint8{0, 0, 0, 0, 0, 0, 0, 0, 254, 255, 255, 7, 254, 255, 255, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
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
			if !matchString("solid") {
				goto l0
			}
			if !p.rules[14]() {
				goto l0
			}
			{
				position1, thunkPosition1 := position, thunkPosition
				if !p.rules[12]() {
					goto l2
				}
				do(0)
				goto l1
			l2:
				position, thunkPosition = position1, thunkPosition1
			}
		l1:
			if !p.rules[10]() {
				goto l0
			}
		l3:
			{
				position4, thunkPosition4 := position, thunkPosition
				{
					position5, thunkPosition5 := position, thunkPosition
					if !p.rules[14]() {
						goto l5
					}
					goto l6
				l5:
					position, thunkPosition = position5, thunkPosition5
				}
			l6:
				if !matchString("facet") {
					goto l4
				}
				if !p.rules[14]() {
					goto l4
				}
				do(1)
				if !matchString("normal") {
					goto l4
				}
				if !p.rules[14]() {
					goto l4
				}
				if !p.rules[11]() {
					goto l4
				}
				if !p.rules[14]() {
					goto l4
				}
				if !p.rules[11]() {
					goto l4
				}
				if !p.rules[14]() {
					goto l4
				}
				if !p.rules[11]() {
					goto l4
				}
				if !p.rules[10]() {
					goto l4
				}
				{
					position7, thunkPosition7 := position, thunkPosition
					if !p.rules[14]() {
						goto l7
					}
					goto l8
				l7:
					position, thunkPosition = position7, thunkPosition7
				}
			l8:
				if !matchString("outer") {
					goto l4
				}
				if !p.rules[14]() {
					goto l4
				}
				if !matchString("loop") {
					goto l4
				}
				if !p.rules[10]() {
					goto l4
				}
				if !p.rules[6]() {
					goto l4
				}
				if !p.rules[6]() {
					goto l4
				}
				if !p.rules[6]() {
					goto l4
				}
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
				if !matchString("endloop") {
					goto l4
				}
				if !p.rules[10]() {
					goto l4
				}
				{
					position11, thunkPosition11 := position, thunkPosition
					if !p.rules[14]() {
						goto l11
					}
					goto l12
				l11:
					position, thunkPosition = position11, thunkPosition11
				}
			l12:
				if !matchString("endfacet") {
					goto l4
				}
				if !p.rules[10]() {
					goto l4
				}
				goto l3
			l4:
				position, thunkPosition = position4, thunkPosition4
			}
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
			if !matchString("endsolid") {
				goto l0
			}
			if !p.rules[14]() {
				goto l0
			}
			{
				position15, thunkPosition15 := position, thunkPosition
				if !p.rules[12]() {
					goto l16
				}
				do(3)
				goto l15
			l16:
				position, thunkPosition = position15, thunkPosition15
			}
		l15:
			if !p.rules[10]() {
				goto l0
			}
			{
				position17, thunkPosition17 := position, thunkPosition
				if !matchDot() {
					goto l17
				}
				goto l0
			l17:
				position, thunkPosition = position17, thunkPosition17
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
		nil,
		/* 2 Facet <- (FacetHeader OuterLoop Vertex Vertex Vertex EndLoop EndFacet) */
		nil,
		/* 3 FacetHeader <- (Space? 'facet' Space { p.Add() } 'normal' Space Vector NewLine) */
		nil,
		/* 4 Vector <- (Number Space Number Space Number) */
		nil,
		/* 5 OuterLoop <- (Space? 'outer' Space 'loop' NewLine) */
		nil,
		/* 6 Vertex <- (Space? 'vertex' { p.Vertex() } Space Number Space Number Space Number NewLine) */
		func() bool {
			position0, thunkPosition0 := position, thunkPosition
			{
				position24, thunkPosition24 := position, thunkPosition
				if !p.rules[14]() {
					goto l24
				}
				goto l25
			l24:
				position, thunkPosition = position24, thunkPosition24
			}
		l25:
			if !matchString("vertex") {
				goto l23
			}
			do(2)
			if !p.rules[14]() {
				goto l23
			}
			if !p.rules[11]() {
				goto l23
			}
			if !p.rules[14]() {
				goto l23
			}
			if !p.rules[11]() {
				goto l23
			}
			if !p.rules[14]() {
				goto l23
			}
			if !p.rules[11]() {
				goto l23
			}
			if !p.rules[10]() {
				goto l23
			}
			return true
		l23:
			position, thunkPosition = position0, thunkPosition0
			return false
		},
		/* 7 EndLoop <- (Space? 'endloop' NewLine) */
		nil,
		/* 8 EndFacet <- (Space? 'endfacet' NewLine) */
		nil,
		/* 9 EndSolid <- (Space? 'endsolid' Space ((Identifier { p.EndName = buffer[begin:end] }) /) NewLine) */
		nil,
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
		/* 12 Identifier <- (< [a-zA-Z]+ >) */
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
