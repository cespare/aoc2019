package main

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcode(prog)
	a := newArcade(ic)
	ic.run()
	ctx.reportPart1(a.countBlocks())

	prog[0] = 2
	ic = newIntcode(prog)
	a = newArcade(ic)
	ic.run()
	if a.countBlocks() > 0 {
		panic("you lose")
	}
	ctx.reportPart1(a.score)
}

type arcade struct {
	ic  *intcode
	out []int64

	m      map[ivec2]int
	score  int64
	paddle ivec2
	ball   ivec2
}

func newArcade(ic *intcode) *arcade {
	a := &arcade{
		ic: ic,
		m:  make(map[ivec2]int),
	}
	ic.in = func(buf []int64) []int64 {
		var v int64
		switch {
		case a.paddle.x < a.ball.x:
			v = 1
		case a.paddle.x > a.ball.x:
			v = -1
		default:
			v = 0
		}
		return append(buf, v)
	}
	ic.out = func(v int64) {
		a.out = append(a.out, v)
		if len(a.out) < 3 {
			return
		}
		a.add(a.out[0], a.out[1], a.out[2])
		a.out = a.out[:0]
	}
	return a
}

func (a *arcade) add(n0, n1, n2 int64) {
	if n0 == -1 && n1 == 0 {
		a.score = n2
		return
	}
	t := int(n2)
	v := ivec2{n0, n1}
	a.m[v] = t
	switch t {
	case 3:
		a.paddle = v
	case 4:
		a.ball = v
	}
}

func (a *arcade) countBlocks() int {
	var blocks int
	for _, t := range a.m {
		if t == 2 {
			blocks++
		}
	}
	return blocks
}
