package renderers

import (
	"context"
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func TestDarts180DrawsAFrame(t *testing.T) {
	matrix := newPrepareMatrix(128, 64)
	for position := range matrix.pixels {
		matrix.pixels[position] = color.RGBA{}
	}
	screen := rgbmatrix.NewScreen(matrix)
	renderer := Darts180(screen)
	dc := gg.NewContextForImage(screen.Canvas)

	if err := renderer.drawFrame(context.Background(), dc, 1.8); err != nil {
		t.Fatal(err)
	}

	brightPixels := 0
	image := dc.Image()
	for y := 0; y < image.Bounds().Dy(); y++ {
		for x := 0; x < image.Bounds().Dx(); x++ {
			rgba := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
			if int(rgba.R)+int(rgba.G)+int(rgba.B) > 180 {
				brightPixels++
			}
		}
	}
	if brightPixels < 100 {
		t.Fatalf("darts 180 frame has too few bright pixels: %d", brightPixels)
	}
}
