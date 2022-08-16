package auth

import (
	"analytics/internal/ports"
	"context"
)

type service struct {
	authClient ports.AuthPort
}

func New(authClient ports.AuthPort) ports.AuthPort {
	return &service{authClient: authClient}
}

func (s *service) Validate(ctx context.Context, accessToken string) (bool, error) {
	isAuthenticated, err := s.authClient.Validate(ctx, accessToken)
	if err != nil {
		return isAuthenticated, err
	}
	return isAuthenticated, nil
}
