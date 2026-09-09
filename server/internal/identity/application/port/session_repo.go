package port

import (
	"context"
	"time"

	"GoatChat/GoatChat/internal/identity/domain"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Save(ctx context.Context, user domain.Session) error
	Revoke(ctx context.Context, sessionID uuid.UUID, revokedAt time.Time) error

	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)

	FindByRefreshTokenHash(ctx context.Context, rtHash domain.RefreshTokenHash) (*domain.Session, error)
	FindByRefreshTokenHashForUpdate(ctx context.Context, rtHash domain.RefreshTokenHash) (*domain.Session, error)
	UpdateRefreshTokenHash(ctx context.Context, sessionID uuid.UUID, newRTHash domain.RefreshTokenHash) error
}
