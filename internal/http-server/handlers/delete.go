package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ya-aan/url-shortener/internal/service"
)

type DeleteService interface {
	DeleteURL(ctx context.Context, alias string) error
}

func Delete(urlService DeleteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		err := urlService.DeleteURL(r.Context(), alias)
		if errors.Is(err, service.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error": "url not found",
			})
			return
		}

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": "internal error",
			})
			return
		}

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": "internal error",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
