package grpc_events_receiver

import (
	epb "analytics/internal/adapters/grpc/events_receiver/event_pb"
	"analytics/internal/domain/models"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const layout string = "2006-01-02 15:04:05.999999999 -0700 MST"

func (s *EventWriterServer) WriteEvent(ctx context.Context, req *epb.EventParameters) (*epb.WriteStatus, error) {
	writeStatus := &epb.WriteStatus{}

	eventTime, err := time.Parse(layout, req.Time)
	if err != nil {
		writeStatus.Success = false
		return writeStatus, status.Error(codes.Internal, "parse event time failed")
	}

	event := models.Event{
		TaskId:          req.TaskId,
		Time:            eventTime,
		Type:            req.Type.String(),
		User:            req.User,
		ApproversNumber: req.ApproversNumber,
	}

	switch event.Type {
	case "create":
		err := s.analytics.CreateTask(ctx, event)
		if err != nil {
			writeStatus.Success = false
			return writeStatus, status.Error(codes.Internal, "unexpected error")
		}
	case "send_mail":
		err := s.analytics.AddMail(ctx, event)
		if err != nil {
			writeStatus.Success = false
			return writeStatus, status.Error(codes.Internal, "unexpected error")
		}
	case "approve":
		err := s.analytics.AddApproveClick(ctx, event)
		if err != nil {
			writeStatus.Success = false
			return writeStatus, status.Error(codes.Internal, "unexpected error")
		}
	case "reject":
		err := s.analytics.AddRejectClick(ctx, event)
		if err != nil {
			writeStatus.Success = false
			return writeStatus, status.Error(codes.Internal, "unexpected error")
		}
	default:
		writeStatus.Success = false
		return writeStatus, status.Error(codes.Internal, "incorrect event type")
	}

	writeStatus.Success = true
	return writeStatus, nil
}
