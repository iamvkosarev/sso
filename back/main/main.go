package main

import (
	"github.com/iamvkosarev/go-shared-utils/cors"
	"github.com/iamvkosarev/go-shared-utils/jwts"
	"github.com/iamvkosarev/go-shared-utils/logs"
	"github.com/iamvkosarev/sso/back/api"
	"github.com/iamvkosarev/sso/back/internal/storage/local_storage"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

const JWT_SECRET_KEY = "JWT_SECRET"

type API interface {
	RegisterHandler() http.Handler
	LoginHandler() http.Handler
}

func main() {
	godotenv.Load()
	var newAPI API = api.NewAPI(
		api.Deps{
			Storage:    local_storage.NewStorage(),
			Cors:       cors.NewCORS([]string{"https://kosarev.app", "http://localhost:63343"}),
			HttpLogger: logs.NewHttpLogger(true, true),
			JWT: jwts.NewJWT(
				os.Getenv(JWT_SECRET_KEY),
				time.Hour*24,
			),
		},
	)
	http.Handle("/register", newAPI.RegisterHandler())
	http.Handle("/login", newAPI.LoginHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
