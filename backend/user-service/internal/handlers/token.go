package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func extractUserIDFromToken(r *http.Request) (uint, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return 0, errors.New("Missing Authorization header")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("Invalid Authorization header")
	}

	tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTSecret, nil
	})
	if err != nil || !tok.Valid {
		return 0, errors.New("Invalid token")
	}

	claims := tok.Claims.(jwt.MapClaims)
	// user_id теперь число
	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("Invalid user_id claim")
	}
	return uint(uidFloat), nil
}
