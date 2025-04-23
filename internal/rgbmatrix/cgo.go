//go:build with_cgo

package rgbmatrix

/*
#cgo CFLAGS: -std=c99 -I${SRCDIR}/../../3rdparty/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE
#cgo LDFLAGS: -lrgbmatrix -L${SRCDIR}/../../3rdparty/rpi-rgb-led-matrix/lib -lstdc++ -lm
*/
import "C"
