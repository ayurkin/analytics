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
		h.Get("/tasks/approved", s.GetApprovedTasksCount)
		h.Get("/task/totalresponsetime", s.GetTaskTotalResponseTime)
	})
	return h
}

//GetRejectedTasksCount
//@ID GetRejectedTasksCount
//@tags analytics
//@Summary Get rejected tasks count
//@Description Get rejected tasks count if access token is valid. Token format: access=token
//@Security ApiKeyAuth
//@Success 200  {string}  1
//@Failure 401 {string} string "unauthorized"
//@Failure 500 {string} string "internal error"
//@Router /tasks/rejected [get]
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

//GetApprovedTasksCount
//@ID GetApprovedTasksCount
//@tags analytics
//@Summary Get approved tasks count
//@Description Get approved tasks count if access token is valid. Token format: access=token
//@Security ApiKeyAuth
//@Success 200  {string}  1
//@Failure 401 {string} string "unauthorized"
//@Failure 500 {string} string "internal error"
//@Router /tasks/approved [get]
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

//GetTaskTotalResponseTime
//@ID GetTaskTotalResponseTime
//@tags analytics
//@Summary Get task total response time
//@Description Get task total response time if access token is valid. Token format: access=token
//@Security ApiKeyAuth
//@Param   id   query  int  true  "Task ID"
//@Success 200 {string}  1
//@Failure 401 {string} string "unauthorized"
//@Failure 500 {string} string "internal error"
//@Router /task/totalresponsetime [get]
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
