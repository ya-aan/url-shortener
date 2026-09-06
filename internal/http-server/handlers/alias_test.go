package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ya-aan/url-shortener/internal/service"
)

func TestHandlers_InvalidAlias(t *testing.T) {
	svc := service.New(nil)
	router := chi.NewRouter()
	router.Post("/url", Save(svc))
	router.Get("/{alias}", Redirect(svc))
	router.Patch("/{alias}", Update(svc))
	router.Delete("/{alias}", Delete(svc))

	for _, tc := range []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodPost, "/url", `{"url":"https://example.com","alias":"bad alias"}`},
		{http.MethodGet, "/bad%20alias", ""},
		{http.MethodPatch, "/bad%20alias", `{"url":"https://example.com"}`},
		{http.MethodDelete, "/bad%20alias", ""},
	} {
		t.Run(tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "invalid alias") {
				t.Fatalf("got status %d, body %s; want 400 invalid alias", rec.Code, rec.Body)
			}
		})
	}
}
