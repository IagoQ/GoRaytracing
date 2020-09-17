package main

import "math"

type World struct {
	shapes []Shape
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

type Sphere struct {
	pos Vector
	r   float64
	m   Material
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
