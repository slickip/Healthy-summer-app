package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// В контексте будем хранить userID под этим ключом
type contextKey string

const (
	JWT_SECRET       = "OMGMYKEY" // замени на свой ключ
	ContextUserIDKey = contextKey("userID")
)

// JWT Middleware
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		log.Println("Auth header:", authHeader)
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(JWT_SECRET), nil
		})
		if err != nil {
			log.Printf("[JWTAuth] parse error: %v", err)
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Printf("[JWTAuth] token invalid, claims=%+v", token.Claims)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("[JWTAuth] cannot cast claims: %+v", token.Claims)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Универсальная проверка user_id
		var userID uint

		switch v := claims["user_id"].(type) {
		case float64:
			userID = uint(v)
		case string:
			uid64, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				log.Printf("Claims: %+v\n", claims)
				return
			}
			userID = uint(uid64)
		default:
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			log.Printf("Claims: %+v\n", claims)
			return
		}

		// Добавляем userID в контекст
		ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
		log.Printf("user_id raw: %+v\n", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
