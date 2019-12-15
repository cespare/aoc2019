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

	m0 := make(map[ivec2]int64, 200e3)
	var v ivec2
	d0 := int64(1)
	for _, s := range wires[0] {
		vd, n, err := parseWireDir(s)
		if err != nil {
			log.Fatal(err)
		}
		for i := int64(0); i < n; i++ {
			v = v.add(vd)
			m0[v] = d0
			d0++
		}
	}

	v = ivec2{0, 0}
	d1 := int64(1)
	var bestPart1, bestPart2 int64
	for _, s := range wires[1] {
		vd, n, err := parseWireDir(s)
		if err != nil {
			log.Fatal(err)
		}
		for i := int64(0); i < n; i++ {
			v = v.add(vd)

			if d0, ok := m0[v]; ok {
				if bestPart1 == 0 || v.len() < bestPart1 {
					bestPart1 = v.len()
				}
				d := d0 + d1
				if bestPart2 == 0 || d < bestPart2 {
					bestPart2 = d
				}
			}

			d1++
		}

	}

	ctx.reportPart1(bestPart1)
	ctx.reportPart2(bestPart2)
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

func (v ivec2) sub(v1 ivec2) ivec2 {
	return ivec2{
		x: v.x - v1.x,
		y: v.y - v1.y,
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
