package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/ya-aan/url-shortener/internal/service"
)

type URLService interface {
	CreateURL(ctx context.Context, rawURL, alias string) (int64, string, error)
}

type SaveRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type SaveResponse struct {
	ID    int64  `json:"id"`
	Alias string `json:"alias"`
}

func Save(urlService URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SaveRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": "invalid json",
			})
			return
		}

		id, alias, err := urlService.CreateURL(
			r.Context(),
			req.URL,
			req.Alias,
		)

		if errors.Is(err, service.ErrInvalidURL) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": "invalid url",
			})
			return
		}

		if errors.Is(err, service.ErrInvalidAlias) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": "invalid alias",
			})
			return
		}

		if errors.Is(err, service.ErrAliasExists) {
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, map[string]string{
				"error": "alias already exists",
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

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, SaveResponse{
			ID:    id,
			Alias: alias,
		})
	}
}
