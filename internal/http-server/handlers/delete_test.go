package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ya-aan/url-shortener/internal/service"
)

type mockDeleteNotFoundService struct{}
type mockDeleteService struct{}

func (m *mockDeleteNotFoundService) DeleteURL(
	ctx context.Context,
	alias string,
) error {
	return service.ErrNotFound
}

func (m *mockDeleteService) DeleteURL(
	ctx context.Context,
	alias string,
) error {
	return nil
}

func TestDelete_Success(t *testing.T) {
	service := &mockDeleteService{}

	router := chi.NewRouter()
	router.Delete("/{alias}", Delete(service))

	req := httptest.NewRequest(
		http.MethodDelete,
		"/google",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			rec.Code,
		)
	}
}

func TestDelete_NotFound(t *testing.T) {
	serviceMock := &mockDeleteNotFoundService{}

	router := chi.NewRouter()
	router.Delete("/{alias}", Delete(serviceMock))

	req := httptest.NewRequest(
		http.MethodDelete,
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
