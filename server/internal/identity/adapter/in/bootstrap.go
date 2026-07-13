package in

import (
	"GoatChat/GoatChat/internal/identity/adapter/out"
	"GoatChat/GoatChat/internal/identity/application"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BootstrapIdentity(r chi.Router, db *pgxpool.Pool) chi.Router {
	repo := out.NewUserPgRepository(db)
	tx := dbtx.NewPgxTx(db)
	svc := application.NewUserService(tx, repo)
	h := NewIdentityHandler(svc)
	router := NewIdentityRouter(h)
	router.route(r)

	return r
}
