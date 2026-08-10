package renderers

import (
	"context"
)

type AfterRenderFunc func()

type Renderer interface {
	Render(ctx context.Context, cb ...AfterRenderFunc) error
}
