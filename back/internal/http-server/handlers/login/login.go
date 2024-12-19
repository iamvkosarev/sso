package login

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/iamvkosarev/go-shared-utils/api/response"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/back/internal/lib/jwt"
	"github.com/iamvkosarev/sso/back/internal/model"
	"github.com/iamvkosarev/sso/back/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

type UserProvider interface {
	GetUser(email string) (model.User, error)
}

type SecretProvider interface {
	GetSecret() (string, error)
}

func NewLoginHandler(log *slog.Logger, userProvider UserProvider, app model.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.login.NewLoginHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode body"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			var validatorErr validator.ValidationErrors
			errors.As(err, &validatorErr)

			log.Error("failed to validate body", err)
			render.JSON(w, r, resp.ValidateErrors(validatorErr))
			return
		}

		user, err := userProvider.GetUser(req.Email)
		if errors.Is(err, storage.ErrorUserNotFound) {
			log.Error("user not exists", sl.Err(err))
			render.JSON(w, r, resp.Error("user not exists"))
			return
		}
		if err != nil {
			log.Error("failed to get user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get user"))
			return
		}
		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(req.Password)); err != nil {
			log.Error("invalid credentials", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid credentials"))
			return
		}

		token, err := jwt.NewToken(user, app)
		if err != nil {
			log.Error("failed to create token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to create token"))
			return
		}

		log.Info(
			"success authorization",
			slog.String("email", user.Email),
		)
		responseOK(w, r, token)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, token string) {
	render.JSON(
		w, r, Response{
			Response: resp.Ok(),
			Token:    token,
		},
	)
}
