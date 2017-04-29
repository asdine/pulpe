package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// RegisterPageHandler register the routes for serving pages.
func RegisterPageHandler(mux *ServeMux, dir string) {
	pattern := filepath.Join(dir, "*.tmpl.html")

	h := pageHandler{
		templates: template.Must(template.ParseGlob(pattern)),
	}

	fmt.Println(h.templates.DefinedTemplates())

	mux.HandleFunc("/", h.handleIndex)
}

type pageHandler struct {
	templates *template.Template
}

func (h *pageHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	h.templates.ExecuteTemplate(w, "index.tmpl.html", nil)
}
