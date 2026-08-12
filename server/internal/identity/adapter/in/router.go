package in

import (
	"github.com/go-chi/chi/v5"
)

const (
	UserRouteGroup = "/users"
	AuthRouteGroup = "/auth"
)

type IdentityRouter struct {
	h *IdentityHandler
}

func NewIdentityRouter(h *IdentityHandler) *IdentityRouter {
	ir := &IdentityRouter{h: h}

	return ir
}

func (ir *IdentityRouter) route(r chi.Router) {
	r.Route(UserRouteGroup, func(r chi.Router) {
		r.Post("/", ir.h.Register)
		r.Get("/{username}", ir.h.FindByUsername)
		r.Put("/{username}", ir.h.UpdateContact)
	})

	r.Route(AuthRouteGroup, func(r chi.Router) {
		r.Post("/login", ir.h.Login)
		r.Post("/logout", ir.h.Logout)
		r.Post("/refresh", ir.h.Refresh)
	})
}
