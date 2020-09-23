package main

import (
	"sort"
)

type bvh struct {
	left    Shape
	right   Shape
	box     boundingbox
	isempty bool
}

func CreateBvh(shapelist []Shape, depth int) Shape {

	if len(shapelist) == 0 {
		return bvh{isempty: true}
	}
	if len(shapelist) == 1 {
		left := shapelist[0]
		right := bvh{isempty: true}
		bb := left.bb()
		return bvh{left, right, bb, false}
	}
	if len(shapelist) == 2 {
		left := shapelist[0]
		right := shapelist[1]
		bb := combinebb(left.bb(), right.bb())
		return bvh{left, right, bb, false}
	}

	//sort shapes by alternating axis
	if depth%3 == 0 {
		sort.Sort(ByX(shapelist))
	}
	if depth%3 == 1 {
		sort.Sort(ByY(shapelist))
	}
	if depth%3 == 2 {
		sort.Sort(ByZ(shapelist))
	}

	half := len(shapelist) / 2

	left := CreateBvh(shapelist[0:half], depth+1)
	right := CreateBvh(shapelist[half:], depth+1)
	bb := combinebb(left.bb(), right.bb())

	return bvh{left, right, bb, false}

}
func (t bvh) bb() boundingbox {
	return t.box
}

func (t bvh) hit(r Ray, rec *HitRec, tmin, tmax float64) bool {

	if t.isempty {
		return false
	}

	if !t.box.hit(r, rec, tmin, tmax) {
		return false
	}
	lefthit := t.left.hit(r, rec, tmin, tmax)
	if lefthit {
		return true
	}
	righthit := t.right.hit(r, rec, tmin, tmax)
	return righthit
}

type boundingbox struct {
	min, max Vector
}

func combinebb(b1, b2 boundingbox) boundingbox {
	minx := min(b1.min.x, b2.min.x)
	miny := min(b1.min.y, b2.min.y)
	minz := min(b1.min.z, b2.min.z)

	maxx := max(b1.max.x, b2.max.x)
	maxy := max(b1.max.y, b2.max.y)
	maxz := max(b1.max.z, b2.max.z)
	return boundingbox{Vector{minx, miny, minz}, Vector{maxx, maxy, maxz}}
}

func (bb boundingbox) hit(r Ray, rec *HitRec, tmin, tmax float64) bool {
	// ugly code, refactor wouldnt be too much prettier, open for better ideas
	var t0, t1 float64
	t0 = min(bb.min.x-(r.orig.x)/r.dir.x, bb.max.x-(r.orig.x)/r.dir.x)
	t1 = max(bb.min.x-(r.orig.x)/r.dir.x, bb.max.x-(r.orig.x)/r.dir.x)
	tmin = max(t0, tmin)
	tmax = max(t1, tmax)
	if tmin >= tmax {
		return false
	}
	t0 = min(bb.min.y-(r.orig.y)/r.dir.y, bb.max.y-(r.orig.y)/r.dir.y)
	t1 = max(bb.min.y-(r.orig.y)/r.dir.y, bb.max.y-(r.orig.y)/r.dir.y)
	tmin = max(t0, tmin)
	tmax = max(t1, tmax)
	if tmin >= tmax {
		return false
	}
	t0 = min(bb.min.z-(r.orig.z)/r.dir.z, bb.max.z-(r.orig.z)/r.dir.z)
	t1 = max(bb.min.z-(r.orig.z)/r.dir.z, bb.max.z-(r.orig.z)/r.dir.z)
	tmin = max(t0, tmin)
	tmax = max(t1, tmax)
	if tmin >= tmax {
		return false
	}
	return true

}

// Order shape slices by axis (from bb min vector)

type ByX []Shape

func (a ByX) Len() int      { return len(a) }
func (a ByX) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByX) Less(i, j int) bool {
	if a[i] == nil || a[j] == nil {
		return false
	}
	return a[i].bb().min.x < a[j].bb().min.x
}

type ByY []Shape

func (a ByY) Len() int      { return len(a) }
func (a ByY) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByY) Less(i, j int) bool {
	if a[i] == nil || a[j] == nil {
		return false
	}
	return a[i].bb().min.y < a[j].bb().min.y
}

type ByZ []Shape

func (a ByZ) Len() int      { return len(a) }
func (a ByZ) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByZ) Less(i, j int) bool {
	if a[i] == nil || a[j] == nil {
		return false
	}
	return a[i].bb().min.z < a[j].bb().min.z
}
