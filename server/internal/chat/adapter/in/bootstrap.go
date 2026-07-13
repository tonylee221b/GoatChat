package in

import (
	"GoatChat/GoatChat/internal/chat/adapter/out"
	chatsqlc "GoatChat/GoatChat/internal/chat/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/chat/application"

	"github.com/go-chi/chi/v5"
)

func BoostrapChat(r chi.Router, db chatsqlc.DBTX) chi.Router {
	repo := out.NewChatPgRepository(db)
	crSvc := application.NewChatService(repo)
	h := NewChatHandler(*crSvc)
	router := NewChatRouter(h)
	router.route(r)

	return r
}
