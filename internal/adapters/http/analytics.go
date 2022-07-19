package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func (s *Server) analyticsHandlers() http.Handler {
	h := chi.NewMux()
	h.Route("/", func(r chi.Router) {
		h.Use(s.ValidateAuth)
		h.Get("/tasks/rejected", s.GetRejectedTasksCount)
		h.Get("/tasks/apporoved", s.GetApprovedTasksCount)
		h.Get("/task/totalresponsetime", s.GetTaskTotalResponseTime)
	})
	return h
}

func (s *Server) GetRejectedTasksCount(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "request", r)
	rejectedTasksCnt, err := s.analytics.GetRejectedTasksCount(ctx)
	if err != nil {
		s.logger.Error("get rejected tasks count failed", err)
		http.Error(w, "incorrect task id", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("%d", rejectedTasksCnt)))
}

func (s *Server) GetApprovedTasksCount(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "request", r)
	approvedTasksCnt, err := s.analytics.GetApprovedTasksCount(ctx)
	if err != nil {
		s.logger.Error("get approved tasks count failed", err)
		http.Error(w, "incorrect task id", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("%d", approvedTasksCnt)))
}

func (s *Server) GetTaskTotalResponseTime(w http.ResponseWriter, r *http.Request) {
	idRaw := r.URL.Query().Get("id")
	taskId, err := strconv.Atoi(idRaw)
	if err != nil {
		s.logger.Error("incorrect task id", err)
		http.Error(w, "incorrect task id", http.StatusBadRequest)
	}
	ctx := context.WithValue(context.Background(), "request", r)
	taskTotalResponseTime, err := s.analytics.GetTotalTaskResponseTime(ctx, int32(taskId))
	if err != nil {
		s.logger.Error(fmt.Sprintf("get total task %d response time failed", taskId))
		http.Error(w, "get total task response time failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(taskTotalResponseTime))
}
