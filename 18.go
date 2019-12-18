package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	v := newVault()
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		row := []byte(strings.TrimSpace(scanner.Text()))
		v.addRow(row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	ctx.reportPart1(v.copy().shortestKeyPath())

	p := v.p[0]
	v.p = []ivec2{
		{p.x - 1, p.y - 1},
		{p.x + 1, p.y - 1},
		{p.x - 1, p.y + 1},
		{p.x + 1, p.y + 1},
	}
	v.m[p.y][p.x] = '#'
	v.m[p.y][p.x-1] = '#'
	v.m[p.y][p.x+1] = '#'
	v.m[p.y-1][p.x] = '#'
	v.m[p.y+1][p.x] = '#'
	ctx.reportPart2(v.shortestKeyPath())
}

type vault struct {
	// immutable after construction
	m [][]byte
	w int
	h int

	// mutable
	p          []ivec2
	keys       map[byte]ivec2
	doors      map[byte]ivec2
	keysByPos  map[ivec2]byte
	doorsByPos map[ivec2]byte
}

func newVault() *vault {
	return &vault{
		keys:       make(map[byte]ivec2),
		doors:      make(map[byte]ivec2),
		keysByPos:  make(map[ivec2]byte),
		doorsByPos: make(map[ivec2]byte),
	}
}

func (v *vault) copy() *vault {
	v1 := &vault{
		m: v.m,
		w: v.w,
		h: v.h,

		p:          append([]ivec2(nil), v.p...),
		keys:       make(map[byte]ivec2),
		doors:      make(map[byte]ivec2),
		keysByPos:  make(map[ivec2]byte),
		doorsByPos: make(map[ivec2]byte),
	}
	for k, v := range v.keys {
		v1.keys[k] = v
	}
	for k, v := range v.doors {
		v1.doors[k] = v
	}
	for k, v := range v.keysByPos {
		v1.keysByPos[k] = v
	}
	for k, v := range v.doorsByPos {
		v1.doorsByPos[k] = v
	}
	return v1
}

func (v *vault) addRow(row []byte) {
	if v.h > 0 && len(row) != v.w {
		panic("mismatched rows")
	}
	v.m = append(v.m, row)
	v.w = len(row)
	y := v.h
	v.h++

	for x, c := range row {
		p := ivec2{int64(x), int64(y)}
		switch {
		case c == '#' || c == '.':
		case c == '@':
			v.p = append(v.p, p)
			row[x] = '.'
		case c >= 'A' && c <= 'Z':
			v.doors[c] = p
			v.doorsByPos[p] = c
			row[x] = '.'
		case c >= 'a' && c <= 'z':
			v.keys[c] = p
			v.keysByPos[p] = c
			row[x] = '.'
		default:
			fmt.Println(string(c))
			panic("bad input byte")
		}
	}
}

func (v *vault) shortestKeyPath() int64 {
	q := newVaultQueue()
	q.add(v, 0)
	for q.Len() > 0 {
		v, d := q.pop()
		if len(v.keys) == 0 {
			return d
		}
		pds := v.reachableKeys()
		for _, pd := range pds {
			v1 := v.copy()
			key, ok := v1.keysByPos[pd.p]
			if !ok {
				panic("no key")
			}
			delete(v1.keys, key)
			delete(v1.keysByPos, pd.p)
			d1 := d + pd.d
			door := key + 'A' - 'a'
			if doorPos, ok := v1.doors[door]; ok {
				delete(v1.doors, door)
				delete(v1.doorsByPos, doorPos)
			}
			v1.p[pd.robotIdx] = pd.p
			q.add(v1, d1)
		}
	}
	panic("no path")
}

type vaultQueueState struct {
	p    string
	keys uint32 // bitset, 'a' is bit 0
}

func (v *vault) queueState() vaultQueueState {
	state := vaultQueueState{p: fmt.Sprintf("%v", v.p)}
	for key := range v.keys {
		state.keys |= uint32(1) << (key - 'a')
	}
	return state
}

