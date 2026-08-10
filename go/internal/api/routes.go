package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

func ListenAndServe(commands chan renderers.Command) {
	log.Fatal(http.ListenAndServe(":8085", NewHandler(commands)))
}

func NewHandler(commands chan renderers.Command) http.Handler {
	return newHandler(commands, defaultCatalog())
}

func newHandler(commands chan renderers.Command, catalog catalog) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /images", catalogHandler(catalog.images))
	mux.HandleFunc("GET /gifs", catalogHandler(catalog.gifs))
	mux.HandleFunc("GET /dashboards", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, catalogResponse{Items: catalog.dashboards()})
	})
	mux.HandleFunc("GET /animations", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, catalogResponse{Items: catalog.animations()})
	})

	// Image handlers
	mux.HandleFunc("PUT /image", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeImage, Name: name}
		w.WriteHeader(http.StatusOK)
	})

	// GIF handlers
	mux.HandleFunc("PUT /gif", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeGIF, Name: name}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("PUT /gif-once", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeGIFOnce, Name: name, IsTemporary: true}
		w.WriteHeader(http.StatusOK)
	})

	// Dashboard handlers
	mux.HandleFunc("PUT /dashboard", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeDashboard, Name: name}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("PUT /animation", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeAnimation, Name: name}
		w.WriteHeader(http.StatusOK)
	})

	return mux
}

type catalogResponse struct {
	Items []string `json:"items"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func catalogHandler(load func() ([]string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		items, err := load()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, catalogResponse{Items: items})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
