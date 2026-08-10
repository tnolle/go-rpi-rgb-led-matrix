package rgbmatrix

type ScanMode int8

const (
	Progressive ScanMode = 0
	Interlaced  ScanMode = 1
)

type RuntimeOptions struct {
	// GPIOSlowdown is the slowdown factor for GPIO access.
	GPIOSlowdown int `toml:"gpio_slowdown"`
}

// DefaultConfig default WS281x configuration
var DefaultConfig = MatrixOptions{
	HardwareMapping:   "regular",
	Rows:              32,
	Cols:              32,
	ChainLength:       1,
	Parallel:          1,
	PWMBits:           11,
	PWMDitherBits:     0,
	PWMLSBNanoseconds: 130,
	Brightness:        100,
	ScanMode:          Progressive,
}

// MatrixOptions rgb-led-matrix configuration
type MatrixOptions struct {
	// Rows the number of rows supported by the display, so 32 or 16.
	Rows int `toml:"rows"`
	// Cols the number of columns supported by the display, so 32 or 64 .
	Cols int `toml:"cols"`

	// ChainLength is the number of displays daisy-chained together
	// (output of one connected to input of next).
	ChainLength int `toml:"chain_length"`

	// Parallel is the number of parallel chains connected to the Pi; in old Pis
	// with 26 GPIO pins, that is 1, in newer Pis with 40 interfaces pins, that
	// can also be 2 or 3. The effective number of pixels in vertical direction is
	// then thus rows * parallel.
	Parallel int `toml:"parallel"`

	// Set PWM bits used for output. Default is 11, but if you only deal with
	// limited comic-colors, 1 might be sufficient. Lower require less CPU and
	// increases refresh-rate.
	PWMBits int `toml:"pwm_bits"`

	// Change the base time-unit for the on-time in the lowest significant bit in
	// nanoseconds.  Higher numbers provide better quality (more accurate color,
	// less ghosting), but have a negative impact on the frame rate.
	PWMLSBNanoseconds int `toml:"pwm_lsb_nanoseconds"` // the DMA channel to use

	// The lower bits can be time-dithered for higher refresh rate.
	// Flag: --led-pwm-dither-bits
	PWMDitherBits int `toml:"pwm_dither_bits"`

	// Brightness is the initial brightness of the panel in percent. Valid range
	// is 1..100
	Brightness int `toml:"brightness"`

	// ScanMode progressive or interlaced
	ScanMode ScanMode `toml:"scan_mode"` // strip color layout

	RowAddressType int `toml:"row_address_type"`

	Multiplexing int `toml:"multiplexing"`

	// Disable the PWM hardware subsystem to create pulses. Typically, you don't
	// want to disable hardware pulsing, this is mostly for debugging and figuring
	// ppm if there is interference with the sound system.
	// This won't do anything if output enable is not connected to GPIO 18 in
	// non-standard wirings.
	DisableHardwarePulsing bool `toml:"disable_hardware_pulsing"`

	ShowRefreshRate    bool `toml:"show_refresh_rate"`
	LimitRefreshRateHz int  `toml:"limit_refresh_rate_hz"`
	InverseColors      bool `toml:"inverse_colors"`

	// Name of GPIO mapping used
	HardwareMapping string `toml:"hardware_mapping"`
}

func (c *MatrixOptions) geometry() (width, height int) {
	return c.Cols * c.ChainLength, c.Rows * c.Parallel
}
