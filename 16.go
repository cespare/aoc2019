package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func init() {
	addSolutions(16, problem16)
}

func problem16(ctx *problemContext) {
	var input []int
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range bytes.TrimSpace(b) {
		input = append(input, int(c-'0'))
	}
	ctx.reportLoad()

	digs := append([]int(nil), input...)
	for i := 0; i < 100; i++ {
		digs = fftRound(digs)
	}
	var part1 strings.Builder
	for _, d := range digs[:8] {
		fmt.Fprintf(&part1, "%d", d)
	}
	ctx.reportPart1(part1.String())

	const mul = 10000
	digs = make([]int, len(input)*mul)
	for i := range digs {
		digs[i] = input[i%len(input)]
	}
	var off int
	for i := 0; i < 7; i++ {
		off *= 10
		off += input[i]
	}
	if off-1 < len(digs)/2 {
		panic("not possible")
	}

	sub := digs[off:]
	for i := 0; i < 100; i++ {
		for j := len(sub) - 2; j >= 0; j-- {
			sub[j] = (sub[j] + sub[j+1]) % 10
		}
	}
	var part2 strings.Builder
	for _, d := range sub[:8] {
		fmt.Fprintf(&part2, "%d", d)
	}
	ctx.reportPart2(part2.String())
}

var fftBase = []int{0, 1, 0, -1}

func fftRound(digs []int) []int {
	out := make([]int, len(digs))
	for i := range digs {
		repeat := i + 1
		r := repeat - 1
		j := 0
		for _, d := range digs {
			if r == 0 {
				j = (j + 1) % len(fftBase)
				r = repeat
			}
			out[i] += fftBase[j] * d
			r--
		}
		out[i] = lastDigit(out[i])
	}
	return out
}

func lastDigit(n int) int {
	if n < 0 {
		n = -n
	}
	return n % 10
}
