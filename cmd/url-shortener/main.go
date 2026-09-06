package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/ya-aan/url-shortener/internal/config"
	"github.com/ya-aan/url-shortener/internal/http-server/handlers"
	"github.com/ya-aan/url-shortener/internal/http-server/middleware"
	"github.com/ya-aan/url-shortener/internal/service"
	"github.com/ya-aan/url-shortener/internal/storage/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	_ = godotenv.Load()
	cfg := config.MustLoad()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	storage, err := postgres.New(ctx, connectionString)
	if err != nil {
		return err
	}
	defer storage.Close()

	urlService := service.New(storage)

	router := chi.NewRouter()
	username := os.Getenv("HTTP_SERVER_USER")
	password := os.Getenv("HTTP_SERVER_PASSWORD")

	router.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth(username, password))

		r.Post("/url", handlers.Save(urlService))
		r.Delete("/{alias}", handlers.Delete(urlService))
		r.Patch("/{alias}", handlers.Update(urlService))
	})
	router.Get("/{alias}", handlers.Redirect(urlService))

	server := &http.Server{
		Addr:              cfg.HTTPServer.Address,
		Handler:           router,
		ReadTimeout:       cfg.HTTPServer.Timeout,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	return serve(ctx, server)
}

func serve(ctx context.Context, server *http.Server) error {
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("failed to shut down HTTP server: %w", err)
		}

		if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
