package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
)

// Wrapper for http server
type Server struct {
	httpServer *http.Server
}

// Create new http server
func NewServer(config *config.ServerConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.Host, config.Port),
			MaxHeaderBytes:    config.MaxHeaderBytes,
			Handler:           handler,
			ReadTimeout:       config.ReadTimeout,
			WriteTimeout:      config.WriteTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			IdleTimeout:       config.IdleTimeout,
		},
	}
}

// Launch http server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown http server with context
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
