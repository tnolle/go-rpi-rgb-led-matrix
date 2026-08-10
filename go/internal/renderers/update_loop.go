package renderers

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func updateLoop(ctx context.Context, commands chan Command, m rgbmatrix.Matrix, frames *display.Hub) {
	if frames != nil {
		m = rgbmatrix.NewObservable(m, frames.Publish)
	}
	s := rgbmatrix.NewScreen(m)
	defer s.Close()
	width, height := m.Geometry()
	log.Printf("renderer loop started: matrix=%dx%d streaming=%t", width, height, frames != nil)

	go func() { commands <- Command{Type: TypeImage, Name: "autodarts"} }()
	//go func() { commands <- Command{Type: TypeDashboard, Name: dashboard.Shopify.String()} }()

	renderCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var lastCommand Command
	resetScreen := func() { commands <- lastCommand }

	for {
		select {
		case <-ctx.Done():
			log.Printf("renderer loop stopping")
			cancel()
			return
		case cmd := <-commands:
			log.Printf("renderer requested: type=%s name=%q temporary=%t", cmd.Type, cmd.Name, cmd.IsTemporary)
			prepareCtx := cmd.Context
			if prepareCtx == nil {
				prepareCtx = ctx
			}
			prepared, err := prepare(prepareCtx, cmd, s)
			if err != nil {
				log.Printf("renderer rejected: type=%s name=%q error=%v", cmd.Type, cmd.Name, err)
				respond(cmd, err)
				continue
			}

			cancel()
			renderCtx, cancel = context.WithCancel(ctx)
			callbacks := []AfterRenderFunc(nil)
			if cmd.IsTemporary {
				callbacks = append(callbacks, resetScreen)
			}
			if err := prepared.start(renderCtx, callbacks...); err != nil {
				log.Printf("renderer failed to start: type=%s name=%q error=%v", cmd.Type, cmd.Name, err)
				respond(cmd, err)
				continue
			}
			if !cmd.IsTemporary {
				lastCommand = cmd
				lastCommand.Context = nil
				lastCommand.Result = nil
			}
			log.Printf("renderer started: type=%s name=%q temporary=%t", cmd.Type, cmd.Name, cmd.IsTemporary)
			respond(cmd, nil)
		}
	}
}

func respond(cmd Command, err error) {
	if cmd.Result == nil {
		return
	}
	select {
	case cmd.Result <- err:
	default:
	}
}

func UpdateLoopWithMatrix(ctx context.Context, commands chan Command, matrix rgbmatrix.Matrix, frames *display.Hub) {
	updateLoop(ctx, commands, matrix, frames)
}

type SoftBloomRingsRenderer struct {
	screen *rgbmatrix.Screen
}

func SoftBloomRings(screen *rgbmatrix.Screen) *SoftBloomRingsRenderer {
	return &SoftBloomRingsRenderer{screen: screen}
}

func (r *SoftBloomRingsRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	cx, cy := w/2, h/2
	start := time.Now()

	const maxRings = 5
	const ringSpacing = 10.0

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx, dy := x-cx, y-cy
					dist := math.Hypot(dx, dy)

					// Create multiple rings using modulus
					progress := dist - t*20
					ringPhase := math.Mod(progress, ringSpacing)

					if ringPhase < 2.5 { // threshold for thickness
						hue := math.Mod(0.6+dist*0.01+t*0.1, 1.0)
						brightness := 1.0 - (ringPhase / 2.5)
						brightness *= 0.8

						r, g, b := hsvToRGB(hue, 1.0, brightness)
						r = math.Max(r, 0.15)
						g = math.Max(g, 0.15)
						b = math.Max(b, 0.15)

						dc.SetRGB(r, g, b)
						dc.SetPixel(int(x), int(y))
					}
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}
