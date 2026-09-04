package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ya-aan/url-shortener/internal/service"
)

type mockRedirectService struct{}
type mockRedirectNotFoundService struct{}

func (m *mockRedirectNotFoundService) GetURL(
	ctx context.Context,
	alias string,
) (string, error) {
	return "", service.ErrNotFound
}

func (m *mockRedirectService) GetURL(
	ctx context.Context,
	alias string,
) (string, error) {
	return "https://google.com", nil
}

func TestRedirect_Success(t *testing.T) {
	service := &mockRedirectService{}
	router := chi.NewRouter()
	router.Get("/{alias}", Redirect(service))

	req := httptest.NewRequest(
		http.MethodGet,
		"/google",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusFound,
			rec.Code,
		)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	serviceMock := &mockRedirectNotFoundService{}

	router := chi.NewRouter()
	router.Get("/{alias}", Redirect(serviceMock))

	req := httptest.NewRequest(
		http.MethodGet,
		"/unknown",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}
