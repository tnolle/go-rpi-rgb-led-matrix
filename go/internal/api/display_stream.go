package api

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
)

const (
	streamWriteTimeout = 5 * time.Second
	streamPongTimeout  = 60 * time.Second
	streamPingInterval = 30 * time.Second
)

var streamUpgrader = websocket.Upgrader{}

type streamMetadata struct {
	Type        string `json:"type"`
	Version     int    `json:"version"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	PixelFormat string `json:"pixel_format"`
}

func displayStreamHandler(hub *display.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if hub == nil {
			writeError(w, http.StatusServiceUnavailable, "display_stream_unavailable", "the display stream is unavailable")
			return
		}

		connection, err := streamUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer connection.Close()

		connection.SetReadLimit(1024)
		_ = connection.SetReadDeadline(time.Now().Add(streamPongTimeout))
		connection.SetPongHandler(func(string) error {
			return connection.SetReadDeadline(time.Now().Add(streamPongTimeout))
		})

		disconnected := make(chan struct{})
		go func() {
			defer close(disconnected)
			for {
				if _, _, err := connection.ReadMessage(); err != nil {
					return
				}
			}
		}()

		frames, unsubscribe := hub.Subscribe()
		defer unsubscribe()
		pings := time.NewTicker(streamPingInterval)
		defer pings.Stop()

		metadataSent := false
		for {
			select {
			case frame, ok := <-frames:
				if !ok {
					return
				}
				if !metadataSent {
					metadata := streamMetadata{
						Type: "display", Version: 1,
						Width: frame.Width, Height: frame.Height, PixelFormat: "rgb24",
					}
					if err := writeWebSocketJSON(connection, metadata); err != nil {
						return
					}
					metadataSent = true
				}
				if err := writeWebSocketMessage(connection, websocket.BinaryMessage, frame.Pixels); err != nil {
					return
				}
			case <-pings.C:
				if err := writeWebSocketMessage(connection, websocket.PingMessage, nil); err != nil {
					return
				}
			case <-disconnected:
				return
			case <-r.Context().Done():
				return
			}
		}
	}
}

func writeWebSocketJSON(connection *websocket.Conn, body any) error {
	if err := connection.SetWriteDeadline(time.Now().Add(streamWriteTimeout)); err != nil {
		return err
	}
	return connection.WriteJSON(body)
}

func writeWebSocketMessage(connection *websocket.Conn, messageType int, body []byte) error {
	if err := connection.SetWriteDeadline(time.Now().Add(streamWriteTimeout)); err != nil {
		return err
	}
	return connection.WriteMessage(messageType, body)
}
