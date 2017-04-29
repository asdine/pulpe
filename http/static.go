package http

import "net/http"

// RegisterStaticHandler register the routes for serving static files.
func RegisterStaticHandler(mux *ServeMux, path string) {
	staticHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(path)))

	mux.HandleFunc("assets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// save the actual path because the static handler strips the "assets" prefix.
		actualPath := r.URL.Path
		staticHandler.ServeHTTP(w, r)
		r.URL.Path = actualPath
	})
}
