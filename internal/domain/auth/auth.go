package auth

import (
	"analytics/internal/ports"
	"context"
)

type Service struct {
	authClient ports.AuthPort
}

func New(authClient ports.AuthPort) *Service {
	return &Service{authClient: authClient}
}

func (s *Service) Validate(ctx context.Context, accessToken string) (bool, error) {
	isAuthenticated, err := s.authClient.Validate(ctx, accessToken)
	if err != nil {
		return isAuthenticated, err
	}
	return isAuthenticated, nil
}
