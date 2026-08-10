//go:generate go run github.com/dmarkham/enumer -type=Animation -transform=kebab

package animation

type Animation int

const (
	Aurora Animation = iota
	BlobbyFusion
	Checkerboard
	ColorWave
	Firefly
	Kaleidoscope
	LavaLamp
	Lightning
	Mandelbrot
	MatrixRain
	Nebula
	Plasma
	RadarSweep
	Ripple
	Spectrum
	Spiral
	Starfield
	Tunnel
	Vortex
	PixelBloom
	RGBFlow
	Glitch
	HypnoticRings
	SpinningGrid
	HexPulse
	SnakeTrail
	ExplosionBurst
	BeatGrid
	AudioOrbit
	AuroraCurtains
	UlamSpiral
	GameOfLife
	VectorFieldFlow
	SierpinskiTriangle
	FluidDream
	FluidRainbow
	OrbitingMetaballs
	MarbleShader
)
