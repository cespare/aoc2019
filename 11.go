package main

import (
	"fmt"
)

func init() {
	addSolutions(11, problem11)
}

func problem11(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	colors := make(map[ivec2]int64)
	paintHull(prog, colors)
	ctx.reportPart1(len(colors))

	colors = map[ivec2]int64{{0, 0}: 1}
	paintHull(prog, colors)
	printHull(colors)
	ctx.reportPart2("^^^")
}

func paintHull(prog []int64, colors map[ivec2]int64) {
	ic := newIntcodeWithMem(prog)
	ic.setChannelMode()
	go ic.run()
	var p ivec2
	dir := ivec2{0, 1}
	for {
		ic.inCh <- colors[p]

		color := <-ic.outCh
		t, ok := <-ic.outCh

		colors[p] = color

		dir = turn(dir, t)
		p = p.add(dir)

		if !ok {
			return
		}
	}
}

func printHull(colors map[ivec2]int64) {
	min := ivec2{9999999, 9999999}
	max := ivec2{-9999999, -9999999}
	for p := range colors {
		if p.x < min.x {
			min.x = p.x
		}
		if p.y < min.y {
			min.y = p.y
		}
		if p.x > max.x {
			max.x = p.x
		}
		if p.y > max.y {
			max.y = p.y
		}
	}
	for y := max.y; y >= min.y; y-- {
		for x := min.x; x <= max.x; x++ {
			color, ok := colors[ivec2{x, y}]
			if !ok {
				fmt.Print("  ")
				continue
			}
			switch color {
			case 0:
				fmt.Print("░░")
			case 1:
				fmt.Print("██")
			default:
				panic("bad color")
			}
		}
		fmt.Println()
	}
}

var (
	north = ivec2{0, 1}
	east  = ivec2{1, 0}
	south = ivec2{0, -1}
	west  = ivec2{-1, 0}
)

func turn(dir ivec2, t int64) ivec2 {
	switch t {
	case 0:
		switch dir {
		case north:
			return west
		case east:
			return north
		case south:
			return east
		case west:
			return south
		default:
			panic("bad dir")
		}
	case 1:
		switch dir {
		case north:
			return east
		case east:
			return south
		case south:
			return west
		case west:
			return north
		default:
			panic("bad dir")
		}
	default:
		panic("bad turn dir")
	}
}
