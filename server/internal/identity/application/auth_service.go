package application

import (
	"context"
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
	pwHasher           port.PasswordHasher
	accessTokenManager port.AccessTokenManager
	refreshTokenTTL    time.Duration
}

func NewAuthService(
	tx dbtx.Tx,
	userRepo port.UserRepository,
	pwHasher port.PasswordHasher,
	accessTokenManager port.AccessTokenManager,
	refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		tx,
		userRepo,
		pwHasher,
		accessTokenManager,
		refreshTokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, username domain.Username, plainPW string) (TokenPair, error) {
	var tp TokenPair

	a.tx.WithinTx(ctx, func(ctx context.Context) error {
		ufd, err := a.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return err
		}

		err = a.pwHasher.Compare(ufd.PasswordHash.Value, plainPW)
		if err != nil {
			return err
		}

		generatedRT, err := generateRefreshToken()
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		refreshExpiresAt := now.Add(a.refreshTokenTTL)

		session, err := domain.NewSession(
			ufd.ID,
			generatedRT.Hash,
			now,
			refreshExpiresAt,
		)
		if err != nil {
			return err
		}

		return nil
	})

	return tp, nil
}
