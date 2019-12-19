package main

import "fmt"

func init() {
	addSolutions(19, problem19)
}

func problem19(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcodeWithMem(prog)
	const print = false
	var tractorPoints int
	for y := int64(0); y < 50; y++ {
		for x := int64(0); x < 50; x++ {
			if inTractorBeam(ic, ivec2{x, y}) {
				if print {
					fmt.Print("#")
				}
				tractorPoints++
			} else {
				if print {
					fmt.Print(".")
				}
			}
		}
		if print {
			fmt.Println()
		}
	}
	ctx.reportPart1(tractorPoints)

	const size = 100
	// 1: Find a reasonable starting point.
	topRight := ivec2{0, 20}
	for !inTractorBeam(ic, topRight) {
		topRight = topRight.add(ivec2{1, 0})
	}
	for {
		// 2: Find the right-most point on this row to place the
		// top-right corner.
		for {
			next := topRight.add(ivec2{1, 0})
			if !inTractorBeam(ic, next) {
				break
			}
			topRight = next
		}
		// 3. Check whether the bottom-left corner is in the beam.
		bottomLeft := topRight.add(ivec2{-(size - 1), size - 1})
		if inTractorBeam(ic, bottomLeft) {
			topLeft := topRight.add(ivec2{-(size - 1), 0})
			ctx.reportPart2(topLeft.x*10000 + topLeft.y)
			return
		}
		topRight = topRight.add(ivec2{0, 1})
	}
}

func inTractorBeam(ic *intcode, p ivec2) bool {
	ic = ic.clone()
	defer ic.free()
	ic.input = []int64{p.x, p.y}
	ic.run()
	if len(ic.output) != 1 {
		panic("bad")
	}
	switch ic.output[0] {
	case 0:
		return false
	case 1:
		return true
	default:
		panic("bad")
	}
}
