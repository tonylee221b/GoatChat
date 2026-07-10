package application

import (
	"context"

	"GoatChat/GoatChat/internal/identity/application/port"
	"GoatChat/GoatChat/internal/identity/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"
)

type UserService struct {
	tx   dbtx.Tx
	repo port.UserRepository
}

func NewUserService(tx dbtx.Tx, repo port.UserRepository) *UserService {
	return &UserService{tx, repo}
}

func (u *UserService) Register(ctx context.Context, username domain.Username) error {
	return u.tx.WithinTx(ctx, func(ctx context.Context) error {
		user, err := domain.NewUser(username)
		if err != nil {
			return err
		}

		if err = u.repo.Save(ctx, *user); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserService) FindByUsername(ctx context.Context, un domain.Username) (*domain.User, error) {
	var user *domain.User

	err := u.tx.WithinTx(ctx, func(ctx context.Context) error {
		ufd, err := u.repo.FindByUsername(ctx, un)
		if err != nil {
			return err
		}

		user = ufd

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
