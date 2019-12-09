package main

func init() {
	addSolutions(4, problem4)
}

func problem4(ctx *problemContext) {
	var combos int
	for n := 372304; n <= 847060; n++ {
		if passwordOK(n) {
			combos++
		}
	}
	ctx.reportPart1(combos)

	combos = 0
	for n := 372304; n <= 847060; n++ {
		if passwordOK2(n) {
			combos++
		}
	}
	ctx.reportPart2(combos)
}

func passwordOK(n int) bool {
	double := false
	var d int
	prev := n % 10
	n /= 10
	for n > 0 {
		d, n = n%10, n/10
		if d > prev {
			return false
		}
		if d == prev {
			double = true
		}
		prev = d
	}
	return double
}

func passwordOK2(n int) bool {
	double := false
	repeated := 1
	var d int
	prev := n % 10
	n /= 10
	for n > 0 {
		d, n = n%10, n/10
		if d > prev {
			return false
		}
		if d == prev {
			repeated++
		} else {
			if repeated == 2 {
				double = true
			}
			repeated = 1
		}
		prev = d
	}
	if repeated == 2 {
		double = true
	}
	return double
}
