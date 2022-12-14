package ports

import (
	"analytics/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type AnalyticsStoragePort interface {
	CreateTask(ctx context.Context, event models.Event) error
	AddMail(ctx context.Context, event models.Event) error
	AddApproveClick(ctx context.Context, event models.Event) error
	AddRejectClick(ctx context.Context, event models.Event) error
	GetTotalTaskResponseTime(ctx context.Context, taskId int32) (string, error)
	GetTasksCount(ctx context.Context, taskType string) (int32, error)
	CheckIdempotency(ctx context.Context, uuid uuid.UUID) (bool, error)
}
