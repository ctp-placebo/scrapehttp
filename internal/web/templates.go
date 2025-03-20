package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", name, err)
	}
	return nil
}
