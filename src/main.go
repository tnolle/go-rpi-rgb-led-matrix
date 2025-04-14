package main

import (
	"image/color"
	"time"

	rgbmatrix "github.com/tnolle/go-rpi-rgb-led-matrix"
)

var P5PanelOptions = &rgbmatrix.MatrixOptions{
	Rows:           32,
	Cols:           64,
	Brightness:     20,
	Multiplexing:   2,
	RowAddressType: 3,

	// Defaults
	PWMBits:           11,
	PWMLSBNanoseconds: 130,
	PWMDitherBits:     0,
	ScanMode:          rgbmatrix.Progressive,
}

var P3PanelOptions = &rgbmatrix.MatrixOptions{
	Rows:       64,
	Cols:       64,
	Brightness: 20,

	// Defaults
	PWMBits:           11,
	PWMLSBNanoseconds: 130,
	PWMDitherBits:     0,
	ScanMode:          rgbmatrix.Progressive,
}

func main() {
	options := *P5PanelOptions
	options.HardwareMapping = "adafruit-hat"
	options.ShowRefreshRate = true
	options.ChainLength = 1
	options.Parallel = 1
	runtimeOptions := &rgbmatrix.RuntimeOptions{
		GPIOSlowdown: -1,
	}

	m, err := rgbmatrix.NewRGBLedMatrix(&options, runtimeOptions)
	if err != nil {
		panic(err)
	}

	c := rgbmatrix.NewCanvas(m)
	defer c.Close() // don't forgot close the Matrix, if not your leds will remain on

	b := 1
	for {
		b = (b + 1) % 2
		time.Sleep(1000 * time.Millisecond)
		//for a := range 64 {
		//	if a%2 == 0 {
		//		DrawColumn(c, a, color.RGBA{255, 0, 0, 0})
		//	}
		//}
		for a := range 32 {
			if a%2 == b {
				DrawRow(c, a, color.RGBA{255, 255, 255, 0})
			}
		}
		c.Render()

		//DrawRow(c, 0, color.RGBA{255, 0, 255, 0})
		//DrawRow(c, 31, color.RGBA{255, 0, 255, 0})
		//DrawColumn(c, 0, color.RGBA{255, 0, 255, 0})
		//DrawColumn(c, 63, color.RGBA{255, 0, 255, 0})
		//
		//for a := range 32 {
		//	c.Set(a, a, color.RGBA{255, 0, 0, 0})
		//	c.Set(32+a, 31-a, color.RGBA{255, 255, 0, 0})
		//}
		//c.Render()

		//for a := range c.Bounds().Dy() {
		//	DrawRow(c, a, color.RGBA{0, 0, 255, 0})
		//	c.Render()
		//	time.Sleep(500 * time.Millisecond)
		//}
		//c.Clear()
		//for a := range c.Bounds().Dx() {
		//	DrawColumn(c, a, color.RGBA{0, 0, 255, 0})
		//	c.Render()
		//	time.Sleep(100 * time.Millisecond)
		//}
		//c.Clear()

		//DrawRow(c, 0, color.RGBA{255, 0, 0, 0})
		//DrawRow(c, 8, color.RGBA{255, 0, 0, 0})
		//DrawRow(c, 16, color.RGBA{255, 0, 0, 0})
		//DrawRow(c, 24, color.RGBA{255, 0, 0, 0})
		//
		//DrawRow(c, 7, color.RGBA{0, 0, 255, 0})
		//DrawRow(c, 15, color.RGBA{0, 0, 255, 0})
		//DrawRow(c, 23, color.RGBA{0, 0, 255, 0})
		//DrawRow(c, 31, color.RGBA{0, 0, 255, 0})

		//DrawRow(c, 25, color.RGBA{0, 0, 255, 0})

		//c.Render()
		//time.Sleep(100 * time.Millisecond)
	}

	// Run forever
	//forever := make(chan struct{})
	//<-forever
}

func DrawRow(c *rgbmatrix.Canvas, row int, color color.Color) {
	bounds := c.Bounds()
	for x := 0; x < bounds.Dx(); x++ {
		c.Set(x, row, color)
	}
}

func DrawColumn(c *rgbmatrix.Canvas, column int, color color.Color) {
	bounds := c.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		c.Set(column, y, color)
	}
}
