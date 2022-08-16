package http

import (
	"analytics/internal/ports"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	analytics ports.AnalyticsPort
	auth      ports.AuthPort
	server    *http.Server
	logger    *zap.SugaredLogger
}

func New(analytics ports.AnalyticsPort, auth ports.AuthPort, logger *zap.SugaredLogger) *Server {
	return &Server{analytics: analytics, auth: auth, server: &http.Server{}, logger: logger}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("failed to listen on port 3000: %v", err)
	}

	s.server.Handler = s.routes()

	if err := s.server.Serve(listen); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve http server over port 3000: %v", err)
	}
	return nil
}

func (s *Server) routes() http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/analytics/v1/healtz", s.healtzHandler)
	r.Get("/analytics/v1/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", "/analytics/v1"))))
	r.Mount("/analytics/v1/", s.analyticsHandlers())
	return r
}

func (s *Server) healtzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
