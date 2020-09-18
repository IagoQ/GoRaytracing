package main

import (
	"math"
	"math/rand"
)

type Vector struct {
	x, y, z float64
}

func (a Vector) add(b Vector) Vector {
	a.x += b.x
	a.y += b.y
	a.z += b.z

	return a
}

func (a Vector) sub(b Vector) Vector {
	a.x -= b.x
	a.y -= b.y
	a.z -= b.z

	return a
}

func (a Vector) scalarMult(b float64) Vector {
	a.x *= b
	a.y *= b
	a.z *= b

	return a
}

func (a Vector) dot(b Vector) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func (a Vector) cross(b Vector) Vector {
	return Vector{
		x: a.y*b.z - a.z*b.y,
		y: a.z*b.x - a.x*b.z,
		z: a.x*b.y - a.y*b.x,
	}
}

func (a Vector) length() float64 {
	return math.Sqrt(a.dot(a))
}
func (a Vector) lengthSquared() float64 {
	return a.dot(a)
}

func (a Vector) normalize() Vector {
	return a.scalarMult(1. / a.length())
}

func randomUnitVector() Vector {
	a := rand.Float64() * math.Pi * 2
	z := (rand.Float64() * 2) - 1
	r := math.Sqrt(1 - z*z)
	return Vector{r * math.Cos(a), r * math.Sin(a), z}
}

func randomDiskVector() Vector {
	for {
		p := Vector{rand.Float64()*2 - 1, rand.Float64()*2 - 1, 0}
		if p.lengthSquared() < 1 {
			return p
		}
	}
}
