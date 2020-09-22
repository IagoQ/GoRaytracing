package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting up...")
	//scene
	height := 500
	width := 500
	aspect := float64(width) / float64(height)
	samples := 10
	maxdepth := 5
	from := Vector{1, 13, 0}
	to := Vector{2, 2, 2}

	//multithreading
	cores := samples
	samplesPerRoutine := samples / cores
	bufferSize := cores

	//camera
	c := CreateCamera(aspect, 50.0, from, to, Vector{0, 1, 0}, 0.1, 10)
	//world

	p1 := Vector{1, 1, 1}
	p2 := Vector{4, 4, 4}
	p3 := Vector{1, 4, 4}
	mat := Matte{1, Color{0.2, 0.8, 0.5}}

	w := CreateWorld()
	w.add(CreateTriangle(p3, p2, p1, mat))
	// w.add(Sphere{p1, 0.5, mat})
	// w.add(Sphere{p2, 0.5, mat})
	// w.add(Sphere{p3, 0.5, mat})
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
