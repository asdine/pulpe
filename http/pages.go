package http

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

// RegisterPageHandler register the routes for serving pages.
func RegisterPageHandler(mux *ServeMux, dir string, lazy bool) {
	pattern := filepath.Join(dir, "*.tmpl.html")

	h := pageHandler{
		lazy: lazy,
		dir:  dir,
	}

	if !lazy {
		h.templates = template.Must(template.ParseGlob(pattern))
	}

	mux.HandleFunc("/join", h.handleRegister)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		h.handleIndex(w, r)
	})
}

type pageHandler struct {
	templates *template.Template
	lazy      bool
	dir       string
}

func (h *pageHandler) render(wr io.Writer, name string, data interface{}) {
	if h.lazy {
		template.Must(template.ParseFiles(filepath.Join(h.dir, name))).Execute(wr, data)
	} else {
		h.templates.ExecuteTemplate(wr, name, data)
	}
}

func (h *pageHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.render(w, "index.tmpl.html", map[string]string{
		"Title": "Pulpe",
	})
}

func (h *pageHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	h.render(w, "register.tmpl.html", map[string]string{
		"Title": "Join",
	})
}
