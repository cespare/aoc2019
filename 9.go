package main

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcodeWithMem(prog, 1)
	ic.run()
	ctx.reportPart1(ic.output)

	ic = newIntcodeWithMem(prog, 2)
	ic.run()
	ctx.reportPart2(ic.output)
}
