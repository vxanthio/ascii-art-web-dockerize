package web

import (
	"html/template"
)

func NewTemplateCache() (map[string]*template.Template, error) {

	templateMap := make(map[string]*template.Template)
	ts, err := template.ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		return nil, err
	}
	templateMap["index.html"] = ts
	return templateMap, nil
}
