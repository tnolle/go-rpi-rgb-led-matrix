package renderers

import (
	"context"
	"fmt"
	"log"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/animation"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/dashboard"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

type preparedRenderer struct {
	renderer Renderer
	async    bool
}

func (p preparedRenderer) start(ctx context.Context, callbacks ...AfterRenderFunc) error {
	if !p.async {
		return p.renderer.Render(ctx, callbacks...)
	}
	go func() {
		if err := p.renderer.Render(ctx, callbacks...); err != nil {
			log.Printf("renderer stopped with error: %v", err)
		}
	}()
	return nil
}

func prepare(ctx context.Context, cmd Command, screen *rgbmatrix.Screen) (preparedRenderer, error) {
	switch cmd.Type {
	case TypePlayground:
		return preparedRenderer{renderer: MarbleShader(screen), async: true}, nil
	case TypeImage:
		renderer, err := Image(screen, cmd.Name)
		return preparedRenderer{renderer: renderer}, err
	case TypeGIF:
		renderer, err := GIFLoop(screen, cmd.Name)
		return preparedRenderer{renderer: renderer}, err
	case TypeGIFOnce:
		renderer, err := GIFOnce(screen, cmd.Name)
		return preparedRenderer{renderer: renderer}, err
	case TypeDashboard:
		return prepareDashboard(ctx, cmd.Name, screen)
	case TypeAnimation:
		return prepareAnimation(cmd.Name, screen)
	default:
		return preparedRenderer{}, fmt.Errorf("%w: renderer type %d", ErrUnknownContent, cmd.Type)
	}
}

func prepareDashboard(ctx context.Context, name string, screen *rgbmatrix.Screen) (preparedRenderer, error) {
	value, err := dashboard.DashboardString(name)
	if err != nil {
		return preparedRenderer{}, fmt.Errorf("%w: dashboard %q", ErrUnknownContent, name)
	}

	var renderer Renderer
	switch value {
	case dashboard.Clock:
		renderer = Clock(screen)
	case dashboard.Autodarts:
		renderer, err = UserCountDashboard(screen)
	case dashboard.Shopify:
		renderer, err = ShopifyDashboard(ctx, screen)
	}
	return preparedRenderer{renderer: renderer, async: true}, err
}

func prepareAnimation(name string, screen *rgbmatrix.Screen) (preparedRenderer, error) {
	value, err := animation.AnimationString(name)
	if err != nil {
		return preparedRenderer{}, fmt.Errorf("%w: animation %q", ErrUnknownContent, name)
	}

	factories := map[animation.Animation]func(*rgbmatrix.Screen) Renderer{
		animation.Aurora:             func(s *rgbmatrix.Screen) Renderer { return Aurora(s) },
		animation.BlobbyFusion:       func(s *rgbmatrix.Screen) Renderer { return BlobbyFusion(s) },
		animation.Checkerboard:       func(s *rgbmatrix.Screen) Renderer { return Checkerboard(s) },
		animation.ColorWave:          func(s *rgbmatrix.Screen) Renderer { return ColorWave(s) },
		animation.Firefly:            func(s *rgbmatrix.Screen) Renderer { return Firefly(s) },
		animation.Kaleidoscope:       func(s *rgbmatrix.Screen) Renderer { return Kaleidoscope(s) },
		animation.LavaLamp:           func(s *rgbmatrix.Screen) Renderer { return LavaLamp(s) },
		animation.Lightning:          func(s *rgbmatrix.Screen) Renderer { return Lightning(s) },
		animation.Mandelbrot:         func(s *rgbmatrix.Screen) Renderer { return Mandelbrot(s) },
		animation.MatrixRain:         func(s *rgbmatrix.Screen) Renderer { return MatrixRain(s) },
		animation.Nebula:             func(s *rgbmatrix.Screen) Renderer { return Nebula(s) },
		animation.Plasma:             func(s *rgbmatrix.Screen) Renderer { return Plasma(s) },
		animation.RadarSweep:         func(s *rgbmatrix.Screen) Renderer { return RadarSweep(s) },
		animation.Ripple:             func(s *rgbmatrix.Screen) Renderer { return Ripple(s) },
		animation.Spectrum:           func(s *rgbmatrix.Screen) Renderer { return Spectrum(s) },
		animation.Spiral:             func(s *rgbmatrix.Screen) Renderer { return Spiral(s) },
		animation.Starfield:          func(s *rgbmatrix.Screen) Renderer { return Starfield(s) },
		animation.Tunnel:             func(s *rgbmatrix.Screen) Renderer { return Tunnel(s) },
		animation.Vortex:             func(s *rgbmatrix.Screen) Renderer { return Vortex(s) },
		animation.PixelBloom:         func(s *rgbmatrix.Screen) Renderer { return PixelBloom(s) },
		animation.RGBFlow:            func(s *rgbmatrix.Screen) Renderer { return RGBFlow(s) },
		animation.Glitch:             func(s *rgbmatrix.Screen) Renderer { return Glitch(s) },
		animation.HypnoticRings:      func(s *rgbmatrix.Screen) Renderer { return HypnoticRings(s) },
		animation.SpinningGrid:       func(s *rgbmatrix.Screen) Renderer { return SpinningGrid(s) },
		animation.HexPulse:           func(s *rgbmatrix.Screen) Renderer { return HexPulse(s) },
		animation.SnakeTrail:         func(s *rgbmatrix.Screen) Renderer { return SnakeTrail(s) },
		animation.ExplosionBurst:     func(s *rgbmatrix.Screen) Renderer { return ExplosionBurst(s) },
		animation.BeatGrid:           func(s *rgbmatrix.Screen) Renderer { return BeatGrid(s) },
		animation.AudioOrbit:         func(s *rgbmatrix.Screen) Renderer { return AudioOrbit(s) },
		animation.AuroraCurtains:     func(s *rgbmatrix.Screen) Renderer { return AuroraCurtains(s) },
		animation.UlamSpiral:         func(s *rgbmatrix.Screen) Renderer { return UlamSpiral(s) },
		animation.GameOfLife:         func(s *rgbmatrix.Screen) Renderer { return GameOfLife(s) },
		animation.VectorFieldFlow:    func(s *rgbmatrix.Screen) Renderer { return VectorFieldFlow(s) },
		animation.SierpinskiTriangle: func(s *rgbmatrix.Screen) Renderer { return SierpinskiTriangle(s) },
		animation.FluidDream:         func(s *rgbmatrix.Screen) Renderer { return FluidDream(s) },
		animation.FluidRainbow:       func(s *rgbmatrix.Screen) Renderer { return FluidRainbow(s) },
		animation.OrbitingMetaballs:  func(s *rgbmatrix.Screen) Renderer { return OrbitingMetaballs(s) },
		animation.MarbleShader:       func(s *rgbmatrix.Screen) Renderer { return MarbleShader(s) },
		animation.SolarSystem:        func(s *rgbmatrix.Screen) Renderer { return SolarSystem(s) },
		animation.Darts_180:          func(s *rgbmatrix.Screen) Renderer { return Darts180(s) },
	}
	factory, ok := factories[value]
	if !ok {
		return preparedRenderer{}, fmt.Errorf("%w: animation %q", ErrUnknownContent, name)
	}
	return preparedRenderer{renderer: factory(screen), async: true}, nil
}
