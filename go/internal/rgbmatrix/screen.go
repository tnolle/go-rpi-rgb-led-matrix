package rgbmatrix

import (
	"context"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"time"

	"github.com/disintegration/imaging"
)

type TransformFunc func(img image.Image) image.Image

type Screen struct {
	Canvas    *Canvas
	Transform TransformFunc
}

func NewScreen(m Matrix) *Screen {
	w, h := m.Geometry()
	return &Screen{
		Canvas: NewCanvas(m),
		Transform: func(img image.Image) image.Image {
			if img.Bounds().Dx() == w && img.Bounds().Dy() == h {
				return img
			}
			return fitCenter(img, w, h, imaging.Lanczos)
		},
	}
}

func (s *Screen) WithTransform(f TransformFunc) *Screen {
	s.Transform = f
	return s
}

func (s *Screen) Fill(color color.Color) {
	draw.Draw(s.Canvas, s.Canvas.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
}

func (s *Screen) ShowImage(ctx context.Context, i image.Image) error {
	if s.Transform != nil {
		i = s.Transform(i)
	}

	draw.Draw(s.Canvas, s.Canvas.Bounds(), i, image.Point{}, draw.Over)

	return s.Canvas.Render()
}

func (s *Screen) ShowImageFor(ctx context.Context, i image.Image, d time.Duration) error {
	if s.Transform != nil {
		i = s.Transform(i)
	}

	draw.Draw(s.Canvas, s.Canvas.Bounds(), i, image.Point{}, draw.Over)
	err := s.Canvas.Render()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(d):
		return nil
	}
}

func (s *Screen) ShowImages(ctx context.Context, imgs []image.Image, ds []time.Duration, loop bool) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		i := 0
		for i < len(imgs) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.ShowImageFor(ctx, imgs[i], ds[i]); err != nil {
					return
				}
			}
			if i++; i == len(imgs) && loop {
				i = 0
			}
		}
	}()
	return done
}

func (s *Screen) LoopGIF(ctx context.Context, gif *gif.GIF) chan struct{} {
	delay := make([]time.Duration, len(gif.Delay))
	images := make([]image.Image, len(gif.Image))
	for i, img := range gif.Image {
		images[i] = img
		delay[i] = time.Duration(gif.Delay[i]*10) * time.Millisecond
	}
	defer s.Canvas.Clear()
	return s.ShowImages(ctx, images, delay, true)
}

func (s *Screen) PlayGIF(ctx context.Context, gif *gif.GIF) chan struct{} {
	delay := make([]time.Duration, len(gif.Delay))
	images := make([]image.Image, len(gif.Image))
	for i, img := range gif.Image {
		images[i] = img
		delay[i] = time.Duration(gif.Delay[i]*10) * time.Millisecond
	}
	defer s.Canvas.Clear()
	return s.ShowImages(ctx, images, delay, false)
}

// DrawText draws the given text onto the image at the specified coordinates with the given color.
func (s *Screen) DrawText(font *BDFFont, text string, x, y int, color color.Color) {
	cursorX := x
	for _, ch := range text {
		glyph := font.Glyphs[ch]
		if glyph == nil {
			cursorX += 6 // fallback spacing
			continue
		}
		for row := 0; row < len(glyph.Bitmap); row++ {
			line := glyph.Bitmap[row]
			// Parse hex row data into bytes
			rowData, _ := hex.DecodeString(line)
			for col := 0; col < glyph.Width; col++ {
				byteIndex := col / 8
				bitIndex := 7 - (col % 8) // BDF stores MSB first

				if byteIndex < len(rowData) && (rowData[byteIndex]&(1<<bitIndex)) != 0 {
					px := cursorX + glyph.XOffset + col
					py := y + row // + glyph.YOffset
					if image.Pt(px, py).In(s.Canvas.Bounds()) {
						s.Canvas.Set(px, py, color)
					}
				}
			}
		}
		cursorX += glyph.DeviceWidth
	}
}

func fitCenter(img image.Image, width, height int, filter imaging.ResampleFilter) *image.NRGBA {
	maxW, maxH := width, height

	if maxW <= 0 || maxH <= 0 {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.NRGBA{}
	}

	if srcW <= maxW && srcH <= maxH {
		dst := image.NewNRGBA(image.Rect(0, 0, maxW, maxH))
		offsetX := (maxW - srcW) / 2
		offsetY := (maxH - srcH) / 2
		for y := 0; y < srcH; y++ {
			for x := 0; x < srcW; x++ {
				dst.Set(offsetX+x, offsetY+y, img.At(srcBounds.Min.X+x, srcBounds.Min.Y+y))
			}
		}
		return dst
	}

	srcAspectRatio := float64(srcW) / float64(srcH)
	maxAspectRatio := float64(maxW) / float64(maxH)

	var newW, newH int
	if srcAspectRatio > maxAspectRatio {
		newW = maxW
		newH = int(float64(newW) / srcAspectRatio)
	} else {
		newH = maxH
		newW = int(float64(newH) * srcAspectRatio)
	}

	resizedImg := imaging.Resize(img, newW, newH, filter)
	resizedBounds := resizedImg.Bounds()
	resizedW := resizedBounds.Dx()
	resizedH := resizedBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, maxW, maxH))
	offsetX := (maxW - resizedW) / 2
	offsetY := (maxH - resizedH) / 2

	for y := 0; y < resizedH; y++ {
		for x := 0; x < resizedW; x++ {
			dst.Set(offsetX+x, offsetY+y, resizedImg.At(resizedBounds.Min.X+x, resizedBounds.Min.Y+y))
		}
	}

	return dst
}

func (s *Screen) Clear() error { return s.Canvas.Clear() }
func (s *Screen) Close() error { return s.Canvas.Close() }
