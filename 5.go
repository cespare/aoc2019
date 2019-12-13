package main

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcode(prog, 1)
	ic.run()
	ctx.reportPart1(ic.output)

	ic = newIntcode(prog, 5)
	ic.run()

	ctx.reportPart2(ic.output)
}
