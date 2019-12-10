package main

import (
	"bufio"
	"log"
	"math"
	"sort"
)

func init() {
	addSolutions(10, problem10)
}

func problem10(ctx *problemContext) {
	var lines [][]byte
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		lines = append(lines, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	m := loadMap(lines)
	var max int
	var station ivec2
	for p := range m {
		if m := countVisible(m, p); m > max {
			max = m
			station = p
		}
	}
	ctx.reportPart1(max)

	delete(m, station)

	targetsByNorm := make(map[ivec2][]targetAsteroid)
	for p := range m {
		rel := p.sub(station)
		n := norm(rel)
		a := targetAsteroid{
			abs: p,
			rel: rel,
			d:   rel.mag(),
			n:   n,
		}
		targetsByNorm[n] = append(targetsByNorm[n], a)
	}
	var targets [][]targetAsteroid
	for _, ts := range targetsByNorm {
		targets = append(targets, ts)
	}
	sort.Slice(targets, func(i, j int) bool {
		a0, a1 := targets[i][0], targets[j][0]
		if a0.n == a1.n {
			return false
		}
		card0, card1 := normToCardinalGroup(a0.n), normToCardinalGroup(a1.n)
		if card0 != card1 {
			return card0 < card1
		}
		slope0, slope1 := a0.n.slope(), a1.n.slope()
		return slope0 < slope1
	})
	for _, ts := range targets {
		sort.Slice(ts, func(i, j int) bool {
			a0, a1 := ts[i], ts[j]
			return a0.rel.mag() < a1.rel.mag()
		})
	}
	var removed int
	for {
		found := false
		for i, ts := range targets {
			if len(ts) == 0 {
				continue
			}
			found = true
			a := ts[0]
			targets[i] = ts[1:]
			removed++
			if removed == 200 {
				ctx.reportPart2(a.abs.x*100 + a.abs.y)
				return
			}
		}
		if !found {
			panic("not enough")
		}
	}
}

func normToCardinalGroup(n ivec2) int {
	switch n {
	case ivec2{0, -1}:
		return 0
	case ivec2{1, 0}:
		return 2
	case ivec2{0, 1}:
		return 4
	case ivec2{-1, 0}:
		return 6
	}
	switch {
	case n.x > 0 && n.y < 0:
		return 1
	case n.x > 0 && n.y > 0:
		return 3
	case n.x < 0 && n.y > 0:
		return 5
	case n.x < 0 && n.y < 0:
		return 7
	}
	panic("unreached")
}

func loadMap(lines [][]byte) map[ivec2]struct{} {
	m := make(map[ivec2]struct{})
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				m[ivec2{int64(x), int64(y)}] = struct{}{}
			}
		}
	}
	return m
}

func countVisible(m map[ivec2]struct{}, p ivec2) int {
	slopes := make(map[ivec2]struct{})
	for p1 := range m {
		if p1 == p {
			continue
		}
		slope := norm(p1.sub(p))
		slopes[slope] = struct{}{}
	}
	return len(slopes)
}

func norm(v ivec2) ivec2 {
	switch {
	case v.x == 0 && v.y > 0:
		return ivec2{0, 1}
	case v.x == 0 && v.y < 0:
		return ivec2{0, -1}
	case v.y == 0 && v.x > 0:
		return ivec2{1, 0}
	case v.y == 0 && v.x < 0:
		return ivec2{-1, 0}
	}
	lim := min(abs(v.x), abs(v.y))
	for n := lim; n >= 2; n-- {
		if v.x%n == 0 && v.y%n == 0 {
			v.x /= n
			v.y /= n
			return v
		}
	}
	return v
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type targetAsteroid struct {
	abs ivec2 // absolute pos
	rel ivec2 // pos from station
	n   ivec2 // norm(rel)
	d   float64
}

func (v ivec2) mag() float64 {
	x := math.Abs(float64(v.x))
	y := math.Abs(float64(v.y))
	return math.Sqrt(x*x + y*y)
}

func (v ivec2) slope() float64 {
	return float64(v.y) / float64(v.x)
}
