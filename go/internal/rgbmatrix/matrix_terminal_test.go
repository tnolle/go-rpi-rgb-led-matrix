package rgbmatrix

import (
	"bytes"
	"image/color"
	"strings"
	"testing"
)

func TestTerminalRender(t *testing.T) {
	var output bytes.Buffer
	matrix := NewTerminalWithWriter(1, 2, &output)
	matrix.Set(0, color.RGBA{R: 255, A: 255})
	matrix.Set(1, color.RGBA{B: 255, A: 255})

	if err := matrix.Render(); err != nil {
		t.Fatal(err)
	}

	frame := output.String()
	if !strings.Contains(frame, "\x1b[38;2;255;0;0m") {
		t.Fatalf("frame does not contain upper pixel color: %q", frame)
	}
	if !strings.Contains(frame, "\x1b[48;2;0;0;255m▀") {
		t.Fatalf("frame does not contain lower pixel color: %q", frame)
	}
	if got := matrix.At(0); got != (color.RGBA{}) {
		t.Fatalf("render did not clear framebuffer: %v", got)
	}

	if err := matrix.Close(); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(output.String(), terminalLeave) {
		t.Fatal("close did not restore terminal state")
	}
}

func TestTerminalGeometry(t *testing.T) {
	matrix := NewTerminalWithWriter(128, 64, &bytes.Buffer{})
	width, height := matrix.Geometry()
	if width != 128 || height != 64 {
		t.Fatalf("unexpected geometry: %dx%d", width, height)
	}
}
