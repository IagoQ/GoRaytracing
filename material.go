package main

import (
	"math"
	"math/rand"
)

// TODO: add texture to materials

type Material interface {
	scatter(*Ray, *HitRec, *Color, *Ray) bool
	emmited(float64, float64, Vector) Color
}

type Diffuselight struct {
	c Color
}

func (dl Diffuselight) emmited(u, v float64, p Vector) Color {
	return dl.c
}

func (dl Diffuselight) scatter(r *Ray, hit *HitRec, c *Color, scattered *Ray) bool {

	// scatterdir := hit.normal.add(randomUnitVector())
	// *scattered = Ray{hit.point, scatterdir}
	// *c = dl.c
	return false

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

	cosTheta := r.dir.normalize().scalarMult(-1).dot(hit.normal)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	if eta*sinTheta > 1.0 {
		reflected := reflect(r.dir.normalize(), hit.normal)
		(*scattered) = Ray{hit.point, reflected}
		return true
	}

	reflectProb := schlick(cosTheta, eta)
	if rand.Float64() < reflectProb {
		reflected := reflect(r.dir.normalize(), hit.normal)
		(*scattered) = Ray{hit.point, reflected}
		return true
	}

	refracted := refract(r.dir.normalize(), hit.normal, eta)
	(*scattered) = Ray{hit.point, refracted}
	return true
}

func (d Dielectric) emmited(u, v float64, p Vector) Color {
	return Color{0, 0, 0}
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

func (f FuzzyMirror) emmited(u, v float64, p Vector) Color {
	return Color{0, 0, 0}
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

func (m Matte) emmited(u, v float64, p Vector) Color {
	return Color{0, 0, 0}
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

func (m Mirror) emmited(u, v float64, p Vector) Color {
	return Color{0, 0, 0}
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

func schlick(cosine, ref float64) float64 {
	r0 := math.Pow((1-ref)/(1+ref), 2)
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
