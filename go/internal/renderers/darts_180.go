package renderers

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

const darts180FrameInterval = time.Second / 30

type Darts180Renderer struct {
	screen   *rgbmatrix.Screen
	confetti []darts180Confetti
}

type darts180Confetti struct {
	x, y  float64
	speed float64
	phase float64
	hue   float64
}

var darts180Digits = map[byte][7]string{
	'0': {"11111", "10001", "10011", "10101", "11001", "10001", "11111"},
	'1': {"00100", "01100", "00100", "00100", "00100", "00100", "01110"},
	'8': {"11111", "10001", "10001", "11111", "10001", "10001", "11111"},
}

func Darts180(screen *rgbmatrix.Screen) *Darts180Renderer {
	width, height := screen.Canvas.Bounds().Dx(), screen.Canvas.Bounds().Dy()
	rng := rand.New(rand.NewSource(180))
	confettiCount := max(18, width*height/240)
	confetti := make([]darts180Confetti, confettiCount)
	for i := range confetti {
		confetti[i] = darts180Confetti{
			x:     rng.Float64() * float64(width),
			y:     rng.Float64() * float64(height),
			speed: 7 + rng.Float64()*13,
			phase: rng.Float64() * 2 * math.Pi,
			hue:   rng.Float64(),
		}
	}
	return &Darts180Renderer{screen: screen, confetti: confetti}
}

