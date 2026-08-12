package in

import (
	"GoatChat/GoatChat/internal/identity/adapter/out"
	"GoatChat/GoatChat/internal/identity/application"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"

	"aidanwoods.dev/go-paseto"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BootstrapIdentity(r chi.Router, db *pgxpool.Pool) chi.Router {
	pwHasher, _ := out.NewBcryptHasher(7)
	userRepo := out.NewUserPgRepository(db)
	sessionRepo := out.NewSessionPgRepository(db)
	tx := dbtx.NewPgxTx(db)

	cfg := NewIdentityConfig()
	sk, _ := paseto.NewV4AsymmetricSecretKeyFromHex(cfg.AccessTokenSecretKey)
	accessTokenManager, _ := out.NewPasetoAccessTokenManager(
		sk,
		cfg.AccessTokenIssuer,
		cfg.AccessTokenAudience,
		cfg.AccessTokenTTL,
	)

	userSvc := application.NewUserService(tx, userRepo, pwHasher)
	authSvc := application.NewAuthService(tx, userRepo, sessionRepo, pwHasher, accessTokenManager, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	h := NewIdentityHandler(userSvc, authSvc)
	router := NewIdentityRouter(h)
	router.route(r)

	return r
}
