package renderers

import (
	"context"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func updateLoop(ctx context.Context, commands chan Command, m rgbmatrix.Matrix) {
	s := rgbmatrix.NewScreen(m)
	defer s.Close()

	//go func() { commands <- Command{Type: TypeImage, Name: "autodarts"} }()
	go func() { commands <- Command{Type: TypeDashboard, Name: string(DasboardShopify)} }()

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
				switch DashboardName(cmd.Name) {
				case DashboardClock:
					go Clock(s).Render(renderCtx)
				case DashboardUserCount:
					go UserCountDashboard(s).Render(renderCtx)
				case DasboardShopify:
					go ShopifyDashboard(s).Render(renderCtx)
				}

			case TypeAnimation:
				switch AnimationName(cmd.Name) {
				case AnimationAurora:
					go Aurora(s).Render(renderCtx)
				case AnimationCheckerboard:
					go Checkerboard(s).Render(renderCtx)
				case AnimationColorWave:
					go ColorWave(s).Render(renderCtx)
				case AnimationBlobbyFusion:
					go BlobbyFusion(s).Render(renderCtx)
				case AnimationFirefly:
					go Firefly(s).Render(renderCtx)
				case AnimationKaleidoscope:
					go Kaleidoscope(s).Render(renderCtx)
				case AnimationLavaLamp:
					go LavaLamp(s).Render(renderCtx)
				case AnimationLightning:
					go Lightning(s).Render(renderCtx)
				case AnimationMandelbrot:
					go Mandelbrot(s).Render(renderCtx)
				case AnimationMatrixRain:
					go MatrixRain(s).Render(renderCtx)
				case AnimationNebula:
					go Nebula(s).Render(renderCtx)
				case AnimationPlasma:
					go Plasma(s).Render(renderCtx)
				case AnimationRadarSweep:
					go RadarSweep(s).Render(renderCtx)
				case AnimationRipple:
					go Ripple(s).Render(renderCtx)
				case AnimationSpectrum:
					go Spectrum(s).Render(renderCtx)
				case AnimationSpiral:
					go Spiral(s).Render(renderCtx)
				case AnimationStarfield:
					go Starfield(s).Render(renderCtx)
				case AnimationTunnel:
					go Tunnel(s).Render(renderCtx)
				case AnimationVortex:
					go Vortex(s).Render(renderCtx)
				case AnimationPixelBloom:
					go PixelBloom(s).Render(renderCtx)
				case AnimationRGBFlow:
					go RGBFlow(s).Render(renderCtx)
				case AnimationGlitch:
					go Glitch(s).Render(renderCtx)
				case AnimationHypnoticRings:
					go HypnoticRings(s).Render(renderCtx)
				case AnimationSpinningGrid:
					go SpinningGrid(s).Render(renderCtx)
				case AnimationHexPulse:
					go HexPulse(s).Render(renderCtx)
				case AnimationSnakeTrail:
					go SnakeTrail(s).Render(renderCtx)
				case AnimationExplosionBurst:
					go ExplosionBurst(s).Render(renderCtx)
				case AnimationBeatGrid:
					go BeatGrid(s).Render(renderCtx)
				case AnimationAudioOrbit:
					go AudioOrbit(s).Render(renderCtx)
				case AnimationAuroraCurtains:
					go AuroraCurtains(s).Render(renderCtx)
				case AnimationUlamSpiral:
					go UlamSpiral(s).Render(renderCtx)
				case AnimationGameOfLife:
					go GameOfLife(s).Render(renderCtx)
				case AnimationVectorFieldFlow:
					go VectorFieldFlow(s).Render(renderCtx)
				case AnimationSierpinskiTriangle:
					go SierpinskiTriangle(s).Render(renderCtx)
				case AnimationFluidDream:
					go FluidDream(s).Render(renderCtx)
				case AnimationFluidRainbow:
					go FluidRainbow(s).Render(renderCtx)
				case AnimationOrbitingMetaballs:
					go OrbitingMetaballs(s).Render(renderCtx)
				case AnimationMarbleShader:
					go MarbleShader(s).Render(renderCtx)
				}

			}
		}
	}
}

func UpdateLoopEmulated(ctx context.Context, commands chan Command, config rgbmatrix.Config) {
	m := rgbmatrix.NewEmulator(64*3, 32*3, 12)
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
