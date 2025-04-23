package renderers

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

type BeatGridRenderer struct {
	screen *rgbmatrix.Screen
}

func BeatGrid(screen *rgbmatrix.Screen) *BeatGridRenderer {
	return &BeatGridRenderer{screen: screen}
}

func (r *BeatGridRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	start := time.Now()

	gridSize := 8.0 // size of each square
	cols := int(width / gridSize)
	rows := int(height / gridSize)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()

			for y := 0; y < rows; y++ {
				for x := 0; x < cols; x++ {
					phase := float64(x+y) * 0.6
					brightness := 0.4 + 0.6*math.Sin(t*4+phase)
					hue := math.Mod(t*0.15+float64(x+y)*0.05, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, brightness)
					dc.SetRGB(r, g, b)
					dc.DrawRectangle(float64(x)*gridSize, float64(y)*gridSize, gridSize-1, gridSize-1)
					dc.Fill()
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type RGBFlowRenderer struct {
	screen *rgbmatrix.Screen
}

func RGBFlow(screen *rgbmatrix.Screen) *RGBFlowRenderer {
	return &RGBFlowRenderer{screen: screen}
}

func (r *RGBFlowRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := float64(dc.Width()), float64(dc.Height())
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()
			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					phase := math.Sin(x*0.2 + y*0.2 + t)
					hue := math.Mod((phase+1)/2+0.5*t, 1.0)
					rVal, gVal, bVal := hsvToRGB(hue, 1.0, 1.0)
					dc.SetRGB(math.Max(rVal, 0.1), math.Max(gVal, 0.1), math.Max(bVal, 0.1))
					dc.SetPixel(int(x), int(y))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type PixelBloomRenderer struct {
	screen *rgbmatrix.Screen
}

func PixelBloom(screen *rgbmatrix.Screen) *PixelBloomRenderer {
	return &PixelBloomRenderer{screen: screen}
}

func (r *PixelBloomRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := float64(dc.Width()), float64(dc.Height())
	cx, cy := w/2, h/2
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()
			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					dist := math.Hypot(dx, dy)
					pulse := math.Sin(dist*0.2 - t*2)
					hue := math.Mod(0.6+dist*0.01+t*0.05, 1.0)
					bright := 0.5 + 0.5*pulse
					rVal, gVal, bVal := hsvToRGB(hue, 1.0, bright)
					dc.SetRGB(math.Max(rVal, 0.1), math.Max(gVal, 0.1), math.Max(bVal, 0.1))
					dc.SetPixel(int(x), int(y))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type GlitchRenderer struct {
	screen *rgbmatrix.Screen
}

func Glitch(screen *rgbmatrix.Screen) *GlitchRenderer {
	return &GlitchRenderer{screen: screen}
}

func (r *GlitchRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.Clear()
			baseHue := rand.Float64()
			for y := 0; y < h; y++ {
				offset := 0
				if rng.Float64() < 0.2 {
					offset = rng.Intn(4) - 2 // -2 to 2 pixel horizontal glitch
				}
				hue := math.Mod(baseHue+float64(y)/float64(h), 1.0)
				rVal, gVal, bVal := hsvToRGB(hue, 1.0, 1.0)
				for x := 0; x < w; x++ {
					tx := (x + offset + w) % w
					dc.SetRGB(rVal, gVal, bVal)
					dc.SetPixel(tx, y)
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type RadarSweepRenderer struct {
	screen *rgbmatrix.Screen
}

func RadarSweep(screen *rgbmatrix.Screen) *RadarSweepRenderer {
	return &RadarSweepRenderer{screen: screen}
}

func (r *RadarSweepRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	cx := width / 2
	cy := height / 2
	start := time.Now()

	decay := make([][]float64, int(height))
	for y := range decay {
		decay[y] = make([]float64, int(width))
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			// Update decay values
			for y := 0; y < int(height); y++ {
				for x := 0; x < int(width); x++ {
					decay[y][x] *= 0.94
				}
			}

			// Draw radar sweep beam
			sweepAngle := now * 1.5
			for d := 0.0; d < math.Min(width, height)/2; d += 0.3 {
				x := cx + math.Cos(sweepAngle)*d
				y := cy + math.Sin(sweepAngle)*d
				if x >= 0 && x < width && y >= 0 && y < height {
					ix, iy := int(x), int(y)
					decay[iy][ix] = 1.0
				}
			}

			// Render all pixels with current decay brightness
			for y := 0; y < int(height); y++ {
				for x := 0; x < int(width); x++ {
					brightness := decay[y][x]
					if brightness > 0.01 {
						r, g, b := hsvToRGB(0.33, 1.0, brightness)
						dc.SetRGB(r, g, b)
						dc.SetPixel(x, y)
					}
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type NebulaRenderer struct {
	screen *rgbmatrix.Screen
}

func Nebula(screen *rgbmatrix.Screen) *NebulaRenderer {
	return &NebulaRenderer{screen: screen}
}

func (r *NebulaRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	start := time.Now()

	cx := w / 2
	cy := h / 2

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					// Normalize coordinates
					nx := x / w
					ny := y / h
					dist := math.Hypot(x-cx, y-cy)

					v := math.Sin(nx*10+now) +
						math.Sin(ny*10-now*1.3) +
						math.Sin((nx+ny)*10+now*1.1) +
						math.Sin(dist*0.25-now*0.7)

					hue := math.Mod(0.6+v*0.05+now*0.01, 1.0)
					brightness := 0.3 + 0.7*(0.5+0.5*math.Sin(v+now))

					r, g, b := hsvToRGB(hue, 1.0, brightness)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type AuroraRenderer struct {
	screen *rgbmatrix.Screen
}

func Aurora(screen *rgbmatrix.Screen) *AuroraRenderer {
	return &AuroraRenderer{screen: screen}
}

func (r *AuroraRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < height; y++ {
				for x := 0.0; x < width; x++ {
					xf := x / width * 2 * math.Pi
					yf := y / height

					wave := math.Sin(xf*2 + now*1.5)
					curve := math.Sin(yf*4*math.Pi + wave + now)

					hue := math.Mod(0.4+0.2*wave+now*0.01, 1.0)
					brightness := 0.3 + 0.7*(0.5+0.5*curve)

					r, g, b := hsvToRGB(hue, 1.0, brightness)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type LavaLampRenderer struct {
	screen *rgbmatrix.Screen
}

func LavaLamp(screen *rgbmatrix.Screen) *LavaLampRenderer {
	return &LavaLampRenderer{screen: screen}
}

func (r *LavaLampRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					xf := x / w * 2 * math.Pi
					yf := y / h * 2 * math.Pi

					value := math.Sin(xf*2+now) +
						math.Sin(yf*3+now*0.7) +
						math.Sin((xf+yf)*2+now*1.3)

					hue := math.Mod((value+3)/6+now*0.02, 1.0)
					brightness := 0.4 + 0.6*math.Sin(value*3+now)

					r, g, b := hsvToRGB(hue, 0.8, brightness)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type ColorWaveRenderer struct {
	screen *rgbmatrix.Screen
}

func ColorWave(screen *rgbmatrix.Screen) *ColorWaveRenderer {
	return &ColorWaveRenderer{screen: screen}
}

func (r *ColorWaveRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())

	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					offset := math.Sin(now*0.5) * 20 // Directional offset oscillates over time
					xx := x + offset
					yy := y + offset

					// Generate evolving pattern using sine waves
					red := 0.5 + 0.5*math.Sin((xx+now*30)*0.1)
					green := 0.5 + 0.5*math.Sin((yy+now*40)*0.1)
					blue := 0.5 + 0.5*math.Sin((xx+yy+now*50)*0.1)

					// Prevent black by ensuring a minimum color threshold
					red = math.Max(red, 0.1)
					green = math.Max(green, 0.1)
					blue = math.Max(blue, 0.1)

					dc.SetRGB(red, green, blue)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type PlasmaRenderer struct {
	screen *rgbmatrix.Screen
}

func Plasma(screen *rgbmatrix.Screen) *PlasmaRenderer {
	return &PlasmaRenderer{screen: screen}
}

func (r *PlasmaRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())

	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					value := math.Sin(x*0.1+now) +
						math.Sin(y*0.1+now) +
						math.Sin((x+y)*0.1+now) +
						math.Sin(math.Hypot(x-w/2, y-h/2)*0.1-now)

					hue := (value + 4) / 8 // Normalize to [0, 1]
					hue = math.Mod(hue, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, 1.0)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)
					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type RippleRenderer struct {
	screen *rgbmatrix.Screen
}

func Ripple(screen *rgbmatrix.Screen) *RippleRenderer {
	return &RippleRenderer{screen: screen}
}

func (r *RippleRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())

	cx := w / 2
	cy := h / 2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					dist := math.Hypot(dx, dy)

					value := math.Sin(dist*0.3 - now*3)
					hue := math.Mod(value*0.25+now*0.1, 1.0)
					brightness := 0.5 + 0.5*math.Sin(dist*0.2-now*2)

					r, g, b := hsvToRGB(hue, 1.0, brightness)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type SpiralRenderer struct {
	screen *rgbmatrix.Screen
}

func Spiral(screen *rgbmatrix.Screen) *SpiralRenderer {
	return &SpiralRenderer{screen: screen}
}

func (r *SpiralRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())

	cx, cy := w/2, h/2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					angle := math.Atan2(dy, dx)

					hue := (angle + now) / (2 * math.Pi)
					hue = math.Mod(hue, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, 1.0)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)
					dc.SetRGB(r, g, b)

					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

func hsvToRGB(h, s, v float64) (r, g, b float64) {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return
}

type TunnelRenderer struct {
	screen *rgbmatrix.Screen
}

func Tunnel(screen *rgbmatrix.Screen) *TunnelRenderer {
	return &TunnelRenderer{screen: screen}
}

func (r *TunnelRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())

	cx := w / 2
	cy := h / 2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					dist := math.Hypot(dx, dy)
					angle := math.Atan2(dy, dx)

					depth := math.Sin(dist*0.1 - now)
					radius := 1.0 / (0.1 + dist*0.05)

					hue := math.Mod((angle/(2*math.Pi))+now*0.1+depth*0.5, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, radius)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type SpectrumRenderer struct {
	screen *rgbmatrix.Screen
}

func Spectrum(screen *rgbmatrix.Screen) *SpectrumRenderer {
	return &SpectrumRenderer{screen: screen}
}

func (r *SpectrumRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := dc.Width()
	height := dc.Height()
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			barCount := 16
			barWidth := width / barCount

			for i := 0; i < barCount; i++ {
				freq := 0.5 + float64(i)*0.1
				amplitude := math.Sin(now*freq+float64(i))*0.5 + 0.5
				barHeight := int(float64(height) * amplitude)

				hue := math.Mod(float64(i)/float64(barCount)+now*0.05, 1.0)

				for y := height - 1; y >= height-barHeight; y-- {
					for x := i * barWidth; x < (i+1)*barWidth; x++ {
						r, g, b := hsvToRGB(hue, 1.0, 1.0)
						r = math.Max(r, 0.1)
						g = math.Max(g, 0.1)
						b = math.Max(b, 0.1)
						dc.SetRGB(r, g, b)
						dc.SetPixel(x, y)
					}
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type StarfieldRenderer struct {
	screen *rgbmatrix.Screen
}

func Starfield(screen *rgbmatrix.Screen) *StarfieldRenderer {
	return &StarfieldRenderer{screen: screen}
}

func (r *StarfieldRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	cx, cy := float64(w)/2, float64(h)/2
	stars := make([][3]float64, 100)
	for i := range stars {
		stars[i] = [3]float64{
			(rand.Float64()*2 - 1) * cx,
			(rand.Float64()*2 - 1) * cy,
			rand.Float64()*1.5 + 0.5,
		}
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.SetRGB(0.0, 0.0, 0.05) // dark background instead of full clear
			dc.Clear()
			for i := range stars {
				stars[i][2] -= 0.02
				if stars[i][2] <= 0.1 {
					stars[i][0] = (rand.Float64()*2 - 1) * cx
					stars[i][1] = (rand.Float64()*2 - 1) * cy
					stars[i][2] = 1.5
				}
				sx := cx + stars[i][0]/stars[i][2]
				sy := cy + stars[i][1]/stars[i][2]
				if sx >= 0 && sx < float64(w) && sy >= 0 && sy < float64(h) {
					brightness := 1.0 - (stars[i][2]-0.5)/1.5
					brightness = math.Max(brightness, 0.1)
					dc.SetRGB(brightness, brightness, brightness)
					dc.SetPixel(int(sx), int(sy))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type FireflyRenderer struct {
	screen *rgbmatrix.Screen
}

func Firefly(screen *rgbmatrix.Screen) *FireflyRenderer {
	return &FireflyRenderer{screen: screen}
}

func (r *FireflyRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	type Firefly struct{ x, y, dx, dy float64 }
	ff := make([]Firefly, 20)
	for i := range ff {
		ff[i] = Firefly{rand.Float64() * float64(w), rand.Float64() * float64(h), rand.Float64() - 0.5, rand.Float64() - 0.5}
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.SetRGB(0, 0, 0.1)
			dc.Clear()
			for i := range ff {
				ff[i].x += ff[i].dx
				ff[i].y += ff[i].dy
				if ff[i].x < 0 || ff[i].x >= float64(w) {
					ff[i].dx *= -1
				}
				if ff[i].y < 0 || ff[i].y >= float64(h) {
					ff[i].dy *= -1
				}
				hue := math.Mod(ff[i].x/float64(w)+ff[i].y/float64(h), 1)
				r, g, b := hsvToRGB(hue, 1, 1)
				dc.SetRGB(r, g, b)
				dc.SetPixel(int(ff[i].x), int(ff[i].y))
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type MatrixRainRenderer struct {
	screen *rgbmatrix.Screen
}

func MatrixRain(screen *rgbmatrix.Screen) *MatrixRainRenderer {
	return &MatrixRainRenderer{screen: screen}
}

func (r *MatrixRainRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	drops := make([]int, w)
	for i := range drops {
		drops[i] = rand.Intn(h)
	}
	trailLength := 10
	lengths := make([]int, w)
	for i := range lengths {
		lengths[i] = 1 + rand.Intn(16) // lengths between 5 and 14
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			for x := 0; x < w; x++ {
				headY := drops[x]
				for t := 0; t < lengths[x]; t++ {
					y := (headY - t + h) % h
					brightness := 1.0 - float64(t)/float64(trailLength)
					rVal, gVal, bVal := 0.0, brightness, 0.0
					rVal = math.Max(rVal, 0.1)
					gVal = math.Max(gVal, 0.1)
					bVal = math.Max(bVal, 0.1)
					dc.SetRGB(rVal, gVal, bVal)
					dc.SetPixel(x, y)
				}
				drops[x] = (drops[x] + 1) % h
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(10 * time.Millisecond)
		}
	}
}

type CheckerboardRenderer struct {
	screen *rgbmatrix.Screen
}

func Checkerboard(screen *rgbmatrix.Screen) *CheckerboardRenderer {
	return &CheckerboardRenderer{screen: screen}
}

func (r *CheckerboardRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()
			size := 8 + int(4*math.Sin(t))
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if (x/size+y/size)%2 == 0 {
						dc.SetRGB(1, 1, 1)
					} else {
						dc.SetRGB(0.1, 0.1, 0.1)
					}
					dc.SetPixel(x, y)
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type VortexRenderer struct {
	screen *rgbmatrix.Screen
}

func Vortex(screen *rgbmatrix.Screen) *VortexRenderer {
	return &VortexRenderer{screen: screen}
}

func (r *VortexRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := float64(dc.Width()), float64(dc.Height())
	cx, cy := w/2, h/2
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.Clear()
			now := time.Since(start).Seconds()
			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					angle := math.Atan2(dy, dx)
					dist := math.Hypot(dx, dy)
					value := math.Sin(dist*0.1 - now + angle)
					hue := math.Mod((value+1)/2+now*0.02, 1)
					r, g, b := hsvToRGB(hue, 1, 1)
					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type LightningRenderer struct {
	screen *rgbmatrix.Screen
}

func Lightning(screen *rgbmatrix.Screen) *LightningRenderer {
	return &LightningRenderer{screen: screen}
}

func (r *LightningRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	var lastFlash time.Time
	var flashX int
	var flashing bool
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Dim background to gradually fade out
			dc.SetRGBA(0, 0, 0, 0.1)
			dc.Clear()

			now := time.Now()
			if flashing && now.Sub(lastFlash) > 100*time.Millisecond {
				flashing = false
			}
			if !flashing && (lastFlash.IsZero() || now.Sub(lastFlash) > time.Duration(rand.Intn(1000))*time.Millisecond) {
				lastFlash = now
				flashX = rand.Intn(w)
				flashing = true
			}
			if flashing {
				x := flashX
				for y := 0; y < h; y++ {
					if x >= 0 && x < w {
						dc.SetRGB(1, 1, 1)
						dc.SetPixel(x, y)
					}
					x += rand.Intn(3) - 1
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type MandelbrotRenderer struct {
	screen *rgbmatrix.Screen
}

func Mandelbrot(screen *rgbmatrix.Screen) *MandelbrotRenderer {
	return &MandelbrotRenderer{screen: screen}
}

func (r *MandelbrotRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()
			zoom := 1.5 + 0.5*math.Sin(now*0.2)
			centerX := -0.5 + 0.2*math.Sin(now*0.1)
			centerY := 0.0 + 0.2*math.Cos(now*0.1)

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					cx := (float64(x)/float64(w))*zoom + centerX - zoom/2
					cy := (float64(y)/float64(h))*zoom + centerY - zoom/2
					zx, zy := 0.0, 0.0
					iter, maxIter := 0, 30
					for zx*zx+zy*zy < 4 && iter < maxIter {
						tmp := zx*zx - zy*zy + cx
						zy, zx = 2*zx*zy+cy, tmp
						iter++
					}
					hue := float64(iter) / float64(maxIter)
					r, g, b := hsvToRGB(hue, 1, 1)
					dc.SetRGB(r, g, b)
					dc.SetPixel(x, y)
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type BlobbyFusionRenderer struct {
	screen *rgbmatrix.Screen
}

func BlobbyFusion(screen *rgbmatrix.Screen) *BlobbyFusionRenderer {
	return &BlobbyFusionRenderer{screen: screen}
}

func (r *BlobbyFusionRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	type Blob struct {
		x, y, radius, dx, dy float64
	}

	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	start := time.Now()

	// Create some initial blobs
	blobs := []Blob{}
	for i := 0; i < 5; i++ {
		blobs = append(blobs, Blob{
			x:      rand.Float64() * w,
			y:      rand.Float64() * h,
			radius: 10 + rand.Float64()*10,
			dx:     rand.Float64()*2 - 1,
			dy:     rand.Float64()*2 - 1,
		})
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			// Move blobs
			for i := range blobs {
				blobs[i].x += blobs[i].dx
				blobs[i].y += blobs[i].dy

				if blobs[i].x < 0 || blobs[i].x > w {
					blobs[i].dx *= -1
				}
				if blobs[i].y < 0 || blobs[i].y > h {
					blobs[i].dy *= -1
				}
			}

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					sum := 0.0
					for _, b := range blobs {
						dx := x - b.x
						dy := y - b.y
						dist := math.Hypot(dx, dy)
						sum += b.radius * b.radius / (dist*dist + 1)
					}

					normalized := math.Min(sum/5.0, 1.0)
					hue := math.Mod(0.6+0.3*normalized+now*0.02, 1.0)
					val := math.Pow(normalized, 1.2)

					r, g, b := hsvToRGB(hue, 0.8, val)
					dc.SetRGB(math.Max(r, 0.1), math.Max(g, 0.1), math.Max(b, 0.1))
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type KaleidoscopeRenderer struct {
	screen *rgbmatrix.Screen
}

func Kaleidoscope(screen *rgbmatrix.Screen) *KaleidoscopeRenderer {
	return &KaleidoscopeRenderer{screen: screen}
}

func (r *KaleidoscopeRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	cx := w / 2
	cy := h / 2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.Clear()
			now := time.Since(start).Seconds()
			segments := 6 // number of mirrored segments

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					angle := math.Atan2(dy, dx) + now*0.5
					dist := math.Hypot(dx, dy)

					// wrap angle into a segment
					angle = math.Mod(angle, 2*math.Pi/float64(segments))
					xx := math.Cos(angle) * dist
					yy := math.Sin(angle) * dist

					value := math.Sin(xx*0.2+now) + math.Cos(yy*0.2+now*1.1)
					hue := math.Mod(0.6+value*0.15+now*0.01, 1.0)
					brightness := 0.3 + 0.7*math.Sin(value+now)

					r, g, b := hsvToRGB(hue, 1.0, brightness)
					dc.SetRGB(math.Max(r, 0.1), math.Max(g, 0.1), math.Max(b, 0.1))
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type HypnoticRingsRenderer struct{ screen *rgbmatrix.Screen }

func HypnoticRings(screen *rgbmatrix.Screen) *HypnoticRingsRenderer {
	return &HypnoticRingsRenderer{screen: screen}
}

func (r *HypnoticRingsRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := float64(dc.Width()), float64(dc.Height())
	cx, cy := w/2, h/2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx, dy := x-cx, y-cy
					dist := math.Hypot(dx, dy)
					value := math.Sin(dist*0.2 - now*2)
					bright := 0.5 + 0.5*value
					hue := math.Mod(dist*0.01+now*0.1, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, bright)
					dc.SetRGB(math.Max(r, 0.1), math.Max(g, 0.1), math.Max(b, 0.1))
					dc.SetPixel(int(x), int(y))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type SpinningGridRenderer struct{ screen *rgbmatrix.Screen }

func SpinningGrid(screen *rgbmatrix.Screen) *SpinningGridRenderer {
	return &SpinningGridRenderer{screen: screen}
}

func (r *SpinningGridRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := float64(dc.Width()), float64(dc.Height())
	cx, cy := w/2, h/2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			angle := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()
			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx := x - cx
					dy := y - cy
					rx := dx*math.Cos(angle) - dy*math.Sin(angle)
					ry := dx*math.Sin(angle) + dy*math.Cos(angle)
					if int(rx)%10 == 0 || int(ry)%10 == 0 {
						hue := math.Mod(angle*0.1+rx*0.01+ry*0.01, 1.0)
						r, g, b := hsvToRGB(hue, 1, 1)
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

type HexPulseRenderer struct{ screen *rgbmatrix.Screen }

func HexPulse(screen *rgbmatrix.Screen) *HexPulseRenderer {
	return &HexPulseRenderer{screen: screen}
}

func (r *HexPulseRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()
			radius := 4.0
			dx := radius * 3 / 2
			dy := radius * math.Sqrt(3)
			for y := 0.0; y < float64(h); y += dy {
				for x := 0.0; x < float64(w); x += dx {
					offset := 0.0
					if int(y/dy)%2 == 1 {
						offset = radius * 0.75
					}
					dist := math.Hypot(x-float64(w)/2+offset, y-float64(h)/2)
					pulse := 0.5 + 0.5*math.Sin(dist*0.2-t*4)
					hue := math.Mod(dist*0.01+t*0.1, 1)
					r, g, b := hsvToRGB(hue, 1.0, pulse)
					dc.SetRGB(r, g, b)
					dc.DrawCircle(x+offset, y, 1)
					dc.Fill()
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type SnakeTrailRenderer struct{ screen *rgbmatrix.Screen }

func SnakeTrail(screen *rgbmatrix.Screen) *SnakeTrailRenderer {
	return &SnakeTrailRenderer{screen: screen}
}

func (r *SnakeTrailRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	type Point struct{ x, y int }

	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	snake := []Point{{w / 2, h / 2}}
	dir := []Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	heading := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			head := snake[0]
			nx := (head.x + dir[heading].x + w) % w
			ny := (head.y + dir[heading].y + h) % h
			snake = append([]Point{{nx, ny}}, snake...)
			if len(snake) > 20 {
				snake = snake[:20]
			}
			if rand.Float64() < 0.3 {
				heading = rand.Intn(4)
			}

			dc.SetRGB(0, 0, 0)
			dc.Clear()
			for i, p := range snake {
				hue := float64(i) / float64(len(snake))
				r, g, b := hsvToRGB(hue, 1, 1)
				dc.SetRGB(r, g, b)
				dc.SetPixel(p.x, p.y)
			}
			r.screen.ShowImage(ctx, dc.Image())
		}
	}
}

type ExplosionBurstRenderer struct{ screen *rgbmatrix.Screen }

func ExplosionBurst(screen *rgbmatrix.Screen) *ExplosionBurstRenderer {
	return &ExplosionBurstRenderer{screen: screen}
}

func (r *ExplosionBurstRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := float64(dc.Width()), float64(dc.Height())
	cx, cy := w/2, h/2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					dx, dy := x-cx, y-cy
					dist := math.Hypot(dx, dy)
					ring := math.Sin(dist*0.5 - now*4)
					brightness := 0.5 + 0.5*ring
					hue := math.Mod(now*0.1+dist*0.02, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, brightness)
					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type AudioOrbitRenderer struct {
	screen *rgbmatrix.Screen
}

func AudioOrbit(screen *rgbmatrix.Screen) *AudioOrbitRenderer {
	return &AudioOrbitRenderer{screen: screen}
}

func (r *AudioOrbitRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	cx := w / 2
	cy := h / 2
	start := time.Now()

	numOrbits := 6
	radius := math.Min(w, h) * 0.3

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.SetRGB(0, 0, 0)
			dc.Clear()
			now := time.Since(start).Seconds()

			for i := 0; i < numOrbits; i++ {
				angle := now*1.5 + float64(i)*math.Pi*2/float64(numOrbits)
				amp := math.Sin(now*3 + float64(i))
				ringRadius := radius + 5*amp
				x := cx + ringRadius*math.Cos(angle)
				y := cy + ringRadius*math.Sin(angle)
				hue := math.Mod(float64(i)/float64(numOrbits)+now*0.1, 1.0)
				r, g, b := hsvToRGB(hue, 1.0, 1.0)
				dc.SetRGB(r, g, b)
				dc.DrawCircle(x, y, 2)
				dc.Fill()
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type AuroraCurtainsRenderer struct {
	screen *rgbmatrix.Screen
}

func AuroraCurtains(screen *rgbmatrix.Screen) *AuroraCurtainsRenderer {
	return &AuroraCurtainsRenderer{screen: screen}
}

func (r *AuroraCurtainsRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			for x := 0.0; x < width; x++ {
				wave := math.Sin(x*0.2 + now*1.5)
				for y := 0.0; y < height; y++ {
					yRatio := y / height
					offset := math.Sin(yRatio*math.Pi*4 + wave*2 + now*0.5)
					brightness := 0.4 + 0.6*(0.5+0.5*offset)

					hue := math.Mod(0.3+0.2*wave+now*0.02, 1.0)
					r, g, b := hsvToRGB(hue, 1.0, brightness)
					r = math.Max(r, 0.1)
					g = math.Max(g, 0.1)
					b = math.Max(b, 0.1)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type UlamSpiralRenderer struct {
	screen *rgbmatrix.Screen
}

func UlamSpiral(screen *rgbmatrix.Screen) *UlamSpiralRenderer {
	return &UlamSpiralRenderer{screen: screen}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func (r *UlamSpiralRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	cx, cy := w/2, h/2
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.SetRGB(0, 0, 0)
			dc.Clear()
			x, y := 0, 0
			dx, dy := 0, -1
			steps := int(math.Max(float64(w), float64(h))) * int(math.Max(float64(w), float64(h)))
			num := 1

			for i := 0; i < steps; i++ {
				sx := cx + x
				sy := cy + y
				if sx >= 0 && sx < w && sy >= 0 && sy < h {
					if isPrime(num) {
						t := time.Since(start).Seconds()
						hue := math.Mod(float64(num)*0.01+t*0.1, 1.0)
						r, g, b := hsvToRGB(hue, 1.0, 1.0)
						dc.SetRGB(r, g, b)
						dc.SetPixel(sx, sy)
					}
				}
				if (x == y) || (x < 0 && x == -y) || (x > 0 && x == 1-y) {
					dx, dy = -dy, dx
				}
				x += dx
				y += dy
				num++
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(150 * time.Millisecond)
		}
	}
}

type GameOfLifeRenderer struct {
	screen *rgbmatrix.Screen
}

func GameOfLife(screen *rgbmatrix.Screen) *GameOfLifeRenderer {
	return &GameOfLifeRenderer{screen: screen}
}

func (r *GameOfLifeRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w, h := dc.Width(), dc.Height()
	grid := make([][]bool, h)
	next := make([][]bool, h)
	for y := range grid {
		grid[y] = make([]bool, w)
		next[y] = make([]bool, w)
		for x := range grid[y] {
			grid[y][x] = rand.Float64() < 0.2
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Draw current grid
			dc.SetRGB(0, 0, 0)
			dc.Clear()
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if grid[y][x] {
						hue := float64((x*y)%360) / 360
						r, g, b := hsvToRGB(hue, 1.0, 1.0)
						dc.SetRGB(r, g, b)
						dc.SetPixel(x, y)
					}
				}
			}
			r.screen.ShowImage(ctx, dc.Image())

			// Compute next state
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					live := 0
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							if dx == 0 && dy == 0 {
								continue
							}
							ny := (y + dy + h) % h
							nx := (x + dx + w) % w
							if grid[ny][nx] {
								live++
							}
						}
					}
					if grid[y][x] {
						next[y][x] = live == 2 || live == 3
					} else {
						next[y][x] = live == 3
					}
				}
			}

			// Swap buffers
			grid, next = next, grid
			time.Sleep(50 * time.Millisecond)
		}
	}
}

type VectorFieldFlowRenderer struct {
	screen *rgbmatrix.Screen
}

func VectorFieldFlow(screen *rgbmatrix.Screen) *VectorFieldFlowRenderer {
	return &VectorFieldFlowRenderer{screen: screen}
}

func (r *VectorFieldFlowRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	type Particle struct{ x, y float64 }

	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	start := time.Now()

	numParticles := 100
	particles := make([]Particle, numParticles)
	for i := range particles {
		particles[i] = Particle{
			x: rand.Float64() * w,
			y: rand.Float64() * h,
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			for i := range particles {
				p := &particles[i]

				// Curl-inspired sine/cos field
				angle := math.Sin(p.y*0.1+now)*math.Pi + math.Cos(p.x*0.1-now)*math.Pi
				dx := math.Cos(angle)
				dy := math.Sin(angle)

				p.x += dx * 0.5
				p.y += dy * 0.5

				// Wrap around
				if p.x < 0 {
					p.x += w
				} else if p.x >= w {
					p.x -= w
				}
				if p.y < 0 {
					p.y += h
				} else if p.y >= h {
					p.y -= h
				}

				hue := math.Mod(float64(i)/float64(numParticles)+now*0.1, 1.0)
				r, g, b := hsvToRGB(hue, 1.0, 1.0)
				dc.SetRGB(r, g, b)
				dc.SetPixel(int(p.x), int(p.y))
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type SierpinskiTriangleRenderer struct {
	screen *rgbmatrix.Screen
}

func SierpinskiTriangle(screen *rgbmatrix.Screen) *SierpinskiTriangleRenderer {
	return &SierpinskiTriangleRenderer{screen: screen}
}

func (ren *SierpinskiTriangleRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	var drawTriangle func(dc *gg.Context, x1, y1, x2, y2, x3, y3 float64, depth int)

	drawTriangle = func(dc *gg.Context, x1, y1, x2, y2, x3, y3 float64, depth int) {
		if depth <= 0 {
			dc.MoveTo(x1, y1)
			dc.LineTo(x2, y2)
			dc.LineTo(x3, y3)
			dc.ClosePath()
			dc.Fill()
			return
		}
		mid := func(x1, y1, x2, y2 float64) (float64, float64) {
			return (x1 + x2) / 2, (y1 + y2) / 2
		}
		ax, ay := mid(x1, y1, x2, y2)
		bx, by := mid(x2, y2, x3, y3)
		cx, cy := mid(x3, y3, x1, y1)

		drawTriangle(dc, x1, y1, ax, ay, cx, cy, depth-1)
		drawTriangle(dc, x2, y2, ax, ay, bx, by, depth-1)
		drawTriangle(dc, x3, y3, bx, by, cx, cy, depth-1)
	}

	dc := gg.NewContextForImage(ren.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			depth := int(2 + math.Floor(math.Sin(now*0.5)*1.5)) // animate between 1 and 3

			dc.SetRGB(0, 0, 0)
			dc.Clear()

			hue := math.Mod(now*0.1, 1.0)
			r, g, b := hsvToRGB(hue, 1.0, 1.0)
			dc.SetRGB(r, g, b)

			// Points of triangle
			x1, y1 := width/2, 0.0
			x2, y2 := 0.0, height
			x3, y3 := width, height

			drawTriangle(dc, x1, y1, x2, y2, x3, y3, depth)

			ren.screen.ShowImage(ctx, dc.Image())
			time.Sleep(200 * time.Millisecond)
		}
	}
}

type FluidDreamRenderer struct {
	screen *rgbmatrix.Screen
}

func FluidDream(screen *rgbmatrix.Screen) *FluidDreamRenderer {
	return &FluidDreamRenderer{screen: screen}
}

func (r *FluidDreamRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < height; y++ {
				for x := 0.0; x < width; x++ {
					nx := x / width
					ny := y / height

					// Smooth coordinate motion
					u := nx + 0.1*math.Sin(ny*10+t*0.3)
					v := ny + 0.1*math.Cos(nx*10-t*0.2)

					// Complex wave field
					value := math.Sin(u*8+v*4+t*0.7) +
						math.Cos(u*10-v*8-t*0.5) +
						math.Sin((u+v)*15-t*0.2)

					value = value / 3.0 // normalize to [-1, 1]

					hue := math.Mod((value+1)/2+t*0.1, 1.0)
					brightness := 0.4 + 0.6*math.Sin(value*2+t*0.8)

					r, g, b := hsvToRGB(hue, 1.0, brightness)
					minVal := 0.15
					r = math.Max(r, minVal)
					g = math.Max(g, minVal)
					b = math.Max(b, minVal)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

// FluidRainbowRenderer creates a smoothly evolving rainbow that flows across the screen like liquid ink.
type FluidRainbowRenderer struct {
	screen *rgbmatrix.Screen
}

func FluidRainbow(screen *rgbmatrix.Screen) *FluidRainbowRenderer {
	return &FluidRainbowRenderer{screen: screen}
}

func (r *FluidRainbowRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	width := float64(dc.Width())
	height := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < height; y++ {
				for x := 0.0; x < width; x++ {
					xf := x / width
					yf := y / height

					// Offset fields for fluid movement
					u := xf + 0.1*math.Sin(yf*10+t*0.4)
					v := yf + 0.1*math.Cos(xf*10-t*0.3)

					// Use sin+cos waves to distort hue across surface
					value := math.Sin(u*6+t) + math.Cos(v*8-t*1.2)
					hue := math.Mod((value+2)/4+t*0.05, 1.0)
					brightness := 0.7 + 0.3*math.Sin(value*2+t*0.8)

					r, g, b := hsvToRGB(hue, 1.0, brightness)
					r = math.Max(r, 0.15)
					g = math.Max(g, 0.15)
					b = math.Max(b, 0.15)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type OrbitingMetaballsRenderer struct {
	screen *rgbmatrix.Screen
}

func OrbitingMetaballs(screen *rgbmatrix.Screen) *OrbitingMetaballsRenderer {
	return &OrbitingMetaballsRenderer{screen: screen}
}

func (r *OrbitingMetaballsRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	type Blob struct {
		radius float64
		speed  float64
		offset float64
	}

	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	start := time.Now()

	centerX := w / 2
	centerY := h / 2

	blobs := []Blob{
		{radius: 12, speed: 1.0, offset: 0},
		{radius: 10, speed: 1.3, offset: math.Pi / 3},
		{radius: 9, speed: 0.8, offset: math.Pi * 2 / 3},
		{radius: 8, speed: 1.5, offset: math.Pi},
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			now := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					sum := 0.0
					for _, blob := range blobs {
						angle := now*blob.speed + blob.offset
						orbX := centerX + math.Cos(angle)*w*0.25
						orbY := centerY + math.Sin(angle)*h*0.25
						dx := x - orbX
						dy := y - orbY
						distSq := dx*dx + dy*dy
						sum += blob.radius * blob.radius / (distSq + 1)
					}
					val := math.Min(sum/4.0, 1.0)
					hue := math.Mod(0.65+0.3*val+now*0.05, 1.0)
					bright := math.Pow(val, 1.4)

					r, g, b := hsvToRGB(hue, 0.8, bright)
					r = math.Max(r, 0.15)
					g = math.Max(g, 0.15)
					b = math.Max(b, 0.15)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}
			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}

type MarbleShaderRenderer struct {
	screen *rgbmatrix.Screen
}

func MarbleShader(screen *rgbmatrix.Screen) *MarbleShaderRenderer {
	return &MarbleShaderRenderer{screen: screen}
}

func (r *MarbleShaderRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	w := float64(dc.Width())
	h := float64(dc.Height())
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			t := time.Since(start).Seconds()
			dc.Clear()

			for y := 0.0; y < h; y++ {
				for x := 0.0; x < w; x++ {
					xf := x / w
					yf := y / h

					noise := math.Sin((xf*10+math.Sin(yf*10+t*0.2))*3 + t)
					marble := math.Sin(xf*20 + noise*2 + t*0.5)
					value := (marble + 1) / 2

					hue := math.Mod(0.5+value*0.3+t*0.02, 1.0)
					brightness := 0.4 + 0.6*value

					r, g, b := hsvToRGB(hue, 0.7, brightness)
					r = math.Max(r, 0.15)
					g = math.Max(g, 0.15)
					b = math.Max(b, 0.15)

					dc.SetRGB(r, g, b)
					dc.SetPixel(int(x), int(y))
				}
			}

			r.screen.ShowImage(ctx, dc.Image())
			time.Sleep(30 * time.Millisecond)
		}
	}
}
