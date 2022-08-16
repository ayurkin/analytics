//go:build approve
// +build approve

package tests

import (
	"analytics/internal/adapters/postgres"
	"analytics/internal/domain/analytics"
	"analytics/internal/domain/models"
	"analytics/internal/ports"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestApproveRun(t *testing.T) {
	suite.Run(t, new(ApproveSuite))
}

type ApproveSuite struct {
	suite.Suite
	pgContainer testcontainers.Container
	analytics   ports.AnalyticsPort
}

func (suite *ApproveSuite) SetupSuite() {
	ctx := context.Background()

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:11",
			ExposedPorts: []string{"5432"},
			Env: map[string]string{
				"POSTGRES_DB":       dbName,
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPass,
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
			SkipReaper: true,
			AutoRemove: true,
		},
		Started: true,
	})
	suite.Require().NoError(err)

	// with a second delay migrations work properly
	time.Sleep(time.Second * 10)

	ip, err := dbContainer.Host(ctx)
	suite.Require().NoError(err)
	port, err := dbContainer.MappedPort(ctx, "5432")
	suite.T().Log(fmt.Sprintf("Postgres container port: %v", port))
	suite.Require().NoError(err)

	cfg := &pgx.ConnConfig{
		Config: pgconn.Config{
			Host:     ip,
			Port:     uint16(port.Int()),
			Database: dbName,
			User:     dbUser,
			Password: dbPass,
		},
	}

	connString := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		"disable",
	)
	err = applyMigrations(connString)
	suite.T().Log("Migrations finished")
	suite.Require().NoError(err)

	db, err := postgres.New(ctx, connString)

	suite.Require().NoError(err)

	analyticsS := analytics.New(db)
	suite.analytics = analyticsS
	suite.pgContainer = dbContainer

	suite.T().Log("Suite setup is done")
	time.Sleep(time.Second * 5)
}

func (suite *ApproveSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(context.Background())
	if err != nil {
		suite.T().Error("Terminate container failed")
	}
	suite.T().Log("Suite stop is done")
}

func (suite *ApproveSuite) Test1CreateTask() {
	ctx := context.Background()

	eventTime, err := time.Parse(layout, "2022-02-02 15:00:00.000000 +0000 UTC")
	suite.Require().NoError(err)

	event := models.Event{
		UUID:            uuid.New(),
		TaskId:          1,
		Time:            eventTime,
		Type:            "create",
		User:            "author@mail.com",
		ApproversNumber: 2,
	}
	err = suite.analytics.CreateTask(ctx, event)
	suite.Require().NoError(err)
}

func (suite *ApproveSuite) Test2AddEmail1() {
	ctx := context.Background()

	eventTime, err := time.Parse(layout, "2022-02-02 15:00:05.000000 +0000 UTC")
	suite.Require().NoError(err)

	event := models.Event{
		UUID:            uuid.New(),
		TaskId:          1,
		Time:            eventTime,
		Type:            "send_mail",
		User:            "addressee1@mail.com",
		ApproversNumber: 0,
	}
	err = suite.analytics.AddMail(ctx, event)
	suite.Require().NoError(err)
}

func (suite *ApproveSuite) Test3AddApproveClick1() {
	ctx := context.Background()

	eventTime, err := time.Parse(layout, "2022-02-02 15:01:05.000000 +0000 UTC")
	suite.Require().NoError(err)

	event := models.Event{
		UUID:            uuid.New(),
		TaskId:          1,
		Time:            eventTime,
		Type:            "approve",
		User:            "addressee1@mail.com",
		ApproversNumber: 0,
	}
	err = suite.analytics.AddApproveClick(ctx, event)
	suite.Require().NoError(err)
}

func (suite *ApproveSuite) Test4AddEmail2() {
	ctx := context.Background()

	eventTime, err := time.Parse(layout, "2022-02-02 15:01:10.000000 +0000 UTC")
	suite.Require().NoError(err)

	event := models.Event{
		UUID:            uuid.New(),
		TaskId:          1,
		Time:            eventTime,
		Type:            "send_mail",
		User:            "addressee2@mail.com",
		ApproversNumber: 0,
	}
	err = suite.analytics.AddMail(ctx, event)
	suite.Require().NoError(err)
}

func (suite *ApproveSuite) Test5AddApproveClick2() {
	ctx := context.Background()

	eventTime, err := time.Parse(layout, "2022-02-02 15:02:10.000000 +0000 UTC")
	suite.Require().NoError(err)

	event := models.Event{
		UUID:            uuid.New(),
		TaskId:          1,
		Time:            eventTime,
		Type:            "approve",
		User:            "addressee2@mail.com",
		ApproversNumber: 0,
	}
	err = suite.analytics.AddApproveClick(ctx, event)
	suite.Require().NoError(err)
}

func (suite *ApproveSuite) Test6CheckApprovedTasksCount() {
	ctx := context.Background()

	approvedTasksCount, err := suite.analytics.GetApprovedTasksCount(ctx)
	suite.Require().NoError(err)

	a := assert.New(suite.T())
	a.Equal(int32(1), approvedTasksCount)
}

func (suite *ApproveSuite) Test7CheckRejectedTasksCount() {
	ctx := context.Background()

	rejectedTasksCount, err := suite.analytics.GetRejectedTasksCount(ctx)
	suite.Require().NoError(err)

	a := assert.New(suite.T())
	a.Equal(int32(0), rejectedTasksCount)
}

func (suite *ApproveSuite) Test8TotalTaskResponseTime() {
	ctx := context.Background()

	totalTaskResponseTime, err := suite.analytics.GetTotalTaskResponseTime(ctx, 1)
	suite.Require().NoError(err)

	a := assert.New(suite.T())
	a.Equal("2 minutes", totalTaskResponseTime)
}
