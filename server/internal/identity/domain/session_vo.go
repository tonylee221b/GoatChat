package domain

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log/slog"
)

const refreshTokenHashSize = sha256.Size

type RefreshTokenHash struct {
	Value [refreshTokenHashSize]byte
}

func NewRefreshTokenHash(plainToken string) (RefreshTokenHash, error) {
	if plainToken == "" {
		return RefreshTokenHash{}, errors.New(ErrInvalidRefreshTokenHash)
	}

	return RefreshTokenHash{Value: sha256.Sum256([]byte(plainToken))}, nil
}

func RestoreRefreshFromByte(raw []byte) (RefreshTokenHash, error) {
	if len(raw) != refreshTokenHashSize {
		slog.Warn("refresh token raw is empty")
		return RefreshTokenHash{}, errors.New(ErrInvalidRefreshTokenHash)
	}

	var r [refreshTokenHashSize]byte
	copy(r[:], raw)

	return RefreshTokenHash{Value: r}, nil
}

func RestoreRefreshFromHex(hexStr string) (RefreshTokenHash, error) {
	raw, err := hex.DecodeString(hexStr)
	if err != nil {
		slog.Warn("failed to decode refresh token hex")
		return RefreshTokenHash{}, errors.New(ErrInvalidRefreshTokenHash)
	}

	return RestoreRefreshFromByte(raw)
}

func (h RefreshTokenHash) Bytes() []byte {
	result := make([]byte, len(h.Value))
	copy(result, h.Value[:])

	return result
}

func (h RefreshTokenHash) String() string {
	return hex.EncodeToString(h.Value[:])
}

func (h RefreshTokenHash) Equals(other RefreshTokenHash) bool {
	return subtle.ConstantTimeCompare(h.Value[:], other.Value[:]) == 1
}

func (h RefreshTokenHash) IsZero() bool {
	return h == RefreshTokenHash{}
}
