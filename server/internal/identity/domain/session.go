package domain

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Session struct is an Aggregate Root
//
// it manages user auth session
type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash RefreshTokenHash
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
}

func NewSession(
	userID uuid.UUID,
	rth RefreshTokenHash,
	now time.Time,
	expiresAt time.Time,
) (*Session, error) {
	if userID == uuid.Nil {
		return nil, errors.New(ErrInvalidUserID)
	}

	if rth.IsZero() {
		return nil, errors.New(ErrInvalidRefreshTokenHash)
	}

	if !expiresAt.After(now) {
		return nil, errors.New(ErrSessionExpired)
	}

	return &Session{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshTokenHash: rth,
		ExpiresAt:        expiresAt.UTC(),
		RevokedAt:        nil,
		CreatedAt:        now.UTC(),
	}, nil
}

func RestoreSession(
	id uuid.UUID,
	userID uuid.UUID,
	rth RefreshTokenHash,
	expiresAt time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
) *Session {
	return &Session{
		ID:               id,
		UserID:           userID,
		RefreshTokenHash: rth,
		ExpiresAt:        expiresAt,
		RevokedAt:        revokedAt,
		CreatedAt:        createdAt,
	}
}

func (s Session) CanRefresh(now time.Time) error {
	if s.RevokedAt != nil {
		slog.Info(ErrSessionRevoked)
		return errors.New(ErrSessionRevoked)
	}

	if !now.Before(s.ExpiresAt) {
		slog.Info(ErrSessionExpired)
		return errors.New(ErrSessionExpired)
	}

	return nil
}

func (s *Session) RotateRefreshToken(newHash RefreshTokenHash, now time.Time) error {
	if err := s.CanRefresh(now); err != nil {
		return err
	}

	if newHash.IsZero() {
		slog.Info("refresh token is empty")
		return errors.New(ErrInvalidRefreshTokenHash)
	}

	s.RefreshTokenHash = newHash

	return nil
}

func (s *Session) Revoke(now time.Time) {
	if s.RevokedAt != nil {
		slog.Info("session is already revoked")
		return
	}

	s.RevokedAt = &now
}
