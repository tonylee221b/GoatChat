package in

import (
	"GoatChat/GoatChat/internal/chat/adapter/out"
	"GoatChat/GoatChat/internal/chat/application"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BoostrapChat(r chi.Router, db *pgxpool.Pool) chi.Router {
	repo := out.NewChatPgRepository(db)
	tx := dbtx.NewPgxTx(db)
	crSvc := application.NewChatService(tx, repo)
	h := NewChatHandler(*crSvc)
	router := NewChatRouter(h)
	router.route(r)

	return r
}
