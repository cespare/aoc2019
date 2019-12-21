package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(21, problem21)
}

func problem21(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	spring := `
NOT J J
AND A J
AND B J
AND C J
NOT J J
AND D J
WALK
`
	ctx.reportPart1(runSpring(prog, spring))

	spring = `
NOT J J
AND A J
AND B J
AND C J
NOT J J
AND D J
OR E T
OR H T
AND T J
RUN
`
	ctx.reportPart2(runSpring(prog, spring))
}

func runSpring(prog []int64, spring string) string {
	ic := newIntcodeWithMem(prog)
	for _, c := range []byte(spring[1:]) {
		ic.input = append(ic.input, int64(c))
	}
	if !ic.run() {
		panic("need more input")
	}
	var out strings.Builder
	for _, n := range ic.output {
		if n >= 0 && n < 256 {
			out.WriteByte(byte(n))
		} else {
			fmt.Fprintf(&out, "%d", n)
		}
	}
	return out.String()
}
