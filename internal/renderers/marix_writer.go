package renderers

import (
	"image/color"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

type MatrixWriter struct {
	screen *rgbmatrix.Screen
	font   *rgbmatrix.BDFFont
	x, y   int
}

func NewMatrixWriter(screen *rgbmatrix.Screen, font *rgbmatrix.BDFFont) *MatrixWriter {
	return &MatrixWriter{screen: screen, font: font}
}

func (w *MatrixWriter) Write(text string, color color.Color) {
	w.screen.DrawText(w.font, text, w.x, w.y, color)
	w.x += len(text) * w.font.Width()
}

func (w *MatrixWriter) WriteLn(text string, color color.Color) {
	w.Write(text, color)
	w.NewLine()
}

func (w *MatrixWriter) NewLine() {
	w.x = 0
	w.y += w.font.Height()
}

func (w *MatrixWriter) SetPosition(x, y int) {
	w.x = x
	w.y = y
}

func (w *MatrixWriter) Flush() {
	w.x = 0
	w.y = 0
	w.screen.Canvas.Render()
}

func (w *MatrixWriter) Clear() {
	w.x = 0
	w.y = 0
	w.screen.Canvas.Clear()
}
