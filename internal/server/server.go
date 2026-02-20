package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"ascii-art-color/internal/middleware"
)

type Server struct {
	*http.Server
}

func New(addr string) *Server {
	mux := http.NewServeMux()
	
	// Register routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/ascii-art", handleAsciiArt)

	// Wrap with logging middleware
	logger := log.New(os.Stdout, "", log.LstdFlags)
	handler := middleware.Logging(logger)(mux)

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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ASCII Art"))
}
