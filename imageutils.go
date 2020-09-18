package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

type Color struct {
	r float64
	g float64
	b float64
}

func (a Color) add(b Color) Color {
	a.r += b.r
	a.g += b.g
	a.b += b.b

	return a
}

func (a Color) sub(b Color) Color {
	a.r -= b.r
	a.g -= b.g
	a.b -= b.b

	return a
}

func (a Color) scalarMult(b float64) Color {
	a.r *= b
	a.g *= b
	a.b *= b

	return a
}

func (a Color) mult(b Color) Color {
	a.r *= b.r
	a.g *= b.g
	a.b *= b.b
	return a
}

func (a Color) sqrt() Color {
	a.r = math.Sqrt(a.r)
	a.g = math.Sqrt(a.g)
	a.b = math.Sqrt(a.b)

	return a
}

func randomColor() Color {
	return Color{rand.Float64(), rand.Float64(), rand.Float64()}
}

//TODO make canvas struct?
func CreateCanvas(width int, height int) [][]Color {
	image := make([][]Color, width)
	for i := range image {
		image[i] = make([]Color, height)
	}
	return image
}

func addCanvas(a, b [][]Color) [][]Color {
	final := CreateCanvas(len(a), len(a[0]))
	//TODO: check sizes first
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			final[i][j] = a[i][j].add(b[i][j])
		}
	}
	return final
}

func multCanvas(a [][]Color, n float64) [][]Color {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			a[i][j] = a[i][j].scalarMult(n)
		}
	}
	return a
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

	f, _ := os.Create(filename)
	png.Encode(f, img)

}
