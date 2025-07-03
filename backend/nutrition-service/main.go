package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/config"
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Address:     "0.0.0.0:8083",
			Timeout:     5 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
	}

	//роутер
	mux := http.NewServeMux()

	//простейший эндпоинт
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong from nutrition-service")
	})

	//создаем сервер
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
