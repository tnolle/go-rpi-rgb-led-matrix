package renderers

import (
	"context"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

const (
	pacmanFrameInterval = time.Second / 30
	pacmanRoundSeconds  = 18.0
	pacmanRunSeconds    = 14.0
	pacmanPowerSecond   = 1.4
	pacmanGhostSpeed    = 0.018
)

type PacmanRenderer struct {
	screen *rgbmatrix.Screen
}

type pacmanPoint struct {
	x, y      float64
	direction float64
}

type pacmanGhost struct {
	start            float64
	release          float64
	red, green, blue float64
}

var pacmanGhosts = []pacmanGhost{
	{start: 0.18, release: 0.0, red: 1.00, green: 0.12, blue: 0.18},
	{start: 0.30, release: 0.4, red: 1.00, green: 0.45, blue: 0.72},
	{start: 0.42, release: 0.8, red: 0.05, green: 0.90, blue: 0.95},
	{start: 0.54, release: 1.2, red: 1.00, green: 0.55, blue: 0.12},
}

var pacmanScoreDigits = map[byte][5]string{
	'0': {"111", "101", "101", "101", "111"},
	'1': {"010", "110", "010", "010", "111"},
	'2': {"111", "001", "111", "100", "111"},
	'4': {"101", "101", "111", "001", "001"},
	'6': {"111", "100", "111", "101", "111"},
	'8': {"111", "101", "111", "101", "111"},
}

func Pacman(screen *rgbmatrix.Screen) *PacmanRenderer {
	return &PacmanRenderer{screen: screen}
}

func (r *PacmanRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	dc := gg.NewContextForImage(r.screen.Canvas)
	start := time.Now()
	ticker := time.NewTicker(pacmanFrameInterval)
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

func (r *PacmanRenderer) drawFrame(ctx context.Context, dc *gg.Context, elapsed float64) error {
	width, height := float64(dc.Width()), float64(dc.Height())
	scale := math.Max(0.65, math.Min(width/128, height/64))
	roundSecond := math.Mod(elapsed, pacmanRoundSeconds)
	progress := math.Min(1, roundSecond/pacmanRunSeconds)

	dc.SetRGB(0.005, 0.005, 0.035)
	dc.Clear()
	r.drawMaze(dc, width, height, scale)
	r.drawPellets(dc, width, height, scale, progress)

	point := pacmanPath(progress, width, height)
	for index, ghost := range pacmanGhosts {
		r.drawLevelGhost(dc, ghost, index, scale, elapsed, roundSecond, width, height)
	}
	if roundSecond < pacmanRunSeconds {
		r.drawPacman(dc, point, scale, elapsed)
	} else {
		r.drawLevelComplete(dc, point, width, height, scale, roundSecond-pacmanRunSeconds)
	}
	r.drawCaptureScores(dc, scale, roundSecond, width, height)
	return r.screen.ShowImage(ctx, dc.Image())
}

func (r *PacmanRenderer) drawMaze(dc *gg.Context, width, height, scale float64) {
	insetX := math.Max(6, width*0.075)
	insetY := math.Max(7, height*0.16)
	radius := math.Max(3, 4*scale)

	dc.SetRGBA(0.05, 0.12, 0.95, 0.30)
	dc.SetLineWidth(math.Max(2.2, 3.5*scale))
	dc.DrawRoundedRectangle(insetX-3*scale, insetY-3*scale, width-2*insetX+6*scale, height-2*insetY+6*scale, radius)
	dc.Stroke()
	dc.SetRGB(0.08, 0.25, 1.0)
	dc.SetLineWidth(math.Max(0.8, 1.2*scale))
	dc.DrawRoundedRectangle(insetX-3*scale, insetY-3*scale, width-2*insetX+6*scale, height-2*insetY+6*scale, radius)
	dc.Stroke()

	// Interior blocks suggest the maze without obstructing the chase loop.
	blockWidth := width * 0.16
	blockHeight := math.Max(4, height*0.10)
	for _, x := range []float64{width * 0.24, width*0.5 - blockWidth/2, width*0.76 - blockWidth} {
		r.drawMazeBlock(dc, x, height*0.31, blockWidth, blockHeight, scale)
		r.drawMazeBlock(dc, x, height*0.59, blockWidth, blockHeight, scale)
	}

	// The center ghost house is the visual anchor of the board.
	houseWidth := width * 0.19
	houseHeight := height * 0.17
	dc.SetRGBA(0.08, 0.22, 1.0, 0.8)
	dc.SetLineWidth(math.Max(0.8, scale))
	dc.DrawRectangle(width/2-houseWidth/2, height/2-houseHeight/2, houseWidth, houseHeight)
	dc.Stroke()
	dc.SetRGB(1.0, 0.35, 0.75)
	dc.DrawLine(width/2-houseWidth*0.22, height/2-houseHeight/2, width/2+houseWidth*0.22, height/2-houseHeight/2)
	dc.Stroke()
}

func (r *PacmanRenderer) drawMazeBlock(dc *gg.Context, x, y, width, height, scale float64) {
	dc.SetRGBA(0.06, 0.18, 1.0, 0.78)
	dc.SetLineWidth(math.Max(0.7, scale))
	dc.DrawRoundedRectangle(x, y, width, height, math.Max(1, 2*scale))
	dc.Stroke()
}

func (r *PacmanRenderer) drawPellets(dc *gg.Context, width, height, scale, progress float64) {
	const pelletCount = 44
	for index := 0; index < pelletCount; index++ {
		pelletProgress := float64(index) / pelletCount
		// Pellets behind Pac-Man stay eaten until the next lap.
		if pelletProgress < progress && progress-pelletProgress < 0.93 {
			continue
		}
		point := pacmanPath(pelletProgress, width, height)
		powerPellet := index == 4 || index == 30
		radius := math.Max(0.55, 0.75*scale)
		if powerPellet {
			radius = math.Max(1.3, 1.8*scale)
		}
		dc.SetRGBA(1.0, 0.78, 0.55, 0.20)
		dc.DrawCircle(point.x, point.y, radius+1)
		dc.Fill()
		dc.SetRGB(1.0, 0.82, 0.62)
		dc.DrawCircle(point.x, point.y, radius)
		dc.Fill()
	}
}

func pacmanPath(progress, width, height float64) pacmanPoint {
	progress = math.Max(0, math.Min(1, progress))
	points := []pacmanPoint{
		{x: width * 0.08, y: height * 0.16},
		{x: width * 0.50, y: height * 0.16},
		{x: width * 0.50, y: height * 0.31},
		{x: width * 0.92, y: height * 0.31},
		{x: width * 0.92, y: height * 0.84},
		{x: width * 0.62, y: height * 0.84},
		{x: width * 0.62, y: height * 0.58},
		{x: width * 0.38, y: height * 0.58},
		{x: width * 0.38, y: height * 0.84},
		{x: width * 0.08, y: height * 0.84},
		{x: width * 0.08, y: height * 0.50},
		{x: width * 0.28, y: height * 0.50},
		{x: width * 0.28, y: height * 0.16},
		{x: width * 0.08, y: height * 0.16},
	}
	totalDistance := 0.0
	for index := 0; index < len(points)-1; index++ {
		totalDistance += math.Hypot(points[index+1].x-points[index].x, points[index+1].y-points[index].y)
	}
	distance := progress * totalDistance
	for index := 0; index < len(points)-1; index++ {
		start, end := points[index], points[index+1]
		dx, dy := end.x-start.x, end.y-start.y
		segmentDistance := math.Hypot(dx, dy)
		if distance <= segmentDistance {
			fraction := distance / segmentDistance
			return pacmanPoint{
				x: start.x + dx*fraction, y: start.y + dy*fraction,
				direction: math.Atan2(dy, dx),
			}
		}
		distance -= segmentDistance
	}
	last := points[len(points)-1]
	previous := points[len(points)-2]
	last.direction = math.Atan2(last.y-previous.y, last.x-previous.x)
	return last
}

func pacmanGhostCaptureSecond(ghost pacmanGhost) float64 {
	pacmanSpeed := 1 / pacmanRunSeconds
	return (ghost.start - ghost.release*pacmanGhostSpeed) / (pacmanSpeed - pacmanGhostSpeed)
}

func (r *PacmanRenderer) drawLevelGhost(dc *gg.Context, ghost pacmanGhost, _ int, scale, elapsed, roundSecond, width, height float64) {
	house := pacmanPoint{x: width / 2, y: height / 2, direction: 0}
	startPoint := pacmanPath(ghost.start, width, height)
	if roundSecond < ghost.release {
		releaseProgress := roundSecond / ghost.release
		point := interpolatePacmanPoint(house, startPoint, releaseProgress)
		point.direction = math.Atan2(startPoint.y-house.y, startPoint.x-house.x)
		r.drawGhost(dc, point, ghost, scale, elapsed, false)
		return
	}

	captureSecond := pacmanGhostCaptureSecond(ghost)
	if roundSecond < captureSecond {
		ghostProgress := ghost.start + (roundSecond-ghost.release)*pacmanGhostSpeed
		point := pacmanPath(ghostProgress, width, height)
		r.drawGhost(dc, point, ghost, scale, elapsed, roundSecond >= pacmanPowerSecond)
		return
	}

	capturePoint := pacmanPath(captureSecond/pacmanRunSeconds, width, height)
	returnProgress := math.Min(1, (roundSecond-captureSecond)/1.6)
	if returnProgress < 1 {
		point := interpolatePacmanPoint(capturePoint, house, returnProgress)
		point.direction = math.Atan2(house.y-point.y, house.x-point.x)
		r.drawGhostEyes(dc, point, scale)
	}
}

func interpolatePacmanPoint(start, end pacmanPoint, progress float64) pacmanPoint {
	return pacmanPoint{
		x: start.x + (end.x-start.x)*progress,
		y: start.y + (end.y-start.y)*progress,
	}
}

func (r *PacmanRenderer) drawCaptureScores(dc *gg.Context, scale, roundSecond, width, height float64) {
	scores := []string{"200", "400", "800", "1600"}
	for index, score := range scores {
		captureSecond := pacmanGhostCaptureSecond(pacmanGhosts[index])
		age := roundSecond - captureSecond
		if age < 0 || age > 0.85 {
			continue
		}
		point := pacmanPath(captureSecond/pacmanRunSeconds, width, height)
		pixelSize := max(1, int(math.Round(scale)))
		textWidth := len(score)*3*pixelSize + (len(score)-1)*pixelSize
		x := int(point.x) - textWidth/2
		y := int(point.y - 6*scale - age*3)
		brightness := 1 - age/0.85
		r.drawTinyNumber(dc, score, x, y, pixelSize, brightness)
	}
}

func (r *PacmanRenderer) drawTinyNumber(dc *gg.Context, value string, x, y, scale int, brightness float64) {
	for digitIndex := range len(value) {
		pattern := pacmanScoreDigits[value[digitIndex]]
		for row, line := range pattern {
			for column, pixel := range line {
				if pixel != '1' {
					continue
				}
				px := float64(x + digitIndex*4*scale + column*scale)
				py := float64(y + row*scale)
				dc.SetRGBA(0.25, 1.0, 1.0, 0.20*brightness)
				dc.DrawRectangle(px-1, py-1, float64(scale+2), float64(scale+2))
				dc.Fill()
				dc.SetRGBA(0.72, 1.0, 1.0, brightness)
				dc.DrawRectangle(px, py, float64(scale), float64(scale))
				dc.Fill()
			}
		}
	}
}

func (r *PacmanRenderer) drawPacman(dc *gg.Context, point pacmanPoint, scale, elapsed float64) {
	radius := math.Max(2.2, 4.2*scale)
	mouth := 0.12 + 0.32*math.Abs(math.Sin(elapsed*10))
	dc.SetRGBA(1.0, 0.82, 0.0, 0.22)
	dc.DrawCircle(point.x, point.y, radius+1.5)
	dc.Fill()
	dc.SetRGB(1.0, 0.84, 0.0)
	dc.MoveTo(point.x, point.y)
	dc.DrawArc(point.x, point.y, radius, point.direction+mouth, point.direction+2*math.Pi-mouth)
	dc.ClosePath()
	dc.Fill()
}

func (r *PacmanRenderer) drawLevelComplete(dc *gg.Context, point pacmanPoint, width, height, scale, elapsed float64) {
	r.drawPacman(dc, point, scale*(1+0.08*math.Sin(elapsed*7)), elapsed)
	centerX, centerY := width/2, height/2

	for ring := 1; ring <= 3; ring++ {
		radius := math.Mod(elapsed*18+float64(ring)*10, math.Min(width, height)*0.48)
		hue := math.Mod(elapsed*0.18+float64(ring)*0.22, 1)
		red, green, blue := hsvToRGB(hue, 0.9, 1)
		dc.SetRGBA(red, green, blue, 0.55*(1-radius/(math.Min(width, height)*0.5)))
		dc.SetLineWidth(math.Max(0.7, scale))
		dc.DrawCircle(centerX, centerY, radius)
		dc.Stroke()
	}

	for spark := 0; spark < 18; spark++ {
		angle := float64(spark)*2*math.Pi/18 + elapsed*0.7
		distance := 8*scale + math.Mod(elapsed*13+float64(spark*7), math.Min(width, height)*0.42)
		x := centerX + math.Cos(angle)*distance
		y := centerY + math.Sin(angle)*distance*0.55
		hue := math.Mod(float64(spark)/18+elapsed*0.12, 1)
		red, green, blue := hsvToRGB(hue, 0.9, 1)
		dc.SetRGB(red, green, blue)
		dc.DrawCircle(x, y, math.Max(0.5, scale*0.7))
		dc.Fill()
	}

	pixelSize := max(1, int(math.Round(scale)))
	bonus := "10000"
	textWidth := len(bonus)*3*pixelSize + (len(bonus)-1)*pixelSize
	r.drawTinyNumber(dc, bonus, int(centerX)-textWidth/2, int(centerY)-3*pixelSize, pixelSize, 0.75+0.25*math.Sin(elapsed*6))
}

func (r *PacmanRenderer) drawGhost(dc *gg.Context, point pacmanPoint, ghost pacmanGhost, scale, elapsed float64, frightened bool) {
	radius := math.Max(1.9, 3.5*scale)
	red, green, blue := ghost.red, ghost.green, ghost.blue
	if frightened {
		red, green, blue = 0.08, 0.20, 0.95
		if math.Mod(elapsed, 0.45) > 0.32 {
			red, green, blue = 0.92, 0.92, 1.0
		}
	}

	dc.SetRGBA(red, green, blue, 0.20)
	dc.DrawCircle(point.x, point.y, radius+1)
	dc.Fill()
	dc.SetRGB(red, green, blue)
	dc.DrawCircle(point.x, point.y-radius*0.20, radius)
	dc.Fill()
	dc.DrawRectangle(point.x-radius, point.y-radius*0.15, radius*2, radius*1.15)
	dc.Fill()

	// Alternating feet preserve the familiar two-frame arcade gait.
	feetPhase := int(elapsed*9) % 2
	for foot := -1; foot <= 1; foot++ {
		offset := float64(foot) * radius * 0.68
		bob := 0.0
		if (foot+feetPhase)%2 == 0 {
			bob = radius * 0.18
		}
		dc.DrawCircle(point.x+offset, point.y+radius*0.88+bob, radius*0.38)
		dc.Fill()
	}

	eyeOffsetX := math.Cos(point.direction) * radius * 0.20
	eyeOffsetY := math.Sin(point.direction) * radius * 0.20
	for _, side := range []float64{-0.38, 0.38} {
		eyeX := point.x + side*radius
		eyeY := point.y - radius*0.25
		dc.SetRGB(1, 1, 1)
		dc.DrawCircle(eyeX, eyeY, radius*0.28)
		dc.Fill()
		dc.SetRGB(0.03, 0.12, 0.55)
		dc.DrawCircle(eyeX+eyeOffsetX, eyeY+eyeOffsetY, radius*0.13)
		dc.Fill()
	}
}

func (r *PacmanRenderer) drawGhostEyes(dc *gg.Context, point pacmanPoint, scale float64) {
	radius := math.Max(1.9, 3.5*scale)
	eyeOffsetX := math.Cos(point.direction) * radius * 0.20
	eyeOffsetY := math.Sin(point.direction) * radius * 0.20
	for _, side := range []float64{-0.38, 0.38} {
		eyeX := point.x + side*radius
		eyeY := point.y - radius*0.25
		dc.SetRGBA(0.35, 0.75, 1.0, 0.20)
		dc.DrawCircle(eyeX, eyeY, radius*0.42)
		dc.Fill()
		dc.SetRGB(1, 1, 1)
		dc.DrawCircle(eyeX, eyeY, radius*0.28)
		dc.Fill()
		dc.SetRGB(0.03, 0.12, 0.55)
		dc.DrawCircle(eyeX+eyeOffsetX, eyeY+eyeOffsetY, radius*0.13)
		dc.Fill()
	}
}
