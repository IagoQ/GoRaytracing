package main

import (
	"math/rand"
)

type Scene struct {
	width, height, samples, depth int
	world                         World
	camera                        Camera
}

func render(s Scene) [][]Color {
	canvas := CreateCanvas(s.width, s.height)
	var pixelColor Color
	for j := 0; j < s.height; j++ {
		for i := 0; i < s.width; i++ {
			pixelColor = Color{0, 0, 0}
			for k := 0; k < s.samples; k++ {
				u := (float64(i) + rand.Float64()) / float64(s.width-1)
				v := (float64(j) + rand.Float64()) / float64(s.height-1)
				r := s.camera.getRay(u, v)
				pixelColor = pixelColor.add(raycolor(r, s.world, s.depth))
			}
			canvas[i][s.height-j-1] = pixelColor.scalarMult(1.0 / float64(s.samples))
		}
	}
	return canvas
}

func raycolor(r Ray, s Shape, depth int) Color {
	var hit HitRec

	if depth <= 0 {
		return Color{0, 0, 0}
	}

	var scattered Ray
	var attenuation Color
	if s.hit(r, &hit, 0.001, 100) {
		if (*hit.material).scatter(&r, &hit, &attenuation, &scattered) {
			return attenuation.mult(raycolor(scattered, s, depth-1))
		}
		return Color{0, 0, 0}
	}

	unitdir := r.dir.normalize()
	t := 0.5 * (unitdir.y + 1.0)
	return Color{1, 1, 1}.scalarMult(t).add(Color{0.5, 0.7, 1.0}.scalarMult(1 - t))
}
