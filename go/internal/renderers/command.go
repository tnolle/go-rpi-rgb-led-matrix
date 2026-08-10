package renderers

import "context"

type ScreenType int

const (
	TypeImage ScreenType = iota
	TypeGIF
	TypeGIFOnce
	TypeDashboard
	TypePlayground
	TypeAnimation
)

func (t ScreenType) String() string {
	switch t {
	case TypeImage:
		return "image"
	case TypeGIF:
		return "gif"
	case TypeGIFOnce:
		return "gif-once"
	case TypeDashboard:
		return "dashboard"
	case TypePlayground:
		return "playground"
	case TypeAnimation:
		return "animation"
	default:
		return "unknown"
	}
}

type Command struct {
	Type        ScreenType
	Name        string
	IsTemporary bool
	Context     context.Context
	Result      chan error
}
