package verify

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/iamvkosarev/go-shared-utils/api/response"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/internal/config"
	"github.com/iamvkosarev/sso/internal/lib/jwt"
	"log/slog"
	"net/http"
	"time"
)

type Response struct {
	resp.Response
	UserID int64 `json:"user_id"`
}

func NewVerifyHandler(log *slog.Logger, app config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.verify.NewVerifyHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		token, err := jwt.GetToken(r)
		if err != nil {
			log.Error("failed to extract token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to extract token"))
			return
		}
		tokenClaims, err := jwt.ParseToken(token, app)
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Error("token is expired", sl.Err(err))
			render.JSON(w, r, resp.Error(resp.ErrorTokenExpired.Error()))
			return
		}
		if err != nil {
			log.Error("failed to parse token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse token"))
			return
		}
		if tokenClaims.Exp.Before(time.Now()) {
			log.Error("token is expired", sl.Err(err))
			render.JSON(w, r, resp.Error(resp.ErrorTokenExpired.Error()))
			return
		}

		responseOK(w, r, tokenClaims.UserID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(
		w, r, Response{
			UserID:   id,
			Response: resp.Ok(),
		},
	)
}
