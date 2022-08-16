package ports

import (
	"analytics/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type AnalyticsPort interface {
	CreateTask(ctx context.Context, event models.Event) error
	AddMail(ctx context.Context, event models.Event) error
	AddApproveClick(ctx context.Context, event models.Event) error
	AddRejectClick(ctx context.Context, event models.Event) error
	GetApprovedTasksCount(ctx context.Context) (int32, error)
	GetRejectedTasksCount(ctx context.Context) (int32, error)
	GetTotalTaskResponseTime(ctx context.Context, taskId int32) (string, error)
	CheckIdempotency(ctx context.Context, uuid uuid.UUID) (bool, error)
}
