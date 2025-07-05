package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/config"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/handlers"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Подключаемся к базе
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "postgres"),
		getEnv("DB_USER", "healthyuser"),
		getEnv("DB_PASSWORD", "healthypass"),
		getEnv("DB_NAME", "healthydb"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Автомиграция
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Создаем структуру с зависимостями
	h := &handlers.Handler{
		DB: db,
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

	// Роуты с зависимостями
	mux.HandleFunc("/api/users/register", h.RegisterHandler)
	mux.HandleFunc("/api/users/login", h.LoginHandler)
	mux.HandleFunc("/api/users/profile", h.ProfileHandler)

	// Сервер
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Printf("Starting server on %s", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

// getEnv возвращает значение переменной окружения или дефолт
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
