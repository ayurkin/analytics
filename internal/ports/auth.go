package ports

import (
	"context"
)

type AuthPort interface {
	Validate(ctx context.Context, accessToken string) (bool, error)
}