func (r *Darts180Renderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	start := time.Now()
	ticker := time.NewTicker(darts180FrameInterval)
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

func (r *Darts180Renderer) drawFrame(ctx context.Context, dc *gg.Context, elapsed float64) error {
	width, height := float64(dc.Width()), float64(dc.Height())
	dc.SetRGB(0.015, 0.005, 0.055)
	dc.Clear()

	r.drawNeonBackdrop(dc, width, height, elapsed)
	r.drawConfetti(dc, width, height, elapsed)
	r.drawMarquee(dc, width, height, elapsed)
	r.drawDarts(dc, width, height, elapsed)
	r.drawScore(dc, width, height, elapsed)

	return r.screen.ShowImage(ctx, dc.Image())
}

func (r *Darts180Renderer) drawNeonBackdrop(dc *gg.Context, width, height, elapsed float64) {
	vanishingX := width / 2
	vanishingY := height * 0.38

	// Horizontal bands and converging lines evoke the polished lanes and
	// geometric carpet of a classic bowling alley.
	for band := 0; band < 5; band++ {
		y := vanishingY + float64(band*band)*height*0.025
		hue := math.Mod(elapsed*0.08+float64(band)*0.12+0.55, 1)
		red, green, blue := hsvToRGB(hue, 0.85, 0.65)
		dc.SetRGBA(red, green, blue, 0.32)
		dc.SetLineWidth(0.6)
		dc.DrawLine(2, y, width-3, y)
		dc.Stroke()
	}
	for lane := -4; lane <= 4; lane++ {
		hue := math.Mod(0.82+float64(lane)*0.035+elapsed*0.04, 1)
		red, green, blue := hsvToRGB(hue, 0.9, 0.8)
		dc.SetRGBA(red, green, blue, 0.30)
		dc.DrawLine(vanishingX+float64(lane)*2, vanishingY, vanishingX+float64(lane)*width*0.15, height)
		dc.Stroke()
	}

	// A sweeping neon wash travels behind the score.
	sweepX := math.Mod(elapsed*35, width+40) - 20
	dc.SetRGBA(0.10, 0.85, 1.0, 0.10)
	dc.DrawRectangle(sweepX-8, 2, 16, height-4)
	dc.Fill()
}

func (r *Darts180Renderer) drawConfetti(dc *gg.Context, width, height, elapsed float64) {
	for _, piece := range r.confetti {
		x := math.Mod(piece.x+math.Sin(elapsed*1.7+piece.phase)*4+width, width)
		y := math.Mod(piece.y+elapsed*piece.speed, height)
		brightness := 0.65 + 0.35*math.Sin(elapsed*4+piece.phase)
		red, green, blue := hsvToRGB(math.Mod(piece.hue+elapsed*0.04, 1), 0.9, brightness)
		dc.SetRGBA(red, green, blue, 0.85)
		dc.SetLineWidth(1)
		dc.DrawLine(x, y, x+math.Sin(piece.phase)*2, y+1)
		dc.Stroke()
	}
}

func (r *Darts180Renderer) drawMarquee(dc *gg.Context, width, height, elapsed float64) {
	spacing := math.Max(6, math.Min(width, height)/8)
	bulb := func(x, y, index float64) {
		chase := math.Mod(index-elapsed*12, 6)
		brightness := 0.28
		if chase < 1.7 {
			brightness = 1
		}
		dc.SetRGBA(1.0, 0.42+0.45*brightness, 0.05, 0.18+0.32*brightness)
		dc.DrawCircle(x, y, 1.8)
		dc.Fill()
		dc.SetRGB(1.0, 0.38+0.55*brightness, 0.08)
		dc.DrawCircle(x, y, 0.65+0.35*brightness)
		dc.Fill()
	}

	index := 0.0
	for x := spacing; x < width-spacing/2; x += spacing {
		bulb(x, 2, index)
		index++
	}
	for y := spacing; y < height-spacing/2; y += spacing {
		bulb(width-2, y, index)
		index++
	}
	for x := width - spacing; x > spacing/2; x -= spacing {
		bulb(x, height-2, index)
		index++
	}
	for y := height - spacing; y > spacing/2; y -= spacing {
		bulb(2, y, index)
		index++
	}
}

func (r *Darts180Renderer) drawDarts(dc *gg.Context, width, height, elapsed float64) {
	centerX := width / 2
	baseY := math.Max(7, height*0.14)
	for index := 0; index < 3; index++ {
		arrival := math.Mod(elapsed*0.9+float64(index)*0.72, 2.4)
		x := centerX + float64(index-1)*math.Max(7, width*0.075)
		y := baseY - math.Max(0, 1-arrival)*7
		angle := -0.35 + float64(index)*0.35
		r.drawDart(dc, x, y, angle, math.Min(width/128, height/64))
	}
}

func (r *Darts180Renderer) drawDart(dc *gg.Context, x, y, angle, scale float64) {
	scale = math.Max(0.7, scale)
	dx, dy := math.Cos(angle), math.Sin(angle)
	length := 8 * scale
	dc.SetRGBA(0.05, 0.95, 1.0, 0.35)
	dc.SetLineWidth(2.2 * scale)
	dc.DrawLine(x-dx*length, y-dy*length, x, y)
	dc.Stroke()
	dc.SetRGB(0.82, 0.96, 1.0)
	dc.SetLineWidth(math.Max(0.7, scale))
	dc.DrawLine(x-dx*length, y-dy*length, x, y)
	dc.Stroke()

	// Pink flights at the rear give the darts the neon house-ball palette.
	rearX, rearY := x-dx*length, y-dy*length
	perpX, perpY := -dy*2.2*scale, dx*2.2*scale
	dc.SetRGBA(1.0, 0.08, 0.62, 0.9)
	dc.MoveTo(rearX, rearY)
	dc.LineTo(rearX-dx*3*scale+perpX, rearY-dy*3*scale+perpY)
	dc.LineTo(rearX-dx*3*scale-perpX, rearY-dy*3*scale-perpY)
	dc.ClosePath()
	dc.Fill()
}

func (r *Darts180Renderer) drawScore(dc *gg.Context, width, height, elapsed float64) {
	scale := int(math.Max(2, math.Min(width/24, height/10)))
	digitWidth := 5 * scale
	gap := scale
	totalWidth := 3*digitWidth + 2*gap
	startX := int(width)/2 - totalWidth/2
	startY := int(height)/2 - (7*scale)/2 + max(2, int(height*0.05))
	pulse := 0.82 + 0.18*math.Sin(elapsed*5)
	hue := math.Mod(0.88+elapsed*0.035, 1)
	red, green, blue := hsvToRGB(hue, 0.82, pulse)

	for index, digit := range []byte("180") {
		x := startX + index*(digitWidth+gap)
		r.drawDigit(dc, digit, x, startY, scale, red, green, blue)
	}
}

func (r *Darts180Renderer) drawDigit(dc *gg.Context, digit byte, x, y, scale int, red, green, blue float64) {
	pattern := darts180Digits[digit]
	for row, line := range pattern {
		for column, pixel := range line {
			if pixel != '1' {
				continue
			}
			px := float64(x + column*scale)
			py := float64(y + row*scale)
			dc.SetRGBA(red, green, blue, 0.18)
			dc.DrawRectangle(px-1, py-1, float64(scale+2), float64(scale+2))
			dc.Fill()
			dc.SetRGB(red, green, blue)
			dc.DrawRectangle(px, py, float64(scale), float64(scale))
			dc.Fill()
			dc.SetRGBA(1, 1, 1, 0.35)
			dc.DrawRectangle(px, py, math.Max(1, float64(scale)/3), math.Max(1, float64(scale)/3))
			dc.Fill()
		}
	}
}
