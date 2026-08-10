package renderers

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

const solarSystemFrameInterval = time.Second / 30

type SolarSystemRenderer struct {
	screen *rgbmatrix.Screen
	stars  []solarSystemStar
}

type solarSystemStar struct {
	x, y  float64
	phase float64
}

type solarSystemPlanet struct {
	orbit    float64
	speed    float64
	phase    float64
	size     float64
	red      float64
	green    float64
	blue     float64
	hasRings bool
	hasMoon  bool
}

var solarSystemPlanets = []solarSystemPlanet{
	{orbit: 0.15, speed: 1.60, phase: 0.2, size: 0.8, red: 0.65, green: 0.62, blue: 0.58},
	{orbit: 0.25, speed: 1.25, phase: 2.4, size: 1.3, red: 0.94, green: 0.67, blue: 0.30},
	{orbit: 0.36, speed: 1.00, phase: 4.1, size: 1.4, red: 0.15, green: 0.48, blue: 1.00, hasMoon: true},
	{orbit: 0.47, speed: 0.80, phase: 1.3, size: 1.1, red: 0.90, green: 0.25, blue: 0.12},
	{orbit: 0.61, speed: 0.45, phase: 5.2, size: 2.7, red: 0.82, green: 0.59, blue: 0.38},
	{orbit: 0.73, speed: 0.34, phase: 3.0, size: 2.3, red: 0.88, green: 0.76, blue: 0.46, hasRings: true},
	{orbit: 0.84, speed: 0.25, phase: 0.8, size: 1.8, red: 0.40, green: 0.86, blue: 0.90},
	{orbit: 0.95, speed: 0.20, phase: 4.8, size: 1.7, red: 0.22, green: 0.38, blue: 0.95},
}

func SolarSystem(screen *rgbmatrix.Screen) *SolarSystemRenderer {
	width, height := screen.Canvas.Bounds().Dx(), screen.Canvas.Bounds().Dy()
	rng := rand.New(rand.NewSource(42))
	starCount := max(12, width*height/180)
	stars := make([]solarSystemStar, starCount)
	for i := range stars {
		stars[i] = solarSystemStar{
			x:     rng.Float64() * float64(width),
			y:     rng.Float64() * float64(height),
			phase: rng.Float64() * 2 * math.Pi,
		}
	}
	return &SolarSystemRenderer{screen: screen, stars: stars}
}

func (r *SolarSystemRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	start := time.Now()
	ticker := time.NewTicker(solarSystemFrameInterval)
	defer ticker.Stop()

	for {
		if err := r.drawFrame(ctx, dc, time.Since(start).Seconds()); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (r *SolarSystemRenderer) drawFrame(ctx context.Context, dc *gg.Context, elapsed float64) error {
	width, height := float64(dc.Width()), float64(dc.Height())
	centerX, centerY := width/2, height/2
	maxOrbitX := math.Max(5, centerX-4)
	maxOrbitY := math.Max(3, centerY-3)

	dc.SetRGB(0, 0, 0.015)
	dc.Clear()
	r.drawStars(dc, elapsed)

	// Draw all paths first, so planets and the sun remain crisp above them.
	for _, planet := range solarSystemPlanets {
		dc.SetRGBA(0.20, 0.27, 0.38, 0.32)
		dc.SetLineWidth(0.45)
		dc.DrawEllipse(centerX, centerY, maxOrbitX*planet.orbit, maxOrbitY*planet.orbit)
		dc.Stroke()
	}

	for _, planet := range solarSystemPlanets {
		angle := elapsed*planet.speed + planet.phase
		x := centerX + math.Cos(angle)*maxOrbitX*planet.orbit
		y := centerY + math.Sin(angle)*maxOrbitY*planet.orbit
		r.drawPlanet(dc, planet, x, y, angle, width, height)
	}

	r.drawSun(dc, centerX, centerY, math.Max(2.2, math.Min(width, height)*0.055), elapsed)
	return r.screen.ShowImage(ctx, dc.Image())
}

func (r *SolarSystemRenderer) drawStars(dc *gg.Context, elapsed float64) {
	for _, star := range r.stars {
		brightness := 0.20 + 0.45*(0.5+0.5*math.Sin(elapsed*1.8+star.phase))
		dc.SetRGB(brightness*0.75, brightness*0.85, brightness)
		dc.SetPixel(int(star.x), int(star.y))
	}
}

func (r *SolarSystemRenderer) drawPlanet(dc *gg.Context, planet solarSystemPlanet, x, y, angle, width, height float64) {
	scale := math.Max(0.65, math.Min(width/128, height/64))
	size := math.Max(0.75, planet.size*scale)

	// A faint highlight makes even one- and two-pixel planets readable.
	dc.SetRGBA(planet.red, planet.green, planet.blue, 0.22)
	dc.DrawCircle(x, y, size+1)
	dc.Fill()

	if planet.hasRings {
		dc.SetRGBA(0.86, 0.72, 0.48, 0.9)
		dc.SetLineWidth(math.Max(0.55, scale*0.7))
		dc.DrawEllipse(x, y, size*2.1, math.Max(0.7, size*0.55))
		dc.Stroke()
	}

	dc.SetRGB(planet.red, planet.green, planet.blue)
	dc.DrawCircle(x, y, size)
	dc.Fill()
	dc.SetRGBA(1, 1, 1, 0.45)
	dc.DrawCircle(x-size*0.3, y-size*0.3, math.Max(0.35, size*0.28))
	dc.Fill()

	if planet.hasMoon {
		moonAngle := angle*8 + 1.4
		moonDistance := size + math.Max(1.4, scale*1.8)
		moonX := x + math.Cos(moonAngle)*moonDistance
		moonY := y + math.Sin(moonAngle)*moonDistance*0.65
		dc.SetRGB(0.78, 0.80, 0.84)
		dc.DrawCircle(moonX, moonY, math.Max(0.45, scale*0.5))
		dc.Fill()
	}
}

func (r *SolarSystemRenderer) drawSun(dc *gg.Context, x, y, radius, elapsed float64) {
	pulse := 1 + 0.08*math.Sin(elapsed*2.2)
	for layer := 3; layer >= 1; layer-- {
		fraction := float64(layer) / 3
		dc.SetRGBA(1.0, 0.35+0.20*fraction, 0.02, 0.10)
		dc.DrawCircle(x, y, radius*pulse+float64(layer)*1.5)
		dc.Fill()
	}
	dc.SetRGB(1.0, 0.42, 0.03)
	dc.DrawCircle(x, y, radius*pulse)
	dc.Fill()
	dc.SetRGB(1.0, 0.88, 0.25)
	dc.DrawCircle(x-radius*0.2, y-radius*0.2, radius*0.58*pulse)
	dc.Fill()
}
