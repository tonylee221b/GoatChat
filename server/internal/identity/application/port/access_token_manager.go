package port

import (
	"time"
)

type AccessTokenClaims struct {
	UserID    string
	SessionID string
	ExpiresAt time.Time
}

type AccessTokenManager interface {
	Issue(userID, sessionID string) (string, error)
	Verify(rawToken string) (AccessTokenClaims, error)
}
