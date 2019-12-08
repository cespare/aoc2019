package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	line, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	var program []int64
	for _, s := range strings.Split(string(bytes.TrimSpace(line)), ",") {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalln("Bad number:", s)
		}
		program = append(program, n)
	}
	ctx.reportLoad()

	ic := newIntcode2(program, 1)
	ic.run()

	ctx.reportPart1(ic.output)

	ic = newIntcode2(program, 5)
	ic.run()

	ctx.reportPart1(ic.output)
}

type intcode2 struct {
	state []int64
	pc    int

	input  []int64
	output []int64

	inCh  chan int64
	outCh chan int64
}

func newIntcode2(prog []int64, input ...int64) *intcode2 {
	return &intcode2{
		input: input,
		state: append([]int64(nil), prog...),
	}
}

func (ic *intcode2) setChannelMode() {
	if len(ic.input) > 0 || len(ic.output) > 0 {
		panic("cannot set channel mode after using input/output slices")
	}
	ic.inCh = make(chan int64, 1)
	ic.outCh = make(chan int64, 1)
}

const (
	opAdd       = 1
	opMul       = 2
	opInput     = 3
	opOutput    = 4
	opJumpTrue  = 5
	opJumpFalse = 6
	opLess      = 7
	opEqual     = 8
	opHalt      = 99
)

func (ic *intcode2) decode(inst int64) (modes []int, opcode int) {
	if inst < 0 {
		panic("bad instruction")
	}
	opcode = int(inst % 100)
	n := inst / 100
	var numParams int
	switch opcode {
	case opAdd:
		numParams = 3
	case opMul:
		numParams = 3
	case opInput:
		numParams = 1
	case opOutput:
		numParams = 1
	case opJumpTrue:
		numParams = 2
	case opJumpFalse:
		numParams = 2
	case opLess:
		numParams = 3
	case opEqual:
		numParams = 3
	case opHalt:
		numParams = 0
	default:
		panic("bad opcode")
	}
	modes = make([]int, numParams)
	var m int
	for i := range modes {
		m, n = int(n%10), n/10
		if m > 1 {
			panic("bad parameter mode")
		}
		modes[i] = m
	}
	if n > 0 {
		panic(fmt.Sprintf("too many parameter modes for opcode %d", opcode))
	}
	return modes, opcode
}

func (ic *intcode2) run() {
	for !ic.step() {
	}
}

func (ic *intcode2) step() (done bool) {
	modes, opcode := ic.decode(ic.state[ic.pc])
	ic.pc++
	switch opcode {
	case opAdd:
		a := ic.get(modes[0])
		b := ic.get(modes[1])
		ic.set(a+b, modes[2])
	case opMul:
		a := ic.get(modes[0])
		b := ic.get(modes[1])
		ic.set(a*b, modes[2])
	case opInput:
		var v int64
		if ic.inCh == nil {
			v = ic.input[0]
			ic.input = ic.input[1:]
		} else {
			v = <-ic.inCh
		}
		ic.set(v, modes[0])
	case opOutput:
		v := ic.get(modes[0])
		if ic.outCh == nil {
			if len(ic.output) >= 1e6 {
				panic("output exploded")
			}
			ic.output = append(ic.output, v)
		} else {
			ic.outCh <- v
		}
	case opJumpTrue:
		v := ic.get(modes[0])
		targ := ic.get(modes[1])
		if v != 0 {
			ic.pc = int(targ)
		}
	case opJumpFalse:
		v := ic.get(modes[0])
		targ := ic.get(modes[1])
		if v == 0 {
			ic.pc = int(targ)
		}
	case opLess:
		a := ic.get(modes[0])
		b := ic.get(modes[1])
		var v int64
		if a < b {
			v = 1
		}
		ic.set(v, modes[2])
	case opEqual:
		a := ic.get(modes[0])
		b := ic.get(modes[1])
		var v int64
		if a == b {
			v = 1
		}
		ic.set(v, modes[2])
	case opHalt:
		if ic.outCh != nil {
			close(ic.outCh)
		}
		return true
	default:
		panic("bad state")
	}
	return false
}

func (ic *intcode2) get(mode int) int64 {
	param := ic.state[ic.pc]
	ic.pc++
	switch mode {
	case 0:
		return ic.state[param]
	case 1:
		return param
	default:
		panic("bad mode")
	}
}

func (ic *intcode2) set(val int64, mode int) {
	param := ic.state[ic.pc]
	ic.pc++
	switch mode {
	case 0:
		ic.state[param] = val
	case 1:
		panic("write to immediate")
	default:
		panic("bad mode")
	}
}
