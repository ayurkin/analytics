package auth_grpc_client

import (
	apb "analytics/internal/adapters/grpc/auth/auth_pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	client apb.AuthApiClient
}

func New(target string) (*GrpcClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("create grpc client connection failed: %v", err)
	}

	client := apb.NewAuthApiClient(conn)
	return &GrpcClient{client: client}, nil
}

func (s *GrpcClient) Validate(ctx context.Context, accessToken string) (bool, error) {
	var isAuthenticated bool
	response, err := s.client.Authenticate(ctx, &apb.AuthRequest{AccessToken: accessToken})
	if err != nil {
		return false, err
	}
	if response.Error != "" {
		return false, fmt.Errorf("validate access token failed: %v", err)
	}
	isAuthenticated = true
	return isAuthenticated, nil
}
