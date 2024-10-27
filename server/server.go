// C:\GoProject\src\eShop\server\server.go

package server

import (
	"context"
	"net/http"
	"time"
)

// Server struct represents an HTTP server with a field for http.Server...
type Server struct {
	httpServer *http.Server // Embedded HTTP server...
}

// Run starts the HTTP server on the specified port with the provided handler...
// Configures the server with a max header size of 1 MB and read/write timeouts of 10 seconds...
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          // 1 MB
		ReadTimeout:    10 * time.Second, // 10 seconds
		WriteTimeout:   10 * time.Second, // 10 seconds
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server, closing all active connections...
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
