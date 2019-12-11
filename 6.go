package main

import (
	"bufio"
	"log"
	"strings"
)

func init() {
	addSolutions(6, problem6)
}

func problem6(ctx *problemContext) {
	allNodes := make(map[string]struct{})
	g := make(orbitGraph)
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ")")
		if len(parts) != 2 {
			log.Fatalf("bad orbit line: %q", scanner.Text())
		}
		g.add(parts[1], parts[0])
		allNodes[parts[0]] = struct{}{}
		allNodes[parts[1]] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	var total int
	cache := make(map[string]int)
	for n := range allNodes {
		total += g.depth(n, cache)
	}
	ctx.reportPart1(total)

	dists := make(map[string]int) // distance from me to each ancestor up to COM
	n := "YOU"
	for d := 0; ; d++ {
		parent, ok := g[n]
		if !ok {
			if n != "COM" {
				panic("bad")
			}
			break
		}
		dists[parent] = d
		n = parent
	}
	n = "SAN"
	for d := 0; ; d++ {
		n = g[n]
		if d1, ok := dists[n]; ok {
			ctx.reportPart2(d + d1)
			return
		}
	}
}

type orbitGraph map[string]string

func (g orbitGraph) add(from, to string) {
	if _, ok := g[from]; ok {
		panic("bad tree")
	}
	g[from] = to
}

func (g orbitGraph) depth(n string, cache map[string]int) int {
	if d, ok := cache[n]; ok {
		return d
	}
	parent, ok := g[n]
	if !ok {
		return 0
	}
	d := g.depth(parent, cache) + 1
	cache[n] = d
	return d
}
