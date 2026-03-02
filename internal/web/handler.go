// Package web implements the HTTP layer of the ASCII Art Web application.
//
// Responsibilities:
//
//   - Registers and exposes HTTP handlers
//   - Renders HTML templates
//   - Handles HTTP errors (404, 500)
//
// Architecture:
//
//   - Application struct holds shared dependencies
//   - Handlers are defined as methods on Application
//   - Templates are retrieved from a pre-built template cache
//
// Example:
//
//	app := &Application{
//	    TemplateCache: templateCache,
//	}
//
//	http.HandleFunc("/", app.Home)
//	http.ListenAndServe(":8080", nil)
package web

import (
	"ascii-art-web/internal/server"
	"html/template"
	"net/http"
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
type PageData struct {
	Result string
	Title  string
}

// Home handles requests to the root route ("/").
//
// Behavior:
//
//   - Retrieves "index.html" from the template cache
//   - Renders the template if found
//   - Returns 404 if the template does not exist
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
	err := ts.Execute(w, PageData{Result: "", Title: "Home"})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func (app *Application) HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	result, status, err := server.GenerateASCII(text, banner)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	ts, found := app.TemplateCache["index.html"]
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageData{Result: result, Title: "Home"}
	if err := ts.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
