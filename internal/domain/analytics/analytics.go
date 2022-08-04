package analytics

import (
	"analytics/internal/domain/models"
	"analytics/internal/ports"
	"context"
	"fmt"
)

type Service struct {
	db ports.AnalyticsStoragePort
}

func New(db ports.AnalyticsStoragePort) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateTask(ctx context.Context, event models.Event) error {
	err := s.db.CreateTask(ctx, event)

	if err != nil {
		return fmt.Errorf("create task failed: %v", err)
	}
	return nil
}

func (s *Service) AddMail(ctx context.Context, event models.Event) error {
	err := s.db.AddMail(ctx, event)

	if err != nil {
		return fmt.Errorf("add mail failed: %v", err)
	}
	return nil
}

func (s *Service) AddApproveClick(ctx context.Context, event models.Event) error {
	err := s.db.AddApproveClick(ctx, event)

	if err != nil {
		return fmt.Errorf("add approve click failed: %v", err)
	}
	return nil
}

func (s *Service) AddRejectClick(ctx context.Context, event models.Event) error {
	err := s.db.AddRejectClick(ctx, event)

	if err != nil {
		return fmt.Errorf("add reject click failed: %v", err)
	}
	return nil
}

func (s *Service) GetTotalTaskResponseTime(ctx context.Context, taskId int32) (string, error) {
	totalTaskResponseTime, err := s.db.GetTotalTaskResponseTime(ctx, taskId)

	if err != nil {
		return "", fmt.Errorf("get total task %d response time failed: %v", taskId, err)
	}
	return totalTaskResponseTime, nil
}

func (s *Service) GetApprovedTasksCount(ctx context.Context) (int32, error) {
	approvedTaskCount, err := s.db.GetTasksCount(ctx, "approved")

	if err != nil {
		return -1, fmt.Errorf("get approved tasks count: %v", err)
	}
	return approvedTaskCount, nil
}

func (s *Service) GetRejectedTasksCount(ctx context.Context) (int32, error) {
	approvedTaskCount, err := s.db.GetTasksCount(ctx, "rejected")

	if err != nil {
		return -1, fmt.Errorf("get rejected tasks count: %v", err)
	}
	return approvedTaskCount, nil
}
