package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Not allowed", http.StatusMethodNotAllowed)
		fmt.Println("Not allowed: " + request.Method)
		return
	}
	var body struct {
		Email    string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, "Bad request", http.StatusBadRequest)
		fmt.Println("Bad request")
	}
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	http.Error(writer, "Server error", http.StatusInternalServerError)
	//	fmt.Println("Server error")
	//	return
	//}

	//_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", body.Email, hashedPassword)
	//if err != nil {
	//	http.Error(writer, "User already exists", http.StatusConflict)
	//	return
	//}

	writer.WriteHeader(http.StatusCreated)
	fmt.Fprint(writer, "User registered successfully")
}

func main() {
	http.Handle("/register", enableCORS(http.HandlerFunc(registerHandler)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
