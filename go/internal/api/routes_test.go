package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

func TestCommandValidation(t *testing.T) {
	handler := newHandler(make(chan renderers.Command), catalog{}, display.NewHub())
	tests := []struct {
		name      string
		path      string
		status    int
		errorCode string
	}{
		{name: "missing name", path: "/animation", status: http.StatusBadRequest, errorCode: "missing_name"},
		{name: "blank name", path: "/animation?name=%20", status: http.StatusBadRequest, errorCode: "missing_name"},
		{name: "unknown animation", path: "/animation?name=unknown", status: http.StatusNotFound, errorCode: "animation_not_found"},
		{name: "path traversal", path: "/image?name=..%2F..%2Fsecret", status: http.StatusNotFound, errorCode: "image_not_found"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := performRequest(handler, http.MethodPut, test.path)
			assertAPIError(t, response, test.status, test.errorCode)
		})
	}
}

func TestCommandSuccessWaitsForRenderer(t *testing.T) {
	commands := make(chan renderers.Command)
	handler := newHandler(commands, catalog{}, display.NewHub())
	go func() {
		command := <-commands
		if command.Type != renderers.TypeAnimation || command.Name != "plasma" {
			t.Errorf("unexpected command: %#v", command)
		}
		command.Result <- nil
	}()

	response := performRequest(handler, http.MethodPut, "/animation?name=plasma")
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d, body: %s", response.Code, response.Body.String())
	}
	var body displayResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Display.Type != "animation" || body.Display.Name != "plasma" || body.Display.Temporary {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestRendererErrorsAreMapped(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		status    int
		errorCode string
	}{
		{name: "invalid asset", err: renderers.ErrInvalidAsset, status: http.StatusUnprocessableEntity, errorCode: "invalid_image"},
		{name: "service unavailable", err: renderers.ErrServiceUnavailable, status: http.StatusServiceUnavailable, errorCode: "service_unavailable"},
		{name: "internal failure", err: errors.New("output failed"), status: http.StatusInternalServerError, errorCode: "renderer_failed"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commands := make(chan renderers.Command)
			root := t.TempDir()
			catalog := catalog{imagesDir: root}
			// The catalog only needs a sorted entry for validation; preparation is
			// represented by the command result in this handler-level test.
			writeTestFile(t, root, "valid.png")
			handler := newHandler(commands, catalog, display.NewHub())
			go func() {
				command := <-commands
				command.Result <- test.err
			}()

			response := performRequest(handler, http.MethodPut, "/image?name=valid")
			assertAPIError(t, response, test.status, test.errorCode)
		})
	}
}

func TestHealthOnlyAcceptsGet(t *testing.T) {
	handler := newHandler(make(chan renderers.Command), catalog{}, display.NewHub())
	response := performRequest(handler, http.MethodPost, "/healthz")
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

func performRequest(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func assertAPIError(t *testing.T, response *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("unexpected status: got %d, want %d; body: %s", response.Code, status, response.Body.String())
	}
	var body errorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != code || body.Error.Message == "" {
		t.Fatalf("unexpected error: %#v", body.Error)
	}
}

func writeTestFile(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
		t.Fatal(err)
	}
}
