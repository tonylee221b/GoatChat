package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	chat_in "GoatChat/GoatChat/internal/chat/adapter/in"
	identity_in "GoatChat/GoatChat/internal/identity/adapter/in"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	APIVersion = "/api/v1"
)

func main() {
	ctx := context.Background()

	// TODO (mgyoo) : 추후 config파일로 분리
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	port := ":" + os.Getenv("PORT")

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

	r.Route(APIVersion, func(r chi.Router) {
		identity_in.BootstrapIdentity(r, pool)
		chat_in.BoostrapChat(r, pool)
	})

	slog.Info("GOAT Chat Server is running...", "PORT", port)
	http.ListenAndServe(port, r)
}
