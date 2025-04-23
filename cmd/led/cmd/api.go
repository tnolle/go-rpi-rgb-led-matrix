package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

var HOST = "http://localhost:8085"

func do(host string, command renderers.Type, name string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s?name=%s", host, command, name), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func ShowDashboard(name string) {
	for _, host := range viper.GetStringSlice("hosts") {
		do(host, renderers.TypeDashboard, name)
	}
}

func ShowImage(name string) {
	for _, h := range viper.GetStringSlice("hosts") {
		do(h, renderers.TypeImage, name)
	}
}

func ShowGIF(name string, once bool) {
	if once {
		for _, h := range viper.GetStringSlice("hosts") {
			do(h, renderers.TypeGIFOnce, name)
		}
		return
	}
	for _, h := range viper.GetStringSlice("hosts") {
		do(h, renderers.TypeGIF, name)
	}
}

func ShowAnimation(name string) {
	for _, h := range viper.GetStringSlice("hosts") {
		do(h, renderers.TypeAnimation, name)
	}
}
