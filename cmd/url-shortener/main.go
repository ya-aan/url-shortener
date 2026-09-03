package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/ya-aan/url-shortener/internal/http-server/handlers"
	"github.com/ya-aan/url-shortener/internal/service"
	"github.com/ya-aan/url-shortener/internal/storage/postgres"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	ctx := context.Background()

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	storage, err := postgres.New(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	urlService := service.New(storage)

	router := chi.NewRouter()

	router.Post("/url", handlers.Save(urlService))
	router.Get("/{alias}", handlers.Redirect(urlService))
	router.Delete("/{alias}", handlers.Delete(urlService))
	router.Patch("/{alias}", handlers.Update(urlService))

	log.Println("server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
