package in

import "github.com/go-chi/chi/v5"

func SetRoutes(h *IdentityHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", h.Signup)
	r.Post("/login", h.Login)

	return r
}
