package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(3, problem3)
}

func problem3(ctx *problemContext) {
	var wires [][]string
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		wire := strings.Split(scanner.Text(), ",")
		wires = append(wires, wire)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if len(wires) != 2 {
		log.Fatalf("Expect 2 wires; got %d", len(wires))
	}
	ctx.reportLoad()

	m0, err := wireToDist(wires[0])
	if err != nil {
		log.Fatal(err)
	}
	m1, err := wireToDist(wires[1])
	if err != nil {
		log.Fatal(err)
	}

	var best int64
	for v := range m0 {
		if _, ok := m1[v]; ok {
			if best == 0 || v.len() < best {
				best = v.len()
			}
		}
	}

	ctx.reportPart1(best)

	best = 0
	for v, d0 := range m0 {
		if d1, ok := m1[v]; ok {
			d := d0 + d1
			if best == 0 || d < best {
				best = d
			}
		}
	}

	ctx.reportPart2(best)
}

func wireToDist(w []string) (map[ivec2]int64, error) {
	m := make(map[ivec2]int64)
	var v ivec2
	d := int64(1)
	for _, s := range w {
		vd, n, err := parseWireDir(s)
		if err != nil {
			return nil, err
		}
		for i := int64(0); i < n; i++ {
			v = v.add(vd)
			m[v] = d
			d++
		}
	}
	return m, nil
}

func parseWireDir(s string) (v ivec2, n int64, err error) {
	if len(s) == 0 {
		return ivec2{}, 0, errors.New("empty dir")
	}
	dir := s[0]
	switch dir {
	case 'U':
		v = ivec2{y: 1}
	case 'D':
		v = ivec2{y: -1}
	case 'L':
		v = ivec2{x: -1}
	case 'R':
		v = ivec2{x: 1}
	default:
		return ivec2{}, 0, fmt.Errorf("unknown dir %c", dir)
	}
	n, err = strconv.ParseInt(s[1:], 10, 64)
	if err != nil {
		return ivec2{}, 0, err
	}
	return v, n, nil
}

type ivec2 struct {
	x, y int64
}

func (v ivec2) add(v1 ivec2) ivec2 {
	return ivec2{
		x: v.x + v1.x,
		y: v.y + v1.y,
	}
}

func (v ivec2) len() int64 {
	return iabs(v.x) + iabs(v.y)
}

func iabs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
