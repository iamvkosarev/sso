package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamvkosarev/go-shared-utils/middleware/logger"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/http-server/handlers/login"
	"github.com/iamvkosarev/sso/internal/http-server/handlers/register"
	"github.com/iamvkosarev/sso/internal/http-server/handlers/verify"
	"github.com/iamvkosarev/sso/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, storage *sqlite.Storage, app config.App) http.Handler {
	mux := chi.NewMux()

	mux.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.URLFormat,
		logger.NewLogger(log),
	)

	mux.Route(
		"/api/v1",
		func(r chi.Router) {
			r.Post("/sso/register", register.NewRegisterHandler(log, storage))
			r.Post("/sso/login", login.NewLoginHandler(log, storage, app))
			r.Get("/sso/verify", verify.NewVerifyHandler(log, app))
		},
	)

	return mux
}
