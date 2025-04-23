package api

import (
	"log"
	"net/http"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers"
)

func ListenAndServe(commands chan renderers.Command) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Image handlers
	http.HandleFunc("PUT /image", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeImage, Name: name}
		w.WriteHeader(http.StatusOK)
	})

	// GIF handlers
	http.HandleFunc("PUT /gif", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeGIF, Name: name}
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("PUT /gif-once", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeGIFOnce, Name: name, IsTemporary: true}
		w.WriteHeader(http.StatusOK)
	})

	// Dashboard handlers
	http.HandleFunc("PUT /dashboard", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeDashboard, Name: name}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("PUT /animation", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		commands <- renderers.Command{Type: renderers.TypeAnimation, Name: name}
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8085", nil))
}
