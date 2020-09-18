package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println("Starting up...")
	//scene
	height := 350
	width := 400
	aspect := float64(width) / float64(height)
	samples := 500
	maxdepth := 20
	from := Vector{13, 2, 3}
	to := Vector{0, 0, 0}

	//multithreading
	cores := 8
	samplesPerRoutine := samples / cores
	bufferSize := cores

	//camera
	c := CreateCamera(aspect, 60.0, from, to, Vector{0, 1, 0}, 0.1, 10)
	//world
	var w World
	w.add(Sphere{Vector{0, 0, 0}, 2, Matte{0.98, Color{1, 0, 1}}})

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
