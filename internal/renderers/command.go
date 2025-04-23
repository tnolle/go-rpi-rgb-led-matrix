package renderers

type Type string

const (
	TypeImage      Type = "image"
	TypeGIF        Type = "gif"
	TypeGIFOnce    Type = "gif-once"
	TypeDashboard  Type = "dashboard"
	TypePlayground Type = "playground"
	TypeAnimation  Type = "animation"
)

type Command struct {
	Type        Type
	Name        string
	IsTemporary bool
}

type DashboardName string

const (
	DashboardClock     DashboardName = "clock"
	DashboardUserCount DashboardName = "user-count"
)

type AnimationName string

const (
	AnimationAurora             AnimationName = "aurora"
	AnimationBlobbyFusion       AnimationName = "blobby-fusion"
	AnimationCheckerboard       AnimationName = "checkerboard"
	AnimationColorWave          AnimationName = "color-wave"
	AnimationFirefly            AnimationName = "firefly"
	AnimationKaleidoscope       AnimationName = "kaleidoscope"
	AnimationLavaLamp           AnimationName = "lava-lamp"
	AnimationLightning          AnimationName = "lightning"
	AnimationMandelbrot         AnimationName = "mandelbrot"
	AnimationMatrixRain         AnimationName = "matrix-rain"
	AnimationNebula             AnimationName = "nebula"
	AnimationPlasma             AnimationName = "plasma"
	AnimationRadarSweep         AnimationName = "radar"
	AnimationRipple             AnimationName = "ripple"
	AnimationSpectrum           AnimationName = "spectrum"
	AnimationSpiral             AnimationName = "spiral"
	AnimationStarfield          AnimationName = "starfield"
	AnimationTunnel             AnimationName = "tunnel"
	AnimationVortex             AnimationName = "vortex"
	AnimationPixelBloom         AnimationName = "pixel-bloom"
	AnimationRGBFlow            AnimationName = "rgb-flow"
	AnimationGlitch             AnimationName = "glitch"
	AnimationHypnoticRings      AnimationName = "hypnotic-rings"
	AnimationSpinningGrid       AnimationName = "spinning-grid"
	AnimationHexPulse           AnimationName = "hex-pulse"
	AnimationSnakeTrail         AnimationName = "snake-trail"
	AnimationExplosionBurst     AnimationName = "explosion-burst"
	AnimationBeatGrid           AnimationName = "beat-grid"
	AnimationAudioOrbit         AnimationName = "audio-orbit"
	AnimationAuroraCurtains     AnimationName = "aurora-curtains"
	AnimationUlamSpiral         AnimationName = "ulam-spiral"
	AnimationGameOfLife         AnimationName = "game-of-life"
	AnimationVectorFieldFlow    AnimationName = "vector-field-flow"
	AnimationSierpinskiTriangle AnimationName = "sierpinski-triangle"
	AnimationFluidDream         AnimationName = "fluid-dream"
	AnimationFluidRainbow       AnimationName = "fluid-rainbow"
	AnimationOrbitingMetaballs  AnimationName = "orgbiting-metaballs"
	AnimationMarbleShader       AnimationName = "marble-shader"
)
