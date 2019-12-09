package main

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcode(prog)
	ic.mem[1] = 12
	ic.mem[2] = 2
	ic.run()
	ctx.reportPart1(ic.mem[0])

	for noun := int64(0); noun < 100; noun++ {
		for verb := int64(0); verb < 100; verb++ {
			ic := newIntcode(prog)
			ic.mem[1] = noun
			ic.mem[2] = verb
			ic.run()
			if ic.mem[0] == 19690720 {
				ctx.reportPart2(100*noun + verb)
				return
			}
		}
	}
	panic("solution not found")
}
