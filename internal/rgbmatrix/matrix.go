package rgbmatrix

import "image/color"

// Matrix is an interface that represent any RGB matrix, very useful for testing
type Matrix interface {
	Geometry() (width, height int)
	At(position int) color.Color
	Set(position int, c color.Color)
	Apply([]color.Color) error
	Render() error
	Close() error
}
