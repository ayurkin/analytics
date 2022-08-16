package events_receiver

import (
	"analytics/internal/domain/models"
	"analytics/internal/ports"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Client struct {
	Reader    *kafka.Reader
	analytics ports.AnalyticsPort
	logger    *zap.SugaredLogger
}

func New(analytics ports.AnalyticsPort, logger *zap.SugaredLogger) *Client {

	c := Client{}
	brokers := []string{"localhost:9092"}
	topic := "team3-topic-analytics"
	groupId := "analytics_group"

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})
	c.logger = logger
	c.analytics = analytics
	return &c
}

func (c *Client) Start(ctx context.Context) {

	for {
		msg, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.logger.Infof("kafka reader stopped by context: %v", err)

				return
			}

			c.logger.Infof("kafka reader failed: %v", err)

			continue
		}
		c.logger.Infof("fetched message: %v value: %s, ofset: %v", msg.Key, msg.Value, msg.Offset)
		err = c.fetchProcessCommit(ctx, msg)
		if err != nil {
			c.logger.Errorf("fetch process commit failed: %v", err)
		}

	}
}

func (c *Client) fetchProcessCommit(ctx context.Context, msg kafka.Message) error {
	var event models.Event
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		c.logger.Errorf("kafka message unmarshall failed: %v, value: %s", err, msg.Value)
		err = c.Reader.CommitMessages(ctx, msg)
		return err
	}

	if event.UUID == uuid.Nil {
		c.logger.Errorf("incorrect event uuid: %s", event.UUID.String())
		err = c.Reader.CommitMessages(ctx, msg)
		return err
	}

	isIdempotent, err := c.analytics.CheckIdempotency(ctx, event.UUID)
	if err != nil {
		return err
	}

	if !isIdempotent {
		c.logger.Infof("%s is not idempotent, will be skipped", event.UUID.String())
		err = c.Reader.CommitMessages(ctx, msg)
		return err
	}

	switch event.Type {
	case "create":
		err = c.analytics.CreateTask(ctx, event)
	case "send_mail":
		err = c.analytics.AddMail(ctx, event)
	case "approve":
		err = c.analytics.AddApproveClick(ctx, event)
	case "reject":
		err = c.analytics.AddRejectClick(ctx, event)
	default:
		c.logger.Errorf("incorrect event type: %s", event.Type)
	}
	if err != nil {
		return err
	}

	err = c.Reader.CommitMessages(ctx, msg)
	return err

}

func (c *Client) Stop() {
	err := c.Reader.Close()
	if err != nil {
		c.logger.Errorf("kafka reader close failed: %v", err)
	}
	c.logger.Info("kafka reader closed")
}
