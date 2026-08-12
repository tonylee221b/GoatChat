package out

import (
	"GoatChat/GoatChat/internal/identity/domain"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) (*BcryptHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		slog.Error("invalid bcrypt cost", "cost", cost)
		return nil, errors.New(domain.ErrInvalidConfig)
	}

	return &BcryptHasher{cost}, nil
}

func (b *BcryptHasher) Hash(plainPW string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(plainPW), b.cost)
	if err != nil {
		slog.Error("failed to hash password", "error", err.Error())
		return "", errors.New(domain.ErrInvalidConfig)
	}

	return string(hashedPW), nil
}

func (b *BcryptHasher) Compare(hashedPW, plainPW string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(plainPW))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			slog.Info("invalid password")
			return errors.New(domain.ErrInvalidCredentials)
		}

		slog.Error("compare password failed", "error", err.Error())
		return errors.New(domain.ErrInvalidConfig)
	}

	return nil
}
