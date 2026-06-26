package main

import (
	"net/http"
	"time"

	identityRouter "GoatChat/GoatChat/internal/identity/adapter/in"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/api/v1/health"))

	r.Get("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Route("/api/v1/", func(r chi.Router) {
		r.Mount("/identity", identityRouter.SetRoutes())
	})

	http.ListenAndServe(":8888", r)
}
