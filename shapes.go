package main

import (
	"math"
)

type World struct {
	shapes []Shape
}

func CreateWorld() World {
	return World{}
}

func (w World) hit(r Ray, rec *HitRec, tmin, tmax float64) bool {
	var temprec HitRec
	hit := false
	closest := tmax
	for _, s := range w.shapes {
		if s.hit(r, &temprec, tmin, closest) {
			hit = true
			closest = temprec.t
			(*rec) = temprec
		}
	}
	return hit
}

func (w *World) add(s Shape) {
	w.shapes = append(w.shapes, s)
}

type Shape interface {
	hit(Ray, *HitRec, float64, float64) bool
}

type Triangle struct {
	p1, p2, p3, normal Vector
	m                  Material
}

func CreateTriangle(p1, p2, p3 Vector, m Material) Triangle {
	var t Triangle
	t.p1 = p1
	t.p2 = p2
	t.p3 = p3
	t.normal = (p2.sub(p1)).cross(p3.sub(p1)).normalize()
	t.m = m
	return t
}

func (tri Triangle) hit(r Ray, rec *HitRec, tmin, tmax float64) bool {
	edge1 := tri.p2.sub(tri.p1)
	edge2 := tri.p3.sub(tri.p1)

	pvec := r.dir.cross(edge2)
	det := edge1.dot(pvec)

	//culling
	if det < 0.00001 {
		return false
	}

	// not culling
	// if math.Abs(det) < 0.00001 {
	// 	return false
	// }

	invDet := 1.0 / det

	tvec := r.orig.sub(tri.p1)
	u := tvec.dot(pvec) * invDet
	if u < 0 || u > 1 {
		return false
	}

	qvec := tvec.cross(edge1)
	v := r.dir.dot(qvec) * invDet
	if v < 0 || v+u > 1 {
		return false
	}

	t := edge2.dot(qvec) * invDet

	if t < tmin || t > tmax {
		return false
	}

	rec.point = r.at(t)
	rec.material = &tri.m
	rec.setFaceNormal(r, tri.normal)
	rec.t = t

	return true
}

type Sphere struct {
	pos Vector
	r   float64
	m   Material
}

func CreateSphere(p Vector, r float64, m Material) Sphere {
	return Sphere{p, r, m}
}
func (s Sphere) hit(r Ray, rec *HitRec, tmin, tmax float64) bool {
	oc := r.orig.sub(s.pos)
	a := r.dir.dot(r.dir)
	b := oc.dot(r.dir)
	c := oc.dot(oc) - s.r*s.r
	dis := b*b - a*c

	if dis > 0 {
		root := math.Sqrt(dis)

		temp := (-b - root) / a
		if temp > tmin && temp < tmax {
			rec.t = temp
			rec.point = r.at(temp)
			outnormal := (rec.point.sub(s.pos).scalarMult(1 / s.r))
			rec.setFaceNormal(r, outnormal)
			rec.material = &s.m
			return true
		}

		temp = (-b + root) / a
		if temp > tmin && temp < tmax {
			rec.t = temp
			rec.point = r.at(temp)
			outnormal := (rec.point.sub(s.pos).scalarMult(1 / s.r))
			rec.setFaceNormal(r, outnormal)
			rec.material = &s.m
			return true
		}
	}
	return false
}
