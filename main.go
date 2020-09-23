package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

//TODO clear up main
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	fmt.Println("Starting up...")

	//profile
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		fmt.Println("Profiling")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//scene
	height := 300
	width := 300
	aspect := float64(width) / float64(height)
	samples := 60
	maxdepth := 5
	from := Vector{1, 13, 0}
	to := Vector{1, 1, 1}

	//multithreading
	cores := 1
	samplesPerRoutine := samples / cores
	bufferSize := cores

	//camera
	c := CreateCamera(aspect, 50.0, from, to, Vector{0, 1, 0}, 0.1, 10)
	//world

	w := CreateWorld()
	sl := make([]Shape, 0)
	for i := 0; i < 20; i++ {
		pos := Vector{rand.Float64()*6 - 3, rand.Float64()*6 - 3, rand.Float64()*6 - 3}
		p1 := Vector{rand.Float64() * 2, rand.Float64() * 2, rand.Float64() * 2}
		p2 := Vector{rand.Float64() * 2, rand.Float64() * 2, rand.Float64() * 2}
		p3 := Vector{rand.Float64() * 2, rand.Float64() * 2, rand.Float64() * 2}
		mat := Matte{1, Color{rand.Float64(), rand.Float64(), rand.Float64()}}
		t := CreateTriangle(p3.add(pos), p2.add(pos), p1.add(pos), mat)
		w.add(t)
		sl = append(sl, t)
	}
	w2 := CreateBvh(sl, 0)
	//render1

	inputChannel := make(chan Scene, bufferSize)
	outputChannel := make(chan [][]Color, bufferSize)

	for i := 0; i < bufferSize; i++ {
		inputChannel <- Scene{width, height, samplesPerRoutine, maxdepth, w2, c}
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

	f, err := os.Create("heap.prof")
	fmt.Println("Profiling")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	GeneratePng(final, "image.png")

}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
