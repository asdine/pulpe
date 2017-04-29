package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RegisterStaticHandler register the routes for serving static files.
func RegisterStaticHandler(router *httprouter.Router, path string) {
	staticHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(path)))

	router.HandlerFunc("GET", "/assets/*assets", func(w http.ResponseWriter, r *http.Request) {
		// save the actual path because the static handler strips the "assets" prefix.
		actualPath := r.URL.Path
		staticHandler.ServeHTTP(w, r)
		r.URL.Path = actualPath
	})
}
