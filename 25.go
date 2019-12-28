package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	addSolutions(25, problem25)
}

func problem25(ctx *problemContext) {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	var term *intcodeTerm
	if len(b) > 0 && b[0] == '{' {
		term = loadIntcodeTerm(bytes.NewReader(b))
	} else {
		prog := readProg(bytes.NewReader(b))
		term = newIntcodeTerm(newIntcode(prog))
	}
	ctx.reportLoad()

	term.run()
}

type intcodeTerm struct {
	ic           *intcode
	out          bytes.Buffer
	stdinScanner *bufio.Scanner
}

func newIntcodeTerm(ic *intcode) *intcodeTerm {
	t := &intcodeTerm{
		ic:           ic,
		stdinScanner: bufio.NewScanner(os.Stdin),
	}
	ic.out = func(v int64) {
		if v < 0 || v > 255 {
			panic("not ASCII")
		}
		t.out.WriteByte(byte(v))
		if v == '\n' {
			io.Copy(os.Stdout, &t.out)
			t.out.Reset()
		}
	}
	ic.in = func(buf []int64) []int64 {
		for {
			if !t.stdinScanner.Scan() {
				panic("stdin terminated")
			}
			cmd := t.stdinScanner.Text()
			if name, ok := trimPrefix(cmd, "save "); ok {
				if err := t.save(name); err == nil {
					log.Printf("Saved to %s", name)
				} else {
					log.Printf("Error saving: %s", err)
				}
				continue
			}
			for _, c := range cmd {
				buf = append(buf, int64(c))
			}
			buf = append(buf, '\n')
			return buf
		}
	}
	return t
}

func loadIntcodeTerm(r io.Reader) *intcodeTerm {
	var ic intcode
	if err := ic.load(r); err != nil {
		panic(err)
	}
	return newIntcodeTerm(&ic)
}

func (t *intcodeTerm) run() {
	t.ic.run()
}

func (t *intcodeTerm) save(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := t.ic.save(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
