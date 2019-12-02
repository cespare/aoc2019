package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	line, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	var ic intcode
	for _, s := range strings.Split(string(bytes.TrimSpace(line)), ",") {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalln("Bad number:", s)
		}
		ic.state = append(ic.state, n)
	}
	ctx.reportLoad()

	ic1 := ic.copy()
	ic1.state[1] = 12
	ic1.state[2] = 2

	for !ic1.step() {
	}
	ctx.reportPart1(ic1.state[0])

	for noun := int64(0); noun < 100; noun++ {
		for verb := int64(0); verb < 100; verb++ {
			ic1 := ic.copy()
			ic1.state[1] = noun
			ic1.state[2] = verb
			for !ic1.step() {
			}
			if ic1.state[0] == 19690720 {
				ctx.reportPart2(100*noun + verb)
				return
			}
		}
	}
	panic("solution not found")
}

type intcode struct {
	state []int64
	pc    int
}

func (ic *intcode) step() (done bool) {
	switch ic.state[ic.pc] {
	case 1:
		i, j, k := ic.state[ic.pc+1], ic.state[ic.pc+2], ic.state[ic.pc+3]
		ic.state[k] = ic.state[i] + ic.state[j]
		ic.pc += 4
		return false
	case 2:
		i, j, k := ic.state[ic.pc+1], ic.state[ic.pc+2], ic.state[ic.pc+3]
		ic.state[k] = ic.state[i] * ic.state[j]
		ic.pc += 4
		return false
	case 99:
		return true
	default:
		panic("bad state")
	}
}

func (ic *intcode) copy() *intcode {
	return &intcode{
		state: append([]int64(nil), ic.state...),
		pc:    ic.pc,
	}
}
