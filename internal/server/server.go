// Package server provides HTTP server functionality for the ASCII art web application.
// It handles routing, middleware integration, and request processing.
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"ascii-art-web/internal/middleware"
	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
	"ascii-art-web/internal/sanitize"
	"ascii-art-web/internal/validation"
)

type Server struct {
	*http.Server
}

func New(addr string) *Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/ascii-art", handleAsciiArt)

	// Chain middleware: Security -> Recovery -> Logging -> Routes
	logger := log.New(os.Stdout, "", log.LstdFlags)
	handler := middleware.SecurityHeaders(
		middleware.Recovery(logger)(
			middleware.Logging(logger)(mux)))

	return &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Home"))
}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if banner == "" {
		banner = "standard"
	}

	if err := validation.ValidateText(text); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateBanner(banner); err != nil {
		http.Error(w, "Not Found: invalid banner", http.StatusNotFound)
		return
	}

	bannerPath := "cmd/ascii-art/testdata/" + banner + ".txt"
	bannerData, err := parser.LoadBanner(os.DirFS("."), bannerPath)
	if err != nil {
		http.Error(w, "Not Found: banner file not found", http.StatusNotFound)
		return
	}

	result, err := renderer.ASCII(text, bannerData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	safeResult := sanitize.HTML(result)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(safeResult))
}
