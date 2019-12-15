package main

import (
	"bufio"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

func init() {
	addSolutions(14, problem14)
}

func problem14(ctx *problemContext) {
	fr := make(fuelReactions)
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		if err := fr.parseLine(scanner.Text()); err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	ctx.reportPart1(fr.calcOre(1))

	moreThan1T := sort.Search(99999999, func(fuel int) bool {
		return fr.calcOre(int64(fuel)) > 1e12
	})
	ctx.reportPart2(moreThan1T - 1)
}

type chemAmount struct {
	n    int64
	chem string
}

func parseChemAmount(s string) (chemAmount, error) {
	var ca chemAmount
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return ca, fmt.Errorf("bad chemAmount %q", s)
	}
	var err error
	ca.n, err = strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return ca, err
	}
	ca.chem = fields[1]
	return ca, nil
}

type reaction struct {
	inputs []chemAmount
	output chemAmount
}

type fuelReactions map[string]reaction // by output chem

func (fr fuelReactions) parseLine(s string) error {
	parts := strings.Split(s, "=>")
	if len(parts) != 2 {
		return fmt.Errorf("bad line %q", s)
	}
	var r reaction
	var err error
	r.output, err = parseChemAmount(strings.TrimSpace(parts[1]))
	if err != nil {
		return err
	}
	for _, is := range strings.Split(parts[0], ",") {
		in, err := parseChemAmount(strings.TrimSpace(is))
		if err != nil {
			return err
		}
		r.inputs = append(r.inputs, in)
	}
	out := r.output.chem
	if _, ok := fr[out]; ok {
		return fmt.Errorf("multiple outputs for %q", out)
	}
	fr[out] = r
	return nil
}

func (fr fuelReactions) calcOre(n int64) int64 {
	extra := make(map[string]int64)
	var q reactionQueue
	q.push(chemAmount{n, "FUEL"})
	var ore int64
	for !q.empty() {
		ca := q.pop()
		if ca.chem == "ORE" {
			ore += ca.n
			continue
		}
		ex := extra[ca.chem]
		if ex >= ca.n {
			extra[ca.chem] = ex - ca.n
			continue
		}
		ca.n -= ex

		rx := fr[ca.chem]
		m := ((ca.n - 1) / rx.output.n) + 1 // multiple of the # of reactions
		extra[ca.chem] = (m * rx.output.n) - ca.n
		for _, in := range rx.inputs {
			need := in
			need.n *= m
			q.push(need)
		}
	}
	return ore
}

type reactionQueue []chemAmount

func (q *reactionQueue) push(ca chemAmount) {
	for i, ca1 := range *q {
		if ca1.chem == ca.chem {
			(*q)[i].n += ca.n
			return
		}
	}
	*q = append(*q, ca)
}

func (q *reactionQueue) pop() chemAmount {
	ca := (*q)[0]
	*q = (*q)[1:]
	return ca
}

func (q *reactionQueue) empty() bool {
	return len(*q) == 0
}
