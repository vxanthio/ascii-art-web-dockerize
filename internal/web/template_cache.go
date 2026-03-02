// Package web implements the HTTP presentation layer of the application.
//
// It is responsible for rendering HTML templates and preparing data
// structures used by HTTP handlers.
package web

import (
	"html/template"
)

// NewTemplateCache parses HTML template files and stores them in a cache.
//
// It loads the base and page templates at startup and returns a map
// keyed by template name for efficient reuse by HTTP handlers.
func NewTemplateCache() (map[string]*template.Template, error) {

	templateMap := make(map[string]*template.Template)
	ts, err := template.ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		return nil, err
	}
	templateMap["index.html"] = ts
	return templateMap, nil
}
