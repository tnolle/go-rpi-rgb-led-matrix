package rgbmatrix

import (
	"errors"
	"image/color"
	"testing"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
)

func TestObservablePublishesSuccessfulRenderAsRGB24(t *testing.T) {
	matrix := &observableMatrix{width: 2, height: 1, pixels: make([]color.Color, 2)}
	var published display.Frame
	observable := NewObservable(matrix, func(frame display.Frame) { published = frame })
	observable.Set(0, color.RGBA{R: 1, G: 2, B: 3, A: 255})
	observable.Set(1, color.RGBA{R: 4, G: 5, B: 6, A: 255})

	if err := observable.Render(); err != nil {
		t.Fatal(err)
	}
	if published.Width != 2 || published.Height != 1 {
		t.Fatalf("unexpected dimensions: %dx%d", published.Width, published.Height)
	}
	want := []byte{1, 2, 3, 4, 5, 6}
	if string(published.Pixels) != string(want) {
		t.Fatalf("unexpected pixels: got %v, want %v", published.Pixels, want)
	}
}

func TestObservableDoesNotPublishFailedRender(t *testing.T) {
	matrix := &observableMatrix{width: 1, height: 1, pixels: make([]color.Color, 1), renderErr: errors.New("failed")}
	published := false
	observable := NewObservable(matrix, func(display.Frame) { published = true })

	if err := observable.Render(); err == nil {
		t.Fatal("expected render error")
	}
	if published {
		t.Fatal("failed render was published")
	}
}

type observableMatrix struct {
	width, height int
	pixels        []color.Color
	renderErr     error
}

func (m *observableMatrix) Geometry() (int, int)            { return m.width, m.height }
func (m *observableMatrix) At(position int) color.Color     { return m.pixels[position] }
func (m *observableMatrix) Set(position int, c color.Color) { m.pixels[position] = c }
func (m *observableMatrix) Apply(pixels []color.Color) error {
	copy(m.pixels, pixels)
	return m.Render()
}
func (m *observableMatrix) Render() error { return m.renderErr }
func (m *observableMatrix) Close() error  { return nil }
