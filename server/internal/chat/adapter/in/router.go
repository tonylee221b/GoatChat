package in

import (
	"github.com/go-chi/chi/v5"
)

type ChatRouter struct {
	h *ChatHandler
}

const (
	ChatRouteGroup     = "/chat"
	ChatroomRouteGruop = "/chatroom"
)

func NewChatRouter(h *ChatHandler) *ChatRouter {
	return &ChatRouter{h}
}

func (cr *ChatRouter) route(r *chi.Mux) {
	r.Route(ChatRouteGroup, func(r chi.Router) {
		r.Post(string(ChatroomRouteGruop), cr.h.CreateChatroom)
	})
}
