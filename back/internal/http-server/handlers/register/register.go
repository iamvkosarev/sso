package register

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/iamvkosarev/go-shared-utils/api/response"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/sso/back/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	UserID int64 `json:"user_id"`
}

type UserSaver interface {
	SaveUser(email string, hashPassword string) (int64, error)
}

func NewRegisterHandler(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.register.NewRegisterHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("empty request body")
			render.JSON(w, r, resp.Error("empty request body"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.String("email", req.Email))

		if err := validator.New().Struct(req); err != nil {
			var validErr validator.ValidationErrors
			errors.As(err, &validErr)

			log.Error("failed to validate request", sl.Err(err))
			render.JSON(w, r, resp.ValidateErrors(validErr))
			return
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to hash password"))
			return
		}

		id, err := userSaver.SaveUser(req.Email, string(passHash))

		if errors.Is(err, storage.ErrorUserExists) {
			log.Error("user already exists", sl.Err(err))
			render.JSON(w, r, resp.Error("user already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save user"))
			return
		}
		log.Info("new user registered", slog.String("email", req.Email), slog.Int64("user_id", id))
		responseOK(w, r, id)
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
