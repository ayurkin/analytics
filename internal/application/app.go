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

type App struct {
	logger              *zap.Logger
	quitCh              chan struct{}
	kafkaEventsReceiver *events_receiver.Client
	httpServer          *http.Server
}

func Start(ctx context.Context, app *App) {
	logger, _ := zap.NewProduction()
	app.logger = logger
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

	app.httpServer = http.New(analyticsS, authS, logger.Sugar())

	app.kafkaEventsReceiver = events_receiver.New(analyticsS, logger.Sugar())

	go func() {
		app.kafkaEventsReceiver.Start(ctx)
	}()
	go func() {
		err := app.httpServer.Start()
		if err != nil {
			logger.Sugar().Fatalf("http server failed: %v", err)
		}
	}()

	logger.Sugar().Info("application has started")

}

func Stop(app *App) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.httpServer.Stop(ctx)
	if err != nil {
		app.logger.Sugar().Errorf("stop http server failed: %v", err)
	}
	app.kafkaEventsReceiver.Stop()
	app.logger.Sugar().Info("app has stopped")
}
