package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ya-aan/url-shortener/internal/service"
)

type mockUpdateNotFoundService struct{}
type mockUpdateService struct{}

func (m *mockUpdateNotFoundService) UpdateURL(
	ctx context.Context,
	alias string,
	newURL string,
) error {
	return service.ErrNotFound
}

func (m *mockUpdateService) UpdateURL(
	ctx context.Context,
	alias string,
	newURL string,
) error {
	return nil
}

func TestUpdate_Success(t *testing.T) {
	service := &mockUpdateService{}

	router := chi.NewRouter()
	router.Patch("/{alias}", Update(service))

	body := `{
		"url": "https://youtube.com"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/google",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	serviceMock := &mockUpdateNotFoundService{}

	router := chi.NewRouter()
	router.Patch("/{alias}", Update(serviceMock))

	body := `{"url":"https://youtube.com"}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/unknown",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdate_BadRequest(t *testing.T) {
	service := &mockUpdateService{}

	router := chi.NewRouter()
	router.Patch("/{alias}", Update(service))

	body := `{"url":"https://youtube.com",`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/google",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
