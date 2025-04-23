package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

var HOST = "http://localhost:8085"

func do(host string, name string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s?name=%s", host, name), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func hosts() []string {
	var h []string
	if selected := viper.GetInt("selectedHost"); selected > 0 {
		h = append(h, viper.GetStringSlice("hosts")[selected-1])
		return h
	}
	for _, host := range viper.GetStringSlice("hosts") {
		h = append(h, host)
	}
	return h
}

func ShowDashboard(name string) {
	for _, h := range hosts() {
		_, err := do(h+"/dashboard", name)
		fmt.Println(h, err)
	}
}

func ShowImage(name string) {
	for _, h := range hosts() {
		_, err := do(h+"/image", name)
		fmt.Println(h, err)
	}
}

func ShowGIF(name string, once bool) {
	if once {
		for _, h := range hosts() {
			_, err := do(h+"/gif-once", name)
			fmt.Println(h, err)
		}
		return
	}
	for _, h := range hosts() {
		_, err := do(h+"/gif", name)
		fmt.Println(h, err)
	}
}

func ShowAnimation(name string) {
	for _, h := range hosts() {
		_, err := do(h+"/animation", name)
		fmt.Println(h, err)
	}
}
