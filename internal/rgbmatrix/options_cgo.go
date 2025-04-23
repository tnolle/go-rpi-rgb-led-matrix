//go:build with_cgo

package rgbmatrix

/*
#include <led-matrix-c.h>

void set_show_refresh_rate(struct RGBLedMatrixOptions *o, int show_refresh_rate) {
  o->show_refresh_rate = show_refresh_rate != 0 ? 1 : 0;
}

void set_disable_hardware_pulsing(struct RGBLedMatrixOptions *o, int disable_hardware_pulsing) {
  o->disable_hardware_pulsing = disable_hardware_pulsing != 0 ? 1 : 0;
}

void set_inverse_colors(struct RGBLedMatrixOptions *o, int inverse_colors) {
  o->inverse_colors = inverse_colors != 0 ? 1 : 0;
}
*/
import "C"

func (r *RuntimeOptions) toC() *C.struct_RGBLedRuntimeOptions {
	return &C.struct_RGBLedRuntimeOptions{
		gpio_slowdown: C.int(r.GPIOSlowdown),
	}
}

func (c *MatrixOptions) toC() *C.struct_RGBLedMatrixOptions {
	o := &C.struct_RGBLedMatrixOptions{
		rows:                  C.int(c.Rows),
		cols:                  C.int(c.Cols),
		chain_length:          C.int(c.ChainLength),
		parallel:              C.int(c.Parallel),
		pwm_bits:              C.int(c.PWMBits),
		pwm_lsb_nanoseconds:   C.int(c.PWMLSBNanoseconds),
		pwm_dither_bits:       C.int(c.PWMDitherBits),
		brightness:            C.int(c.Brightness),
		scan_mode:             C.int(c.ScanMode),
		row_address_type:      C.int(c.RowAddressType),
		multiplexing:          C.int(c.Multiplexing),
		limit_refresh_rate_hz: C.int(c.LimitRefreshRateHz),
		hardware_mapping:      C.CString(c.HardwareMapping),
	}

	if c.ShowRefreshRate == true {
		C.set_show_refresh_rate(o, C.int(1))
	} else {
		C.set_show_refresh_rate(o, C.int(0))
	}

	if c.DisableHardwarePulsing == true {
		C.set_disable_hardware_pulsing(o, C.int(1))
	} else {
		C.set_disable_hardware_pulsing(o, C.int(0))
	}

	if c.InverseColors == true {
		C.set_inverse_colors(o, C.int(1))
	} else {
		C.set_inverse_colors(o, C.int(0))
	}

	return o
}
