package main

import (
	"fmt"
	"math/rand"
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
	samples := 50
	maxdepth := 5
	from := Vector{1, 13, 0}
	to := Vector{1, 1, 1}

	//multithreading
	cores := samples
	samplesPerRoutine := samples / cores
	bufferSize := cores

	//camera
	c := CreateCamera(aspect, 50.0, from, to, Vector{0, 1, 0}, 0.1, 10)
	//bvh
	sl := make([]Shape, 0)
	// populate bvh
	for i := 0; i < 1000; i++ {
		pos := Vector{rand.Float64()*6 - 3, rand.Float64()*6 - 3, rand.Float64()*6 - 3}
		p1 := Vector{rand.Float64() * 2, rand.Float64() * 2, rand.Float64() * 2}
		p2 := Vector{rand.Float64() * 2, rand.Float64() * 2, rand.Float64() * 2}
		p3 := Vector{rand.Float64() * 2, rand.Float64() * 2, rand.Float64() * 2}

		mat := Matte{1, Color{rand.Float64(), rand.Float64(), rand.Float64()}}
		t := CreateTriangle(p3.add(pos), p2.add(pos), p1.add(pos), mat)
		sl = append(sl, t)
	}
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
