package rgbmatrix

import (
	"image/color"
	"testing"
)

func TestMemoryMatrix(t *testing.T) {
	matrix := NewMemory(2, 1)
	width, height := matrix.Geometry()
	if width != 2 || height != 1 {
		t.Fatalf("unexpected geometry: %dx%d", width, height)
	}

	matrix.Set(0, color.RGBA{R: 1, G: 2, B: 3, A: 255})
	if got := matrix.At(0); got != (color.RGBA{R: 1, G: 2, B: 3, A: 255}) {
		t.Fatalf("unexpected pixel: %v", got)
	}
	if err := matrix.Render(); err != nil {
		t.Fatal(err)
	}
	if got := matrix.At(0); got != (color.RGBA{}) {
		t.Fatalf("render did not clear drawing buffer: %v", got)
	}
}
