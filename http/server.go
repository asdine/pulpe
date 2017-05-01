package http

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/blankrobot/pulpe"
)

// NewServer returns a new instance of Server.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
	}
}

// Server represents an HTTP server.
type Server struct {
	// Handler to serve.
	Handler http.Handler

	// Bind address to open.
	Addr string

	server http.Server
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.server.Handler = s.Handler

	// Start HTTP server.
	go func() {
		err := s.server.Serve(ln)
		if err != http.ErrServerClosed {
			log.Print(err)
		}
	}()

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// NewServeMux instantiates a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux: http.NewServeMux(),
	}
}

// ServeMux is a wrapper around a http.Handler.
type ServeMux struct {
	*http.ServeMux
}

// ServeHTTP delegates a request to the underlying ServeMux.
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rw := NewResponseWriter(w)
	s.ServeMux.ServeHTTP(rw, r)

	log.Printf(
		"%s %s %s %d %d %s",
		clientIP(r),
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

// Connector creates a session from and a request.
type Connector func(*http.Request) pulpe.Session

// NewCookieConnector returns a connector that creates a session and loads the session token from a cookie.
func NewCookieConnector(client pulpe.Client) Connector {
	return func(r *http.Request) pulpe.Session {
		session := client.Connect()

		cookie, err := r.Cookie("pulpesid")
		if err == nil {
			session.SetAuthToken(cookie.Value)
		}

		return session
	}
}
