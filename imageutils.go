package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type Color struct {
	r float64
	g float64
	b float64
}

func CreateCanvas(width int, height int) [][]Color {
	image := make([][]Color, width)
	for i := range image {
		image[i] = make([]Color, height)
	}
	return image
}

func GeneratePng(canvas [][]Color, filename string) {
	height := len(canvas[0])
	width := len(canvas)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := uint8(canvas[x][y].r * 255.999)
			g := uint8(canvas[x][y].g * 255.999)
			b := uint8(canvas[x][y].b * 255.999)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	f, _ := os.Create("image.png")
	png.Encode(f, img)

}
