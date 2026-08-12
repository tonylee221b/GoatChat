package in

import (
	"github.com/go-chi/chi/v5"
)

type ChatRouter struct {
	h *ChatHandler
}

const (
	ChatRouteGroup     = "/chat"
	ChatroomRouteGroup = "/chatroom"
)

func NewChatRouter(h *ChatHandler) *ChatRouter {
	return &ChatRouter{h}
}

func (cr *ChatRouter) route(r chi.Router) {
	r.Route(ChatRouteGroup, func(r chi.Router) {
		r.Post(string(ChatroomRouteGroup), cr.h.CreateChatroom)
		r.Patch(string(ChatroomRouteGroup), cr.h.UpdateChatroom)
		r.Delete(string(ChatroomRouteGroup), cr.h.DeleteChatroom)
	})
}
