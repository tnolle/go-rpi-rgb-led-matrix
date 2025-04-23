package renderers

import (
	"context"
	"fmt"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/fs"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

type ImageRenderer struct {
	screen *rgbmatrix.Screen
	path   string
}

func Image(screen *rgbmatrix.Screen, path string) *ImageRenderer {
	return &ImageRenderer{screen: screen, path: path}
}

func (r *ImageRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	img, err := fs.LoadPNG(fmt.Sprintf("images/pngs/%s.png", r.path))
	if err != nil {
		return err
	}
	return r.screen.ShowImage(ctx, img)
}

type GIFOnceRenderer struct {
	screen *rgbmatrix.Screen
	path   string
}

func GIFOnce(screen *rgbmatrix.Screen, path string) *GIFOnceRenderer {
	return &GIFOnceRenderer{screen: screen, path: path}
}

func (r *GIFOnceRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	img, err := fs.LoadGIF(fmt.Sprintf("images/gifs/%s.gif", r.path))
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-r.screen.PlayGIF(ctx, img):
			if len(cb) == 1 {
				cb[0]()
			}
		}
	}()
	return nil
}

type GIFLoopRenderer struct {
	screen *rgbmatrix.Screen
	path   string
}

func GIFLoop(screen *rgbmatrix.Screen, path string) *GIFLoopRenderer {
	return &GIFLoopRenderer{screen: screen, path: path}
}

func (r *GIFLoopRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	img, err := fs.LoadGIF(fmt.Sprintf("images/gifs/%s.gif", r.path))
	if err != nil {
		return err
	}
	r.screen.LoopGIF(ctx, img)
	return nil
}
