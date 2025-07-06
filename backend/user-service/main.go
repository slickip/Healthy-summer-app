package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/config"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/db"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/handlers"
)

func main() {
	// Загружаем конфиг
	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Address:     "0.0.0.0:8081",
			Timeout:     5 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
	}

	database := db.New()

	// Создаем структуру с зависимостями
	h := &handlers.Handler{
		DB: database,
	}

	// Роутер
	mux := http.NewServeMux()

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong from user-service")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Браузер отправляет preflight-запрос методом OPTIONS
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Передаем управление mux
			h.ServeHTTP(w, r)
		})
	}

	// Роуты с зависимостями
	mux.HandleFunc("/api/users/register", h.RegisterHandler)
	mux.HandleFunc("/api/users/login", h.LoginHandler)
	mux.HandleFunc("/api/users/profile", h.ProfileHandler)

	// Сервер
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      corsHandler(mux),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Printf("Starting server on %s", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
