package in

import (
	"github.com/go-chi/chi/v5"
)

const (
	UserRouteGroup = "/users"
)

type IdentityRouter struct {
	h *IdentityHandler
}

func NewIdentityRouter(h *IdentityHandler) *IdentityRouter {
	ir := &IdentityRouter{h: h}

	return ir
}

func (ir *IdentityRouter) route(r *chi.Mux) {
	r.Route(UserRouteGroup, func(r chi.Router) {
		r.Post("/register", ir.h.Register)
		r.Get("/{username}", ir.h.FindByUsername)
	})
}
