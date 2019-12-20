package main

import (
	"bufio"
	"log"
)

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	var maze plutoMaze
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		maze.addRow([]byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	maze.finalize()
	ctx.reportLoad()

	ctx.reportPart1(maze.shortestPath())
	ctx.reportPart2(maze.shortestPathRecursive())
}

type plutoMaze struct {
	m            [][]byte
	w            int64
	h            int64
	start        ivec2
	end          ivec2
	innerPortals map[ivec2]ivec2
	outerPortals map[ivec2]ivec2
}

func (m *plutoMaze) addRow(row []byte) {
	if w := int64(len(row)); w > m.w {
		m.w = w
	}
	m.m = append(m.m, row)
	m.h++
}

func (m *plutoMaze) finalize() {
	for y, row := range m.m {
		for int64(len(row)) < m.w {
			row = append(row, ' ')
		}
		m.m[y] = row
	}
	innerLabels := make(map[string]ivec2)
	outerLabels := make(map[string]ivec2)
	labelCells := make(map[ivec2]struct{})
	for y, row := range m.m {
		for x, c := range row {
			x, y := int64(x), int64(y)
			if !mazeLabel(c) {
				continue
			}
			p0 := ivec2{x, y}
			if _, ok := labelCells[p0]; ok {
				continue
			}
			p1 := ivec2{x, y + 1}
			var p2 ivec2
			if mazeLabel(m.at(p1)) {
				p2 = ivec2{x, y - 1}
				if m.at(p2) != '.' {
					p2 = ivec2{x, y + 2}
					if m.at(p2) != '.' {
						panic("can't find labeled cell")
					}
				}
			} else {
				p1 = ivec2{x + 1, y}
				if !mazeLabel(m.at(p1)) {
					panic("bad label")
				}
				p2 = ivec2{x - 1, y}
				if m.at(p2) != '.' {
					p2 = ivec2{x + 2, y}
					if m.at(p2) != '.' {
						panic("can't find labeled cell")
					}
				}
			}
			label := string([]byte{m.at(p0), m.at(p1)})
			if p0.x == 0 || p0.x == m.w-2 || p0.y == 0 || p0.y == m.h-2 {
				if _, ok := outerLabels[label]; ok {
					panic("duplicate label")
				}
				outerLabels[label] = p2
			} else {
				if _, ok := innerLabels[label]; ok {
					panic("duplicate label")
				}
				innerLabels[label] = p2
			}
			labelCells[p0] = struct{}{}
			labelCells[p1] = struct{}{}
		}
	}

	cell, ok := outerLabels["AA"]
	if !ok {
		panic("bad start")
	}
	m.start = cell
	delete(outerLabels, "AA")

	cell, ok = outerLabels["ZZ"]
	if !ok {
		panic("bad end")
	}
	m.end = cell
	delete(outerLabels, "ZZ")

	m.innerPortals = make(map[ivec2]ivec2)
	m.outerPortals = make(map[ivec2]ivec2)
	for label, cellOuter := range outerLabels {
		cellInner, ok := innerLabels[label]
		if !ok {
			panic("unmatched outer label")
		}
		m.outerPortals[cellOuter] = cellInner
		m.innerPortals[cellInner] = cellOuter
		delete(innerLabels, label)
	}
	if len(innerLabels) > 0 {
		panic("unmatched inner label")
	}
}

func mazeLabel(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func (m *plutoMaze) at(p ivec2) byte {
	if p.y < 0 || p.y >= m.h || p.x < 0 || p.x >= m.w {
		return ' '
	}
	return m.m[p.y][p.x]
}

func (m *plutoMaze) shortestPath() int64 {
	type posDepth struct {
		p ivec2
		d int64
	}
	visited := map[ivec2]struct{}{m.start: {}}
	q := []posDepth{{m.start, 0}}
	for len(q) > 0 {
		pd := q[0]
		q = q[1:]
		for _, p1 := range m.neighbors(pd.p) {
			if _, ok := visited[p1]; ok {
				continue
			}
			visited[p1] = struct{}{}
			d1 := pd.d + 1
			if p1 == m.end {
				return d1
			}
			q = append(q, posDepth{p1, d1})
		}
	}
	panic("no path")
}

func (m *plutoMaze) neighbors(p ivec2) []ivec2 {
	var ns []ivec2
	for _, d := range []ivec2{
		{-1, 0},
		{0, -1},
		{1, 0},
		{0, 1},
	} {
		p1 := p.add(d)
		if m.at(p1) == '.' {
			ns = append(ns, p1)
		}
	}
	if p1, ok := m.innerPortals[p]; ok {
		ns = append(ns, p1)
	}
	if p1, ok := m.outerPortals[p]; ok {
		ns = append(ns, p1)
	}
	return ns
}

func (m *plutoMaze) shortestPathRecursive() int64 {
	type posDepth struct {
		pr posRec
		d  int64
	}
	start := posRec{m.start, 0}
	visited := map[posRec]struct{}{start: {}}
	q := []posDepth{{start, 0}}
	for len(q) > 0 {
		pd := q[0]
		q = q[1:]
		for _, pr1 := range m.neighborsRecursive(pd.pr) {
			if _, ok := visited[pr1]; ok {
				continue
			}
			visited[pr1] = struct{}{}
			d1 := pd.d + 1
			if pr1.r == 0 && pr1.p == m.end {
				return d1
			}
			q = append(q, posDepth{pr1, d1})
		}
	}
	panic("no path")
}

type posRec struct {
	p ivec2
	r int64
}

func (m *plutoMaze) neighborsRecursive(pr posRec) []posRec {
	var ns []posRec
	for _, d := range []ivec2{
		{-1, 0},
		{0, -1},
		{1, 0},
		{0, 1},
	} {
		p1 := pr.p.add(d)
		if m.at(p1) == '.' {
			ns = append(ns, posRec{p1, pr.r})
		}
	}
	if p1, ok := m.innerPortals[pr.p]; ok {
		ns = append(ns, posRec{p1, pr.r + 1})
	}
	if pr.r > 0 {
		if p1, ok := m.outerPortals[pr.p]; ok {
			ns = append(ns, posRec{p1, pr.r - 1})
		}
	}
	return ns
}
