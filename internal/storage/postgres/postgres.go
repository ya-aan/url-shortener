package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ya-aan/url-shortener/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, connectionString string) (*Storage, error) {
	db, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveURL(ctx context.Context, url, alias string) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		ctx,
		`INSERT INTO urls (url, alias) VALUES ($1, $2) RETURNING id`,
		url,
		alias,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save url: %w", err)
	}

	return id, nil
}

func (s *Storage) AliasExists(ctx context.Context, alias string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM urls WHERE alias = $1)`,
		alias,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check alias existence: %w", err)
	}
	return exists, nil
}

func (s *Storage) GetURLByAlias(ctx context.Context, alias string) (string, error) {
	var url string
	err := s.db.QueryRow(
		ctx,
		`SELECT url FROM urls WHERE alias = $1`,
		alias,
	).Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	if err != nil {

		return "", fmt.Errorf("failed to get url by alias: %w", err)
	}
	return url, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	result, err := s.db.Exec(
		ctx,
		`DELETE FROM urls WHERE alias = $1`,
		alias,
	)
	if err != nil {
		return fmt.Errorf("failed to delete url: %w", err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) UpdateURL(ctx context.Context, alias, newURL string) error {
	result, err := s.db.Exec(
		ctx,
		`UPDATE urls SET url = $1 WHERE alias = $2`,
		newURL,
		alias,
	)
	if err != nil {
		return fmt.Errorf("failed to update url: %w", err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
