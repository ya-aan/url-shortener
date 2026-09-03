package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ya-aan/url-shortener/internal/storage/postgres"
)

func main() {
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

	fmt.Println("connected to postgres")

	_ = storage
}
