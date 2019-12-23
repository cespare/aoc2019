package main

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcode(prog)
	ic.setInput(1)
	var lastVal int64
	ic.setOutputLastVal(&lastVal)
	ic.run()
	ctx.reportPart1(lastVal)

	ic = newIntcode(prog)
	ic.setInput(5)
	ic.setOutputLastVal(&lastVal)
	ic.run()
	ctx.reportPart2(lastVal)
}