func (v *vault) queueItem(d int64) *vaultQueueItem {
	return &vaultQueueItem{
		v:     v,
		idx:   -1,
		state: v.queueState(),
		d:     d,
	}
}

type vaultQueue struct {
	q []*vaultQueueItem
	m map[vaultQueueState]*vaultQueueItem
}

func newVaultQueue() *vaultQueue {
	return &vaultQueue{m: make(map[vaultQueueState]*vaultQueueItem)}
}

type vaultQueueItem struct {
	v     *vault
	idx   int
	state vaultQueueState
	d     int64
}

func (q *vaultQueue) add(v *vault, d int64) {
	item := v.queueItem(d)
	prev, ok := q.m[item.state]
	if !ok {
		heap.Push(q, item)
		return
	}
	if d < prev.d {
		prev.d = d
		heap.Fix(q, prev.idx)
	}
}

func (q *vaultQueue) pop() (*vault, int64) {
	item := heap.Pop(q).(*vaultQueueItem)
	return item.v, item.d
}

func (q *vaultQueue) Len() int           { return len(q.q) }
func (q *vaultQueue) Less(i, j int) bool { return q.q[i].d < q.q[j].d }
func (q *vaultQueue) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].idx = i
	q.q[j].idx = j
}
func (q *vaultQueue) Push(x interface{}) {
	item := x.(*vaultQueueItem)
	item.idx = len(q.q)
	q.q = append(q.q, item)
	q.m[item.state] = item
}
func (q *vaultQueue) Pop() interface{} {
	item := q.q[len(q.q)-1]
	q.q = q.q[:len(q.q)-1]
	delete(q.m, item.state)
	item.idx = -1
	return item
}

func (v *vault) reachableKeys() []robotPosDistance {
	var keys []robotPosDistance
	for robotIdx, p := range v.p {
		visited := map[ivec2]struct{}{p: {}}
		q := []robotPosDistance{{robotIdx, p, 0}}
		for len(q) > 0 {
			pd := q[0]
			q = q[1:]
			for _, p1 := range v.reachable(pd.p, visited) {
				visited[p1] = struct{}{}
				pd1 := robotPosDistance{robotIdx, p1, pd.d + 1}
				if _, ok := v.keysByPos[p1]; ok {
					keys = append(keys, pd1)
				} else {
					q = append(q, pd1)
				}
			}
		}
	}
	return keys
}

type robotPosDistance struct {
	robotIdx int
	p        ivec2
	d        int64
}

func (v *vault) reachable(p ivec2, visited map[ivec2]struct{}) []ivec2 {
	neighbors := make([]ivec2, 0, 2)
	for _, d := range []ivec2{
		{-1, 0},
		{0, -1},
		{1, 0},
		{0, 1},
	} {
		p1 := p.add(d)
		if p1.x < 0 || p1.x >= int64(v.w) || p1.y < 0 || p1.y >= int64(v.h) {
			continue
		}
		if _, ok := visited[p1]; ok {
			continue
		}
		if _, ok := v.doorsByPos[p1]; ok {
			continue
		}
		if v.m[p1.y][p1.x] == '#' {
			continue
		}
		neighbors = append(neighbors, p1)
	}
	return neighbors
}

func (v *vault) at(p ivec2) byte {
	for _, p1 := range v.p {
		if p1 == p {
			return '@'
		}
	}
	if c, ok := v.doorsByPos[p]; ok {
		return c
	}
	if c, ok := v.keysByPos[p]; ok {
		return c
	}
	return v.m[p.y][p.x]
}

func (v *vault) print() {
	out := make([]byte, v.w+1)
	out[len(out)-1] = '\n'
	for y, row := range v.m {
		for x := range row {
			out[x] = v.at(ivec2{int64(x), int64(y)})
		}
		os.Stdout.Write(out)
	}
}
