package main

import (
	"math/rand"
)

func main() {

	//image
	height := 350
	width := 600
	aspect := float64(width) / float64(height)
	canvas := CreateCanvas(width, height)
	samples := 50
	maxdepth := 20

	//camera
	c := CreateCamera(2.0*aspect, 2.0, 1.0, Vector{0, 0, 0})

	//world
	var w World
	s := Sphere{Vector{0, 100.5, -1}, 100, Matte{1, Color{0.8, 0.8, 0}}}
	s1 := Sphere{Vector{0, 0, -1}, 0.5, Matte{1, Color{0.7, 0.3, 0.3}}}
	s2 := Sphere{Vector{-1.2, 0, -1}, 0.5, Dielectric{1.5, Color{0.8, 0.8, 0.8}}}
	s3 := Sphere{Vector{1.2, 0, -1}, 0.5, FuzzyMirror{0.1, Color{0.8, 0.8, 0.8}}}
	w.add(s)
	w.add(s1)
	w.add(s2)
	w.add(s3)

	//render
	var pixelColor Color
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			pixelColor = Color{0, 0, 0}
			for k := 0; k < samples; k++ {
				u := (float64(i) + rand.Float64()) / float64(width-1)
				v := (float64(j) + rand.Float64()) / float64(height-1)
				r := c.getRay(u, v)
				pixelColor = pixelColor.add(raycolor(r, w, maxdepth))
			}
			canvas[i][j] = pixelColor.scalarMult(1.0 / float64(samples))
		}
	}

	GeneratePng(canvas, "image.png")

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
