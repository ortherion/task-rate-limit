package http

import (
	"github.com/go-chi/chi/v5"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
	"gitlab.com/g6834/team17/task-service/internal/utils"
	"net/http"
	"strconv"
)

func (s *Server) taskHandlers() http.Handler {
	h := chi.NewMux()

	h.Use(s.ValidateAuth())

	h.Group(func(r chi.Router) {
		h.Post("/", s.CreateTask)
		h.Delete("/{id}", s.DeleteTask)
		h.Put("/approve/{id}", s.Approve)
		h.Put("/reject/{id}", s.Reject)
	})

	return h
}

func (s *Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	var taskDTO *models.TaskDTO
	err := utils.ReadJson(r, &taskDTO)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorBadRequest(err))
		return
	}

	taskID, err := s.task.Create(ctx, taskDTO)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	s.presenter.JSON(w, r, nil)
	s.logger.Info().Msgf("task: %d created\r\n", taskID)
}

func (s *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	ID := chi.URLParam(r, "id")

	taskID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	err = s.task.Delete(ctx, taskID)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	s.presenter.JSON(w, r, nil)
	s.logger.Info().Msgf("task: %d deleted\r\n", taskID)
}

func (s *Server) Approve(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	ID := chi.URLParam(r, "id")

	taskID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	err = s.task.Sign(ctx, taskID)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	err = s.task.CheckTaskStatus(ctx, taskID)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	s.presenter.JSON(w, r, nil)
	s.logger.Info().Msgf("task: %d approved \r\n")
}

func (s *Server) Reject(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	ID := chi.URLParam(r, "id")

	taskID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}

	err = s.task.Reject(ctx, taskID)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}
	err = s.task.CheckTaskStatus(ctx, taskID)
	if err != nil {
		s.presenter.Error(w, r, models.ErrorInternal(err))
		return
	}
	s.presenter.JSON(w, r, nil)
	s.logger.Info().Msgf("task: %d rejected\r\n", taskID)
}
