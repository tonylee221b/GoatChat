package port

import (
	"context"

	"GoatChat/GoatChat/internal/identity/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	FindByUsername(ctx context.Context, username domain.Username) (*domain.User, error)
}
