package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var HOST = "http://localhost:8085"

var httpClient = &http.Client{Timeout: 5 * time.Second}

func do(endpoint string, name string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	query := u.Query()
	query.Set("name", name)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPut, u.String(), nil)
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, res.Body)
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("server returned %s", res.Status)
	}
	return nil
}

func hosts() []string {
	var h []string
	if selected := viper.GetInt("selectedHost"); selected > 0 {
		configured := viper.GetStringSlice("hosts")
		if selected <= len(configured) {
			h = append(h, strings.TrimRight(configured[selected-1], "/"))
		}
		if len(h) > 0 {
			return h
		}
	}
	for _, host := range viper.GetStringSlice("hosts") {
		h = append(h, strings.TrimRight(host, "/"))
	}
	if len(h) == 0 {
		h = append(h, HOST)
	}
	return h
}

func ShowDashboard(name string) {
	for _, h := range hosts() {
		err := do(h+"/dashboard", name)
		fmt.Println(h, err)
	}
}

func ShowImage(name string) {
	for _, h := range hosts() {
		err := do(h+"/image", name)
		fmt.Println(h, err)
	}
}

func ShowGIF(name string, once bool) {
	if once {
		for _, h := range hosts() {
			err := do(h+"/gif-once", name)
			fmt.Println(h, err)
		}
		return
	}
	for _, h := range hosts() {
		err := do(h+"/gif", name)
		fmt.Println(h, err)
	}
}

func ShowAnimation(name string) {
	for _, h := range hosts() {
		err := do(h+"/animation", name)
		fmt.Println(h, err)
	}
}

type catalogResponse struct {
	Items []string `json:"items"`
}

func fetchCatalog(host, endpoint string) ([]string, error) {
	res, err := httpClient.Get(strings.TrimRight(host, "/") + endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = res.Status
		}
		return nil, fmt.Errorf("server returned %s: %s", res.Status, message)
	}

	var body catalogResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode catalog: %w", err)
	}
	return body.Items, nil
}
