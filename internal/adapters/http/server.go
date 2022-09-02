package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/config"
	"gitlab.com/g6834/team17/task-service/internal/ports"
	"net/http"
	"time"
)

type Server struct {
	auth      ports.Auth
	task      ports.Task
	presenter ports.Presenter
	logger    *zerolog.Logger
	srv       *http.Server
}

func New(cfg *config.Config, logger *zerolog.Logger, auth ports.Auth, task ports.Task, presenter ports.Presenter) (*Server, error) {
	s := new(Server)
	s.auth = auth
	s.task = task
	s.presenter = presenter
	s.logger = logger
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%v:%v", cfg.Rest.Host, cfg.Rest.Port),
		Handler: s.routes(),
	}

	return s, nil
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", s.healthzHandler)
	r.Mount("/task", s.taskHandlers())

	return r
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
