//go:build with_cgo

package rgbmatrix

/*
#cgo CFLAGS: -std=c99 -I${SRCDIR}/../../3rdparty/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE
#cgo LDFLAGS: -lrgbmatrix -L${SRCDIR}/../../3rdparty/rpi-rgb-led-matrix/lib -lstdc++ -lm
#include <led-matrix-c.h>

void led_matrix_swap(struct RGBLedMatrix *matrix, struct LedCanvas *offscreen_canvas,
                     int width, int height, const uint32_t pixels[]) {
  int i, x, y;
  uint32_t color;
  for (x = 0; x < width; ++x) {
    for (y = 0; y < height; ++y) {
      i = x + (y * width);
      color = pixels[i];

      led_canvas_set_pixel(offscreen_canvas, x, y,
        (color >> 16) & 255, (color >> 8) & 255, color & 255);
    }
  }
  offscreen_canvas = led_matrix_swap_on_vsync(matrix, offscreen_canvas);
}
*/
import "C"
import (
	"fmt"
	"image/color"
	"os"
	"unsafe"
)

// RGBLedMatrix matrix representation for ws281x
type RGBLedMatrix struct {
	Config *MatrixOptions

	height int
	width  int
	matrix *C.struct_RGBLedMatrix
	buffer *C.struct_LedCanvas
	leds   []C.uint32_t
}

const MatrixEmulatorENV = "MATRIX_EMULATOR"

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewRGBLedMatrix(opts *MatrixOptions, rOpts *RuntimeOptions) (c Matrix, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("error creating matrix: %v", r)
			}
		}
	}()

	w, h := opts.geometry()
	m := C.led_matrix_create_from_options_and_rt_options(opts.toC(), rOpts.toC())
	b := C.led_matrix_create_offscreen_canvas(m)
	c = &RGBLedMatrix{
		Config: opts,
		width:  w, height: h,
		matrix: m,
		buffer: b,
		leds:   make([]C.uint32_t, w*h),
	}
	if m == nil {
		return nil, fmt.Errorf("unable to allocate memory")
	}

	return c, nil
}

func isMatrixEmulator() bool {
	if os.Getenv(MatrixEmulatorENV) == "1" {
		return true
	}

	return false
}

// Initialize initialize library, must be called once before other functions are
// called.
func (c *RGBLedMatrix) Initialize() error {
	return nil
}

// Geometry returns the width and the height of the matrix
func (c *RGBLedMatrix) Geometry() (width, height int) {
	return c.width, c.height
}

// Apply set all the pixels to the values contained in leds
func (c *RGBLedMatrix) Apply(leds []color.Color) error {
	for position, l := range leds {
		c.Set(position, l)
	}

	return c.Render()
}

// Render update the display with the data from the LED buffer
func (c *RGBLedMatrix) Render() error {
	w, h := c.Config.geometry()

	C.led_matrix_swap(
		c.matrix,
		c.buffer,
		C.int(w), C.int(h),
		(*C.uint32_t)(unsafe.Pointer(&c.leds[0])),
	)

	c.leds = make([]C.uint32_t, w*h)
	return nil
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *RGBLedMatrix) At(position int) color.Color {
	return uint32ToColor(c.leds[position])
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *RGBLedMatrix) Set(position int, color color.Color) {
	c.leds[position] = C.uint32_t(colorToUint32(color))
}

// Close finalizes the ws281x interface
func (c *RGBLedMatrix) Close() error {
	C.led_matrix_delete(c.matrix)
	return nil
}

func colorToUint32(c color.Color) uint32 {
	if c == nil {
		return 0
	}

	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, _ := c.RGBA()
	return (red>>8)<<16 | (green>>8)<<8 | blue>>8
}

func uint32ToColor(u C.uint32_t) color.Color {
	return color.RGBA{
		uint8(u>>16) & 255,
		uint8(u>>8) & 255,
		uint8(u>>0) & 255,
		0,
	}
}
