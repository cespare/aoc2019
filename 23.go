package main

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	nw := newNetwork(prog, 50)
	reportedPart1 := false
	for {
		if part2 := nw.step(); part2 >= 0 {
			ctx.reportPart2(part2)
			return
		}
		if !reportedPart1 && nw.firstNAT != nil {
			ctx.reportPart1(nw.firstNAT.y)
			reportedPart1 = true
		}
	}
}

type network struct {
	cs []*netComputer

	nat         *packet
	firstNAT    *packet
	prevNATYVal int64
}

func newNetwork(prog []int64, n int) *network {
	nw := &network{
		cs: make([]*netComputer, n),
	}
	for i := range nw.cs {
		ic := newIntcode(prog)
		addr := int64(i)
		c := &netComputer{
			nw:   nw,
			addr: addr,
			ic:   ic,
		}
		ic.setInput(addr)
		ic.in = c.in
		ic.out = c.out
		nw.cs[i] = c
	}
	return nw
}

func (nw *network) step() (part2 int64) {
	idle := true
	for _, c := range nw.cs {
		if c.idleIn < 100 || c.cyclesSinceOut < 100 {
			idle = false
			break
		}
	}
	if nw.nat == nil {
		idle = false
	}
	if idle {
		if nw.nat.y == nw.prevNATYVal {
			return nw.nat.y
		}
		nw.prevNATYVal = nw.nat.y
		nw.cs[0].ic.setInput(nw.nat.x, nw.nat.y)
	}
	for _, c := range nw.cs {
		if idle {
			c.idleIn = 0
		}
		c.cyclesSinceOut++
		c.step()
	}
	return -1
}

type netComputer struct {
	nw   *network
	addr int64
	ic   *intcode

	outBuf         []int64
	idleIn         int
	cyclesSinceOut int
}

type packet struct {
	x int64
	y int64
}

func (c *netComputer) in(buf []int64) []int64 {
	c.idleIn++
	return append(buf, -1)
}

func (c *netComputer) out(v int64) {
	c.cyclesSinceOut = 0
	c.outBuf = append(c.outBuf, v)
	if len(c.outBuf) < 3 {
		return
	}
	addr, x, y := c.outBuf[0], c.outBuf[1], c.outBuf[2]
	c.outBuf = c.outBuf[:0]
	if addr == 255 {
		p := &packet{x: x, y: y}
		c.nw.nat = p
		if c.nw.firstNAT == nil {
			c.nw.firstNAT = p
		}
		return
	}
	peer := c.nw.cs[addr]
	peer.ic.setInput(x, y)
}

func (c *netComputer) step() {
	if c.ic.step(false) == stateHalt {
		panic("unexpected halt")
	}
}
