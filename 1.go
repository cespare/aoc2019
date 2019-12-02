package main

import (
	"bufio"
	"log"
	"strconv"
)

func init() {
	addSolutions(1, problem1)
}

func problem1(ctx *problemContext) {
	var masses []int64
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		n, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatalln("Bad number:", scanner.Text())
		}
		masses = append(masses, n)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	var total int64
	for _, m := range masses {
		total += m/3 - 2
	}
	ctx.reportPart1(total)

	var total2 int64
	for _, m := range masses {
		total2 += fuelForMass(m)
	}
	ctx.reportPart2(total2)
}

func fuelForMass(mass int64) int64 {
	f := mass/3 - 2
	if f < 0 {
		return 0
	}
	return f + fuelForMass(f)
}
