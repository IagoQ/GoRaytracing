package main

import (
	"fmt"
	"math/rand"
)

func main() {

	//image
	height := 350
	width := 400
	aspect := float64(width) / float64(height)
	canvas := CreateCanvas(width, height)
	samples := 150
	maxdepth := 20

	from := Vector{13, 2, 3}
	to := Vector{0, 0, 0}

	//camera
	c := CreateCamera(aspect, 20.0, from, to, Vector{0, 1, 0}, 0.1, 10)

	//world
	w := randomScene()

	//render
	var pixelColor Color
	for j := 0; j < height; j++ {
		fmt.Println(float64(j) / float64(height))
		for i := 0; i < width; i++ {
			pixelColor = Color{0, 0, 0}
			for k := 0; k < samples; k++ {
				u := (float64(i) + rand.Float64()) / float64(width-1)
				v := (float64(j) + rand.Float64()) / float64(height-1)
				r := c.getRay(u, v)
				pixelColor = pixelColor.add(raycolor(r, w, maxdepth))
			}
			canvas[i][height-j-1] = pixelColor.scalarMult(1.0 / float64(samples))
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

func randomScene() World {
	var w World
	groundMaterial := Matte{0.98, Color{0.5, 0.5, 0.5}}
	w.add(Sphere{Vector{0, 1000, 0}, 1000, groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			matc := rand.Float64()
			center := Vector{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.sub(Vector{4, 0.2, 0}).length() > 0.9 {
				if matc < 0.6 {
					smat := Matte{0.98, randomColor()}
					w.add(Sphere{center, rand.Float64() / 4.0, smat})
				} else if matc < 0.8 {
					mcolor := Color{1, 1, 1}.sub(randomColor().scalarMult(0.3))
					mmat := FuzzyMirror{rand.Float64() / 8, mcolor}
					w.add(Sphere{center, rand.Float64() / 2.0, mmat})
				} else {
					w.add(Sphere{center, rand.Float64() / 2, Dielectric{1.5, Color{1, 1, 1}}})
				}
			}
		}
	}
	w.add(Sphere{Vector{0, 1, 0}, 1, Dielectric{1.5, Color{1, 1, 1}}})
	w.add(Sphere{Vector{-4, 1, 0}, 1, Matte{0.99, Color{0.4, 0.2, 0.1}}})
	w.add(Sphere{Vector{4, 1, 0}, 1, Mirror{Color{0.9, 0.9, 0.9}}})
	return w
}
