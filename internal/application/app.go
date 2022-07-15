package application

import (
	grpc_events_receiver "analytics/internal/adapters/grpc/events_receiver"
	"analytics/internal/adapters/postgres"
	"analytics/internal/domain/analytics"
	"context"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func Start(ctx context.Context) {
	logger, _ = zap.NewProduction()

	pgconn := "postgresql://app:secret@localhost:5432/app?sslmode=disable"

	db, err := postgres.New(ctx, pgconn)
	if err != nil {
		logger.Sugar().Fatalf("db init failed: %s", err)
	}

	analyticsS := analytics.New(db)

	totalTime, err := analyticsS.GetTotalTaskResponseTime(ctx, 321)
	if err != nil {
		logger.Sugar().Fatalf("cannot get total time: %s", err)
	}
	logger.Sugar().Infof("totalTime: %s", totalTime)

	approvedTasksCount, err := analyticsS.GetApprovedTasksCount(ctx)
	if err != nil {
		logger.Sugar().Fatalf("cannot get approvedTasksCount: %s", err)
	}

	logger.Sugar().Infof("approvedTasksCount: %d", approvedTasksCount)

	rejectedTasksCount, err := analyticsS.GetRejectedTasksCount(ctx)
	if err != nil {
		logger.Sugar().Fatalf("cannot get rejectedTasksCount: %s", err)
	}
	logger.Sugar().Infof("rejectedTasksCount: %d", rejectedTasksCount)

	grpcEventsReceiver := grpc_events_receiver.New(analyticsS)
	err = grpcEventsReceiver.Start()

	if err != nil {
		logger.Sugar().Fatalf("cannot start grpcEventsReceiver: %s", err)
	}

	logger.Sugar().Info("app has started")
}

func Stop() {
	logger.Sugar().Info("app has stopped")
}
