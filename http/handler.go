package http

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/blankrobot/pulpe"
	"github.com/julienschmidt/httprouter"
)

// NewHandler instantiates a new Handler.
func NewHandler(c pulpe.Client) *Handler {
	router := httprouter.New()

	client := client{c}
	registerCardHandler(router, &client)
	registerListHandler(router, &client)
	registerBoardHandler(router, &client)
	registerUserHandler(router, &client)

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

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, status int, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

type client struct {
	client pulpe.Client
}

func (c *client) session(w http.ResponseWriter, r *http.Request) pulpe.Session {
	session := c.client.Connect()

	cookie, err := r.Cookie("pulpesid")
	if err == nil {
		session.SetAuthToken(cookie.Value)
	}

	return session
}
