package out

import (
	"context"
	"errors"
	"log/slog"
	"time"

	identitysqlc "GoatChat/GoatChat/internal/identity/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/identity/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionPgRepository struct {
	pool *pgxpool.Pool
}

func NewSessionPgRepository(pool *pgxpool.Pool) *SessionPgRepository {
	return &SessionPgRepository{pool: pool}
}

func (r *SessionPgRepository) Save(ctx context.Context, session domain.Session) error {
	q := r.queries(ctx)

	_, err := q.CreateSession(ctx, identitysqlc.CreateSessionParams{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash.Bytes(),
		ExpiresAt:        session.ExpiresAt,
		RevokedAt:        session.RevokedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info(domain.ErrSessionExpired)
			return errors.New(domain.ErrSessionExpired)
		}

		slog.Error("somethien went wrong with db", "error", err.Error())
		return errors.New(domain.ErrDB)
	}

	slog.Debug("sessions found", "session id", session.ID)
	return nil
}

func (r *SessionPgRepository) Revoke(ctx context.Context, sessionID uuid.UUID, revokedAt time.Time) error {
	q := r.queries(ctx)

	err := q.RevokeSession(ctx, identitysqlc.RevokeSessionParams{
		ID:        sessionID,
		RevokedAt: &revokedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info(domain.ErrSessionNotFound)
			return errors.New(domain.ErrSessionNotFound)
		}

		slog.Error("something went wrong with db", "error", err.Error())
		return err
	}

	return nil
}

func (r *SessionPgRepository) FindByRefreshTokenHash(ctx context.Context, rtHash domain.RefreshTokenHash) (*domain.Session, error) {
	q := r.queries(ctx)

	sfd, err := q.FindByRefreshTokenHash(ctx, identitysqlc.FindByRefreshTokenHashParams{
		RefreshTokenHash: rtHash.Bytes(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info(domain.ErrSessionNotFound)
			return nil, errors.New(domain.ErrSessionNotFound)
		}

		slog.Error("something went wrong with db", "error", err.Error())
		return nil, errors.New(domain.ErrDB)
	}

	return toSessionDomain(sfd), nil
}

func (r *SessionPgRepository) FindByRefreshTokenHashForUpdate(ctx context.Context, rtHash domain.RefreshTokenHash) (*domain.Session, error) {
	q := r.queries(ctx)

	sfd, err := q.FindByRefreshTokenHashForUpdate(ctx, identitysqlc.FindByRefreshTokenHashForUpdateParams{
		RefreshTokenHash: rtHash.Bytes(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info(domain.ErrSessionNotFound)
			return nil, errors.New(domain.ErrSessionNotFound)
		}

		slog.Error("something went wrong with db", "error", err.Error())
		return nil, errors.New(domain.ErrDB)
	}

	return toSessionDomain(sfd), nil
}

func (r *SessionPgRepository) UpdateRefreshTokenHash(ctx context.Context, sessionID uuid.UUID, newRTHash domain.RefreshTokenHash) error {
	q := r.queries(ctx)

	_, err := q.UpdateRefreshTokenHash(ctx, identitysqlc.UpdateRefreshTokenHashParams{
		ID:                  sessionID,
		NewRefreshTokenHash: newRTHash.Bytes(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info(domain.ErrSessionNotFound)
			return errors.New(domain.ErrSessionNotFound)
		}

		slog.Error("something went wrong with db", "error", err.Error())
		return errors.New(domain.ErrDB)
	}

	return nil
}

func (r *SessionPgRepository) queries(ctx context.Context) *identitysqlc.Queries {
	return identitysqlc.New(r.db(ctx))
}

func (r *SessionPgRepository) db(ctx context.Context) identitysqlc.DBTX {
	if pgxTx := dbtx.FromContext(ctx); pgxTx != nil {
		return pgxTx
	}

	return r.pool
}

func toSessionDomain(sfd identitysqlc.Session) *domain.Session {
	rrt, _ := domain.RestoreRefreshFromByte(sfd.RefreshTokenHash)

	return domain.RestoreSession(sfd.ID, sfd.UserID, rrt, sfd.ExpiresAt, sfd.RevokedAt, sfd.CreatedAt)
}
