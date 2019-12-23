package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

func readProg(r io.Reader) []int64 {
	line, err := ioutil.ReadAll(r)
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
	return prog
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
	opSetRel    = 9
	opHalt      = 99
)

func decodeIntcode(inst int64) (modes []int, opcode int) {
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
	case opSetRel:
		numParams = 1
	case opHalt:
		numParams = 0
	default:
		panic("bad opcode")
	}
	modes = make([]int, numParams)
	var m int
	for i := range modes {
		m, n = int(n%10), n/10
		if m > 2 {
			panic("bad parameter mode")
		}
		modes[i] = m
	}
	if n > 0 {
		panic(fmt.Sprintf("too many parameter modes for opcode %d", opcode))
	}
	return modes, opcode
}

const maxIntcodeMem = 1e6

type intcodeOpt uint

const (
	optRun intcodeOpt = 1 << iota
	optRunUntilInput
)

type intcode struct {
	prog    []int64
	mem     []int64
	pc      int64
	relBase int64

	opt   intcodeOpt
	in    func([]int64) []int64 // append-style
	inBuf []int64
	out   func(int64)
}

func newIntcode(prog []int64) *intcode {
	return &intcode{
		prog: copyInt64s(prog),
		mem:  make([]int64, 8),
		out: func(int64) {
			panic("encountered output instruction; no out function set")
		},
	}
}

func (ic *intcode) setInput(ns ...int64) {
	if ic.opt&optRun != 0 {
		panic("setInitialInput called after run was called")
	}
	ic.inBuf = ns
}

func (ic *intcode) setOutputAll(p *[]int64) {
	ic.out = func(v int64) { *p = append(*p, v) }
}

func (ic *intcode) setOutputLastVal(p *int64) {
	ic.out = func(v int64) { *p = v }
}

func (ic *intcode) setInputChan(ch chan int64) {
	ic.in = func(buf []int64) []int64 { return append(buf, <-ch) }
}

func (ic *intcode) setOutputChan(ch chan int64) {
	ic.out = func(v int64) { ch <- v }
}

func (ic *intcode) runUntilInput() (halt bool) {
	if ic.opt&optRun != 0 {
		panic("runUntilInput called after run was previously called")
	}
	if ic.in != nil {
		panic("runUntilInput called after in was already set")
	}
	ic.opt |= optRunUntilInput
	for {
		switch ic.step(true) {
		case stateRunning:
		case stateSuspend:
			return false
		case stateHalt:
			return true
		}
	}
}

func (ic *intcode) run() {
	if ic.opt&(optRun|optRunUntilInput) != 0 {
		panic("run called after run or runUntilInput was already called")
	}
	if ic.in == nil {
		ic.in = func([]int64) []int64 {
			panic("encountered input instruction; no in function set")
		}
	}
	ic.opt |= optRun
	for {
		if ic.step(false) == stateHalt {
			return
		}
	}
}

type intcodeState int

const (
	stateRunning intcodeState = iota
	stateSuspend
	stateHalt
)

func (ic *intcode) step(waitForInput bool) intcodeState {
	modes, opcode := decodeIntcode(ic.prog[ic.pc])
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
		if len(ic.inBuf) == 0 {
			if waitForInput {
				ic.pc--
				return stateSuspend
			}
			ic.inBuf = ic.in(ic.inBuf)
			if len(ic.inBuf) == 0 {
				panic("ic.in gave no input")
			}
		}
		v := ic.inBuf[0]
		ic.inBuf = ic.inBuf[1:]
		ic.set(v, modes[0])
	case opOutput:
		ic.out(ic.get(modes[0]))
	case opJumpTrue:
		v := ic.get(modes[0])
		targ := ic.get(modes[1])
		if v != 0 {
			ic.pc = targ
		}
	case opJumpFalse:
		v := ic.get(modes[0])
		targ := ic.get(modes[1])
		if v == 0 {
			ic.pc = targ
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
	case opSetRel:
		ic.relBase += ic.get(modes[0])
	case opHalt:
		return stateHalt
	default:
		panic("bad state")
	}
	return stateRunning
}

func (ic *intcode) get(mode int) int64 {
	param := ic.prog[ic.pc]
	ic.pc++
	off := param
	switch mode {
	case 0:
	case 1:
		return param
	case 2:
		off += ic.relBase
	default:
		panic("bad mode")
	}
	if off < int64(len(ic.prog)) {
		return ic.prog[off]
	}
	off -= int64(len(ic.prog))
	ic.allocMem(off)
	return ic.mem[off]
}

func (ic *intcode) set(val int64, mode int) {
	param := ic.prog[ic.pc]
	ic.pc++
	off := param
	switch mode {
	case 0:
	case 1:
		panic("write to immediate")
	case 2:
		off += ic.relBase
	default:
		panic("bad mode")
	}
	if off < int64(len(ic.prog)) {
		ic.prog[off] = val
		return
	}
	off -= int64(len(ic.prog))
	ic.allocMem(off)
	ic.mem[off] = val
}

func (ic *intcode) allocMem(off int64) {
	if off < int64(len(ic.mem)) {
		return
	}
	if off >= maxIntcodeMem {
		panic(fmt.Sprintf("intcode program required more than %d bytes of extra memory", int64(maxIntcodeMem)))
	}
	n := int64(len(ic.mem)) * 10
	if n < off*2 {
		n = off * 2
	}
	if n > maxIntcodeMem {
		n = maxIntcodeMem
	}
	mem := make([]int64, n)
	copy(mem, ic.mem)
	ic.mem = mem
}

var intcodePool sync.Pool

func (ic *intcode) free() { intcodePool.Put(ic) }

func (ic *intcode) clone() *intcode {
	x := intcodePool.Get()
	if x == nil {
		ic1 := *ic
		ic1.prog = copyInt64s(ic.prog)
		ic1.mem = copyInt64s(ic.mem)
		ic1.inBuf = copyInt64s(ic.inBuf)
		return &ic1
	}

	cloneInt64s := func(p *[]int64, ns []int64) {
		if cap(*p) < len(ns) {
			*p = make([]int64, len(ns))
		} else {
			*p = (*p)[:len(ns)]
		}
		copy(*p, ns)
	}

	ic1 := x.(*intcode)
	cloneInt64s(&ic1.prog, ic.prog)
	cloneInt64s(&ic1.mem, ic.mem)
	ic1.pc = ic.pc
	ic1.relBase = ic.relBase
	ic1.opt = ic.opt
	ic1.in = ic.in
	cloneInt64s(&ic1.inBuf, ic.inBuf)
	ic1.out = ic.out
	return ic1
}

func copyInt64s(s []int64) []int64 {
	return append([]int64(nil), s...)
}
