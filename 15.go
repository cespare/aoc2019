package main

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	var oxygen ivec2
	visited := map[ivec2]struct{}{{0, 0}: {}}
	q := []pathAndPosition{{nil, ivec2{0, 0}}}
	for len(q) > 0 {
		pp := q[0]
		q = q[1:]
	dirLoop:
		for d1 := int64(1); d1 <= 4; d1++ {
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

			path := append(copyInt64s(pp.path), d1)
			result := evalRepairPath(prog, path)
			switch result {
			case 0:
				continue dirLoop
			case 1:
			case 2:
				if oxygen != (ivec2{0, 0}) {
					panic("multiple oxygens")
				}
				ctx.reportPart1(len(path))
				oxygen = pos
			default:
				panic("bad output")
			}
			visited[pos] = struct{}{}
			q = append(q, pathAndPosition{path, pos})
		}

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

type pathAndPosition struct {
	path []int64
	pos  ivec2
}

func evalRepairPath(prog, path []int64) int64 {
	ic := newIntcodeWithMem(prog, copyInt64s(path)...)
	ic.setSuspendMode()
	ic.run()
	if len(ic.output) != len(path) {
		panic("mismatched output")
	}
	for _, d := range ic.output[:len(ic.output)-1] {
		if d == 0 {
			panic("unexpectedly got wall")
		}
	}
	return ic.output[len(ic.output)-1]
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
