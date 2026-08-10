package renderers

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/autodarts"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type UserCountDashboardRenderer struct {
	*MatrixWriter
	UserCount  int
	MatchCount int
	c          *autodarts.AutodartsWSClient
}

func UserCountDashboard(screen *rgbmatrix.Screen) (*UserCountDashboardRenderer, error) {
	config, err := rgbmatrix.LoadConfigFile("./config.toml")
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	font, err := rgbmatrix.LoadBDF(config.Dashboards.Font)
	if err != nil {
		return nil, fmt.Errorf("load dashboard font: %w", err)
	}
	c := autodarts.NewAutodartsWSClient()
	if err := c.Connect(); err != nil {
		return nil, fmt.Errorf("%w: connect to Autodarts: %v", ErrServiceUnavailable, err)
	}
	return &UserCountDashboardRenderer{MatrixWriter: NewMatrixWriter(screen, font), c: c}, nil
}

func (r *UserCountDashboardRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	r.UserCount = autodarts.NewAutodartsAPIClient().GetOnlineUsers()
	r.MatchCount = autodarts.NewAutodartsAPIClient().GetNumMatches()
	p := message.NewPrinter(language.English)
	onUsers := r.c.OnOnlineUsersChange(ctx)
	onMatch := r.c.OnMatchCountChange(ctx)

	blue := color.RGBA{0, 0, 200, 255}
	white := color.RGBA{200, 200, 200, 255}
	cyan := color.RGBA{0, 200, 200, 255}

	draw := func() {
		r.WriteLn("Autodarts", blue)
		r.MatrixWriter.y += 5
		r.Write("Users:   ", white)
		r.WriteLn(p.Sprintf("%5d", r.UserCount), cyan)
		r.Write("Matches: ", white)
		r.WriteLn(p.Sprintf("%5d", r.MatchCount), cyan)
		r.NewLine()
		r.WriteLn(time.Now().Format("15:04"), white)
		r.Flush()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-onUsers:
			r.UserCount = msg.Online
			draw()
		case msg := <-onMatch:
			r.MatchCount = msg.Count
			draw()
		}
	}
}
