package http

import (
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/blankrobot/pulpe"
	"github.com/julienschmidt/httprouter"
)

// NewHandler instantiates a new Handler.
func NewHandler(router *httprouter.Router) *Handler {
	return &Handler{
		router: router,
	}
}

// Handler is a collection of all the service handlers.
type Handler struct {
	staticHandler http.Handler
	assetsPath    string
	indexPath     string
	router        *httprouter.Router
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rw := NewResponseWriter(w)
	switch {
	case strings.HasPrefix(r.URL.Path, "/api"):
		h.router.ServeHTTP(rw, r)
	case h.staticHandler != nil && strings.HasPrefix(r.URL.Path, "/assets/"):
		// save the actual path because the static handler strips the "assets" prefix.
		actualPath := r.URL.Path
		h.staticHandler.ServeHTTP(rw, r)
		r.URL.Path = actualPath
	default:
		if h.assetsPath != "" {
			http.ServeFile(rw, r, h.indexPath)
		} else {
			http.NotFound(rw, r)
		}
	}

	// TODO use httprouter and define all the routes
	log.Printf(
		"%s %s %s %d %d %s",
		r.RemoteAddr,
		r.Method,
		r.URL,
		rw.status,
		rw.len,
		time.Since(start),
	)
}

// EnableStatic sets the assets directory to be served.
// By default, no assets are served.
func (h *Handler) EnableStatic(assetsPath string) {
	h.staticHandler = http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath)))
	h.assetsPath = assetsPath
	h.indexPath = path.Join(h.assetsPath, "index.html")
}

// NewResponseWriter instantiates a ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

// ResponseWriter is a wrapper around http.ResponseWriter.
// It allows to capture informations about the response.
type ResponseWriter struct {
	http.ResponseWriter

	status int
	len    int
}

// WriteHeader stores the status before calling the underlying
// http.ResponseWriter WriteHeader.
func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.len = len(data)
	return w.ResponseWriter.Write(data)
}

// Connector creates a session from a client and a request.
type Connector func(http.ResponseWriter, *http.Request) pulpe.Session

// NewCookieConnector returns a connector that creates a session and loads the session token from a cookie.
func NewCookieConnector(client pulpe.Client) Connector {
	return func(w http.ResponseWriter, r *http.Request) pulpe.Session {
		session := client.Connect()

		cookie, err := r.Cookie("pulpesid")
		if err == nil {
			session.SetAuthToken(cookie.Value)
		}

		return session
	}
}
