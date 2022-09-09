package main

import (
	"math"
)

type IShape interface {
	hit(Ray, *HitRec, float64, float64) bool
	bb() boundingbox
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

	// culling
	// if det < 0.00001 {
	// 	return false
	// }

	// not culling
	if math.Abs(det) < 0.00001 {
		return false
	}

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

func (t Triangle) bb() boundingbox {
	// the wiggle room +- 0.00010 prevents some bugs with triangles that share points
	minx := min(t.p1.x, min(t.p2.x, t.p3.x)) - 0.0001
	miny := min(t.p1.y, min(t.p2.y, t.p3.y)) - 0.0001
	minz := min(t.p1.z, min(t.p2.z, t.p3.z)) - 0.0001
	maxx := max(t.p1.x, max(t.p2.x, t.p3.x)) + 0.0001
	maxy := max(t.p1.y, max(t.p2.y, t.p3.y)) + 0.0001
	maxz := max(t.p1.z, max(t.p2.z, t.p3.z)) + 0.0001
	return boundingbox{Vector{minx, miny, minz}, Vector{maxx, maxy, maxz}}
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

func (s Sphere) bb() boundingbox {
	return boundingbox{s.pos.sub(Vector{s.r, s.r, s.r}), s.pos.add(Vector{s.r, s.r, s.r})}
}
