package main

import (
	"fmt"
	"time"
)

//TODO clear up main

func main() {
	fmt.Println("Starting up...")

	//profile
	// startprofile("cpu.prof")
	// defer stopprofile()

	//scene
	height := 200
	width := 200
	aspect := float64(width) / float64(height)
	samples := 10000
	maxdepth := 10
	from := Vector{2, 4, 9}
	to := Vector{2, 2, 2}

	//multithreading
	cores := samples
	samplesPerRoutine := samples / cores
	bufferSize := cores

	//camera
	c := CreateCamera(aspect, 50.0, from, to, Vector{0, 1, 0}, 0.1, 10)

	//bvh
	w := CreateBvh(sl, 0, len(sl))

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
