package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kirill-dolgii/url-shortner/internal/domain/models"
	"github.com/kirill-dolgii/url-shortner/internal/lib/api/response"
	"github.com/kirill-dolgii/url-shortner/internal/lib/logger/sl"
	"github.com/kirill-dolgii/url-shortner/internal/storage"
)

var (
	ErrEmptyAlias  = errors.New("empty alias")
	ErrUrlNotFound = errors.New("url not found")
)

type UrlProvider interface {
	GetUrlByAlias(alias string) (models.Url, error)
}

func New(logger *slog.Logger, urlProvider UrlProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Error("empty alias", sl.Err(ErrEmptyAlias))
			render.JSON(w, r, response.Error("empty alias"))
			return
		}

		url, err := urlProvider.GetUrlByAlias(alias)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				logger.Error("url not found", sl.Err(ErrUrlNotFound))
				render.JSON(w, r, response.Error("not found"))
				return
			}
			logger.Error("url lookup failed", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		http.Redirect(w, r, url.FullUrl, http.StatusFound)
	}
}
