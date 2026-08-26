package in

import (
	"os"
	"time"
)

type identityConfig struct {
	AccessTokenSecretKey string
	AccessTokenIssuer    string
	AccessTokenAudience  string
	AccessTokenTTL       time.Duration // seconds
	RefreshTokenTTL      time.Duration // seconds
}

func NewIdentityConfig() *identityConfig {
	atTTL := os.Getenv("ACCESS_TOKEN_TTL_IN_SECONDS")
	atTTLDuration, _ := time.ParseDuration(atTTL)
	rtTTL := os.Getenv("REFRESH_TOKEN_TTL_IN_SECONDS")
	rtTTLDuration, _ := time.ParseDuration(rtTTL)

	return &identityConfig{
		AccessTokenSecretKey: os.Getenv("SECRET_KEY"),
		AccessTokenIssuer:    os.Getenv("ISSUER"),
		AccessTokenAudience:  os.Getenv("AUDIENCE"),
		AccessTokenTTL:       atTTLDuration,
		RefreshTokenTTL:      rtTTLDuration,
	}
}
