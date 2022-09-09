package main

import (
	"fmt"
	"math/rand"
	"time"
)

//TODO clear up main

func main() {
	fmt.Println("Starting up...")

	rand.Seed(time.Now().Unix())

	//profile
	// startprofile("cpu.prof")
	// defer stopprofile()

	//scene
	height := 1080
	width := 1920
	aspect := float64(width) / float64(height)
	samples := 100
	maxdepth := 2
	from := Vector{0, 10, 50}
	to := Vector{0, 0, 0}

	//multithreading
	cores := 11
	samplesPerRoutine := samples / cores

	//camera
	c := CreateCamera(aspect, 20.0, from, to, Vector{0, 1, 0}, 0.1, 50)

	w, err := LoadShape("model.obj", Matte{
		reflectance: 1,
		c: Color{
			r: 0,
			g: 0.5,
			b: 1,
		},
	})
	if err != nil {
		panic(err)
	}

	inputChannel := make(chan Scene, cores)
	outputChannel := make(chan [][]Color, cores)

	for i := 0; i < cores; i++ {
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
	for i := 0; i < cores; i++ {
		final = addCanvas(final, <-outputChannel)
	}
	final = multCanvas(final, 1.0/float64(cores))
	close(outputChannel)
	endtime := time.Now()
	time := endtime.Sub(starttime)
	fmt.Println("ended in: " + time.String())

	GeneratePng(final, "image.png")
}

func createshapes() []IShape {
	sl := make([]IShape, 0)
	// populate bvh
	for i := 0; i < 10; i++ {
		p1 := Vector{rand.Float64()*6 - 3, rand.Float64()*6 - 3, rand.Float64()*6 - 3}
		p2 := Vector{rand.Float64()*6 - 3, rand.Float64()*6 - 3, rand.Float64()*6 - 3}
		p3 := Vector{rand.Float64()*6 - 3, rand.Float64()*6 - 3, rand.Float64()*6 - 3}

		var mat Material
		mat = Matte{1, Color{rand.Float64(), rand.Float64(), rand.Float64()}}
		switch rand.Int() % 3 {
		case 0:
			mat = Mirror{Color{rand.Float64(), rand.Float64(), rand.Float64()}}
		case 1:
			mat = Matte{1, Color{rand.Float64(), rand.Float64(), rand.Float64()}}
		case 2:
			mat = FuzzyMirror{0.5, Color{rand.Float64(), rand.Float64(), rand.Float64()}}
		}
		t := CreateTriangle(p1, p2, p3, mat)
		sl = append(sl, t)
	}

	return sl
}
