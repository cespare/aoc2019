package main

import (
	"bufio"
	"log"
	"math/big"
	"strconv"
	"strings"
)

func init() {
	addSolutions(22, problem22)
}

func problem22(ctx *problemContext) {
	var instrs []shufInstr
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		instrs = append(instrs, parseShufInstr(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	d := newDeck(10007)
	for _, instr := range instrs {
		d.apply(instr)
	}
	it := d.iter()
	for i := int64(0); i < d.n; i++ {
		if c := it.next(); c == 2019 {
			ctx.reportPart1(i)
			break
		}
	}

	d = newDeck(119315717514047)
	for _, instr := range instrs {
		d.apply(instr)
	}
	d.repeat(101741582076661)

	it = d.iter()
	for c := int64(0); c < 2020; c++ {
		it.next()
	}
	ctx.reportPart2(it.next())
}

func trimPrefix(s, prefix string) (trimmed string, hadPrefix bool) {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return s, false
}

type shufOp int

const (
	shufDealNewStack shufOp = iota
	shufCutN
	shufDealWithIncr
)

type shufInstr struct {
	op shufOp
	n  int64
}

func parseShufInstr(s string) shufInstr {
	if s == "deal into new stack" {
		return shufInstr{op: shufDealNewStack}
	}
	if s, ok := trimPrefix(s, "cut "); ok {
		n := mustParseInt64(s)
		return shufInstr{shufCutN, n}
	}
	if s, ok := trimPrefix(s, "deal with increment "); ok {
		n := mustParseInt64(s)
		return shufInstr{shufDealWithIncr, n}
	}
	panic("bad instr: " + s)
}

func mustParseInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

type deck struct {
	n     int64
	bn    *big.Int
	start int64
	jump  int64
}

func newDeck(n int64) *deck {
	return &deck{
		n:     n,
		bn:    big.NewInt(n),
		start: 0,
		jump:  1,
	}
}

func (d *deck) apply(instr shufInstr) {
	switch instr.op {
	case shufDealNewStack:
		d.jump = d.n - d.jump
		d.cutN(1)
	case shufCutN:
		d.cutN(instr.n)
	case shufDealWithIncr:
		// jump is the inverse of the "deal with" increment.
		jump := d.modInverse(instr.n)
		d.jump = d.modMul(d.jump, jump)
	default:
		panic("unreached")
	}
}

func (d *deck) cutN(n int64) {
	d.start = d.modAdd(d.start, d.modMul(d.jump, n))
}

// repeat reapplies the current transformation k times total (so k=1 means no
// change).
func (d *deck) repeat(k int64) {
	s1 := d.start
	m1 := d.jump

	// State is represented as two numbers, s ("start") and m ("jump").
	// s is the starting card. m is the jump between cards. If you start
	// with an unshuffled deck and deal with increment d, then m is the
	// modular inverse of d. (All the math here is mod(decksize), which is a
	// prime.)
	//
	// m2 = m1^2 mod n
	// m3 = m1^3 mod n
	// ...
	// mk = m1^k mod n
	//
	// s is harder.
	//
	// s2 = (s1 + s1*m1) mod n
	// s3 = (s2 + s1*m2) mod n
	// s4 = (s3 + s1*m3) mod n
	// ...
	// sk = s1 * [ m1^0 + m1^1 + ... + m1^(k-1) ]
	//
	// the series is geometric:
	//
	// sk = s1 * [ (1 - m1^k) / (1 - m1) ]
	//
	// and the division may be computed with a modular inverse:
	//
	// sk = s1 * (1 - m1^k) * modinverse(1 - m1)

	d.jump = d.modExp(m1, k)

	x := d.modAdd(1, -d.jump) // 1 - m1^k
	y := d.modAdd(1, -m1)     // 1 - m1
	d.start = d.modMul(d.modMul(s1, x), d.modInverse(y))
}

// All the math is mod d.n, and so all the inputs/outputs fit in int64s.

func (d *deck) modInverse(g int64) int64 {
	bg := big.NewInt(g)
	h := big.NewInt(0)
	if h.ModInverse(bg, d.bn) == nil {
		panic("no inverse")
	}
	return h.Int64()
}

func (d *deck) modExp(x, y int64) int64 {
	r := big.NewInt(x)
	return r.Exp(r, big.NewInt(y), d.bn).Int64()
}

func (d *deck) modAdd(x, y int64) int64 {
	sum := big.NewInt(x)
	sum.Add(sum, big.NewInt(y))
	return sum.Mod(sum, d.bn).Int64()
}

func (d *deck) modMul(x, y int64) int64 {
	r := big.NewInt(x)
	r.Mul(r, big.NewInt(y))
	return r.Mod(r, d.bn).Int64()
}

type deckIterator struct {
	d *deck
	n int64
}

func (d *deck) iter() *deckIterator {
	return &deckIterator{
		d: d,
		n: d.start,
	}
}

func (it *deckIterator) next() int64 {
	n := it.n
	it.n = it.d.modAdd(it.n, it.d.jump)
	return n
}
