package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	line, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	var prog []int64
	for _, s := range strings.Split(string(bytes.TrimSpace(line)), ",") {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalln("Bad number:", s)
		}
		prog = append(prog, n)
	}
	ctx.reportLoad()

	ic := newIntcodeWithMem(prog)
	a := newArcade()
	a.run(ic)
	ctx.reportPart1(a.countBlocks())

	prog[0] = 2
	ic = newIntcodeWithMem(prog)
	ic.setSuspendMode()
	a = newArcade()
	for {
		halted := a.run(ic)
		if a.countBlocks() == 0 {
			ctx.reportPart2(a.score)
			return
		}
		if halted {
			panic("you lose")
		}
		var in int64
		switch {
		case a.paddle.x < a.ball.x:
			in = 1
		case a.paddle.x > a.ball.x:
			in = -1
		default:
			in = 0
		}
		ic.input = append(ic.input, in)
	}
}

type arcade struct {
	m      map[ivec2]int
	score  int64
	paddle ivec2
	ball   ivec2
}

func newArcade() *arcade {
	return &arcade{m: make(map[ivec2]int)}
}

func (a *arcade) run(ic *intcode) (halted bool) {
	halted = ic.run()
	a.addOutput(ic)
	return halted
}

func (a *arcade) addOutput(ic *intcode) {
	out := ic.output
	for i := 0; i < len(out); i += 3 {
		a.add(out[i], out[i+1], out[i+2])
	}
	ic.output = out[:0]
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
