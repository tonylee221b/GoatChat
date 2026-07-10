package in

import (
	"GoatChat/GoatChat/internal/chat/adapter/out"
	chatsqlc "GoatChat/GoatChat/internal/chat/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/chat/application/service"

	"github.com/go-chi/chi/v5"
)

func BoostrapChat(db chatsqlc.DBTX) chi.Router {
	r := chi.NewRouter()

	repo := out.NewChatPgRepository(db)
	crSvc := service.NewChatService(repo)
	h := NewChatHandler(*crSvc)
	router := NewChatRouter(h)
	router.route(r)

	return r
}
