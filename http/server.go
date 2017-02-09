package http

import (
	"net"
	"net/http"
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"
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

	srv *graceful.Server
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	srv := graceful.Server{
		Server: &http.Server{Handler: s.Handler},
	}

	// Start HTTP server.
	go func() {
		srv.Serve(ln)
	}()

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.srv != nil {
		s.srv.Stop(time.Second)
	}

	return nil
}
