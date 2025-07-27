package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/config"
	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/db"
	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/handlers"
	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/middleware"
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Address:     "0.0.0.0:8083",
			Timeout:     5 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
	}

	database := db.New()

	mux := http.NewServeMux()

	// CORS middleware
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	// Handlers
	mealHandler := &handlers.MealHandler{DB: database}
	foodHandler := &handlers.FoodHandler{DB: database}
	waterHandler := &handlers.WaterHandler{DB: database}

	// Ping
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong from nutrition-service")
	})

	// Routes
	mux.Handle("/api/meals", middleware.JWTAuth(http.HandlerFunc(mealHandler.Mealhandler)))
	mux.Handle("/api/foods", middleware.JWTAuth(http.HandlerFunc(foodHandler.FoodHandler)))
	mux.Handle("/api/water", middleware.JWTAuth(http.HandlerFunc(waterHandler.WaterHandler)))

	// Server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      corsHandler(mux),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Printf("Starting nutrition-service on %s", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
