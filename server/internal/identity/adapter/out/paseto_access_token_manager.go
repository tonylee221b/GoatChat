package out

import (
	"errors"
	"log/slog"
	"time"

	"GoatChat/GoatChat/internal/identity/application/port"
	"GoatChat/GoatChat/internal/identity/domain"

	"aidanwoods.dev/go-paseto"
)

const (
	ClaimKeySessionID = "sid"
	ClaimKeyTokenType = "type"
)

type PasetoAccessTokenManager struct {
	secretKey paseto.V4AsymmetricSecretKey
	publicKey paseto.V4AsymmetricPublicKey

	issuer   string
	audience string
	ttl      time.Duration
}

func NewPasetoAccessTokenManager(
	sk paseto.V4AsymmetricSecretKey,
	issuer, audience string,
	ttl time.Duration,
) (*PasetoAccessTokenManager, error) {
	if issuer == "" {
		slog.Error(domain.ErrTokenIssuerRequired)
		return nil, errors.New(domain.ErrInvalidConfig)
	}

	if audience == "" {
		slog.Error(domain.ErrTokenAudienceRequired)
		return nil, errors.New(domain.ErrInvalidConfig)
	}

	if ttl <= 0 {
		slog.Error(domain.ErrTokenTTLNotPositive)
		return nil, errors.New(domain.ErrInvalidConfig)
	}

	return &PasetoAccessTokenManager{
		secretKey: sk,
		publicKey: sk.Public(),
		issuer:    issuer,
		audience:  audience,
		ttl:       ttl,
	}, nil
}

func (p *PasetoAccessTokenManager) Issue(userID, sessionID string) (string, error) {
	if userID == "" {
		slog.Info(domain.ErrTokenUserIDRequired)
		return "", errors.New(domain.ErrTokenUserIDRequired)
	}

	if sessionID == "" {
		slog.Info(domain.ErrTokenSessionIDRequired)
		return "", errors.New(domain.ErrTokenSessionIDRequired)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(p.ttl)

	token := paseto.NewToken()

	token.SetSubject(userID)
	token.SetIssuer(p.issuer)
	token.SetAudience(p.audience)
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(expiresAt)

	token.SetString(ClaimKeySessionID, sessionID)
	token.SetString(ClaimKeyTokenType, "access")

	rawToken := token.V4Sign(p.secretKey, nil)

	return rawToken, nil
}

func (p *PasetoAccessTokenManager) Verify(rawToken string) (port.AccessTokenClaims, error) {
	if rawToken == "" {
		slog.Info(domain.ErrInvalidAccessToken)
		return port.AccessTokenClaims{}, errors.New(domain.ErrInvalidCredentials)
	}

	parser := paseto.NewParserForValidNow()
	parser.AddRule(
		paseto.IssuedBy(p.issuer),
		paseto.ForAudience(p.audience),
	)

	token, err := parser.ParseV4Public(p.publicKey, rawToken, nil)
	if err != nil {
		slog.Info("failed to parse access token")
		return port.AccessTokenClaims{}, errors.New(domain.ErrInvalidCredentials)
	}

	userID, err := token.GetSubject()
	if err != nil || userID == "" {
		slog.Info(domain.ErrTokenMissingClaim, "claim field", "subject (user id)")
		return port.AccessTokenClaims{}, errors.New(domain.ErrInvalidCredentials)
	}

	sessionID, err := token.GetString(ClaimKeySessionID)
	if err != nil || sessionID == "" {
		slog.Info(domain.ErrTokenMissingClaim, "claim field", "sid (session id)")
		return port.AccessTokenClaims{}, errors.New(domain.ErrInvalidCredentials)
	}

	tokenType, err := token.GetString(ClaimKeyTokenType)
	if err != nil {
		if tokenType != "access" {
			slog.Info("token type is not access")
		}
		return port.AccessTokenClaims{}, errors.New(domain.ErrInvalidCredentials)
	}

	expiresAt, err := token.GetExpiration()
	if err != nil {
		slog.Info(domain.ErrTokenMissingClaim, "claim field", "exp")
		return port.AccessTokenClaims{}, errors.New(domain.ErrInvalidCredentials)
	}

	slog.Debug("access token verification complete")
	return port.AccessTokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		ExpiresAt: expiresAt,
	}, nil
}
