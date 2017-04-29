package http

import (
	"log"
	"net/http"
	"time"

	"github.com/blankrobot/pulpe"
)

// NewHandler instantiates a new Handler.
func NewHandler(handler http.Handler) *Handler {
	return &Handler{
		handler: handler,
	}
}

// Handler is a wrapper around a http.Handler.
type Handler struct {
	handler http.Handler
}

// ServeHTTP delegates a request to the given handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rw := NewResponseWriter(w)
	h.handler.ServeHTTP(rw, r)

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
