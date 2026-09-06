package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ya-aan/url-shortener/internal/service"
)

type RedirectService interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func Redirect(urlService RedirectService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")

		url, err := urlService.GetURL(r.Context(), alias)

		if errors.Is(err, service.ErrInvalidAlias) {
			http.Error(w, "invalid alias", http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
