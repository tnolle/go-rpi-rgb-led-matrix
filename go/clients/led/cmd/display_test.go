package cmd

import (
	"image/color"
	"testing"
	"time"
)

func TestDisplayStreamURL(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{host: "http://localhost:8085", want: "ws://localhost:8085/display/stream"},
		{host: "https://led.example", want: "wss://led.example/display/stream"},
		{host: "ws://localhost:8085/", want: "ws://localhost:8085/display/stream"},
		{host: "localhost:8085", want: "ws://localhost:8085/display/stream"},
	}

	for _, test := range tests {
		t.Run(test.host, func(t *testing.T) {
			got, err := displayStreamURL(test.host)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("unexpected URL: got %q, want %q", got, test.want)
			}
		})
	}
}

func TestRGB24Colors(t *testing.T) {
	pixels := rgb24Colors([]byte{1, 2, 3, 4, 5, 6})
	want := []color.Color{
		color.RGBA{R: 1, G: 2, B: 3, A: 255},
		color.RGBA{R: 4, G: 5, B: 6, A: 255},
	}
	if len(pixels) != len(want) {
		t.Fatalf("unexpected pixel count: got %d, want %d", len(pixels), len(want))
	}
	for position := range want {
		if pixels[position] != want[position] {
			t.Fatalf("pixel %d: got %v, want %v", position, pixels[position], want[position])
		}
	}
}

func TestDisplayFrameIntervalIsThirtyFPS(t *testing.T) {
	if fps := int(time.Second / displayFrameInterval); fps != 30 {
		t.Fatalf("unexpected frame rate: %d", fps)
	}
}
