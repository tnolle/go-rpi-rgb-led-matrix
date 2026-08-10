package display

import "testing"

func TestHubGivesNewSubscribersTheLatestFrame(t *testing.T) {
	hub := NewHub()
	hub.Publish(Frame{Width: 1, Height: 1, Pixels: []byte{1, 2, 3}})

	frames, unsubscribe := hub.Subscribe()
	defer unsubscribe()
	frame := <-frames
	if frame.Width != 1 || frame.Height != 1 || string(frame.Pixels) != string([]byte{1, 2, 3}) {
		t.Fatalf("unexpected frame: %#v", frame)
	}
}

func TestHubReplacesAStalePendingFrame(t *testing.T) {
	hub := NewHub()
	frames, unsubscribe := hub.Subscribe()
	defer unsubscribe()

	hub.Publish(Frame{Width: 1, Height: 1, Pixels: []byte{1, 1, 1}})
	hub.Publish(Frame{Width: 1, Height: 1, Pixels: []byte{2, 2, 2}})

	frame := <-frames
	if string(frame.Pixels) != string([]byte{2, 2, 2}) {
		t.Fatalf("subscriber received stale frame: %v", frame.Pixels)
	}
}

func TestHubOwnsPublishedPixels(t *testing.T) {
	hub := NewHub()
	pixels := []byte{1, 2, 3}
	hub.Publish(Frame{Width: 1, Height: 1, Pixels: pixels})
	pixels[0] = 9

	frames, unsubscribe := hub.Subscribe()
	defer unsubscribe()
	if frame := <-frames; frame.Pixels[0] != 1 {
		t.Fatalf("published frame was mutated: %v", frame.Pixels)
	}
}
