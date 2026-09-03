package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/ya-aan/url-shortener/internal/service"
)

type UpdateService interface {
	UpdateURL(ctx context.Context, alias, newURL string) error
}

type UpdateRequest struct {
	URL string `json:"url"`
}

func Update(urlService UpdateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")

		var req UpdateRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": "invalid json",
			})
			return
		}

		err := urlService.UpdateURL(r.Context(), alias, req.URL)

		if errors.Is(err, service.ErrInvalidURL) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": "invalid url",
			})
			return
		}

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

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{
			"status": "updated",
		})
	}
}
