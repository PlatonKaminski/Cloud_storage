package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"cloud_storage/config"
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler, cfg config.HTTPConfig) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
