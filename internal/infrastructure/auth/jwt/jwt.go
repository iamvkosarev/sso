package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamvkosarev/sso/internal/domain/entity"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

type TokenClaims struct {
	UserID int64     `json:"user_id"`
	Exp    time.Time `json:"exp"`
}

func NewToken(user entity.User, secret string, tokenTTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}

	token.Claims = claims
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (TokenClaims, error) {
	var tc TokenClaims
	token, err := jwt.Parse(
		tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return tc, entity.ErrTokenExpired
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return tc, entity.ErrTokenIsInvalid
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return tc, entity.ErrTokenIsInvalid
		}
		return tc, fmt.Errorf("can't parse token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return tc, entity.ErrTokenIsInvalid
	}
	if uid, ok := claims["user_id"].(float64); ok {
		tc.UserID = int64(uid)
	}
	if exp, ok := claims["exp"].(float64); ok {
		tc.Exp = time.Unix(int64(exp), 0)
	}
	return tc, nil
}

func GetTokenFormContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", entity.ErrNoMetadata
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", entity.ErrNoAuthHeader
	}

	parts := strings.SplitN(values[0], " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", entity.ErrInvalidAuthHeader
	}

	return parts[1], nil
}
