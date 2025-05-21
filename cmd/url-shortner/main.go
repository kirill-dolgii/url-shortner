package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	// "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kirill-dolgii/url-shortner/internal/config/appconfig"
	"github.com/kirill-dolgii/url-shortner/internal/config/dbconfig"
	"github.com/kirill-dolgii/url-shortner/internal/http-server/handlers/redirect"
	"github.com/kirill-dolgii/url-shortner/internal/http-server/handlers/save"
	mwLogger "github.com/kirill-dolgii/url-shortner/internal/http-server/middleware/logger"
	"github.com/kirill-dolgii/url-shortner/internal/lib/logger/sl"
	"github.com/kirill-dolgii/url-shortner/internal/storage/postgres"
	"github.com/phsym/console-slog"
)

const (
	envDev   = "develop"
	envLocal = "local"
	envProd  = "production"
)

func main() {
	// init config
	config := appconfig.MustLoad()
	// logger
	logger, err := setupLogger(config.Name)
	if err != nil {
		log.Fatal("failed to initialize logger")
	}
	logger.Info("logger initialized")

	// db config
	dbConfig, err := dbconfig.LoadDBConfig()
	if err != nil {
		logger.Error("db config load failed", sl.Err(err))
	}
	logger.Info("db config initialized")

	st, err := postgres.InitDB(dbConfig)
	if err != nil {
		logger.Error("db connection failed", sl.Err(err))
	}
	logger.Info("db connection initialized")

	if err != nil {
		logger.Error("error occurred", sl.Err(err))
	}

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, st))
	router.Post("/{alias}", redirect.New(logger, st))
	logger.Info("starting server", slog.String("address", config.HttpConfig.Address))

	// start server
	srv := &http.Server{
		Addr:         config.HttpConfig.Address,
		Handler:      router,
		ReadTimeout:  config.HttpConfig.Timeout,
		WriteTimeout: config.HttpConfig.Timeout,
		IdleTimeout:  config.HttpConfig.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server start failed", sl.Err(err))
	}
}

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger
	switch env {
	case envDev:
		logger = slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}))
	case envLocal:
		logger = slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		logger = slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelWarn}))
	}
	slog.SetDefault(logger)
	return logger, nil
}
