package display

// Frame is a complete RGB24 snapshot of an LED matrix. Pixels are stored in
// row-major order, with three bytes (red, green, blue) per pixel.
type Frame struct {
	Width  int
	Height int
	Pixels []byte
}

func (f Frame) clone() Frame {
	pixels := make([]byte, len(f.Pixels))
	copy(pixels, f.Pixels)
	f.Pixels = pixels
	return f
}
