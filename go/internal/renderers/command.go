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

type Command struct {
	Type        ScreenType
	Name        string
	IsTemporary bool
	Context     context.Context
	Result      chan error
}
