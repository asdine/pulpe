package http

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/blankrobot/pulpe"
)

// HTTP errors
const (
	ErrInvalidJSON = pulpe.Error("invalid json")
)

// NewHandler instantiates a new Handler.
func NewHandler(c pulpe.Client) *Handler {
	return &Handler{
		CardHandler:  NewCardHandler(c),
		ListHandler:  NewListHandler(c),
		BoardHandler: NewBoardHandler(c),
	}
}

// Handler is a collection of all the service handlers.
type Handler struct {
	CardHandler   *CardHandler
	ListHandler   *ListHandler
	BoardHandler  *BoardHandler
	StaticHandler http.Handler
	assetsPath    string
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rw := NewResponseWriter(w)
	switch {
	case strings.HasPrefix(r.URL.Path, "/v1/cards"):
		h.CardHandler.ServeHTTP(rw, r)
	case strings.HasPrefix(r.URL.Path, "/v1/lists"):
		h.ListHandler.ServeHTTP(rw, r)
	case strings.HasPrefix(r.URL.Path, "/v1/boards"):
		h.BoardHandler.ServeHTTP(rw, r)
	case h.assetsPath != "" && strings.HasPrefix(r.URL.Path, "/assets/"):
		h.StaticHandler.ServeHTTP(rw, r)
	default:
		if h.assetsPath != "" {
			http.ServeFile(rw, r, path.Join(h.assetsPath, "index.html"))
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

// SetStatic sets the assets directory to be served.
// By default, no assets are served.
func (h *Handler) SetStatic(assetsPath string) {
	h.StaticHandler = http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath)))
	h.assetsPath = assetsPath
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

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = pulpe.ErrInternal
	}

	// Write generic error response.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

// errorResponse is a generic response for sending an error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// NotFound writes an API error message to the response.
func NotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}` + "\n"))
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, status int, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}
