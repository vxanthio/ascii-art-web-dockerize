// Package server provides HTTP server functionality for ASCII art web application.
package server

import (
	"net/http"
	"os"

	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
	"ascii-art-web/internal/validation"
)

func Start(addr string) error {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ascii-art", handleAsciiArt)
	return http.ListenAndServe(addr, nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Home"))
}

// GenerateASCII validates input and generates ASCII art.
// Returns the result string and an error with appropriate HTTP status code.
func GenerateASCII(text, banner string) (string, int, error) {
	if banner == "" {
		banner = "standard"
	}

	if err := validation.ValidateText(text); err != nil {
		return "", http.StatusBadRequest, err
	}

	if err := validation.ValidateBanner(banner); err != nil {
		return "", http.StatusNotFound, err
	}

	bannerPath := "cmd/ascii-art/testdata/" + banner + ".txt"
	bannerData, err := parser.LoadBanner(os.DirFS("."), bannerPath)
	if err != nil {
		return "", http.StatusNotFound, err
	}

	result, err := renderer.ASCII(text, bannerData)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("input")
	banner := r.FormValue("font")

	result, status, err := GenerateASCII(text, banner)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
