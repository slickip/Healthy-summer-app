package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/config"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	JWTConfig config.JWTConfig
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// 1) Находим пользователя по email
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	// 2) Проверяем пароль
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3) Генерируем access токен
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, h.JWTConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating access token: %v", err), http.StatusInternalServerError)
		return
	}

	// 4) Генерируем refresh токен
	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, h.JWTConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	// 5) Возвращаем токены клиенту
	w.Header().Set("Content-Type", "application/json")
	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.JWTConfig.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
	}
	json.NewEncoder(w).Encode(response)
}

// RefreshTokenHandler обрабатывает запросы на обновление токена
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Парсим refresh токен
	claims, err := middleware.ParseToken(req.RefreshToken, h.JWTConfig)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Проверяем, что это refresh токен
	if claims.Type != "refresh" {
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	// Проверяем, что пользователь существует
	var user models.User
	if err := h.DB.First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Генерируем новый access токен
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, h.JWTConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating access token: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем новый access токен
	w.Header().Set("Content-Type", "application/json")
	response := RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(h.JWTConfig.AccessExpiry.Seconds()),
		TokenType:   "Bearer",
	}
	json.NewEncoder(w).Encode(response)
}
