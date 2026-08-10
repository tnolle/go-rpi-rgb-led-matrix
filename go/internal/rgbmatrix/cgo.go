//go:build with_cgo

package rgbmatrix

/*
#cgo CFLAGS: -std=c99 -I${SRCDIR}/../../../third_party/rpi-rgb-led-matrix/include -DSHOW_REFRESH_RATE
#cgo LDFLAGS: -lrgbmatrix -L${SRCDIR}/../../../third_party/rpi-rgb-led-matrix/lib -lstdc++ -lm
*/
import "C"
