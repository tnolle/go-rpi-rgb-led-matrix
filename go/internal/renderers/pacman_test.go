package renderers

import (
	"context"
	"image/color"
	"math"
	"testing"

	"github.com/fogleman/gg"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func TestPacmanDrawsAFrame(t *testing.T) {
	matrix := newPrepareMatrix(128, 64)
	for position := range matrix.pixels {
		matrix.pixels[position] = color.RGBA{}
	}
	screen := rgbmatrix.NewScreen(matrix)
	renderer := Pacman(screen)
	dc := gg.NewContextForImage(screen.Canvas)

	if err := renderer.drawFrame(context.Background(), dc, 2.5); err != nil {
		t.Fatal(err)
	}

	brightPixels := 0
	image := dc.Image()
	for y := 0; y < image.Bounds().Dy(); y++ {
		for x := 0; x < image.Bounds().Dx(); x++ {
			rgba := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
			if int(rgba.R)+int(rgba.G)+int(rgba.B) > 160 {
				brightPixels++
			}
		}
	}
	if brightPixels < 100 {
		t.Fatalf("pacman frame has too few bright pixels: %d", brightPixels)
	}
}

func TestPacmanPathLoopsAroundMatrix(t *testing.T) {
	directions := make(map[float64]bool)
	for step := 0; step < 40; step++ {
		progress := float64(step) / 40
		point := pacmanPath(progress, 128, 64)
		if point.x < 0 || point.x >= 128 || point.y < 0 || point.y >= 64 {
			t.Fatalf("path point is outside matrix: %#v", point)
		}
		directions[point.direction] = true
		if point.direction != 0 && point.direction != math.Pi && point.direction != math.Pi/2 && point.direction != -math.Pi/2 {
			t.Fatalf("path direction is not orthogonal: %.2f", point.direction)
		}
	}
	if len(directions) != 4 {
		t.Fatalf("maze path uses %d directions, want 4", len(directions))
	}
}

func TestPacmanCatchesGhostsWithoutPassingThem(t *testing.T) {
	previousCapture := 0.0
	for index, ghost := range pacmanGhosts {
		captureSecond := pacmanGhostCaptureSecond(ghost)
		if captureSecond <= pacmanPowerSecond {
			t.Fatalf("ghost %d is caught before the power pellet: %.2f", index, captureSecond)
		}
		if captureSecond <= previousCapture || captureSecond >= pacmanRunSeconds {
			t.Fatalf("ghost %d has invalid capture time: %.2f", index, captureSecond)
		}

		justBefore := captureSecond - 0.01
		pacmanProgress := justBefore / pacmanRunSeconds
		ghostProgress := ghost.start + (justBefore-ghost.release)*pacmanGhostSpeed
		if ghostProgress <= pacmanProgress {
			t.Fatalf("ghost %d moved through Pac-Man before capture", index)
		}

		atCapturePacman := captureSecond / pacmanRunSeconds
		atCaptureGhost := ghost.start + (captureSecond-ghost.release)*pacmanGhostSpeed
		if math.Abs(atCapturePacman-atCaptureGhost) > 0.000001 {
			t.Fatalf("ghost %d capture positions differ: pacman=%f ghost=%f", index, atCapturePacman, atCaptureGhost)
		}
		previousCapture = captureSecond
	}
}

func TestPacmanLevelEndsWhereTheNextOneStarts(t *testing.T) {
	start := pacmanPath(0, 128, 64)
	finish := pacmanPath(1, 128, 64)
	if start.x != finish.x || start.y != finish.y {
		t.Fatalf("level path jumps when it resets: start=%#v finish=%#v", start, finish)
	}
}
