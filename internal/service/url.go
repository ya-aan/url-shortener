package service

import (
	"context"
	"errors"
	"math/rand"
	"net/url"

	"regexp"

	"github.com/ya-aan/url-shortener/internal/storage"
)

var ErrAliasExists = errors.New("alias already exists")
var ErrInvalidURL = errors.New("invalid url")
var ErrNotFound = errors.New("url not found")

var ErrInvalidAlias = errors.New("invalid alias")

var aliasPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

type URLStorage interface {
	SaveURL(ctx context.Context, url, alias string) (int64, error)
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
) (int64, string, error) {

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return 0, "", ErrInvalidURL
	}

	if alias == "" {
		alias = generateAlias(6)
	}

	if !aliasPattern.MatchString(alias) {
		return 0, "", ErrInvalidAlias
	}
	id, err := s.storage.SaveURL(ctx, rawURL, alias)
	if errors.Is(err, storage.ErrAliasExists) {
		return 0, "", ErrAliasExists
	}
	if err != nil {
		return 0, "", err
	}

	return id, alias, nil

}

func (s *Service) GetURL(ctx context.Context, alias string) (string, error) {
	if !aliasPattern.MatchString(alias) {
		return "", ErrInvalidAlias
	}

	url, err := s.storage.GetURLByAlias(ctx, alias)
	if errors.Is(err, storage.ErrNotFound) {
		return "", ErrNotFound
	}

	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *Service) DeleteURL(ctx context.Context, alias string) error {
	if !aliasPattern.MatchString(alias) {
		return ErrInvalidAlias
	}

	err := s.storage.DeleteURL(ctx, alias)

	if errors.Is(err, storage.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func (s *Service) UpdateURL(
	ctx context.Context,
	alias string,
	newURL string,
) error {
	if !aliasPattern.MatchString(alias) {
		return ErrInvalidAlias
	}

	parsedURL, err := url.ParseRequestURI(newURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrInvalidURL
	}

	err = s.storage.UpdateURL(ctx, alias, newURL)

	if errors.Is(err, storage.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func generateAlias(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)

	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
