package renderers

import (
	"context"
	"errors"
)

var (
	ErrInvalidAsset       = errors.New("invalid asset")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrUnknownContent     = errors.New("unknown content")
)

type AfterRenderFunc func()

type Renderer interface {
	Render(ctx context.Context, cb ...AfterRenderFunc) error
}
