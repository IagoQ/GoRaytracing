package main

type Camera struct {
	viewportHeight, viewportWidth, focalLength float64
	pos, h, v, corner                          Vector
}

func CreateCamera(ViewW, ViewH, Focal float64, Pos Vector) Camera {
	vpw := ViewW
	vph := ViewH
	fl := Focal
	h := Vector{ViewW, 0, 0}
	v := Vector{0, ViewH, 0}
	lowerleft := Pos.sub(h.scalarMult(0.5)).sub(v.scalarMult(0.5)).sub(Vector{0, 0, fl})
	return Camera{vpw, vph, fl, Pos, h, v, lowerleft}
}

func (c Camera) getRay(u, v float64) Ray {
	return Ray{c.pos, c.corner.add(c.h.scalarMult(u)).add(c.v.scalarMult(v)).sub(c.pos)}
}
