package main

import (
	"fmt"
)

func main() {

	height := 1000
	width := 1000

	image := CreateCanvas(width, height)
	fmt.Println(len(image))
	fmt.Println(len(image[0]))

	for j := height - 1; j >= 0; j -= 1 {
		for i := 0; i < width; i++ {
			image[i][j] = Color{float64(i) / float64(width), float64(j) / float64(height), 0}
		}

	}
	GeneratePng(image, "image.ppm")

}
