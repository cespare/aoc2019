package main

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcode(prog)
	ic.setInput(1)
	var val int64
	ic.setOutputLastVal(&val)
	ic.run()
	ctx.reportPart1(val)

	ic = newIntcode(prog)
	ic.setInput(2)
	ic.setOutputLastVal(&val)
	ic.run()
	ctx.reportPart2(val)
}
