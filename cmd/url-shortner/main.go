package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/kirill-dolgii/url-shortner/internal/config/appconfig"
	"github.com/kirill-dolgii/url-shortner/internal/config/dbconfig"
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

	_, err = postgres.InitDB(dbConfig)
	if err != nil {
		logger.Error("db connection failed", sl.Err(err))
	}
	logger.Info("db connection initialized")

	if err != nil {
		logger.Error("error occurred", sl.Err(err))
	}

	// start server
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
