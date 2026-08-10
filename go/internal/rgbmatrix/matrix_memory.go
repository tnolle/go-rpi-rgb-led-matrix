package rgbmatrix

import (
	"image/color"
	"sync"
)

// Memory is a headless matrix output. It accepts the same drawing operations
// as a real matrix while leaving presentation to wrappers such as Observable.
type Memory struct {
	width  int
	height int
	pixels []color.RGBA
	mu     sync.Mutex
}

func NewMemory(width, height int) *Memory {
	return &Memory{
		width:  width,
		height: height,
		pixels: make([]color.RGBA, width*height),
	}
}

func (m *Memory) Geometry() (width, height int) {
	return m.width, m.height
}

func (m *Memory) At(position int) color.Color {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pixels[position]
}

func (m *Memory) Set(position int, c color.Color) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pixels[position] = rgba(c)
}

func (m *Memory) Apply(pixels []color.Color) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for position := range m.pixels {
		if position < len(pixels) {
			m.pixels[position] = rgba(pixels[position])
		} else {
			m.pixels[position] = color.RGBA{}
		}
	}
	clear(m.pixels)
	return nil
}

func (m *Memory) Render() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.pixels)
	return nil
}

func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.pixels)
	return nil
}
