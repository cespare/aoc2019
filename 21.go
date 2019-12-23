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
OR A J
AND B J
AND C J
NOT J J
AND D J
WALK
`
	ctx.reportPart1(runSpring(prog, spring))

	spring = `
OR A J
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
	ic := newIntcode(prog)
	input := make([]int64, len(spring)-1)
	for i := 1; i < len(spring); i++ {
		input[i-1] = int64(spring[i])
	}
	ic.setInput(input...)
	var out strings.Builder
	ic.out = func(n int64) {
		if n >= 0 && n < 256 {
			out.WriteByte(byte(n))
		} else {
			fmt.Fprintf(&out, "%d", n)
		}
	}
	ic.run()
	return out.String()
}
