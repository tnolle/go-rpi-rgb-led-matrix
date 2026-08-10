package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/api"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/keycloak"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	commands := make(chan renderers.Command)
	frames := display.NewHub()
	config := rgbmatrix.LoadConfig()
	keycloak.Init(config.Auth.ClientID, config.Auth.ClientSecret)

	go api.ListenAndServe(commands, frames)
	renderers.UpdateLoopTerminal(ctx, commands, config, frames)
}
