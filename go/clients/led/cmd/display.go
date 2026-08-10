package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

const displayFrameInterval = time.Second / 30

var displayHost string

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Display a server's LED matrix in this terminal",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		host, err := selectedDisplayHost(displayHost)
		if err != nil {
			return err
		}
		return displayServer(host)
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
	displayCmd.Flags().StringVar(&displayHost, "host", "", "server URL (defaults to the selected host)")
}

type displayMetadata struct {
	Type        string `json:"type"`
	Version     int    `json:"version"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	PixelFormat string `json:"pixel_format"`
}

func selectedDisplayHost(override string) (string, error) {
	if override = strings.TrimSpace(override); override != "" {
		return strings.TrimRight(override, "/"), nil
	}

	configured := hosts()
	if len(configured) != 1 {
		return "", fmt.Errorf("display requires one server; select one with `led hosts select <index>` or pass --host")
	}
	return configured[0], nil
}

func displayStreamURL(host string) (string, error) {
	if !strings.Contains(host, "://") {
		host = "http://" + host
	}
	u, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("parse server URL: %w", err)
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported server URL scheme %q", u.Scheme)
	}
	if u.Host == "" {
		return "", errors.New("server URL has no host")
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/display/stream"
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func displayServer(host string) error {
	streamURL, err := displayStreamURL(host)
	if err != nil {
		return err
	}

	connection, _, err := websocket.DefaultDialer.Dial(streamURL, nil)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", streamURL, err)
	}
	defer connection.Close()

	if err := connection.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("set metadata deadline: %w", err)
	}
	metadata, err := readDisplayMetadata(connection)
	if err != nil {
		return err
	}
	if err := connection.SetReadDeadline(time.Time{}); err != nil {
		return fmt.Errorf("clear metadata deadline: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	frames := make(chan []byte, 1)
	streamErrors := make(chan error, 1)
	go receiveDisplayFrames(connection, metadata, frames, streamErrors)

	matrix := rgbmatrix.NewTerminal(metadata.Width, metadata.Height)
	defer matrix.Close()

	ticker := time.NewTicker(displayFrameInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-streamErrors:
			return err
		case <-ticker.C:
			select {
			case frame := <-frames:
				if err := matrix.Apply(rgb24Colors(frame)); err != nil {
					return err
				}
			default:
			}
		}
	}
}

func readDisplayMetadata(connection *websocket.Conn) (displayMetadata, error) {
	messageType, message, err := connection.ReadMessage()
	if err != nil {
		return displayMetadata{}, fmt.Errorf("read display metadata: %w", err)
	}
	if messageType != websocket.TextMessage {
		return displayMetadata{}, fmt.Errorf("expected display metadata, received WebSocket message type %d", messageType)
	}

	var metadata displayMetadata
	if err := json.Unmarshal(message, &metadata); err != nil {
		return displayMetadata{}, fmt.Errorf("decode display metadata: %w", err)
	}
	if metadata.Type != "display" || metadata.Version != 1 {
		return displayMetadata{}, fmt.Errorf("unsupported display stream type %q version %d", metadata.Type, metadata.Version)
	}
	if metadata.PixelFormat != "rgb24" {
		return displayMetadata{}, fmt.Errorf("unsupported pixel format %q", metadata.PixelFormat)
	}
	if metadata.Width <= 0 || metadata.Height <= 0 {
		return displayMetadata{}, fmt.Errorf("invalid display geometry %dx%d", metadata.Width, metadata.Height)
	}
	return metadata, nil
}

func receiveDisplayFrames(connection *websocket.Conn, metadata displayMetadata, frames chan []byte, streamErrors chan error) {
	expectedLength := metadata.Width * metadata.Height * 3
	for {
		messageType, frame, err := connection.ReadMessage()
		if err != nil {
			streamErrors <- fmt.Errorf("read display stream: %w", err)
			return
		}
		if messageType != websocket.BinaryMessage {
			continue
		}
		if len(frame) != expectedLength {
			streamErrors <- fmt.Errorf("invalid display frame length: got %d, want %d", len(frame), expectedLength)
			return
		}

		select {
		case frames <- frame:
		default:
			select {
			case <-frames:
			default:
			}
			select {
			case frames <- frame:
			default:
			}
		}
	}
}

func rgb24Colors(frame []byte) []color.Color {
	pixels := make([]color.Color, len(frame)/3)
	for position := range pixels {
		offset := position * 3
		pixels[position] = color.RGBA{
			R: frame[offset], G: frame[offset+1], B: frame[offset+2], A: 255,
		}
	}
	return pixels
}
