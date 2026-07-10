package main

import (
	"GoatChat/GoatChat/internal/chat/adapter/in"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"GoatChat/GoatChat/internal/identity/adapter/in"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const APIVersion = "/api/v1"

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

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
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	// TODO (mgyoo) : 추후 config파일로 분리
	dbConn, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	r2 := in.BoostrapChat(dbConn)
	r.Mount("/", r2)
	ir := in.BootstrapIdentity(pool)
	r.Mount(APIVersion, ir)

	http.ListenAndServe(":8888", r)
}
