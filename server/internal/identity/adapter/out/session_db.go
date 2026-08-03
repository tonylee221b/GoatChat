package out

import (
	"context"

	"GoatChat/GoatChat/internal/identity/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionPgRepository struct {
	pool *pgxpool.Pool
}

func NewSessionPgRepository(pool *pgxpool.Pool) *SessionPgRepository {
	return &SessionPgRepository{pool: pool}
}

func (r *SessionPgRepository) Save(ctx context.Context, session domain.Session) error {
	return nil
}
