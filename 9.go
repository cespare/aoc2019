package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
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

	ic := newIntcodeWithMem(prog, 1)
	ic.run()
	ctx.reportPart1(ic.output)

	ic = newIntcodeWithMem(prog, 2)
	ic.run()
	ctx.reportPart2(ic.output)
}
