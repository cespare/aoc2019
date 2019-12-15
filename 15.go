package main

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	type qstate struct {
		ic    *intcode
		steps int
		pos   ivec2
	}

	var oxygen ivec2
	visited := map[ivec2]struct{}{{0, 0}: {}}
	ic := newIntcodeWithMem(prog)
	ic.setSuspendMode()
	q := []qstate{{ic, 0, ivec2{0, 0}}}
	for len(q) > 0 {
		pp := q[0]
		q = q[1:]
		for d1 := int64(1); d1 <= 4; d1++ {
			steps := pp.steps + 1
			pos := pp.pos
			switch d1 {
			case dirNorth:
				pos = pos.add(north)
			case dirEast:
				pos = pos.add(east)
			case dirSouth:
				pos = pos.add(south)
			case dirWest:
				pos = pos.add(west)
			default:
				panic("bad")
			}
			if _, ok := visited[pos]; ok {
				continue
			}

			ic := pp.ic.clone()
			result := evalRepairPath(ic, d1)
			switch result {
			case 0:
				continue
			case 1:
			case 2:
				if oxygen != (ivec2{0, 0}) {
					panic("multiple oxygens")
				}
				ctx.reportPart1(steps)
				oxygen = pos
			default:
				panic("bad output")
			}
			visited[pos] = struct{}{}
			q = append(q, qstate{ic, steps, pos})
		}
		pp.ic.free()

	}
	if oxygen == (ivec2{0, 0}) {
		panic("no path found")
	}

	ctx.reportPart2(findDepth(visited, oxygen))
}

const (
	dirNorth = 1
	dirSouth = 2
	dirWest  = 3
	dirEast  = 4
)

func evalRepairPath(ic *intcode, d int64) int64 {
	ic.input = append(ic.input, d)
	ic.run()
	if len(ic.output) != 1 {
		panic("bad")
	}
	r := ic.output[0]
	ic.output = ic.output[:0]
	return r
}

func findDepth(m map[ivec2]struct{}, start ivec2) int64 {
	type posAndDepth struct {
		pos   ivec2
		depth int64
	}
	visited := map[ivec2]struct{}{start: {}}
	var maxDepth int64
	q := []posAndDepth{{pos: start, depth: 0}}
	for len(q) > 0 {
		pd := q[0]
		q = q[1:]
		if pd.depth < maxDepth {
			panic("waddafak")
		}
		if pd.depth > maxDepth {
			maxDepth = pd.depth
		}
		for _, d := range []ivec2{north, east, south, west} {
			p1 := pd.pos.add(d)
			if _, ok := visited[p1]; ok {
				continue
			}
			if _, ok := m[p1]; !ok {
				continue
			}
			visited[p1] = struct{}{}
			q = append(q, posAndDepth{p1, pd.depth + 1})
		}
	}
	return maxDepth
}
