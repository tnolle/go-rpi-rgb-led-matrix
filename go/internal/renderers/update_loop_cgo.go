//go:build with_cgo

package renderers

import (
	"context"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
)

func UpdateLoop(ctx context.Context, commands chan Command, config rgbmatrix.Config) {
	m, err := rgbmatrix.NewRGBLedMatrix(&config.Options, &config.RuntimeOptions)
	if err != nil {
		panic(err)
	}
	updateLoop(ctx, commands, m)
}
