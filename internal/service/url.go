package service

import (
	"context"
	"errors"
	"math/rand"
	"net/url"
)

var ErrAliasExists = errors.New("alias already exists")
var ErrInvalidURL = errors.New("invalid url")

type URLStorage interface {
	SaveURL(ctx context.Context, url, alias string) (int64, error)
	AliasExists(ctx context.Context, alias string) (bool, error)
	GetURLByAlias(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
	UpdateURL(ctx context.Context, alias, newURL string) error
}

type Service struct {
	storage URLStorage
}

func New(storage URLStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateURL(
	ctx context.Context,
	rawURL string,
	alias string,
) (int64, error) {

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return 0, ErrInvalidURL
	}

	if alias == "" {
		alias = generateAlias(6)
	}

	exists, err := s.storage.AliasExists(ctx, alias)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrAliasExists
	}

	id, err := s.storage.SaveURL(ctx, rawURL, alias)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetURL(ctx context.Context, alias string) (string, error) {
	url, err := s.storage.GetURLByAlias(ctx, alias)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *Service) DeleteURL(ctx context.Context, alias string) error {
	return s.storage.DeleteURL(ctx, alias)
}

func (s *Service) UpdateURL(
	ctx context.Context,
	alias string,
	newURL string,
) error {
	return s.storage.UpdateURL(ctx, alias, newURL)
}

func generateAlias(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)

	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
