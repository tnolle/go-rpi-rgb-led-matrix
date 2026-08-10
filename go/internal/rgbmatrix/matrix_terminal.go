package rgbmatrix

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"os"
	"sync"
)

const (
	terminalEnter = "\x1b[?1049h\x1b[?25l\x1b[?7l\x1b[2J\x1b[H"
	terminalLeave = "\x1b[0m\x1b[?7h\x1b[?25h\x1b[?1049l"
)

// Terminal renders two vertical matrix pixels in each terminal cell. The upper
// pixel becomes the foreground color, the lower pixel the background color,
// and the upper-half block character displays both at once.
type Terminal struct {
	width  int
	height int
	pixels []color.RGBA
	out    io.Writer

	mu      sync.Mutex
	started bool
	closed  bool
}

func NewTerminal(width, height int) *Terminal {
	return NewTerminalWithWriter(width, height, os.Stdout)
}

func NewTerminalWithWriter(width, height int, out io.Writer) *Terminal {
	return &Terminal{
		width:  width,
		height: height,
		pixels: make([]color.RGBA, width*height),
		out:    out,
	}
}

func (t *Terminal) Geometry() (width, height int) {
	return t.width, t.height
}

func (t *Terminal) At(position int) color.Color {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.pixels[position]
}

func (t *Terminal) Set(position int, c color.Color) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.pixels[position] = rgba(c)
}

func (t *Terminal) Apply(pixels []color.Color) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for position := range t.pixels {
		if position < len(pixels) {
			t.pixels[position] = rgba(pixels[position])
		} else {
			t.pixels[position] = color.RGBA{}
		}
	}
	return t.renderLocked()
}

func (t *Terminal) Render() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.renderLocked()
}

func (t *Terminal) renderLocked() error {
	if t.closed {
		return fmt.Errorf("terminal matrix is closed")
	}

	var frame bytes.Buffer
	if !t.started {
		frame.WriteString(terminalEnter)
		t.started = true
	} else {
		frame.WriteString("\x1b[H")
	}

	for y := 0; y < t.height; y += 2 {
		for x := 0; x < t.width; x++ {
			upper := t.pixels[x+y*t.width]
			lower := color.RGBA{}
			if y+1 < t.height {
				lower = t.pixels[x+(y+1)*t.width]
			}

			fmt.Fprintf(
				&frame,
				"\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				upper.R, upper.G, upper.B,
				lower.R, lower.G, lower.B,
			)
		}
		frame.WriteString("\x1b[0m")
		if y+2 < t.height {
			frame.WriteString("\r\n")
		}
	}

	if _, err := t.out.Write(frame.Bytes()); err != nil {
		return fmt.Errorf("render terminal matrix: %w", err)
	}

	clear(t.pixels)
	return nil
}

func (t *Terminal) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}
	t.closed = true
	if !t.started {
		return nil
	}
	if _, err := io.WriteString(t.out, terminalLeave); err != nil {
		return fmt.Errorf("restore terminal: %w", err)
	}
	return nil
}

func rgba(c color.Color) color.RGBA {
	if c == nil {
		return color.RGBA{}
	}
	return color.RGBAModel.Convert(c).(color.RGBA)
}
