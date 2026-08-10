package renderers

import (
	"context"
	"fmt"
	"image"
	"image/gif"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/fs"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

type ImageRenderer struct {
	screen *rgbmatrix.Screen
	image  image.Image
}

func Image(screen *rgbmatrix.Screen, path string) (*ImageRenderer, error) {
	img, err := fs.LoadPNG(fmt.Sprintf("images/pngs/%s.png", path))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidAsset, err)
	}
	return &ImageRenderer{screen: screen, image: img}, nil
}

func (r *ImageRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	return r.screen.ShowImage(ctx, r.image)
}

type GIFOnceRenderer struct {
	screen *rgbmatrix.Screen
	gif    *gif.GIF
}

func GIFOnce(screen *rgbmatrix.Screen, path string) (*GIFOnceRenderer, error) {
	img, err := fs.LoadGIF(fmt.Sprintf("images/gifs/%s.gif", path))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidAsset, err)
	}
	return &GIFOnceRenderer{screen: screen, gif: img}, nil
}

func (r *GIFOnceRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-r.screen.PlayGIF(ctx, r.gif):
			if len(cb) == 1 {
				cb[0]()
			}
		}
	}()
	return nil
}

type GIFLoopRenderer struct {
	screen *rgbmatrix.Screen
	gif    *gif.GIF
}

func GIFLoop(screen *rgbmatrix.Screen, path string) (*GIFLoopRenderer, error) {
	img, err := fs.LoadGIF(fmt.Sprintf("images/gifs/%s.gif", path))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidAsset, err)
	}
	return &GIFLoopRenderer{screen: screen, gif: img}, nil
}

func (r *GIFLoopRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	r.screen.LoopGIF(ctx, r.gif)
	return nil
}
