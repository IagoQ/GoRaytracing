package main

import (
	"math"
	"math/rand"
)

type Material interface {
	scatter(*Ray, *HitRec, *Color, *Ray) bool
}

type Dielectric struct {
	index float64
	c     Color
}

func (d Dielectric) scatter(r *Ray, hit *HitRec, c *Color, scattered *Ray) bool {
	(*c) = d.c
	var eta float64
	if hit.front {
		eta = 1.0 / d.index
	} else {
		eta = d.index
	}
	refracted := refract(r.dir.normalize(), hit.normal, eta)
	(*scattered) = Ray{hit.point, refracted}
	return true
}

type FuzzyMirror struct {
	fuzzy float64
	c     Color
}

func (f FuzzyMirror) scatter(r *Ray, hit *HitRec, c *Color, scattered *Ray) bool {

	v := r.dir.normalize()
	n := hit.normal

	reflected := v.sub(n.scalarMult(2.0 * v.dot(n)))

	(*scattered) = Ray{hit.point, reflected.add(randomUnitVector().scalarMult(f.fuzzy))}
	*c = f.c
	return scattered.dir.dot(hit.normal) > 0
}

type Matte struct {
	reflectance float64
	c           Color
}

func (m Matte) scatter(r *Ray, hit *HitRec, c *Color, scattered *Ray) bool {

	if rand.Float64() < m.reflectance {
		scatterdir := hit.normal.add(randomUnitVector())
		*scattered = Ray{hit.point, scatterdir}
		*c = m.c
		return true
	}
	*c = Color{0, 0, 0}
	return false

}

type Mirror struct {
	c Color
}

func (m Mirror) scatter(r *Ray, hit *HitRec, c *Color, scattered *Ray) bool {

	reflected := reflect(r.dir.normalize(), hit.normal)

	(*scattered) = Ray{hit.point, reflected}
	*c = m.c
	return scattered.dir.dot(hit.normal) > 0
}

func reflect(v Vector, n Vector) Vector {
	return v.sub(n.scalarMult(2.0 * v.dot(n)))
}

func refract(uv, n Vector, index float64) Vector {
	cosTheta := uv.scalarMult(-1).dot(n)
	out_perp := uv.add(n.scalarMult(cosTheta)).scalarMult(index)
	out_parallel := n.scalarMult(-1 * math.Sqrt(math.Abs(1.0-out_perp.lengthSquared())))
	return out_perp.add(out_parallel)
}
