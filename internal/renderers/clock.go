package renderers

import (
	"context"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

type ClockRenderer struct {
	screen *rgbmatrix.Screen
}

func Clock(screen *rgbmatrix.Screen) *ClockRenderer {
	return &ClockRenderer{screen: screen}
}

func (r *ClockRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			dc.SetRGB(0, 0, 0)
			dc.Clear()

			// Current time
			t := time.Now()

			// Outer circle
			dc.SetRGB(.5, .5, .5)
			dc.DrawCircle(32, 32, 30)
			for a := range 12 {
				angle := float64(a) / 12 * 360
				p1 := pointOnCircle(gg.Point{X: 32, Y: 32}, 24, angle-90)
				p2 := pointOnCircle(gg.Point{X: 32, Y: 32}, 30, angle-90)
				dc.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
			}
			dc.Stroke()

			seconds := (float64(t.Nanosecond()) / 1e9)
			seconds = (float64(t.Second()) + seconds) / 60
			minutes := float64(t.Minute()) / 60
			hours := (float64(t.Hour()%12) + minutes) / 12

			// Hour hand
			dc.SetRGB(1, 1, 1)
			p := pointOnCircle(gg.Point{X: 32, Y: 32}, 16, hours*360-90)
			dc.DrawLine(32, 32, p.X, p.Y)
			dc.Stroke()

			// Minute hand
			dc.SetRGB(1, 1, 1)
			p = pointOnCircle(gg.Point{X: 32, Y: 32}, 24, minutes*360-90)
			dc.DrawLine(32, 32, p.X, p.Y)
			dc.Stroke()

			// Second hand
			dc.SetRGB(1, 0, 0)
			p = pointOnCircle(gg.Point{X: 32, Y: 32}, 26, seconds*360-90)
			dc.DrawLine(32, 32, p.X, p.Y)
			dc.Stroke()

			r.screen.ShowImage(ctx, dc.Image())
		}
	}

	return nil
}

func pointOnCircle(center gg.Point, radius, angle float64) gg.Point {
	x := center.X + radius*math.Cos(angle*math.Pi/180)
	y := center.Y + radius*math.Sin(angle*math.Pi/180)
	return gg.Point{X: x, Y: y}
}
