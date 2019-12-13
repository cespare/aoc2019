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
		ic := newIntcode(prog, int64(phase), val)
		ic.run()
		if len(ic.output) != 1 {
			panic("bad output")
		}
		val = ic.output[0]
	}
	return val
}

func evalFeedbackAmps(prog []int64, phases []int) int64 {
	ics := make([]*intcode, len(phases))
	vals := make([]int64, len(phases))
	for i := range ics {
		ic := newIntcode(prog)
		ic.setChannelMode()
		go ic.run()
		ics[i] = ic
	}
	var wg sync.WaitGroup
	for i, ic := range ics {
		i := i
		ic0 := ic
		j := i + 1
		if j == len(ics) {
			j = 0
		}
		ic1 := ics[j]
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ic0.outCh {
				vals[i] = v
				ic1.inCh <- v
			}
		}()
	}
	for i, phase := range phases {
		ics[i].inCh <- int64(phase)
	}
	ics[0].inCh <- 0
	wg.Wait()
	return vals[len(vals)-1]
}
