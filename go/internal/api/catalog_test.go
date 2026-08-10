package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

func TestAssetCatalogEndpoints(t *testing.T) {
	root := t.TempDir()
	images := filepath.Join(root, "pngs")
	gifs := filepath.Join(root, "gifs")
	if err := os.MkdirAll(images, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(gifs, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join(images, "zebra.png"),
		filepath.Join(images, "alpha.PNG"),
		filepath.Join(images, "ignored.jpg"),
		filepath.Join(gifs, "party.gif"),
	} {
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	handler := newHandler(make(chan renderers.Command), catalog{imagesDir: images, gifsDir: gifs}, display.NewHub())
	assertCatalog(t, handler, "/images", []string{"alpha", "zebra"})
	assertCatalog(t, handler, "/gifs", []string{"party"})
}

func TestMissingAssetDirectoryReturnsEmptyCatalog(t *testing.T) {
	handler := newHandler(make(chan renderers.Command), catalog{
		imagesDir: filepath.Join(t.TempDir(), "missing"),
		gifsDir:   filepath.Join(t.TempDir(), "missing"),
	}, display.NewHub())
	assertCatalog(t, handler, "/images", []string{})
	assertCatalog(t, handler, "/gifs", []string{})
}

func TestRendererCatalogEndpoints(t *testing.T) {
	handler := newHandler(make(chan renderers.Command), catalog{}, display.NewHub())

	request := httptest.NewRequest(http.MethodGet, "/dashboards", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}

	var body catalogResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	want := []string{"autodarts", "clock", "shopify"}
	if !reflect.DeepEqual(body.Items, want) {
		t.Fatalf("unexpected dashboards: got %v, want %v", body.Items, want)
	}
}

func assertCatalog(t *testing.T, handler http.Handler, path string, want []string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status for %s: %d", path, response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: %q", got)
	}

	var body catalogResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body.Items, want) {
		t.Fatalf("unexpected items for %s: got %v, want %v", path, body.Items, want)
	}
}
