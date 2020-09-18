package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	fmt.Println("Starting up...")
	//scene
	height := 700
	width := 700
	aspect := float64(width) / float64(height)
	samples := 100
	maxdepth := 20
	from := Vector{13, 2, 3}
	to := Vector{0, 0, 0}

	//multithreading
	cores := samples
	samplesPerRoutine := samples / cores
	bufferSize := cores

	//camera
	c := CreateCamera(aspect, 50.0, from, to, Vector{0, 1, 0}, 0.1, 10)
	//world
	w := randomScene()
	//render

	inputChannel := make(chan Scene, bufferSize)
	outputChannel := make(chan [][]Color, bufferSize)

	for i := 0; i < bufferSize; i++ {
		inputChannel <- Scene{width, height, samplesPerRoutine, maxdepth, w, c}
	}
	close(inputChannel)

	fmt.Println("Starting render")
	starttime := time.Now()
	for i := 0; i < cores; i++ {
		go func() {
			outputChannel <- render(<-inputChannel)
		}()
	}

	final := CreateCanvas(width, height)
	for i := 0; i < bufferSize; i++ {
		final = addCanvas(final, <-outputChannel)
	}
	final = multCanvas(final, 1.0/float64(bufferSize))
	close(outputChannel)
	endtime := time.Now()
	time := endtime.Sub(starttime)
	fmt.Println("ended in: " + time.String())

	GeneratePng(final, "image.png")

}

func randomScene() World {
	var w World
	groundMaterial := Matte{0.98, Color{0.5, 0.5, 0.5}}
	w.add(Sphere{Vector{0, 1000, 0}, 1000, groundMaterial})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			matc := rand.Float64()
			radius := rand.Float64() / 2.0
			center := Vector{float64(a) + 0.9*rand.Float64(), radius, float64(b) + 0.9*rand.Float64()}

			if center.sub(Vector{4, 0.2, 0}).length() > 0.9 {
				if matc < 0.4 {
					smat := Matte{0.98, randomColor()}
					w.add(Sphere{center, radius, smat})
				} else if matc < 0.8 {
					mcolor := Color{1, 1, 1}.sub(randomColor().scalarMult(0.3))
					mmat := FuzzyMirror{rand.Float64() / 8, mcolor}
					w.add(Sphere{center, radius, mmat})
				} else {
					w.add(Sphere{center, radius, Dielectric{1.5, Color{1, 1, 1}}})
				}
			}
		}
	}
	w.add(Sphere{Vector{0, 1, 0.2}, 1, Dielectric{1.5, Color{1, 1, 1}}})
	w.add(Sphere{Vector{-4, 1, 0}, 1, Matte{0.99, Color{0.4, 0.2, 0.1}}})
	w.add(Sphere{Vector{4, 1, -0.2}, 1, Mirror{Color{0.9, 0.9, 0.9}}})
	return w
}
