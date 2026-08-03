package port

import (
	"context"

	"GoatChat/GoatChat/internal/identity/domain"
)

type SessionRepository interface {
	Save(ctx context.Context, user domain.Session) error
}
