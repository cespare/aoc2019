package main

func init() {
	addSolutions(17, problem17)
}

func problem17(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	ic := newIntcodeWithMem(prog)
	ic.run()
	var pic [][]byte
	var row []byte
	for _, n := range ic.output {
		c := byte(n)
		if c == '\n' {
			if len(row) > 0 {
				pic = append(pic, row)
			}
			row = nil
		} else {
			row = append(row, c)
		}
	}
	if len(row) > 0 {
		pic = append(pic, row)
	}
	// for _, row := range pic {
	// 	fmt.Println(string(row))
	// }
	var part1 int
	for y, row := range pic {
		for x, c := range row {
			if y == 0 || y == len(pic)-1 || x == 0 || x == len(row)-1 {
				continue
			}
			if c == '#' && pic[y][x-1] == '#' && pic[y][x+1] == '#' && pic[y-1][x] == '#' && pic[y+1][x] == '#' {
				part1 += x * y
			}
		}
	}
	ctx.reportPart1(part1)

	prog[0] = 2
	start, startDir := vacuumStart(pic)
	var scaffold []scaffoldPiece
	p := start
	for {
		dirs := scaffoldDirs(pic, p)
		if len(dirs) > 2 {
			panic("ambiguous")
		}
		if len(scaffold) > 0 {
			prevDir := scaffold[len(scaffold)-1].dir
			if reverseDir(dirs[0]) == prevDir {
				dirs = dirs[1:]
			} else if len(dirs) == 2 && reverseDir(dirs[1]) == prevDir {
				dirs = dirs[:1]
			}
		}
		if len(dirs) == 2 {
			panic("ambiguous")
		}
		if len(dirs) == 0 {
			break
		}
		piece := scaffoldPiece{dir: dirs[0]}
		dv := dirVec(piece.dir)
		for {
			nextp := p.add(dv)
			if !inBounds(pic, nextp) || at(pic, nextp) != '#' {
				break
			}
			piece.len++
			p = nextp
		}
		scaffold = append(scaffold, piece)
	}
	_ = startDir
	// prevDir := startDir
	// for _, piece := range scaffold {
	// 	turn := turnDir(prevDir, piece.dir)
	// 	prevDir = piece.dir
	// 	fmt.Println(turn, piece.len)
	// }

	// L 10	A
	// L 10
	// R 6
	//
	// L 10	A
	// L 10
	// R 6
	//
	// R 12	B
	// L 12
	// L 12
	//
	// R 12 B
	// L 12
	// L 12
	//
	// L 6	C
	// L 10
	// R 12
	// R 12
	//
	// R 12	B
	// L 12
	// L 12
	//
	// L 6	C
	// L 10
	// R 12
	// R 12
	//
	// R 12	B
	// L 12
	// L 12
	//
	// L 6	C
	// L 10
	// R 12
	// R 12
	//
	// L 10	A
	// L 10
	// R 6

	// Program: A,A,B,B,C,B,C,B,C,A
	// A: L,10,L,10,R,6
	// B: R,12,L,12,L,12
	// C: L,6,L,10,R,12,R,12

	ascii := `
A,A,B,B,C,B,C,B,C,A
L,10,L,10,R,6
R,12,L,12,L,12
L,6,L,10,R,12,R,12
n
`
	ascii = ascii[1:]
	input := make([]int64, len(ascii))
	for i, c := range ascii {
		input[i] = int64(c)
	}
	ic = newIntcodeWithMem(prog, input...)
	ic.run()
	ctx.reportPart2(ic.output[len(ic.output)-1])
}

func turnDir(dir0, dir1 int) string {
	switch [2]int{dir0, dir1} {
	case [2]int{dirNorth, dirWest}:
		return "L"
	case [2]int{dirNorth, dirEast}:
		return "R"
	case [2]int{dirSouth, dirWest}:
		return "R"
	case [2]int{dirSouth, dirEast}:
		return "L"
	case [2]int{dirEast, dirNorth}:
		return "L"
	case [2]int{dirEast, dirSouth}:
		return "R"
	case [2]int{dirWest, dirNorth}:
		return "R"
	case [2]int{dirWest, dirSouth}:
		return "L"
	default:
		panic("bad turn")
	}
}

func reverseDir(dir int) int {
	switch dir {
	case dirNorth:
		return dirSouth
	case dirEast:
		return dirWest
	case dirSouth:
		return dirNorth
	case dirWest:
		return dirEast
	default:
		panic("bad dir")
	}
}

func dirStr(dir int) string {
	switch dir {
	case dirNorth:
		return "north"
	case dirEast:
		return "east"
	case dirSouth:
		return "south"
	case dirWest:
		return "west"
	default:
		panic("bad dir")
	}
}

func dirVec(dir int) ivec2 {
	switch dir {
	case dirNorth:
		return ivec2{0, -1}
	case dirEast:
		return ivec2{1, 0}
	case dirSouth:
		return ivec2{0, 1}
	case dirWest:
		return ivec2{-1, 0}
	default:
		panic("bad dir")
	}
}

func at(pic [][]byte, p ivec2) byte {
	return pic[p.y][p.x]
}

// func (v ivec2) mul(n int64) ivec2 {
// 	return ivec2{v.x * n, v.y * n}
// }

func vacuumStart(pic [][]byte) (pos ivec2, dir int) {
	for y, row := range pic {
		for x, c := range row {
			pos = ivec2{int64(x), int64(y)}
			switch c {
			case '^':
				return pos, dirNorth
			case 'v':
				return pos, dirSouth
			case '<':
				return pos, dirWest
			case '>':
				return pos, dirEast
			}
		}
	}
	panic("not found")
}

func inBounds(pic [][]byte, p ivec2) bool {
	return p.x >= 0 && p.x < int64(len(pic[0])) && p.y >= 0 && p.y < int64(len(pic))
}

func scaffoldDirs(pic [][]byte, p ivec2) []int {
	w := int64(len(pic[0]))
	h := int64(len(pic))
	var dirs []int
	if p.x > 0 && pic[p.y][p.x-1] == '#' {
		dirs = append(dirs, dirWest)
	}
	if p.x < w-1 && pic[p.y][p.x+1] == '#' {
		dirs = append(dirs, dirEast)
	}
	if p.y > 0 && pic[p.y-1][p.x] == '#' {
		dirs = append(dirs, dirNorth)
	}
	if p.y < h-1 && pic[p.y+1][p.x] == '#' {
		dirs = append(dirs, dirSouth)
	}
	return dirs
}

type scaffoldPiece struct {
	dir int
	len int
}
