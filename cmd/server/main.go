package main

import (
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/http-server/router"
	"github.com/iamvkosarev/sso/internal/storage/sqlite"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	cfg := config.MustLoad()
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("error setting up logger: %v\n", err)
	}

	storage := sqlite.New(cfg.StoragePath, logger)
	if storage == nil {
		logger.Error("error initializing sqlite storage", sl.Err(err))
		return
	}

	logger.Info("init server", slog.String("address", cfg.HTTPServer.Address))

	http.ListenAndServe(cfg.HTTPServer.Address, router.New(logger, storage, cfg.App))
}
