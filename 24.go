package main

import (
	"bufio"
	"fmt"
	"log"
)

func init() {
	addSolutions(24, problem24)
}

func problem24(ctx *problemContext) {
	var m [][]byte
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		m = append(m, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	g := newBugGrid(m)
	seen := map[uint64]struct{}{g.hash(): struct{}{}}
	for {
		g.step()
		h := g.hash()
		if _, ok := seen[h]; ok {
			ctx.reportPart1(g.biodiv())
			break
		}
		seen[h] = struct{}{}
	}

	gr := newBugGridRec(m)
	for i := 0; i < 200; i++ {
		gr.step()
	}
	ctx.reportPart2(gr.count())
}

type bugGrid struct {
	m       [][]byte
	w       int
	h       int
	scratch [][]byte
}

func newBugGrid(m [][]byte) *bugGrid {
	w := len(m[0])
	h := len(m)
	scratch := make([][]byte, h)
	for y := range scratch {
		scratch[y] = make([]byte, w)
	}
	g := &bugGrid{
		w:       w,
		h:       h,
		scratch: scratch,
	}
	for _, row := range m {
		g.m = append(g.m, append([]byte(nil), row...))
	}
	return g
}

func (g *bugGrid) print() {
	for _, row := range g.m {
		fmt.Println(string(row))
	}
}

func (g *bugGrid) hash() uint64 {
	var h uint64
	var i uint
	for _, row := range g.m {
		for _, c := range row {
			if c == '#' {
				h |= (1 << i)
			}
			i++
		}
	}
	return h
}

func (g *bugGrid) step() {
	for y, row := range g.m {
		for x, c := range row {
			adj := g.adjBugs(x, y)
			switch {
			case c == '#' && adj != 1:
				c = '.'
			case c == '.' && (adj == 1 || adj == 2):
				c = '#'
			}
			g.scratch[y][x] = c
		}
	}
	g.scratch, g.m = g.m, g.scratch
}

func (g *bugGrid) adjBugs(x, y int) int {
	var bugs int
	if y > 0 && g.m[y-1][x] == '#' {
		bugs++
	}
	if y < g.h-1 && g.m[y+1][x] == '#' {
		bugs++
	}
	if x > 0 && g.m[y][x-1] == '#' {
		bugs++
	}
	if x < g.w-1 && g.m[y][x+1] == '#' {
		bugs++
	}
	return bugs
}

func (g *bugGrid) biodiv() int64 {
	n := int64(1)
	var score int64
	for _, row := range g.m {
		for _, c := range row {
			if c == '#' {
				score += n
			}
			n *= 2
		}
	}
	return score
}

type gridCoord struct {
	depth int16
	x     int8
	y     int8
}

type bugGridRec struct {
	bugs map[gridCoord]struct{}
}

func newBugGridRec(b [][]byte) *bugGridRec {
	g := &bugGridRec{bugs: make(map[gridCoord]struct{})}
	for y, row := range b {
		for x, c := range row {
			if x == 2 && y == 2 {
				continue
			}
			if c == '#' {
				coord := gridCoord{x: int8(x), y: int8(y)}
				g.bugs[coord] = struct{}{}
			}
		}
	}
	return g
}

func (g *bugGridRec) count() int {
	return len(g.bugs)
}

func (g *bugGridRec) step() {
	m := make(map[gridCoord]struct{})
	visited := make(map[gridCoord]struct{})
	visit := func(c gridCoord) {
		if _, ok := visited[c]; ok {
			return
		}
		visited[c] = struct{}{}
		var adj int
		for _, n := range c.neighbors() {
			if _, ok := g.bugs[n]; ok {
				adj++
			}
		}
		switch adj {
		case 1:
			m[c] = struct{}{}
		case 2:
			if _, ok := g.bugs[c]; !ok {
				m[c] = struct{}{}
			}
		}
	}
	for c := range g.bugs {
		visit(c)
		for _, n := range c.neighbors() {
			visit(n)
		}
	}
	g.bugs = m
}

func (g *bugGridRec) print(depth int) {
	for y := int8(0); y < 5; y++ {
		row := make([]byte, 5)
		for x := int8(0); x < 5; x++ {
			c := byte('.')
			if x == 2 && y == 2 {
				c = '?'
			}
			if _, ok := g.bugs[gridCoord{x: x, y: y, depth: int16(depth)}]; ok {
				c = '#'
			}
			row[x] = c
		}
		fmt.Println(string(row))
	}
}

func (c gridCoord) neighbors() []gridCoord {
	ns := make([]gridCoord, 0, 8)

	// same depth
	if c.y > 0 && !(c.x == 2 && c.y == 3) {
		ns = append(ns, gridCoord{x: c.x, y: c.y - 1, depth: c.depth})
	}
	if c.y < 4 && !(c.x == 2 && c.y == 1) {
		ns = append(ns, gridCoord{x: c.x, y: c.y + 1, depth: c.depth})
	}
	if c.x > 0 && !(c.x == 3 && c.y == 2) {
		ns = append(ns, gridCoord{x: c.x - 1, y: c.y, depth: c.depth})
	}
	if c.x < 4 && !(c.x == 1 && c.y == 2) {
		ns = append(ns, gridCoord{x: c.x + 1, y: c.y, depth: c.depth})
	}

	// depth-1
	if c.y == 0 {
		ns = append(ns, gridCoord{x: 2, y: 1, depth: c.depth - 1})
	}
	if c.y == 4 {
		ns = append(ns, gridCoord{x: 2, y: 3, depth: c.depth - 1})
	}
	if c.x == 0 {
		ns = append(ns, gridCoord{x: 1, y: 2, depth: c.depth - 1})
	}
	if c.x == 4 {
		ns = append(ns, gridCoord{x: 3, y: 2, depth: c.depth - 1})
	}

	// depth+1
	if c.x == 1 && c.y == 2 {
		for y := int8(0); y < 5; y++ {
			ns = append(ns, gridCoord{x: 0, y: y, depth: c.depth + 1})
		}
	}
	if c.x == 3 && c.y == 2 {
		for y := int8(0); y < 5; y++ {
			ns = append(ns, gridCoord{x: 4, y: y, depth: c.depth + 1})
		}
	}
	if c.x == 2 && c.y == 1 {
		for x := int8(0); x < 5; x++ {
			ns = append(ns, gridCoord{x: x, y: 0, depth: c.depth + 1})
		}
	}
	if c.x == 2 && c.y == 3 {
		for x := int8(0); x < 5; x++ {
			ns = append(ns, gridCoord{x: x, y: 4, depth: c.depth + 1})
		}
	}

	for _, n := range ns {
		if n.x == 2 && n.y == 2 {
			fmt.Println(ns)
			fmt.Println(c)
			panic("HI")
		}
	}
	return ns
}
