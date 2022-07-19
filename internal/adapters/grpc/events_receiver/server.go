package grpc_events_receiver

import (
	epb "analytics/internal/adapters/grpc/events_receiver/event_pb"
	"analytics/internal/ports"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type EventWriterServer struct {
	analytics ports.AnalyticsPort
	epb.UnimplementedEventWriterServer
}

var grpcServer *grpc.Server

func New(analytics ports.AnalyticsPort) *EventWriterServer {
	return &EventWriterServer{analytics: analytics}
}

func (s *EventWriterServer) Start() error {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		return fmt.Errorf("failed to listen on port 9000: %v", err)
	}

	eventWriterServer := s
	grpcServer = grpc.NewServer()
	epb.RegisterEventWriterServer(grpcServer, eventWriterServer)

	if err := grpcServer.Serve(listen); err != nil {
		return fmt.Errorf("failed to server gRPC server over port 9000: %v", err)
	}

	return nil
}

func (s *EventWriterServer) Stop() {
	grpcServer.GracefulStop()
}
