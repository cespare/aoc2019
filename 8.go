package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func init() {
	addSolutions(8, problem8)
}

func problem8(ctx *problemContext) {
	line, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.TrimSpace(line)
	for i, c := range data {
		data[i] = c - '0'
	}
	ctx.reportLoad()

	sif := decodeSIF(data, 25, 6)
	fewestZeros := 999999999999
	var result int
	for layer := range sif {
		zeros := sif.countDigit(layer, 0)
		if zeros < fewestZeros {
			fewestZeros = zeros
			result = sif.countDigit(layer, 1) * sif.countDigit(layer, 2)
		}
	}
	ctx.reportPart1(result)

	im := sif.flatten()

	for _, row := range im {
		for _, c := range row {
			switch c {
			case 0:
				fmt.Print("\x1b[37m██\x1b[0m")
			case 1:
				fmt.Print("\x1b[30m██\x1b[0m")
			case 2:
				fmt.Print("\x1b[32m██\x1b[0m")
			default:
				panic("unexpected")
			}
		}
		fmt.Println()
	}

	ctx.reportPart2("^^^")
}

type sifImage [][][]uint8 // layers -> rows -> cells

func decodeSIF(data []byte, w, h int) sifImage {
	var sif sifImage
	for len(data) > 0 {
		layer := make([][]uint8, h)
		for y := range layer {
			layer[y] = data[:w]
			data = data[w:]
		}
		sif = append(sif, layer)
	}
	return sif
}

func (s sifImage) countDigit(layer int, digit uint8) int {
	var n int
	for _, row := range s[layer] {
		for _, c := range row {
			if c == digit {
				n++
			}
		}
	}
	return n
}

func (s sifImage) flatten() [][]uint8 {
	im := make([][]uint8, len(s[0]))
	for y := range im {
		im[y] = make([]uint8, len(s[0][0]))
		for x := range im[y] {
			im[y][x] = 2 // start transparent
		}
	}
	for layer := len(s) - 1; layer >= 0; layer-- {
		for y, row := range s[layer] {
			for x, c := range row {
				if c != 2 {
					im[y][x] = c
				}
			}
		}
	}
	return im
}
