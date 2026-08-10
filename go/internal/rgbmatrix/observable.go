package rgbmatrix

import (
	"image/color"
	"sync"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
)

// Observable mirrors matrix writes so that each successful render can be
// published as a complete frame before the underlying matrix clears its
// drawing buffer.
type Observable struct {
	matrix  Matrix
	publish func(display.Frame)
	width   int
	height  int
	pixels  []color.RGBA
	mu      sync.Mutex
}

func NewObservable(matrix Matrix, publish func(display.Frame)) *Observable {
	width, height := matrix.Geometry()
	return &Observable{
		matrix:  matrix,
		publish: publish,
		width:   width,
		height:  height,
		pixels:  make([]color.RGBA, width*height),
	}
}

func (o *Observable) Geometry() (width, height int) {
	return o.width, o.height
}

func (o *Observable) At(position int) color.Color {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.pixels[position]
}

func (o *Observable) Set(position int, c color.Color) {
	o.mu.Lock()
	defer o.mu.Unlock()
	pixel := rgba(c)
	o.pixels[position] = pixel
	o.matrix.Set(position, pixel)
}

func (o *Observable) Apply(pixels []color.Color) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	for position := range o.pixels {
		pixel := color.RGBA{}
		if position < len(pixels) {
			pixel = rgba(pixels[position])
		}
		o.pixels[position] = pixel
		o.matrix.Set(position, pixel)
	}
	return o.renderLocked()
}

func (o *Observable) Render() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.renderLocked()
}

func (o *Observable) renderLocked() error {
	pixels := make([]byte, len(o.pixels)*3)
	for position, pixel := range o.pixels {
		offset := position * 3
		pixels[offset] = pixel.R
		pixels[offset+1] = pixel.G
		pixels[offset+2] = pixel.B
	}

	if err := o.matrix.Render(); err != nil {
		return err
	}
	o.publish(display.Frame{Width: o.width, Height: o.height, Pixels: pixels})
	clear(o.pixels)
	return nil
}

func (o *Observable) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.matrix.Close()
}
