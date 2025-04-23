package renderers

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"time"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/rgbmatrix"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ShopifyDashboardRenderer struct {
	*MatrixWriter
	Data   shopifyData
	config rgbmatrix.Config
}

func ShopifyDashboard(screen *rgbmatrix.Screen) *ShopifyDashboardRenderer {
	config := rgbmatrix.LoadConfig()
	font, _ := rgbmatrix.LoadBDF(fmt.Sprintf(config.Dashboards.Font))
	return &ShopifyDashboardRenderer{MatrixWriter: NewMatrixWriter(screen, font), config: config}
}

func (r *ShopifyDashboardRenderer) Render(ctx context.Context, _ ...AfterRenderFunc) error {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	orange := color.RGBA{200, 100, 0, 255}
	white := color.RGBA{200, 200, 200, 255}
	green := color.RGBA{0, 200, 0, 255}
	blue := color.RGBA{0, 200, 200, 255}

	p := message.NewPrinter(language.English)

	requestCtx, cancelRequest := context.WithCancel(ctx)

	draw := func() {
		r.WriteLn("Shopify", green)
		r.MatrixWriter.y += 5

		r.Write("Sales:", white)
		r.WriteLn(p.Sprintf("%14d €", r.Data.TotalSales), green)
		r.Write("Orders:", white)
		r.WriteLn(p.Sprintf("%13d", r.Data.TotalOrders), blue)
		r.Write("Sales (M):", white)
		r.WriteLn(p.Sprintf("%10d €", r.Data.MonthlySales), green)
		r.Write("Sales (d):", white)
		r.WriteLn(p.Sprintf("%10d €", r.Data.TodaySales), green)
		r.Write("Orders (d):", white)
		r.WriteLn(p.Sprintf("%9d", r.Data.TodayOrders), blue)
		r.Flush()
	}

	r.Data, _ = getTotalSales(ctx, r.config)
	draw()

	for {
		select {
		case <-ctx.Done():
			cancelRequest()
			return nil
		case <-t.C:
			cancelRequest()
			requestCtx, cancelRequest = context.WithTimeout(ctx, 1*time.Second)
			go func() {
				defer cancelRequest()
				newData, err := getTotalSales(requestCtx, r.config)
				if err != nil {
					return
				}
				if newData != r.Data {
					r.Data = newData
					r.screen.Fill(orange)
					draw()
					time.Sleep(200 * time.Millisecond)
					r.screen.Fill(color.RGBA{0, 0, 0, 0})
					draw()
				}
			}()
		}
	}
}

type shopifyData struct {
	TotalSales   int
	TotalOrders  int
	MonthlySales int
	TodaySales   int
	TodayOrders  int
}

func getTotalSales(ctx context.Context, config rgbmatrix.Config) (d shopifyData, err error) {
	if d.TotalSales, err = fetch(ctx, config.Shopify.TotalSales); err != nil {
		return
	}
	if d.TotalOrders, err = fetch(ctx, config.Shopify.TotalOrders); err != nil {
		return
	}
	if d.MonthlySales, err = fetch(ctx, config.Shopify.MonthlySales); err != nil {
		return
	}
	if d.TodaySales, err = fetch(ctx, config.Shopify.TodaySales); err != nil {
		return
	}
	if d.TodayOrders, err = fetch(ctx, config.Shopify.TodayOrders); err != nil {
		return
	}
	return
}

type result struct {
	Number int `json:"number"`
}

func fetch(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	var data result
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Number, nil
}
