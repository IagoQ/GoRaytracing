package main

import (
	"math"
)

type Camera struct {
	vph, vpw, fov                 float64
	pos, dh, dv, corner           Vector
	aspect, aperture, focus, lens float64
	w, u, v                       Vector
}

func CreateCamera(aspect, fov float64, Pos, look, up Vector, aperture, focusDistance float64) Camera {
	var c Camera

	theta := degToRad(fov)
	c.vph = 2.0 * math.Tan(theta/2.0)
	c.vpw = aspect * c.vph

	c.w = Pos.sub(look).normalize()
	c.u = up.cross(c.w).normalize()
	c.v = c.w.cross(c.u)

	c.dh = c.u.scalarMult(c.vpw * focusDistance)
	c.dv = c.v.scalarMult(c.vph * focusDistance)
	c.corner = Pos.sub(c.dh.scalarMult(0.5)).sub(c.dv.scalarMult(0.5)).sub(c.w.scalarMult(focusDistance))

	c.fov = fov
	c.pos = Pos
	c.aspect = aspect
	c.aperture = aperture
	c.lens = aperture / 2
	c.focus = focusDistance

	return c
}

func (c Camera) getRay(s, t float64) Ray {
	rd := randomDiskVector().scalarMult(c.aperture / 2)
	offset := c.u.scalarMult(rd.x).add(c.v.scalarMult(rd.y))

	return Ray{c.pos.add(offset), c.corner.add(c.dh.scalarMult(s)).add(c.dv.scalarMult(t)).sub(c.pos).sub(offset)}
}

func degToRad(d float64) float64 {
	return math.Pi * d / 180
}
