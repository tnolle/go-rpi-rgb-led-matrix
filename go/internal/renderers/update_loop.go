package renderers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/animation"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/dashboard"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func updateLoop(ctx context.Context, commands chan Command, m rgbmatrix.Matrix) {
	s := rgbmatrix.NewScreen(m)
	defer s.Close()

	go func() { commands <- Command{Type: TypeImage, Name: "autodarts"} }()
	//go func() { commands <- Command{Type: TypeDashboard, Name: dashboard.Shopify.String()} }()

	renderCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var lastCommand Command
	resetScreen := func() { commands <- lastCommand }

	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case cmd := <-commands:
			cancel()
			renderCtx, cancel = context.WithCancel(ctx)

			if !cmd.IsTemporary {
				lastCommand = cmd
			}
			switch cmd.Type {
			case TypePlayground:
				go MarbleShader(s).Render(renderCtx)

			case TypeImage:
				Image(s, cmd.Name).Render(renderCtx)
			case TypeGIF:
				GIFLoop(s, cmd.Name).Render(renderCtx)
			case TypeGIFOnce:
				GIFOnce(s, cmd.Name).Render(renderCtx, resetScreen)

			case TypeDashboard:
				v, err := dashboard.DashboardString(cmd.Name)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					continue
				}

				switch v {
				case dashboard.Clock:
					go Clock(s).Render(renderCtx)
				case dashboard.Autodarts:
					go UserCountDashboard(s).Render(renderCtx)
				case dashboard.Shopify:
					go ShopifyDashboard(s).Render(renderCtx)
				}

			case TypeAnimation:
				v, err := animation.AnimationString(cmd.Name)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					continue
				}

				switch v {
				case animation.Aurora:
					go Aurora(s).Render(renderCtx)
				case animation.Checkerboard:
					go Checkerboard(s).Render(renderCtx)
				case animation.ColorWave:
					go ColorWave(s).Render(renderCtx)
				case animation.BlobbyFusion:
					go BlobbyFusion(s).Render(renderCtx)
				case animation.Firefly:
					go Firefly(s).Render(renderCtx)
				case animation.Kaleidoscope:
					go Kaleidoscope(s).Render(renderCtx)
				case animation.LavaLamp:
					go LavaLamp(s).Render(renderCtx)
				case animation.Lightning:
					go Lightning(s).Render(renderCtx)
				case animation.Mandelbrot:
					go Mandelbrot(s).Render(renderCtx)
				case animation.MatrixRain:
					go MatrixRain(s).Render(renderCtx)
				case animation.Nebula:
					go Nebula(s).Render(renderCtx)
				case animation.Plasma:
					go Plasma(s).Render(renderCtx)
				case animation.RadarSweep:
					go RadarSweep(s).Render(renderCtx)
				case animation.Ripple:
					go Ripple(s).Render(renderCtx)
				case animation.Spectrum:
					go Spectrum(s).Render(renderCtx)
				case animation.Spiral:
					go Spiral(s).Render(renderCtx)
				case animation.Starfield:
					go Starfield(s).Render(renderCtx)
				case animation.Tunnel:
					go Tunnel(s).Render(renderCtx)
				case animation.Vortex:
					go Vortex(s).Render(renderCtx)
				case animation.PixelBloom:
					go PixelBloom(s).Render(renderCtx)
				case animation.RGBFlow:
					go RGBFlow(s).Render(renderCtx)
				case animation.Glitch:
					go Glitch(s).Render(renderCtx)
				case animation.HypnoticRings:
					go HypnoticRings(s).Render(renderCtx)
				case animation.SpinningGrid:
					go SpinningGrid(s).Render(renderCtx)
				case animation.HexPulse:
					go HexPulse(s).Render(renderCtx)
				case animation.SnakeTrail:
					go SnakeTrail(s).Render(renderCtx)
				case animation.ExplosionBurst:
					go ExplosionBurst(s).Render(renderCtx)
				case animation.BeatGrid:
					go BeatGrid(s).Render(renderCtx)
				case animation.AudioOrbit:
					go AudioOrbit(s).Render(renderCtx)
				case animation.AuroraCurtains:
					go AuroraCurtains(s).Render(renderCtx)
				case animation.UlamSpiral:
					go UlamSpiral(s).Render(renderCtx)
				case animation.GameOfLife:
					go GameOfLife(s).Render(renderCtx)
				case animation.VectorFieldFlow:
					go VectorFieldFlow(s).Render(renderCtx)
				case animation.SierpinskiTriangle:
					go SierpinskiTriangle(s).Render(renderCtx)
				case animation.FluidDream:
					go FluidDream(s).Render(renderCtx)
				case animation.FluidRainbow:
					go FluidRainbow(s).Render(renderCtx)
				case animation.OrbitingMetaballs:
					go OrbitingMetaballs(s).Render(renderCtx)
				case animation.MarbleShader:
					go MarbleShader(s).Render(renderCtx)
				}

			}
		}
	}
}

func UpdateLoopEmulated(ctx context.Context, commands chan Command, config rgbmatrix.Config) {
	fmt.Printf("Emulating a %dx%d matrix\n", config.Options.Cols*config.Options.ChainLength, config.Options.Rows*config.Options.Parallel)
	m := rgbmatrix.NewEmulator(config.Options.Cols*config.Options.ChainLength, config.Options.Rows*config.Options.Parallel, 8)
	m.Run(func() { updateLoop(ctx, commands, m) })
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
