package main

import (
	"math/rand"
	"sort"
	"sync"
)

// bvh definition

type bvh struct {
	left  *IShape
	right *IShape
	box   boundingbox
}

func CreateBvh(shapelist []IShape, start, end int) IShape {
	b := bvh{}

	span := end - start

	switch span {
	case 1:
		b.left = &shapelist[start]
		b.right = &shapelist[start]
	case 2:
		b.left = &shapelist[start]
		b.right = &shapelist[start+1]
	default:

		// sort shapes by random axis
		axisr := rand.Float64()
		if axisr < 0.33 {
			sort.SliceStable(shapelist[start:end], func(i, j int) bool { return shapelist[i].bb().min.x < shapelist[j].bb().min.x })
		} else if axisr < 0.66 {
			sort.SliceStable(shapelist[start:end], func(i, j int) bool { return shapelist[i].bb().min.y < shapelist[j].bb().min.y })
		} else {
			sort.SliceStable(shapelist[start:end], func(i, j int) bool { return shapelist[i].bb().min.z < shapelist[j].bb().min.z })
		}

		mid := start + span/2
		wg := sync.WaitGroup{}
		wg.Add(2)
		var l, r IShape
		go func() {
			l = CreateBvh(shapelist, start, mid)
			wg.Done()
		}()
		go func() {
			r = CreateBvh(shapelist, mid, end)
			wg.Done()
		}()
		wg.Wait()
		b.left = &l
		b.right = &r
	}

	b.box = combinebb((*b.left).bb(), (*b.right).bb())
	return b

}

func (b bvh) hit(r Ray, rec *HitRec, tmin, tmax float64) bool {
	if !b.box.hit(r, rec, tmin, tmax) {
		return false
	}
	lefthit := (*b.left).hit(r, rec, tmin, tmax)
	righthit := (*b.right).hit(r, rec, tmin, tmax)

	return lefthit || righthit
}
func (b bvh) bb() boundingbox {
	return b.box
}

// bounding box defintion

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

	// x axis
	invDet := 1.0 / r.dir.x
	t0 := (bb.min.x - r.orig.x) * invDet
	t1 := (bb.max.x - r.orig.x) * invDet
	if invDet < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}

	//y axis
	invDet = 1.0 / r.dir.y
	t0 = (bb.min.y - r.orig.y) * invDet
	t1 = (bb.max.y - r.orig.y) * invDet
	if invDet < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}
	//z axis
	invDet = 1.0 / r.dir.z
	t0 = (bb.min.z - r.orig.z) * invDet
	t1 = (bb.max.z - r.orig.z) * invDet
	if invDet < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}
	return true

}
