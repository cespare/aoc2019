package main

import (
	"bytes"
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

	ic := newIntcode(program, 1)
	ic.run()
	ctx.reportPart1(ic.output)

	ic = newIntcode(program, 5)
	ic.run()

	ctx.reportPart2(ic.output)
}
