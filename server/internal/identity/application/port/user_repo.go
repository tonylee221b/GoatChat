package port

import (
	"context"

	"GoatChat/GoatChat/internal/identity/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	Update(ctx context.Context, user domain.User) error
	ExistsByUsername(ctx context.Context, username domain.Username) (bool, error)
	FindByUsername(ctx context.Context, username domain.Username) (*domain.User, error)
}
