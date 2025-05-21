package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/kirill-dolgii/url-shortner/internal/lib/api/response"
	"github.com/kirill-dolgii/url-shortner/internal/lib/logger/sl"
	"github.com/kirill-dolgii/url-shortner/internal/storage"
)

var (
	ErrUrlExists = errors.New("url exists")
)

type SaveRequest struct {
	URL   string `json:"url" validate""require,url`
	Alias string `json:"alias,omitempty"`
}

type SaveResponse struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type UrlSaver interface {
	SaveURL(url, alias string) (int64, error)
}

func New(logger *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"
		logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SaveRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			valErr := err.(validator.ValidationErrors)
			logger.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(valErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			err = errors.New("empty alias")
			logger.Error("empty alias", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrUrlExists) {
			logger.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error(ErrUrlExists.Error()))
			return
		}

		if err != nil {
			logger.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		logger.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, SaveResponse{
			Response: response.OK(),
			Alias:    req.Alias,
		})
	}
}
