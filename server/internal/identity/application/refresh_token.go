package application

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"

	"GoatChat/GoatChat/internal/identity/domain"
)

const refreshTokenByteSize = 32

type GeneratedRefreshToken struct {
	Plain string
	Hash  domain.RefreshTokenHash
}

func generateRefreshToken() (GeneratedRefreshToken, error) {
	randomBytes := make([]byte, refreshTokenByteSize)

	if _, err := rand.Read(randomBytes); err != nil {
		slog.Error(domain.ErrGenerateRefreshToken, "error", err.Error())
		return GeneratedRefreshToken{}, errors.New(domain.ErrGenerateRefreshToken)
	}

	plain := base64.RawURLEncoding.EncodeToString(randomBytes)
	hash, err := domain.NewRefreshTokenHash(plain)
	if err != nil {
		return GeneratedRefreshToken{}, err
	}

	return GeneratedRefreshToken{
		Plain: plain,
		Hash:  hash,
	}, nil
}
