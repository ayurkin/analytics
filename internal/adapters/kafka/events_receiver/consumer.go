package events_receiver

import (
	"analytics/internal/domain/models"
	"analytics/internal/ports"
	"context"
	"encoding/json"
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
	topic := "test"
	groupId := "test_group_id"

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

func (c *Client) Start(quit chan struct{}) {

	for {
		err := c.fetchProcessCommit()
		if err != nil {
			c.logger.Errorf("kafka fetch process commit failed: %v", err)
		}
		select {
		case <-quit:
			return
		default:
		}
	}
}

func (c *Client) fetchProcessCommit() error {
	ctx := context.Background()
	msg, err := c.Reader.FetchMessage(ctx)
	if err != nil {
		c.logger.Fatalf("kafka reader failed to fetch message: %v", err)
	}
	c.logger.Infof("fetched message: %v value: %s, ofset: %v", msg.Key, msg.Value, msg.Offset)

	var event models.Event
	err = json.Unmarshal(msg.Value, &event)
	if err != nil {
		c.logger.Errorf("kafka message unmarshall failed: %v, value: %s", err, msg.Value)
		err = c.Reader.CommitMessages(ctx, msg)
		return err
	}

	if event.UUID.String() == "00000000-0000-0000-0000-000000000000" {
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
		err := c.analytics.CreateTask(ctx, event)
		if err != nil {
			return err
		}
	case "send_mail":
		err := c.analytics.AddMail(ctx, event)
		if err != nil {
			return err
		}
	case "approve":
		err := c.analytics.AddApproveClick(ctx, event)
		if err != nil {
			return err
		}
	case "reject":
		err := c.analytics.AddRejectClick(ctx, event)
		if err != nil {
			return err
		}
	default:
		c.logger.Errorf("incorrect event type: %s", event.Type)
	}
	err = c.Reader.CommitMessages(ctx, msg)
	return err

}

func (c *Client) Stop(quit chan struct{}) {
	close(quit)
}
