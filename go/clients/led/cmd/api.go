package cmd

import (
	"encoding/json"
	"errors"
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
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return decodeAPIError(res)
	}
	_, _ = io.Copy(io.Discard, res.Body)
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

func SetDashboard(name string) error {
	return setOnHosts("dashboard", name)
}

func SetImage(name string) error {
	return setOnHosts("image", name)
}

func SetGIF(name string, once bool) error {
	endpoint := "gif"
	if once {
		endpoint = "gif-once"
	}
	return setOnHosts(endpoint, name)
}

func SetAnimation(name string) error {
	return setOnHosts("animation", name)
}

func setOnHosts(endpoint, name string) error {
	var errs []error
	for _, h := range hosts() {
		if err := do(h+"/"+endpoint, name); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", h, err))
			continue
		}
		fmt.Println(h, "OK")
	}
	return errors.Join(errs...)
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiErrorResponse struct {
	Error apiError `json:"error"`
}

func decodeAPIError(res *http.Response) error {
	var body apiErrorResponse
	if err := json.NewDecoder(io.LimitReader(res.Body, 4096)).Decode(&body); err == nil && body.Error.Message != "" {
		return errors.New(body.Error.Message)
	}
	return fmt.Errorf("server returned %s", res.Status)
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
