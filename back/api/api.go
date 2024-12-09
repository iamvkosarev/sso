package api

import (
	"encoding/json"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/jwts"
	"github.com/iamvkosarev/go-shared-utils/logs"
	"github.com/iamvkosarev/sso/back/internal/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

type CORS interface {
	EnableCORS(next http.Handler, methods []string) http.Handler
}

type Deps struct {
	Storage    Storage
	Cors       CORS
	HttpLogger logs.HttpLogger
	JWT        jwts.JWT
}

type API struct {
	deps Deps
}

type Storage interface {
	AddUser(user model.User) error
	GetUser(email string) (*model.User, error)
}

func NewAPI(deps Deps) *API {
	return &API{
		deps: deps,
	}
}

func (api *API) RegisterHandler() http.Handler {
	return api.deps.Cors.EnableCORS(
		http.HandlerFunc(api.registerHandler),
		[]string{"POST"},
	)
}

func (api *API) registerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		api.deps.HttpLogger.Error(writer, "Not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		api.deps.HttpLogger.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.Match([]byte(body.Email)) {
		api.deps.HttpLogger.Error(writer, "Bad email", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		api.deps.HttpLogger.Error(writer, "Server error", http.StatusInternalServerError)
		return
	}

	err = api.deps.Storage.AddUser(
		model.User{
			Email:        body.Email,
			HashPassword: string(hashedPassword),
		},
	)

	if err != nil {
		api.deps.HttpLogger.Error(writer, "User already exists", http.StatusConflict)
		return
	}

	api.deps.HttpLogger.Success(
		writer,
		"User registered successfully",
		http.StatusCreated,
	)
}

func (api *API) LoginHandler() http.Handler {
	return api.deps.Cors.EnableCORS(
		http.HandlerFunc(api.loginHandler),
		[]string{"POST"},
	)
}

func (api *API) loginHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		api.deps.HttpLogger.Error(
			writer, fmt.Sprintf("'%v' not allowed", request.Method),
			http.StatusMethodNotAllowed,
		)
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		api.deps.HttpLogger.InternalError(writer, "Bad request", http.StatusBadRequest, err)
		return
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.Match([]byte(body.Email)) {
		api.deps.HttpLogger.Error(writer, "Bad email", http.StatusBadRequest)
		return
	}

	user, err := api.deps.Storage.GetUser(body.Email)
	if err != nil {
		api.deps.HttpLogger.InternalError(writer, "User not found", http.StatusNotFound, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(body.Password))
	if err != nil {
		api.deps.HttpLogger.InternalError(writer, "Server error", http.StatusInternalServerError, err)
		return
	}
	token, err := api.deps.JWT.GenerateJWT(user.ID)
	if err != nil {
		api.deps.HttpLogger.InternalError(writer, "Failed to generate token", http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(
		map[string]string{
			"token": token,
		},
	)
}
