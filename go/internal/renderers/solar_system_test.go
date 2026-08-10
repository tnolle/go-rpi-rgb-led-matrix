package renderers

import (
	"context"
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func TestSolarSystemDrawsAFrame(t *testing.T) {
	matrix := newPrepareMatrix(128, 64)
	for position := range matrix.pixels {
		matrix.pixels[position] = color.RGBA{}
	}
	screen := rgbmatrix.NewScreen(matrix)
	renderer := SolarSystem(screen)
	dc := gg.NewContextForImage(screen.Canvas)

	if err := renderer.drawFrame(context.Background(), dc, 1.25); err != nil {
		t.Fatal(err)
	}

	litPixels := 0
	image := dc.Image()
	for y := 0; y < image.Bounds().Dy(); y++ {
		for x := 0; x < image.Bounds().Dx(); x++ {
			rgba := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
			if int(rgba.R)+int(rgba.G)+int(rgba.B) > 50 {
				litPixels++
			}
		}
	}
	if litPixels < 20 {
		t.Fatalf("solar system frame has too few lit pixels: %d", litPixels)
	}
}
