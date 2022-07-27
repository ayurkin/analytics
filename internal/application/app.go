package application

import (
	auth_grpc_client "analytics/internal/adapters/grpc/auth"
	"analytics/internal/adapters/http"
	"analytics/internal/adapters/kafka/events_receiver"
	"analytics/internal/adapters/postgres"
	"analytics/internal/domain/analytics"
	"analytics/internal/domain/auth"
	"context"
	"go.uber.org/zap"
	"time"
)

var (
	logger              *zap.Logger
	quitCh              chan struct{}
	kafkaEventsReceiver *events_receiver.Client
	httpServer          *http.Server
)

func Start(ctx context.Context) {
	logger, _ = zap.NewProduction()

	pgconn := "postgresql://app:secret@localhost:5432/app?sslmode=disable"

	db, err := postgres.New(ctx, pgconn)
	if err != nil {
		logger.Sugar().Fatalf("db init failed: %v", err)
	}

	analyticsS := analytics.New(db)

	grpcTarget := "auth.team3.svc.cluster.local:9000"
	authGrpcClient, err := auth_grpc_client.New(grpcTarget)
	if err != nil {
		logger.Sugar().Fatalf("auth grpc client init failed: %v", err)
	}

	authS := auth.New(authGrpcClient)

	httpServer = http.New(analyticsS, authS, logger.Sugar())

	kafkaEventsReceiver = events_receiver.New(analyticsS, logger.Sugar())

	quitCh = make(chan struct{})
	go func() {
		kafkaEventsReceiver.Start(quitCh)
	}()
	go func() {
		err := httpServer.Start()
		if err != nil {
			logger.Sugar().Fatalf("http server failed: %v", err)
		}
	}()

	logger.Sugar().Info("application has started")

}

func Stop() {
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	err := httpServer.Stop(ctx)
	if err != nil {
		logger.Sugar().Errorf("stop http server failed: %v", err)
	}
	kafkaEventsReceiver.Stop(quitCh)
	logger.Sugar().Info("app has stopped")
}
