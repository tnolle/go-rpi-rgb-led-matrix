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

func UserCountDashboard(screen *rgbmatrix.Screen) *UserCountDashboardRenderer {
	font, _ := rgbmatrix.LoadBDF(fmt.Sprintf("fonts/7x14.bdf"))
	c := autodarts.NewAutodartsWSClient()
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	return &UserCountDashboardRenderer{MatrixWriter: NewMatrixWriter(screen, font)}
}

func (r *UserCountDashboardRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	go func() {
		r.UserCount = autodarts.NewAutodartsAPIClient().GetOnlineUsers()
		r.MatchCount = autodarts.NewAutodartsAPIClient().GetNumMatches()
		ws := autodarts.NewAutodartsWSClient()
		ws.Connect()
		p := message.NewPrinter(language.English)
		onUsers := ws.OnOnlineUsersChange(ctx)
		onMatch := ws.OnMatchCountChange(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-onUsers:
				r.UserCount = msg.Online
			case msg := <-onMatch:
				r.MatchCount = msg.Count
			default:
				r.Write("Users:   ", color.RGBA{200, 200, 200, 255})
				r.WriteLn(p.Sprintf("%5d", r.UserCount), color.RGBA{0, 200, 200, 0})
				r.Write("Matches: ", color.RGBA{200, 200, 200, 255})
				r.WriteLn(p.Sprintf("%5d", r.MatchCount), color.RGBA{0, 200, 200, 0})
				r.NewLine()
				r.WriteLn(time.Now().Format("15:04"), color.RGBA{200, 200, 200, 255})
				r.Flush()
			}
		}
	}()
	return nil
}
