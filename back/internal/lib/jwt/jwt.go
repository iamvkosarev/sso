package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamvkosarev/sso/back/internal/model"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrorFailedToExtractToken = errors.New("failed to extract token")
	ErrTokenExpired           = errors.New("token is expired")
)

type TokenClaims struct {
	UserID int64     `json:"user_id"`
	Exp    time.Time `json:"exp"`
}

func NewToken(user model.User, app model.App) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(app.TokenTTL).Unix(),
	}

	token.Claims = claims
	return token.SignedString([]byte(app.Secret))
}

func ParseToken(tokenString string, app model.App) (TokenClaims, error) {
	var tc TokenClaims
	token, err := jwt.Parse(
		tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(app.Secret), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return tc, ErrTokenExpired
		}

		return tc, fmt.Errorf("can't parse token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return tc, ErrInvalidToken
	}
	if uid, ok := claims["user_id"].(float64); ok {
		tc.UserID = int64(uid)
	}
	if exp, ok := claims["exp"].(float64); ok {
		tc.Exp = time.Unix(int64(exp), 0)
	}
	return tc, nil
}

func GetToken(r *http.Request) (string, error) {
	token := tokenFromHeader(r)
	if token != "" {
		return token, nil
	}
	token = tokenFromCookie(r)
	if token != "" {
		return token, nil
	}
	return "", ErrorFailedToExtractToken
}

func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
