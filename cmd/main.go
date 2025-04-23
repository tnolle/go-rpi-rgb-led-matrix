package main

import (
	"context"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/api"
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

	// Config
	config := rgbmatrix.LoadConfig()

	// Keycloak
	keycloak.Init(config.Auth.ClientID, config.Auth.ClientSecret)

	// Start REST API and connect
	go api.ListenAndServe(commands)

	// Run the update loop
	renderers.UpdateLoop(ctx, commands, config)
}
