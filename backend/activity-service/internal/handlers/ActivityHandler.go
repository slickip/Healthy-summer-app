package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/slickip/Healthy-summer-app/backend/activity-service/internal/models"
	"gorm.io/gorm"
)

type ActivityHandler struct {
	DB *gorm.DB
}

// Request payload
type CreateActivityRequest struct {
	UserID    uint   `json:"user_id"` // если используешь JWT, можно брать из токена
	Type      string `json:"type"`
	Duration  int    `json:"duration"`
	Intensity string `json:"intensity"`
	Calories  int    `json:"calories"`
	Location  string `json:"location"`
}

func (h *ActivityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	activity := models.Activity{
		UserID:    req.UserID,
		Type:      req.Type,
		Duration:  req.Duration,
		Intensity: req.Intensity,
		Calories:  req.Calories,
		Location:  req.Location,
	}

	if err := h.DB.Create(&activity).Error; err != nil {
		http.Error(w, "Failed to create activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}
