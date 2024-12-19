package main

import (
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/back/internal/config"
	"github.com/iamvkosarev/sso/back/internal/http-server/router"
	"github.com/iamvkosarev/sso/back/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

func main() {
	cfg := config.MustLoad()
	log, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		fmt.Printf("error setting up logger: %v\n", err)
	}

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error initializing sqlite storage", sl.Err(err))
		return
	}

	log.Info("init server", slog.String("address", cfg.HTTPServerPort.Address))

	http.ListenAndServe(cfg.HTTPServerPort.Address, router.New(log, storage, cfg.App))
}
