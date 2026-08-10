package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

func TestDisplayStreamSendsMetadataThenLatestFrame(t *testing.T) {
	hub := display.NewHub()
	pixels := []byte{1, 2, 3, 4, 5, 6}
	hub.Publish(display.Frame{Width: 2, Height: 1, Pixels: pixels})
	server := httptest.NewServer(newHandler(make(chan renderers.Command), catalog{}, hub))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/display/stream"
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	_ = connection.SetReadDeadline(time.Now().Add(time.Second))

	messageType, message, err := connection.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if messageType != websocket.TextMessage {
		t.Fatalf("metadata message type: got %d, want text", messageType)
	}
	var metadata streamMetadata
	if err := json.Unmarshal(message, &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata.Type != "display" || metadata.Version != 1 || metadata.Width != 2 || metadata.Height != 1 || metadata.PixelFormat != "rgb24" {
		t.Fatalf("unexpected metadata: %#v", metadata)
	}

	messageType, message, err = connection.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if messageType != websocket.BinaryMessage || string(message) != string(pixels) {
		t.Fatalf("unexpected frame: type=%d pixels=%v", messageType, message)
	}
}
