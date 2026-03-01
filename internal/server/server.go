// Package server provides HTTP server functionality for ASCII art web application.
package server

import (
	"html/template"
	"net/http"
	"os"

	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
	"ascii-art-web/internal/validation"
	"ascii-art-web/internal/web"
)

func Start(addr string) error {
	// Initialize template cache
	cache := make(map[string]*template.Template)
	ts, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		return err
	}
	cache["index.html"] = ts

	// Create application with Vasiliki's handlers
	app := &web.Application{
		TemplateCache: cache,
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// Use Vasiliki's handlers
	http.HandleFunc("/", app.Home)
	http.HandleFunc("/ascii-art", app.HandleAsciiArt)
	
	return http.ListenAndServe(addr, nil)
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

