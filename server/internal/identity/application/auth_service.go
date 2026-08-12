package application

import (
	"context"
	"log/slog"
	"time"

	"GoatChat/GoatChat/internal/identity/application/port"
	"GoatChat/GoatChat/internal/identity/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"
)

type TokenPair struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type AuthService struct {
	tx                 dbtx.Tx
	userRepo           port.UserRepository
	sessionRepo        port.SessionRepository
	pwHasher           port.PasswordHasher
	accessTokenManager port.AccessTokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(
	tx dbtx.Tx,
	userRepo port.UserRepository,
	sessionRepo port.SessionRepository,
	pwHasher port.PasswordHasher,
	accessTokenManager port.AccessTokenManager,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		tx,
		userRepo,
		sessionRepo,
		pwHasher,
		accessTokenManager,
		accessTokenTTL,
		refreshTokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, username domain.Username, plainPW string) (TokenPair, error) {
	ufd, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return TokenPair{}, err
	}

	err = a.pwHasher.Compare(ufd.PasswordHash.Value, plainPW)
	if err != nil {
		return TokenPair{}, err
	}

	generatedRT, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	now := time.Now().UTC()
	refreshExpiresAt := now.UTC().Add(a.refreshTokenTTL)

	session, err := domain.NewSession(
		ufd.ID,
		generatedRT.Hash,
		now,
		refreshExpiresAt,
	)
	if err != nil {
		slog.Info("failed to create a new session", "error", err.Error())
		return TokenPair{}, err
	}

	issuedAccessToken, err := a.accessTokenManager.Issue(ufd.ID.String(), session.ID.String())
	if err != nil {
		return TokenPair{}, nil
	}

	if err := a.sessionRepo.Save(ctx, *session); err != nil {
		// TODO: Revoke Access Token created above
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:           issuedAccessToken,
		AccessTokenExpiresAt:  now.Add(a.accessTokenTTL),
		RefreshToken:          generatedRT.Plain,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, plainRT string) error {
	refreshTokenHash, err := domain.NewRefreshTokenHash(plainRT)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()

	err = a.tx.WithinTx(ctx, func(ctx context.Context) error {
		session, err := a.sessionRepo.FindByRefreshTokenHashForUpdate(ctx, refreshTokenHash)
		if err != nil {
			return err
		}

		session.Revoke(now)

		if err := a.sessionRepo.Revoke(ctx, session.ID, now); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Refresh(ctx context.Context, plainRT string) (TokenPair, error) {
	currentHash, err := domain.NewRefreshTokenHash(plainRT)
	if err != nil {
		return TokenPair{}, err
	}

	newRT, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	now := time.Now().UTC()

	var (
		accessToken      string
		sessionExpiresAt time.Time
	)

	err = a.tx.WithinTx(ctx, func(ctx context.Context) error {
		session, err := a.sessionRepo.FindByRefreshTokenHash(ctx, currentHash)
		if err != nil {
			return err
		}

		if err := session.RotateRefreshToken(newRT.Hash, now); err != nil {
			return err
		}

		accessToken, err = a.accessTokenManager.Issue(session.UserID.String(), session.ID.String())
		if err != nil {
			return err
		}

		if err := a.sessionRepo.UpdateRefreshTokenHash(ctx, session.ID, session.RefreshTokenHash); err != nil {
			return err
		}

		sessionExpiresAt = session.ExpiresAt

		return nil
	})
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  now.Add(a.accessTokenTTL),
		RefreshToken:          newRT.Plain,
		RefreshTokenExpiresAt: sessionExpiresAt,
	}, nil
}
