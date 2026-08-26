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
	tx       dbtx.Tx
	repo     port.UserRepository
	pwHasher port.PasswordHasher
}

func NewUserService(tx dbtx.Tx, userRepo port.UserRepository, pwHasher port.PasswordHasher) *UserService {
	return &UserService{tx, userRepo, pwHasher}
}

func (u *UserService) Register(ctx context.Context, username domain.Username, plainPW string) error {
	isExist, err := u.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return err
	}

	if isExist {
		slog.Info("user already exists", "username", username)
		return errors.New(domain.ErrUserAlreadyExist)
	}

	hashedPW, err := u.pwHasher.Hash(plainPW)
	if err != nil {
		return err
	}

	pwHash, err := domain.NewPasswordHash(hashedPW)
	if err != nil {
		return err
	}

	user := domain.NewUser(username, pwHash)
	return u.repo.Save(ctx, *user)
}

func (u *UserService) FindByUsername(ctx context.Context, un domain.Username) (domain.User, error) {
	ufd, err := u.repo.FindByUsername(ctx, un)
	if err != nil {
		return domain.User{}, err
	}

	return *ufd, nil
}

func (u *UserService) UpdateContact(ctx context.Context, username domain.Username, contact domain.Contact) (domain.User, error) {
	var user domain.User

	err := u.tx.WithinTx(ctx, func(ctx context.Context) error {
		ufd, err := u.repo.FindByUsername(ctx, username)
		if err != nil {
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
