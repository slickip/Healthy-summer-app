package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID int64 `json:"user_id"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "postgres://healthyuser:healthypass@localhost:5433/healthydb?sslmode=disable")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Error checking user", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict) //409
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var newID int64
	err = db.QueryRow(
		`INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING id`,
		req.Email, string(hashedPassword),
	).Scan(&newID)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		UserID: newID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
