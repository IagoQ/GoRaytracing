package main

import (
	"math/rand"
)

type Scene struct {
	width, height, samples, depth int
	world                         IShape
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

func raycolor(r Ray, s IShape, depth int) Color {
	var hit HitRec
	min, max := 0.001, 100.0

	// bounce limit reached
	if depth <= 0 {
		return Color{0, 0, 0}
	}

	// didnt hit anything
	if !s.hit(r, &hit, min, max) {
		return Color{1, 1, 1}
	}

	var scattered Ray
	var attenuation Color
	emmited := (*hit.material).emmited(hit.u, hit.v, hit.point)

	if !(*hit.material).scatter(&r, &hit, &attenuation, &scattered) {
		return emmited
	}
	return emmited.add(attenuation.mult(raycolor(scattered, s, depth-1)))

}
