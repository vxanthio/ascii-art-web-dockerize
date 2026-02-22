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

	ts, found := app.TemplateCache["index.html"]
	if !found {
		http.Error(w, "The template does not exist", http.StatusNotFound)
		return
	}

	err := ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to connect to the internal service", http.StatusInternalServerError)
	}
}
