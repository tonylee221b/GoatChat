package out

import (
	"errors"
	"log/slog"
	"time"

	"GoatChat/GoatChat/internal/identity/application/port"

	"aidanwoods.dev/go-paseto"
)

const (
	ErrIssuerRequired   = "issuer is required"
	ErrAudienceRequired = "audience is required"
	ErrTTLNotPositive   = "ttl must be positive"
	ErrClockRequired    = "clock is required"

	ErrUserIDRequired    = "user id is required"
	ErrSessionIDRequired = "session id is required"

	ErrInvalidAccessToken = "invalid access token"

	ErrMissingClaim = "claim is missing"
)

const (
	ClaimKeySessionID = "sid"
	ClaimKeyTokenType = "type"
)

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}

type PasetoAccessTokenManager struct {
	secretKey paseto.V4AsymmetricSecretKey
	publicKey paseto.V4AsymmetricPublicKey

	issuer   string
	audience string
	ttl      time.Duration
	clock    Clock
}

func NewPasetoAccessTokenManager(
	sk paseto.V4AsymmetricSecretKey,
	issuer, audience string,
	ttl time.Duration,
	clock Clock,
) (*PasetoAccessTokenManager, error) {
	if issuer == "" {
		slog.Error(ErrIssuerRequired)
		return nil, errors.New(ErrIssuerRequired)
	}

	if audience == "" {
		slog.Error(ErrAudienceRequired)
		return nil, errors.New(ErrAudienceRequired)
	}

	if ttl <= 0 {
		slog.Error(ErrAudienceRequired)
		return nil, errors.New(ErrTTLNotPositive)
	}

	if clock == nil {
		slog.Error(ErrClockRequired)
		return nil, errors.New(ErrClockRequired)
	}

	return &PasetoAccessTokenManager{
		secretKey: sk,
		publicKey: sk.Public(),
		issuer:    issuer,
		audience:  audience,
		ttl:       ttl,
		clock:     clock,
	}, nil
}

func (p *PasetoAccessTokenManager) Issue(userID, sessionID string) (string, error) {
	if userID == "" {
		slog.Info(ErrUserIDRequired)
		return "", errors.New(ErrUserIDRequired)
	}

	if sessionID == "" {
		slog.Info(ErrSessionIDRequired)
		return "", errors.New(ErrSessionIDRequired)
	}

	now := p.clock.Now()
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
		slog.Info(ErrInvalidAccessToken)
		return port.AccessTokenClaims{}, errors.New(ErrInvalidAccessToken)
	}

	parser := paseto.NewParserForValidNow()
	parser.AddRule(
		paseto.IssuedBy(p.issuer),
		paseto.ForAudience(p.audience),
	)

	token, err := parser.ParseV4Public(p.publicKey, rawToken, nil)
	if err != nil {
		slog.Info("failed to parse access token")
		return port.AccessTokenClaims{}, errors.New(ErrInvalidAccessToken)
	}

	userID, err := token.GetSubject()
	if err != nil || userID == "" {
		slog.Info(ErrMissingClaim, "claim field", "subject (user id)")
		return port.AccessTokenClaims{}, errors.New(ErrInvalidAccessToken)
	}

	sessionID, err := token.GetString(ClaimKeySessionID)
	if err != nil || sessionID == "" {
		slog.Info(ErrMissingClaim, "claim field", "sid (session id)")
		return port.AccessTokenClaims{}, errors.New(ErrInvalidAccessToken)
	}

	tokenType, err := token.GetString(ClaimKeyTokenType)
	if err != nil {
		if tokenType != "access" {
			slog.Info("token type is not access")
		}
		return port.AccessTokenClaims{}, errors.New(ErrInvalidAccessToken)
	}

	expiresAt, err := token.GetExpiration()
	if err != nil {
		slog.Info(ErrMissingClaim, "claim field", "exp")
		return port.AccessTokenClaims{}, errors.New(ErrInvalidAccessToken)
	}

	return port.AccessTokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		ExpiresAt: expiresAt,
	}, nil
}
