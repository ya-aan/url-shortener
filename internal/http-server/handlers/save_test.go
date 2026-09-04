package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ya-aan/url-shortener/internal/service"
)

type mockURLService struct{}

type mockConflictService struct{}

func (m *mockConflictService) CreateURL(
	ctx context.Context,
	rawURL string,
	alias string,
) (int64, string, error) {
	return 0, "", service.ErrAliasExists
}

func (m *mockURLService) CreateURL(
	ctx context.Context,
	rawURL string,
	alias string,
) (int64, string, error) {
	return 1, "google", nil
}

func TestSave_Success(t *testing.T) {
	service := &mockURLService{}
	handler := Save(service)

	body := `{
		"url": "https://google.com",
		"alias": "google"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/url",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestSave_BadRequest(t *testing.T) {
	service := &mockURLService{}
	handler := Save(service)

	body := `{"url":"https://google.com",`

	req := httptest.NewRequest(
		http.MethodPost,
		"/url",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestSave_Conflict(t *testing.T) {
	serviceMock := &mockConflictService{}
	handler := Save(serviceMock)

	body := `{
		"url": "https://google.com",
		"alias": "google"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/url",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusConflict,
			rec.Code,
		)
	}
}
