// Package handlers implements the HTTP layer of the ASCII Art Web application.
//
// Responsibilities:
//
//   - Registers and exposes HTTP handlers
//   - Renders HTML templates
//   - Validates user input and generates ASCII art
//   - Handles HTTP errors (400, 404, 405, 500)
//
// Architecture:
//
//   - Application struct holds shared dependencies
//   - Handlers are defined as methods on Application
//   - Templates are retrieved from a pre-built template cache
//
// Example:
//
//	cache, err := handlers.NewTemplateCache()
//	app := &handlers.Application{TemplateCache: cache}
//	http.HandleFunc("/", app.Home)
//	http.ListenAndServe(":8080", nil)
package handlers

import (
	"html/template"
	"net/http"

	"ascii-art-web/internal/banners"
	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
	"ascii-art-web/internal/validation"
)

// Application represents the web application container.
//
// It holds shared dependencies required by HTTP handlers.
//
// Fields:
//
//   - TemplateCache: a map of parsed HTML templates
//     Key   → template filename (e.g. "index.html")
//     Value → parsed *template.Template
type Application struct {
	TemplateCache map[string]*template.Template
}

// PageData represents the dynamic data passed to HTML templates.
//
// Fields:
//
//   - Result: the generated ASCII art string rendered inside the page.
//     This value is displayed in the <pre> block of index.html.
//   - Title: the page title displayed in the browser tab.
//   - Error: a human-readable error message displayed to the user.
//     Empty string means no error.
type PageData struct {
	Result string
	Title  string
	Error  string
}

// GenerateASCII validates input and generates ASCII art.
//
// Returns the ASCII art result string, an HTTP status code, and any error.
// The banner defaults to "standard" if empty.
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

	bannerData, err := parser.LoadBanner(banners.FS, banner+".txt")
	if err != nil {
		return "", http.StatusNotFound, err
	}

	result, err := renderer.ASCII(text, bannerData)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

// Home handles requests to the root route ("/").
//
// Behavior:
//
//   - Only accepts GET requests
//   - Returns 404 for any path other than "/"
//   - Returns 405 for non-GET methods
//   - Retrieves "index.html" from the template cache
//   - Returns 404 if the template does not exist in cache
//   - Returns 500 if template execution fails
//
// Parameters:
//
//   - w: http.ResponseWriter used to send data to the client
//   - r: *http.Request containing request metadata
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ts, found := app.TemplateCache["index.html"]
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ts.Execute(w, PageData{Title: "Home"}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// HandleAsciiArt handles POST requests to generate ASCII art.
//
// It validates the HTTP method, parses form input (text and banner),
// generates the ASCII art, and renders the result using the index.html template.
//
// On validation or generation error, it re-renders the page with an error message
// visible to the user instead of returning a bare HTTP error response.
//
// On template failure, it responds with the appropriate HTTP status code.
func (app *Application) HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ts, found := app.TemplateCache["index.html"]
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	result, status, err := GenerateASCII(text, banner)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		ts.Execute(w, PageData{Title: "Home", Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ts.Execute(w, PageData{Result: result, Title: "Home"}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
