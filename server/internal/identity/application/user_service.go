package application

import (
	"context"
	"errors"
	"log/slog"

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
		isExist, err := u.repo.ExistsByUsername(ctx, username)
		if err != nil {
			return err
		}

		if isExist {
			slog.Info("user already exists", "username", username)
			return errors.New(domain.ErrUserAlreadyExist)
		}

		user := domain.NewUser(username)

		if err = u.repo.Save(ctx, *user); err != nil {
			slog.Error("db error. user cannot be saved", "error", err)
			return err
		}

		return nil
	})
}

func (u *UserService) FindByUsername(ctx context.Context, un domain.Username) (domain.User, error) {
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
		slog.Error("tx error", "error", err)
		return domain.User{}, err
	}

	return *user, nil
}

func (u *UserService) UpdateContact(ctx context.Context, username domain.Username, contact domain.Contact) (domain.User, error) {
	var user domain.User

	err := u.tx.WithinTx(ctx, func(ctx context.Context) error {
		ufd, err := u.repo.FindByUsername(ctx, username)
		if err != nil {
			slog.Info("user not found", "username", username)
			return err
		}

		err = ufd.UpdateContact(contact)
		if err != nil {
			return err
		}

		err = u.repo.Update(ctx, *ufd)
		if err != nil {
			return err
		}

		user = *ufd

		return nil
	})
	if err != nil {
		slog.Error("tx error", "error", err)
		return domain.User{}, err
	}

	return user, nil
}
