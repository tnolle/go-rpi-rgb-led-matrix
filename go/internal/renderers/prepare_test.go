package renderers

import (
	"context"
	"errors"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/animation"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func TestEveryAnimationCanBePrepared(t *testing.T) {
	screen := rgbmatrix.NewScreen(newPrepareMatrix(16, 16))
	for _, name := range animation.AnimationStrings() {
		t.Run(name, func(t *testing.T) {
			prepared, err := prepare(context.Background(), Command{Type: TypeAnimation, Name: name}, screen)
			if err != nil {
				t.Fatal(err)
			}
			if prepared.renderer == nil || !prepared.async {
				t.Fatalf("animation was not prepared correctly: %#v", prepared)
			}
		})
	}
}

func TestPrepareRejectsCorruptImage(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "images", "pngs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "broken.png"), []byte("not a png"), 0o600); err != nil {
		t.Fatal(err)
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	screen := rgbmatrix.NewScreen(newPrepareMatrix(16, 16))
	_, err = prepare(context.Background(), Command{Type: TypeImage, Name: "broken"}, screen)
	if !errors.Is(err, ErrInvalidAsset) {
		t.Fatalf("unexpected error: %v", err)
	}
}

type prepareMatrix struct {
	width  int
	height int
	pixels []color.Color
}

func newPrepareMatrix(width, height int) *prepareMatrix {
	return &prepareMatrix{width: width, height: height, pixels: make([]color.Color, width*height)}
}

func (m *prepareMatrix) Geometry() (int, int)        { return m.width, m.height }
func (m *prepareMatrix) At(position int) color.Color { return m.pixels[position] }
func (m *prepareMatrix) Set(position int, c color.Color) {
	m.pixels[position] = c
}
func (m *prepareMatrix) Apply(pixels []color.Color) error {
	copy(m.pixels, pixels)
	return nil
}
func (m *prepareMatrix) Render() error { return nil }
func (m *prepareMatrix) Close() error  { return nil }
