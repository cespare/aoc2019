package main

func init() {
	addSolutions(12, problem12)
}

func problem12(ctx *problemContext) {
	pos := []ivec3{
		{1, 2, -9},
		{-1, -9, -4},
		{17, 6, 8},
		{12, 4, 2},
	}
	vel := make([]ivec3, 4)
	ctx.reportLoad()

	m := moons{
		pos: append([]ivec3(nil), pos...),
		vel: append([]ivec3(nil), vel...),
	}
	for i := 0; i < 1000; i++ {
		m.applyGrav()
		m.move()
	}
	var energy int64
	for i := range m.pos {
		energy += m.energy(i)
	}
	ctx.reportPart1(energy)

	px := [4]int64{pos[0].x, pos[1].x, pos[2].x, pos[3].x}
	py := [4]int64{pos[0].y, pos[1].y, pos[2].y, pos[3].y}
	pz := [4]int64{pos[0].z, pos[1].z, pos[2].z, pos[3].z}
	vx := [4]int64{vel[0].x, vel[1].x, vel[2].x, vel[3].x}
	vy := [4]int64{vel[0].y, vel[1].y, vel[2].y, vel[3].y}
	vz := [4]int64{vel[0].z, vel[1].z, vel[2].z, vel[3].z}

	cx := cycleLength(px, vx)
	cy := cycleLength(py, vy)
	cz := cycleLength(pz, vz)
	ctx.reportPart2(lcm(lcm(cx, cy), cz))
}

type ivec3 struct {
	x int64
	y int64
	z int64
}

func (v ivec3) add(v1 ivec3) ivec3 {
	return ivec3{
		v.x + v1.x,
		v.y + v1.y,
		v.z + v1.z,
	}
}

type moons struct {
	pos []ivec3
	vel []ivec3
}

func (m moons) applyGrav() {
	for i, p0 := range m.pos {
		for j := i + 1; j < len(m.pos); j++ {
			p1 := m.pos[j]
			v0, v1 := &m.vel[i], &m.vel[j]
			adjustForGrav(p0.x, p1.x, &v0.x, &v1.x)
			adjustForGrav(p0.y, p1.y, &v0.y, &v1.y)
			adjustForGrav(p0.z, p1.z, &v0.z, &v1.z)
		}
	}
}

func adjustForGrav(p0, p1 int64, v0, v1 *int64) {
	switch {
	case p0 < p1:
		*v0++
		*v1--
	case p0 > p1:
		*v0--
		*v1++
	}
}

func (m moons) move() {
	for i, p := range m.pos {
		m.pos[i] = p.add(m.vel[i])
	}
}

func (m moons) energy(i int) int64 {
	p := m.pos[i]
	v := m.vel[i]
	return (iabs(p.x) + iabs(p.y) + iabs(p.z)) * (iabs(v.x) + iabs(v.y) + iabs(v.z))
}

func cycleLength(p, v [4]int64) int64 {
	var initial [8]int64
	copy(initial[:4], p[:])
	copy(initial[4:], v[:])
	var state [8]int64
	for i := int64(0); ; i++ {
		copy(state[:4], p[:])
		copy(state[4:], v[:])
		if i > 0 && state == initial {
			return i
		}
		applyGravSingle(&p, &v)
		for j, vv := range v {
			p[j] += vv
		}
	}
}

func applyGravSingle(p, v *[4]int64) {
	for i, p0 := range p {
		for j := i + 1; j < 4; j++ {
			p1 := p[j]
			adjustForGrav(p0, p1, &v[i], &v[j])
		}
	}
}

func gcd(m, n int64) int64 {
	if m <= 0 || n <= 0 {
		panic("only positive handled")
	}
	switch {
	case m == n:
		return n
	case m > n:
		return gcd(m-n, n)
	default:
		return gcd(m, n-m)
	}

}

func lcm(m, n int64) int64 {
	return m * n / gcd(m, n)
}
