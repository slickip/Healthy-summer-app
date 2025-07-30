package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/config"
	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/db"
	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/handlers"
	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/middleware"
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Address:     "0.0.0.0:8084",
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
	challengeHandler := &handlers.ChallengeHandler{DB: database}

	// Ping
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong from social-service")
	})

	// Routes
	mux.Handle("/api/challenges", middleware.JWTAuth(http.HandlerFunc(challengeHandler.ChallengeHandler)))
	mux.Handle("/api/challenges/join", middleware.JWTAuth(http.HandlerFunc(challengeHandler.JoinChallenge)))
	mux.Handle("/api/challenges/my", middleware.JWTAuth(http.HandlerFunc(challengeHandler.MyChallenges)))
	mux.Handle("/api/challenges/leaderboard", middleware.JWTAuth(http.HandlerFunc(challengeHandler.ChallengeLeaderboard)))

	mux.Handle("/api/messages", middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			challengeHandler.SendMessage(w, r)
		case http.MethodGet:
			challengeHandler.GetMessages(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/social/feed/friends", middleware.JWTAuth(http.HandlerFunc(challengeHandler.FriendsFeed)))

	// Server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      corsHandler(mux),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Printf("Starting social-service on %s", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
