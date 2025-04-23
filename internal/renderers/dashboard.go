package renderers

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
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
				r.WriteLn("Autodarts", color.RGBA{0, 0, 200, 255})
				r.MatrixWriter.y += 5
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

type ShopifyDashboardRenderer struct {
	*MatrixWriter
	Data   shopifyData
	config rgbmatrix.Config
}

func ShopifyDashboard(screen *rgbmatrix.Screen) *ShopifyDashboardRenderer {
	font, _ := rgbmatrix.LoadBDF(fmt.Sprintf("fonts/7x14.bdf"))
	config := rgbmatrix.LoadConfig()
	return &ShopifyDashboardRenderer{MatrixWriter: NewMatrixWriter(screen, font), config: config}
}

func (r *ShopifyDashboardRenderer) Render(ctx context.Context, cb ...AfterRenderFunc) error {
	go func() {
		r.Data = getTotalSales(r.config)
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		p := message.NewPrinter(language.English)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				r.Data = getTotalSales(r.config)
			default:
				white := color.RGBA{200, 200, 200, 255}
				green := color.RGBA{0, 200, 0, 255}
				blue := color.RGBA{0, 200, 200, 255}

				r.WriteLn("Shopify", green)
				r.MatrixWriter.y += 5

				r.Write("Total Sales:    ", white)
				r.WriteLn(p.Sprintf("%9d €", r.Data.TotalSales), green)
				r.Write("Total Orders:   ", white)
				r.WriteLn(p.Sprintf("%9d", r.Data.TotalOrders), blue)
				r.Write("Monthly Sales:  ", white)
				r.WriteLn(p.Sprintf("%9d €", r.Data.MonthlySales), green)
				r.Write("Today's Sales:  ", white)
				r.WriteLn(p.Sprintf("%9d €", r.Data.TodaySales), green)
				r.Write("Today's Orders: ", white)
				r.WriteLn(p.Sprintf("%9d", r.Data.TodayOrders), blue)
				r.Flush()
			}
		}
	}()
	return nil
}

type shopifyData struct {
	TotalSales   int
	TotalOrders  int
	MonthlySales int
	TodaySales   int
	TodayOrders  int
}

func getTotalSales(config rgbmatrix.Config) (d shopifyData) {
	var err error
	if d.TotalSales, err = fetch(config.Shopify.TotalSales); err != nil {
		fmt.Println("Error fetching total sales:", err)
		return
	}
	if d.TotalOrders, err = fetch(config.Shopify.TotalOrders); err != nil {
		fmt.Println("Error fetching total orders:", err)
		return
	}
	if d.MonthlySales, err = fetch(config.Shopify.MonthlySales); err != nil {
		fmt.Println("Error fetching monthly sales:", err)
		return
	}
	if d.TodaySales, err = fetch(config.Shopify.TodaySales); err != nil {
		fmt.Println("Error fetching today's sales:", err)
		return
	}
	if d.TodayOrders, err = fetch(config.Shopify.TodayOrders); err != nil {
		fmt.Println("Error fetching today's orders:", err)
		return
	}
	return
}

type result struct {
	Number int `json:"number"`
}

func fetch(url string) (int, error) {
	res, err := http.Get(url)
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
