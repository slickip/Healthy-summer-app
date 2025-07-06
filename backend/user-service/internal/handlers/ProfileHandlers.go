package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/models"

	"gorm.io/gorm"
)

var jwtSecret = []byte("OMGMYKEY")

type ProfileResponse struct {
	UserID      uint   `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

// ProfileHandler маршрутизирует GET и PUT
func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getProfile(w, r)
	case http.MethodPut:
		h.updateProfile(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getProfile отдает данные профиля
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.ContextUserIDKey)
	if uidVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := uidVal.(uint)

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp := ProfileResponse{
		UserID:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// updateProfile обновляет профиль
func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.ContextUserIDKey)
	if uidVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := uidVal.(uint)

	var req struct {
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.DisplayName) == "" {
		http.Error(w, "Display name cannot be empty", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Обновляем имя
	user.DisplayName = req.DisplayName
	if err := h.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	resp := ProfileResponse{
		UserID:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
