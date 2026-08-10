package main

import (
	"context"
	"flag"
	"log"
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
	showDisplay := flag.Bool("display", false, "render the LED matrix in this terminal")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	commands := make(chan renderers.Command)
	frames := display.NewHub()
	config := rgbmatrix.LoadConfig()
	keycloak.Init(config.Auth.ClientID, config.Auth.ClientSecret)
	width := config.Options.Cols * config.Options.ChainLength
	height := config.Options.Rows * config.Options.Parallel

	var matrix rgbmatrix.Matrix = rgbmatrix.NewMemory(width, height)
	if *showDisplay {
		log.Printf("server starting: matrix=%dx%d output=terminal", width, height)
		matrix = rgbmatrix.NewTerminal(width, height)
	} else {
		log.Printf("server starting: matrix=%dx%d output=headless", width, height)
	}

	go api.ListenAndServe(commands, frames)
	renderers.UpdateLoopWithMatrix(ctx, commands, matrix, frames)
}
