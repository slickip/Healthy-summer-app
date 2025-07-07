package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// тот же секрет, что и в LoginHandler
const JWTSecret = "OMGMYKEY"

// ключ, под которым кладём userID в контекст
type ctxKey string

const ContextUserIDKey = ctxKey("userID")

// ParseToken достаёт и верифицирует JWT, возвращает userID
func ParseToken(r *http.Request) (uint, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return 0, errors.New("Missing Authorization header")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("Invalid Authorization header")
	}
	tkn, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(JWTSecret), nil
	})
	if err != nil || !tkn.Valid {
		return 0, errors.New("Invalid token")
	}
	claims := tkn.Claims.(jwt.MapClaims)
	var userID uint
	switch v := claims["user_id"].(type) {
	case float64:
		userID = uint(v)
	case string:
		uid64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, errors.New("Invalid user_id claim")
		}
		userID = uint(uid64)
	default:
		return 0, errors.New("Invalid user_id claim")
	}

}

// JWTAuth — middleware, кладёт userID в контекст или возвращает 401
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := ParseToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
