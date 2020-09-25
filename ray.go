package main

type Ray struct {
	orig, dir Vector
}

type HitRec struct {
	point, normal Vector
	t, u, v       float64
	front         bool
	material      *Material
}

func (h *HitRec) setFaceNormal(r Ray, n Vector) {
	h.front = r.dir.dot(n) < 0
	if h.front {
		h.normal = n
	} else {
		h.normal = n.scalarMult(-1)
	}
}

func (a Ray) at(t float64) Vector {
	return a.orig.add(a.dir.scalarMult(t))
}
