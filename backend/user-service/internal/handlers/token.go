package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
)

func extractUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("Missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, errors.New("Invalid Authorization header")
	}
	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok || userIDStr == "" {
		return 0, errors.New("User ID not found in token")
	}

	uid64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("Invalid user ID")
	}
	return uint(uid64), nil
}
