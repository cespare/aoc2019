package main

import (
	"sync"

	"github.com/cespare/permute"
)

func init() {
	addSolutions(7, problem7)
}

func problem7(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	phases := make([]int, 5)
	for i := range phases {
		phases[i] = i
	}
	perm := permute.Ints(phases)
	maxSignal := int64(-999)
	for perm.Permute() {
		if sig := evalAmps(prog, phases); sig > maxSignal {
			maxSignal = sig
		}
	}
	ctx.reportPart1(maxSignal)

	for i := range phases {
		phases[i] = 5 + i
	}
	perm = permute.Ints(phases)
	maxSignal = int64(-999)
	for perm.Permute() {
		if sig := evalFeedbackAmps(prog, phases); sig > maxSignal {
			maxSignal = sig
		}
	}
	ctx.reportPart2(maxSignal)
}

func evalAmps(prog []int64, phases []int) int64 {
	var val int64
	for _, phase := range phases {
		ic := newIntcode(prog)
		ic.setInput(int64(phase), val)
		ic.setOutputLastVal(&val)
		ic.run()
	}
	return val
}

func evalFeedbackAmps(prog []int64, phases []int) int64 {
	chs := make([]chan int64, len(phases))
	for i := range chs {
		chs[i] = make(chan int64, 1)
	}
	var wg sync.WaitGroup
	for i, ch := range chs {
		ic := newIntcode(prog)
		ic.setInput(int64(phases[i]))
		ic.setInputChan(ch)
		ic.setOutputChan(chs[(i+1)%len(phases)])
		wg.Add(1)
		go func() {
			ic.run()
			wg.Done()
		}()
	}
	chs[0] <- 0
	wg.Wait()
	return <-chs[0]
}
