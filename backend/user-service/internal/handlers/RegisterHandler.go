package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type RegisterResponse struct {
	UserID       uint   `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		http.Error(w, "Email, password and display name are required", http.StatusBadRequest)
		return
	}

	// Проверка, что пользователь уже есть
	var count int64
	if err := h.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error checking user: %v", err), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error hashing password: %v", err), http.StatusInternalServerError)
		return
	}

	// Создаём пользователя
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	// Генерируем токены для нового пользователя
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, h.JWTConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating access token: %v", err), http.StatusInternalServerError)
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, h.JWTConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.JWTConfig.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
