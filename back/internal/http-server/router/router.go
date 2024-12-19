package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamvkosarev/go-shared-utils/middleware/logger"
	"github.com/iamvkosarev/sso/back/internal/http-server/handlers/login"
	"github.com/iamvkosarev/sso/back/internal/http-server/handlers/register"
	"github.com/iamvkosarev/sso/back/internal/http-server/handlers/verify"
	"github.com/iamvkosarev/sso/back/internal/model"
	"github.com/iamvkosarev/sso/back/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, storage *sqlite.Storage, app model.App) http.Handler {
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.URLFormat)
	mux.Use(logger.NewLogger(log))

	mux.Group(
		func(r chi.Router) {
			r.Post("/api/sso/register", register.NewRegisterHandler(log, storage))
			r.Post("/api/sso/login", login.NewLoginHandler(log, storage, app))
			r.Post("/api/sso/verify", verify.NewVerifyHandler(log, app))
		},
	)

	return mux
}
