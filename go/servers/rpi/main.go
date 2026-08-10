//go:build linux && with_cgo

package main

import (
	"context"
	"log"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/api"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/keycloak"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func main() {
	// Main context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Communication channel
	commands := make(chan renderers.Command)
	frames := display.NewHub()

	// Config
	config := rgbmatrix.LoadConfig()
	width := config.Options.Cols * config.Options.ChainLength
	height := config.Options.Rows * config.Options.Parallel
	log.Printf("server starting: matrix=%dx%d output=raspberry-pi", width, height)

	// Keycloak
	keycloak.Init(config.Auth.ClientID, config.Auth.ClientSecret)

	// Start REST API and connect
	go api.ListenAndServe(commands, frames)

	// Run the update loop
	renderers.UpdateLoop(ctx, commands, config, frames)
}
