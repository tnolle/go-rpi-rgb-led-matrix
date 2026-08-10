package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/display"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

func ListenAndServe(commands chan renderers.Command, frames *display.Hub) {
	server := &http.Server{
		Addr:              ":8085",
		Handler:           NewHandler(commands, frames),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("HTTP API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func NewHandler(commands chan renderers.Command, frames *display.Hub) http.Handler {
	return newHandler(commands, defaultCatalog(), frames)
}

func newHandler(commands chan renderers.Command, catalog catalog, frames *display.Hub) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("GET /display/stream", displayStreamHandler(frames))

	mux.HandleFunc("PUT /image", commandHandler(commands, commandSpec{
		kind: "image", commandType: renderers.TypeImage, available: catalog.images,
	}))
	mux.HandleFunc("PUT /gif", commandHandler(commands, commandSpec{
		kind: "gif", commandType: renderers.TypeGIF, available: catalog.gifs,
	}))
	mux.HandleFunc("PUT /gif-once", commandHandler(commands, commandSpec{
		kind: "gif", commandType: renderers.TypeGIFOnce, temporary: true, available: catalog.gifs,
	}))
	mux.HandleFunc("PUT /dashboard", commandHandler(commands, commandSpec{
		kind: "dashboard", commandType: renderers.TypeDashboard,
		available: func() ([]string, error) { return catalog.dashboards(), nil },
	}))
	mux.HandleFunc("PUT /animation", commandHandler(commands, commandSpec{
		kind: "animation", commandType: renderers.TypeAnimation,
		available: func() ([]string, error) { return catalog.animations(), nil },
	}))

	return mux
}

type catalogResponse struct {
	Items []string `json:"items"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type displayResponse struct {
	Display displayCommand `json:"display"`
}

type displayCommand struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Temporary bool   `json:"temporary"`
}

type commandSpec struct {
	kind        string
	commandType renderers.ScreenType
	temporary   bool
	available   func() ([]string, error)
}

func catalogHandler(load func() ([]string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		items, err := load()
		if err != nil {
			log.Printf("load catalog: %v", err)
			writeError(w, http.StatusInternalServerError, "catalog_failed", "the catalog could not be loaded")
			return
		}
		writeJSON(w, http.StatusOK, catalogResponse{Items: items})
	}
}

func commandHandler(commands chan renderers.Command, spec commandSpec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" {
			writeError(w, http.StatusBadRequest, "missing_name", "query parameter \"name\" is required")
			return
		}

		items, err := spec.available()
		if err != nil {
			log.Printf("validate %s %q: %v", spec.kind, name, err)
			writeError(w, http.StatusInternalServerError, "catalog_failed", "the catalog could not be loaded")
			return
		}
		if !contains(items, name) {
			writeError(w, http.StatusNotFound, spec.kind+"_not_found", spec.kind+" \""+name+"\" does not exist")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		result := make(chan error, 1)
		command := renderers.Command{
			Type:        spec.commandType,
			Name:        name,
			IsTemporary: spec.temporary,
			Context:     ctx,
			Result:      result,
		}

		select {
		case commands <- command:
		case <-ctx.Done():
			writeCommandContextError(w, ctx.Err())
			return
		}

		select {
		case err := <-result:
			if err != nil {
				writeRendererError(w, spec.kind, name, err)
				return
			}
			writeJSON(w, http.StatusOK, displayResponse{Display: displayCommand{
				Type: spec.kind, Name: name, Temporary: spec.temporary,
			}})
		case <-ctx.Done():
			writeCommandContextError(w, ctx.Err())
		}
	}
}

func writeRendererError(w http.ResponseWriter, kind, name string, err error) {
	log.Printf("start %s %q: %v", kind, name, err)
	switch {
	case errors.Is(err, renderers.ErrInvalidAsset):
		writeError(w, http.StatusUnprocessableEntity, "invalid_"+kind, kind+" \""+name+"\" could not be decoded")
	case errors.Is(err, renderers.ErrServiceUnavailable):
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", kind+" \""+name+"\" is temporarily unavailable")
	case errors.Is(err, renderers.ErrUnknownContent):
		writeError(w, http.StatusNotFound, kind+"_not_found", kind+" \""+name+"\" does not exist")
	default:
		writeError(w, http.StatusInternalServerError, "renderer_failed", "the renderer could not be started")
	}
}

func writeCommandContextError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		writeError(w, http.StatusGatewayTimeout, "renderer_timeout", "the renderer did not start in time")
		return
	}
	writeError(w, http.StatusRequestTimeout, "request_cancelled", "the request was cancelled")
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: apiError{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
