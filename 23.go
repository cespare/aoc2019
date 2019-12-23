package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	prog := readProg(ctx.f)
	ctx.reportLoad()

	start := make(chan struct{})
	nw := newNetwork(50)
	for addr := int64(0); addr < 50; addr++ {
		c := nw.addComputer(prog, addr)
		go func() {
			<-start
			c.ic.run()
			log.Panicf("computer %d halted", c.addr)
		}()
	}
	close(start)
	nw.monitor()
}

type network struct {
	cs      []*netComputer
	natMu   sync.Mutex
	nat     *packet
	lastNAT *packet
}

func newNetwork(n int) *network {
	nw := &network{cs: make([]*netComputer, n)}
	return nw
}

func (nw *network) monitor() {
	for {
		log.Println("Waiting for idle...")
	idleLoop:
		for {
			nw.natMu.Lock()
			natSet := nw.nat != nil
			nw.natMu.Unlock()
			if !natSet {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			for _, c := range nw.cs {
				if atomic.LoadInt64(&c.idle) < 1000 {
					time.Sleep(time.Millisecond)
					continue idleLoop
				}
			}
			break
		}
		nw.natMu.Lock()
		if nw.nat != nil && nw.lastNAT != nil && nw.nat.y == nw.lastNAT.y {
			log.Panicln("part 2:", nw.nat.y)
		}
		nw.lastNAT = nw.nat
		nat := nw.nat
		nw.nat = nil
		nw.natMu.Unlock()
		log.Printf("{%d, %d}", nat.x, nat.y)
		c0 := nw.cs[0]
		c0.mu.Lock()
		c0.in = append(c0.in, nat.x, nat.y)
		c0.mu.Unlock()
	}
}

type netComputer struct {
	nw   *network
	addr int64
	mu   sync.Mutex
	in   []int64
	out  []int64
	idle int64
	ic   *intcode
}

func (nw *network) addComputer(prog []int64, addr int64) *netComputer {
	c := &netComputer{
		nw:   nw,
		addr: addr,
	}
	nw.cs[addr] = c
	c.ic = newIntcode(prog)
	c.ic.setInput(addr)
	c.ic.in = func(buf []int64) []int64 {
		if addr >= 0 {
			buf = append(buf, addr)
			addr = -1
			return buf
		}
		c.mu.Lock()
		if len(c.in) > 0 {
			atomic.SwapInt64(&c.idle, 0)
			buf = append(buf, c.in...)
			c.in = c.in[:0]
			c.mu.Unlock()
			return buf
		}
		c.mu.Unlock()
		time.Sleep(500 * time.Microsecond)
		atomic.AddInt64(&c.idle, 1)
		return append(buf, -1)
	}
	c.ic.out = func(v int64) {
		c.out = append(c.out, v)
		if len(c.out) < 3 {
			return
		}
		addr, x, y := c.out[0], c.out[1], c.out[2]
		// log.Printf("computer %d sending {%d, %d} to %d", c.addr, x, y, addr)
		c.out = c.out[:0]
		if addr == 255 {
			c.nw.natMu.Lock()
			c.nw.nat = &packet{x, y}
			c.nw.natMu.Unlock()
			return
		}
		peer := c.nw.cs[addr]
		peer.mu.Lock()
		peer.in = append(peer.in, x, y)
		peer.mu.Unlock()
	}
	return c
}

type packet struct {
	x int64
	y int64
}
